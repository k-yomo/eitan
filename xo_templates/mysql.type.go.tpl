{{- $short := (shortname .Name "err" "res" "sqlstr" "db" "XOLog") -}}
{{- $table := (.Table.TableName) -}}
{{- if .Comment -}}
// {{ .Comment }}
{{- else -}}
// {{ .Name }} represents a row from '{{ $table }}'.
{{- end }}
type {{ .Name }} struct {
{{- range .Fields }}
	{{ .Name }} {{ retype .Type }} `db:"{{ .Col.ColumnName }}"` // {{ .Col.ColumnName }}
{{- end }}
{{- if .PrimaryKey }}

	// xo fields
	_exists, _deleted bool
{{ end }}
}

{{ if .PrimaryKey }}
// Exists determines if the {{ .Name }} exists in the database.
func ({{ $short }} *{{ .Name }}) Exists() bool {
	return {{ $short }}._exists
}

// GetAll{{ .Name }}s gets all {{ .Name }}s
func GetAll{{ .Name }}s (ctx context.Context, db Queryer) ([]*{{ .Name }}, error) {
	// sql query
	const sqlstr = `SELECT ` +
		`{{ colnames .Fields }} ` +
		`FROM {{ $table }}`

	// log and trace
	XOLog(ctx, sqlstr)
    closeSpan := startSQLSpan(ctx, "GetAll{{ .Name }}s", sqlstr)
    defer closeSpan()

	var {{ $short }}s []*{{ .Name }}
	rows, err := db.QueryContext(ctx, sqlstr)
	if err != nil {
	    return nil, err
	}

	for rows.Next() {
	    {{ $short }} := {{ .Name }}{_exists: true}
	    if err := rows.Scan({{ fieldnames .Fields (print "&" $short) }}); err != nil {
	        return nil, err
	    }
	    {{ $short }}s = append({{ $short }}s, &{{ $short }})
	}
	return {{ $short }}s, nil
}

// Get{{ .Name }} gets a {{ .Name }} by primary key
func Get{{ .Name }}(ctx context.Context, db Queryer, key {{ .PrimaryKey.Type }}) (*{{ .Name }}, error) {
	// sql query
	const sqlstr = `SELECT ` +
		`{{ colnames .Fields }} ` +
		`FROM {{ $table }} ` +
		`WHERE {{ colname .PrimaryKey.Col }} = ?`

	// log and trace
	XOLog(ctx, sqlstr, key)
	closeSpan := startSQLSpan(ctx, "Get{{ .Name }}", sqlstr, key)
	defer closeSpan()

	{{ $short }} := {{ .Name }}{_exists: true}
	err := db.QueryRowxContext(ctx, sqlstr, key).Scan({{ fieldnames .Fields (print "&" $short) }})
	if err != nil {
		return nil, err
	}
	return &{{ $short }}, nil
}

// Get{{ .Name }}s gets {{ .Name }} list by primary keys
func Get{{ .Name }}s(ctx context.Context, db Queryer, keys []{{ .PrimaryKey.Type }}) ([]*{{ .Name }}, error) {
	// sql query
	sqlstr, args, err := sqlx.In(`SELECT ` +
		`{{ colnames .Fields }} ` +
		`FROM {{ $table }} ` +
		`WHERE {{ colname .PrimaryKey.Col }} IN (?)`, keys)
	if err != nil {
		return nil, err
	}

	// log and trace
	XOLog(ctx, sqlstr, args)
	closeSpan := startSQLSpan(ctx, "Get{{ .Name }}s", sqlstr, args)
	defer closeSpan()

	rows, err := db.QueryContext(ctx, sqlstr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// load results
	var res []*{{ .Name }}
	for rows.Next() {
		{{ $short }} := {{ .Name }}{
			_exists: true,
		}

		// scan
		err = rows.Scan({{ fieldnames .Fields (print "&" $short) }})
		if err != nil {
			return nil, err
		}

		res = append(res, &{{ $short }})
	}

	return res, nil
}

func Query{{ .Name }}(ctx context.Context, q sqlx.QueryerContext, sqlstr string, args ...interface{}) (*{{ .Name }}, error) {
	// log and trace
	XOLog(ctx, sqlstr, args)
	closeSpan := startSQLSpan(ctx, "Query{{ .Name }}", sqlstr, args)
	defer closeSpan()

    var dest {{.Name }}
    err := sqlx.GetContext(ctx, q, &dest, sqlstr, args...)
    return &dest, err
}

func Query{{ .Name }}s(ctx context.Context, q sqlx.QueryerContext, sqlstr string, args ...interface{}) ([]*{{ .Name }}, error) {
	// log and trace
	XOLog(ctx, sqlstr, args)
	closeSpan := startSQLSpan(ctx, "Query{{ .Name }}s", sqlstr, args)
	defer closeSpan()

    var dest []*{{.Name }}
    err := sqlx.SelectContext(ctx, q, &dest, sqlstr, args...)
    return dest, err
}

// Deleted provides information if the {{ .Name }} has been deleted from the database.
func ({{ $short }} *{{ .Name }}) Deleted() bool {
	return {{ $short }}._deleted
}

// Insert inserts the {{ .Name }} to the database.
func ({{ $short }} *{{ .Name }}) Insert(ctx context.Context, db Execer) error {
     if t, ok := tx.GetTx(ctx); ok {
      db = t
    }
	// if already exist, bail
	if {{ $short }}._exists {
		return errors.New("insert failed: already exists")
	}


{{ if .Table.ManualPk  }}
	// sql insert query, primary key must be provided
	const sqlstr = `INSERT INTO {{ $table }} (` +
		`{{ colnames .Fields }}` +
		`) VALUES (` +
		`{{ colvals .Fields }}` +
		`)`

	// log and trace
	XOLog(ctx, sqlstr, {{ fieldnames .Fields $short }})
	closeSpan := startSQLSpan(ctx, "{{ .Name }}_Insert", sqlstr, {{ fieldnames .Fields $short }})
	defer closeSpan()

	// run query
	_, err := db.ExecContext(ctx, sqlstr, {{ fieldnames .Fields $short }})
	if err != nil {
		return err
	}

	// set existence
	{{ $short }}._exists = true
{{ else }}
	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO {{ $table }} (` +
		`{{ colnames .Fields .PrimaryKey.Name }}` +
		`) VALUES (` +
		`{{ colvals .Fields .PrimaryKey.Name }}` +
		`)`

	// log and trace
	XOLog(ctx, sqlstr, {{ fieldnames .Fields $short .PrimaryKey.Name }})
	closeSpan := startSQLSpan(ctx, "{{ .Name }}_Insert", sqlstr, {{ fieldnames .Fields $short .PrimaryKey.Name }})
	defer closeSpan()

	// run query
	res, err := db.ExecContext(ctx, sqlstr, {{ fieldnames .Fields $short .PrimaryKey.Name }})
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	{{ $short }}.{{ .PrimaryKey.Name }} = {{ .PrimaryKey.Type }}(id)
	{{ $short }}._exists = true
{{ end }}

	return nil
}

{{ if ne (fieldnamesmulti .Fields $short .PrimaryKeyFields) "" }}
	// Update updates the {{ .Name }} in the database.
	func ({{ $short }} *{{ .Name }}) Update(ctx context.Context, db Execer) error {
        if t, ok := tx.GetTx(ctx); ok {
            db = t
        }
		// if doesn't exist, bail
		if !{{ $short }}._exists {
			return errors.New("update failed: does not exist")
		}
		// if deleted, bail
		if {{ $short }}._deleted {
			return errors.New("update failed: marked for deletion")
		}
		{{ if gt ( len .PrimaryKeyFields ) 1 }}
			// sql query with composite primary key
			const sqlstr = `UPDATE {{ $table }} SET ` +
				`{{ colnamesquerymulti .Fields ", " 0 .PrimaryKeyFields }}` +
				` WHERE {{ colnamesquery .PrimaryKeyFields " AND " }}`

			// log and trace
			XOLog(ctx, sqlstr, {{ fieldnamesmulti .Fields $short .PrimaryKeyFields }}, {{ fieldnames .PrimaryKeyFields $short}})
			closeSpan := startSQLSpan(ctx, "{{ .Name }}_Update", sqlstr, {{ fieldnamesmulti .Fields $short .PrimaryKeyFields }}, {{ fieldnames .PrimaryKeyFields $short}})
            defer closeSpan()

			// run query
			_, err := db.ExecContext(ctx, sqlstr, {{ fieldnamesmulti .Fields $short .PrimaryKeyFields }}, {{ fieldnames .PrimaryKeyFields $short}})
			return err
		{{- else }}
			// sql query
			const sqlstr = `UPDATE {{ $table }} SET ` +
				`{{ colnamesquery .Fields ", " .PrimaryKey.Name }}` +
				` WHERE {{ colname .PrimaryKey.Col }} = ?`

			// log and trace
			XOLog(ctx, sqlstr, {{ fieldnames .Fields $short .PrimaryKey.Name }}, {{ $short }}.{{ .PrimaryKey.Name }})
			closeSpan := startSQLSpan(ctx, "{{ .Name }}_Update", sqlstr, {{ fieldnames .Fields $short .PrimaryKey.Name }}, {{ $short }}.{{ .PrimaryKey.Name }})
            defer closeSpan()

			// run query
			_, err := db.ExecContext(ctx, sqlstr, {{ fieldnames .Fields $short .PrimaryKey.Name }}, {{ $short }}.{{ .PrimaryKey.Name }})
			return err
		{{- end }}
	}
{{ else }}
	// Update statements omitted due to lack of fields other than primary key
{{ end }}
// Delete deletes the {{ .Name }} from the database.
func ({{ $short }} *{{ .Name }}) Delete(ctx context.Context, db Execer) error {
     if t, ok := tx.GetTx(ctx); ok {
      db = t
    }
	// if doesn't exist, bail
	if !{{ $short }}._exists {
		return nil
	}

	// if deleted, bail
	if {{ $short }}._deleted {
		return nil
	}

	{{ if gt ( len .PrimaryKeyFields ) 1 }}
		// sql query with composite primary key
		const sqlstr = `DELETE FROM {{ $table }} WHERE {{ colnamesquery .PrimaryKeyFields " AND " }}`

		// log and trace
		XOLog(ctx, sqlstr, {{ fieldnames .PrimaryKeyFields $short }})
        closeSpan := startSQLSpan(ctx, "{ .Name }}Delete", sqlstr, {{ fieldnames .PrimaryKeyFields $short }})
        defer closeSpan()

		// run query
		_, err := db.ExecContext(ctx, sqlstr, {{ fieldnames .PrimaryKeyFields $short }})
		if err != nil {
			return err
		}
	{{- else }}
		// sql query
		const sqlstr = `DELETE FROM {{ $table }} WHERE {{ colname .PrimaryKey.Col }} = ?`

		// log and trace
		XOLog(ctx, sqlstr, {{ $short }}.{{ .PrimaryKey.Name }})
        closeSpan := startSQLSpan(ctx, "{ .Name }}_Delete", sqlstr, {{ $short }}.{{ .PrimaryKey.Name }})
        defer closeSpan()

		// run query
		_, err := db.ExecContext(ctx, sqlstr, {{ $short }}.{{ .PrimaryKey.Name }})
		if err != nil {
			return err
		}
	{{- end }}

	// set deleted
	{{ $short }}._deleted = true

	return nil
}
{{- end }}

// InsertOrUpdate inserts or updates the {{ .Name }} to the database.
func ({{ $short }} *{{ .Name }}) InsertOrUpdate(ctx context.Context, db Executor) error {
     if t, ok := tx.GetTx(ctx); ok {
      db = t
    }
    _, err := Get{{ .Name }}(ctx, db, {{ $short }}.{{ .PrimaryKey.Name }})
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == sql.ErrNoRows {
		return {{ $short }}.Insert(ctx, db)
	} else {
        {{ $short }}._exists = true
		return {{ $short }}.Update(ctx, db)
	}
}

// InsertOrUpdate inserts or updates the {{ .Name }} to the database.
func ({{ $short }} *{{ .Name }}) InsertIfNotExist(ctx context.Context, db Executor) error {
     if t, ok := tx.GetTx(ctx); ok {
      db = t
    }
    _, err := Get{{ .Name }}(ctx, db, {{ $short }}.{{ .PrimaryKey.Name }})
    if err != nil {
        if err == sql.ErrNoRows {
            return {{ $short }}.Insert(ctx, db)
        }
        return err
    }

	return nil
}


