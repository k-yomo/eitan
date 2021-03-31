package sqlutil

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"time"
)

func PtrToNullString(strp *string) sql.NullString {
	if strp != nil {
		return sql.NullString{String: *strp, Valid: true}
	}
	return sql.NullString{}
}

func NullStrToPtr(nullStr sql.NullString) *string {
	if nullStr.Valid {
		str := nullStr.String
		return &str
	}
	return nil
}

func PtrToNullInt64(intp *int) sql.NullInt64 {
	if intp != nil {
		return sql.NullInt64{Int64: int64(*intp), Valid: true}
	}
	return sql.NullInt64{}
}

func NullInt64ToPtr(nullInt64 sql.NullInt64) *int {
	if nullInt64.Valid {
		num := int(nullInt64.Int64)
		return &num
	}
	return nil
}

func NullFloat64ToPtr(nullFloat64 sql.NullFloat64) *float64 {
	if nullFloat64.Valid {
		return &nullFloat64.Float64
	}
	return nil
}

func PtrToFloat64(float64p *float64) sql.NullFloat64 {
	if float64p != nil {
		return sql.NullFloat64{Float64: *float64p, Valid: true}
	}
	return sql.NullFloat64{}
}

func NullTimeToPtr(nullTime mysql.NullTime) *time.Time {
	if nullTime.Valid {
		t := nullTime.Time
		return &t
	}
	return nil
}
func PtrToNullTime(timep *time.Time) mysql.NullTime {
	if timep != nil {
		return mysql.NullTime{Time: *timep, Valid: true}
	}
	return mysql.NullTime{}
}
