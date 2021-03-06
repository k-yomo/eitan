// Executor is the common interface for database operations that can be used with
// types from schema '{{ schema .Schema }}'.
//
// This should work with database/sql.DB and database/sql.Tx.
type Executor interface {
	Queryer
	Execer
}

type Queryer interface {
	sqlx.Queryer
	sqlx.QueryerContext
}

type Execer interface {
	sqlx.Execer
	sqlx.ExecerContext
}

// XOLog provides the log func used by generated queries.
var XOLog = func(ctx context.Context, sqlstr string, params ...interface{}) { }

// ScannerValuer is the common interface for types that implement both the
// database/sql.Scanner and sql/driver.Valuer interfaces.
type ScannerValuer interface {
	sql.Scanner
	driver.Valuer
}

// StringSlice is a slice of strings.
type StringSlice []string

// quoteEscapeRegex is the regex to match escaped characters in a string.
var quoteEscapeRegex = regexp.MustCompile(`([^\\]([\\]{2})*)\\"`)

// Scan satisfies the sql.Scanner interface for StringSlice.
func (ss *StringSlice) Scan(src interface{}) error {
	buf, ok := src.([]byte)
	if !ok {
		return errors.New("invalid StringSlice")
	}

	// change quote escapes for csv parser
	str := quoteEscapeRegex.ReplaceAllString(string(buf), `$1""`)
	str = strings.Replace(str, `\\`, `\`, -1)

	// remove braces
	str = str[1:len(str)-1]

	// bail if only one
	if len(str) == 0 {
		*ss = StringSlice([]string{})
		return nil
	}

	// parse with csv reader
	cr := csv.NewReader(strings.NewReader(str))
	slice, err := cr.Read()
	if err != nil {
		return err
	}

	*ss = StringSlice(slice)

	return nil
}

// Value satisfies the driver.Valuer interface for StringSlice.
func (ss StringSlice) Value() (driver.Value, error) {
	v := make([]string, len(ss))
	for i, s := range ss {
		v[i] = `"` + strings.Replace(strings.Replace(s, `\`, `\\\`, -1), `"`, `\"`, -1) + `"`
	}
	return "{" + strings.Join(v, ",") + "}", nil
}

// Slice is a slice of ScannerValuers.
type Slice []ScannerValuer

// operation is an operation to database either QUERY or EXEC
type operation int8

const (
	// Exec operation for like sql.DB.Exec()
	Exec	operation = 1
	// Query operation for sql.DB.Query()
	Query	operation = 2
)

func newXOError(op operation, method, table string, err error) error {
	var errcode uint16
	myerr, ok := err.(*mysql.MySQLError)
	if ok {
		errcode = myerr.Number
	}
	xoerr := xoError{
		op:        op,
		method:    method,
		table:     table,
		err:       err,
		myErrCode: errcode,
	}

	if err == sql.ErrNoRows {
		return &xoNotFoundError{xoError: xoerr}
	}

	return &xoerr
}

type xoError struct {
	err       error
	op        operation
	method    string
	table     string
	myErrCode uint16
}

func (e xoError) Error() string {
	return fmt.Sprintf("xo error in %s(%s): %v", e.method, e.table, e.err)
}

func (e xoError) DBErrorCode() uint16 {
	return e.myErrCode
}

func (e xoError) DBError() bool {
	return true
}

func (e xoError) DBOperation() string {
	switch e.op {
	case Exec:
		return "EXEC"
	case Query:
		return "QUERY"
	default:
		return "UNKNOWN"
	}
}

func (e xoError) DBTableName() string {
	return e.table
}

func (e xoError) RawError() error {
	return e.err
}

func (e xoError) Timeout() bool {
	type timeout interface {
		Timeout() bool
	}
	if t, ok := e.err.(timeout); ok {
		return t.Timeout()
	}
	return false
}

type xoNotFoundError struct {
	xoError
}

func (e xoNotFoundError) NotFound() bool {
	return true
}

func startSQLSpan(ctx context.Context, funcName, sqlstr string, params ...interface{}) (closeSpan func()) {
	tracer := otel.Tracer("sql")
	ctx, span := tracer.Start(
		ctx,
		funcName,
		trace.WithSpanKind(trace.SpanKindClient),
	)
	span.SetAttributes(
		attribute.String("sql.query", sqlstr),
		attribute.Array("sql.params", params),
	)
	return func() { span.End() }
}
