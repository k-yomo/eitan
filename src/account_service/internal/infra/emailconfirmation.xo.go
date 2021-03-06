// Package infra contains the types for schema 'accountdb'.
package infra

// GENERATED BY XO. DO NOT EDIT.

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/k-yomo/eitan/src/pkg/tx"
)

// EmailConfirmation represents a row from 'email_confirmations'.
type EmailConfirmation struct {
	Email            string         `db:"email"`             // email
	ConfirmationCode string         `db:"confirmation_code"` // confirmation_code
	ConfirmedAt      mysql.NullTime `db:"confirmed_at"`      // confirmed_at
	CreatedAt        time.Time      `db:"created_at"`        // created_at

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the EmailConfirmation exists in the database.
func (ec *EmailConfirmation) Exists() bool {
	return ec._exists
}

// GetAllEmailConfirmations gets all EmailConfirmations
func GetAllEmailConfirmations(ctx context.Context, db Queryer) ([]*EmailConfirmation, error) {
	// sql query
	const sqlstr = `SELECT ` +
		`email, confirmation_code, confirmed_at, created_at ` +
		`FROM email_confirmations`

	// log and trace
	XOLog(ctx, sqlstr)
	closeSpan := startSQLSpan(ctx, "GetAllEmailConfirmations", sqlstr)
	defer closeSpan()

	var ecs []*EmailConfirmation
	rows, err := db.QueryContext(ctx, sqlstr)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		ec := EmailConfirmation{_exists: true}
		if err := rows.Scan(&ec.Email, &ec.ConfirmationCode, &ec.ConfirmedAt, &ec.CreatedAt); err != nil {
			return nil, err
		}
		ecs = append(ecs, &ec)
	}
	return ecs, nil
}

// GetEmailConfirmation gets a EmailConfirmation by primary key
func GetEmailConfirmation(ctx context.Context, db Queryer, key string) (*EmailConfirmation, error) {
	// sql query
	const sqlstr = `SELECT ` +
		`email, confirmation_code, confirmed_at, created_at ` +
		`FROM email_confirmations ` +
		`WHERE email = ?`

	// log and trace
	XOLog(ctx, sqlstr, key)
	closeSpan := startSQLSpan(ctx, "GetEmailConfirmation", sqlstr, key)
	defer closeSpan()

	ec := EmailConfirmation{_exists: true}
	err := db.QueryRowxContext(ctx, sqlstr, key).Scan(&ec.Email, &ec.ConfirmationCode, &ec.ConfirmedAt, &ec.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &ec, nil
}

// GetEmailConfirmations gets EmailConfirmation list by primary keys
func GetEmailConfirmations(ctx context.Context, db Queryer, keys []string) ([]*EmailConfirmation, error) {
	// sql query
	sqlstr, args, err := sqlx.In(`SELECT `+
		`email, confirmation_code, confirmed_at, created_at `+
		`FROM email_confirmations `+
		`WHERE email IN (?)`, keys)
	if err != nil {
		return nil, err
	}

	// log and trace
	XOLog(ctx, sqlstr, args)
	closeSpan := startSQLSpan(ctx, "GetEmailConfirmations", sqlstr, args)
	defer closeSpan()

	rows, err := db.QueryContext(ctx, sqlstr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// load results
	var res []*EmailConfirmation
	for rows.Next() {
		ec := EmailConfirmation{
			_exists: true,
		}

		// scan
		err = rows.Scan(&ec.Email, &ec.ConfirmationCode, &ec.ConfirmedAt, &ec.CreatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &ec)
	}

	return res, nil
}

func QueryEmailConfirmation(ctx context.Context, q sqlx.QueryerContext, sqlstr string, args ...interface{}) (*EmailConfirmation, error) {
	// log and trace
	XOLog(ctx, sqlstr, args)
	closeSpan := startSQLSpan(ctx, "QueryEmailConfirmation", sqlstr, args)
	defer closeSpan()

	var dest EmailConfirmation
	err := sqlx.GetContext(ctx, q, &dest, sqlstr, args...)
	return &dest, err
}

func QueryEmailConfirmations(ctx context.Context, q sqlx.QueryerContext, sqlstr string, args ...interface{}) ([]*EmailConfirmation, error) {
	// log and trace
	XOLog(ctx, sqlstr, args)
	closeSpan := startSQLSpan(ctx, "QueryEmailConfirmations", sqlstr, args)
	defer closeSpan()

	var dest []*EmailConfirmation
	err := sqlx.SelectContext(ctx, q, &dest, sqlstr, args...)
	return dest, err
}

// Deleted provides information if the EmailConfirmation has been deleted from the database.
func (ec *EmailConfirmation) Deleted() bool {
	return ec._deleted
}

// Insert inserts the EmailConfirmation to the database.
func (ec *EmailConfirmation) Insert(ctx context.Context, db Execer) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	// if already exist, bail
	if ec._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key must be provided
	const sqlstr = `INSERT INTO email_confirmations (` +
		`email, confirmation_code, confirmed_at, created_at` +
		`) VALUES (` +
		`?, ?, ?, ?` +
		`)`

	// log and trace
	XOLog(ctx, sqlstr, ec.Email, ec.ConfirmationCode, ec.ConfirmedAt, ec.CreatedAt)
	closeSpan := startSQLSpan(ctx, "EmailConfirmation_Insert", sqlstr, ec.Email, ec.ConfirmationCode, ec.ConfirmedAt, ec.CreatedAt)
	defer closeSpan()

	// run query
	_, err := db.ExecContext(ctx, sqlstr, ec.Email, ec.ConfirmationCode, ec.ConfirmedAt, ec.CreatedAt)
	if err != nil {
		return err
	}

	// set existence
	ec._exists = true

	return nil
}

// Update updates the EmailConfirmation in the database.
func (ec *EmailConfirmation) Update(ctx context.Context, db Execer) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	// if doesn't exist, bail
	if !ec._exists {
		return errors.New("update failed: does not exist")
	}
	// if deleted, bail
	if ec._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE email_confirmations SET ` +
		`confirmation_code = ?, confirmed_at = ?, created_at = ?` +
		` WHERE email = ?`

	// log and trace
	XOLog(ctx, sqlstr, ec.ConfirmationCode, ec.ConfirmedAt, ec.CreatedAt, ec.Email)
	closeSpan := startSQLSpan(ctx, "EmailConfirmation_Update", sqlstr, ec.ConfirmationCode, ec.ConfirmedAt, ec.CreatedAt, ec.Email)
	defer closeSpan()

	// run query
	_, err := db.ExecContext(ctx, sqlstr, ec.ConfirmationCode, ec.ConfirmedAt, ec.CreatedAt, ec.Email)
	return err
}

// Delete deletes the EmailConfirmation from the database.
func (ec *EmailConfirmation) Delete(ctx context.Context, db Execer) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	// if doesn't exist, bail
	if !ec._exists {
		return nil
	}

	// if deleted, bail
	if ec._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM email_confirmations WHERE email = ?`

	// log and trace
	XOLog(ctx, sqlstr, ec.Email)
	closeSpan := startSQLSpan(ctx, "{ .Name }}_Delete", sqlstr, ec.Email)
	defer closeSpan()

	// run query
	_, err := db.ExecContext(ctx, sqlstr, ec.Email)
	if err != nil {
		return err
	}

	// set deleted
	ec._deleted = true

	return nil
}

// InsertOrUpdate inserts or updates the EmailConfirmation to the database.
func (ec *EmailConfirmation) InsertOrUpdate(ctx context.Context, db Executor) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	_, err := GetEmailConfirmation(ctx, db, ec.Email)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == sql.ErrNoRows {
		return ec.Insert(ctx, db)
	} else {
		ec._exists = true
		return ec.Update(ctx, db)
	}
}

// InsertOrUpdate inserts or updates the EmailConfirmation to the database.
func (ec *EmailConfirmation) InsertIfNotExist(ctx context.Context, db Executor) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	_, err := GetEmailConfirmation(ctx, db, ec.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return ec.Insert(ctx, db)
		}
		return err
	}

	return nil
}

// GetEmailConfirmationByConfirmationCode retrieves a row from 'email_confirmations' as a EmailConfirmation.
// Generated from index 'confirmation_code'.
func GetEmailConfirmationByConfirmationCode(ctx context.Context, db Queryer, confirmationCode string) (*EmailConfirmation, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`email, confirmation_code, confirmed_at, created_at ` +
		`FROM email_confirmations ` +
		`WHERE confirmation_code = ?`

	// log and trace
	XOLog(ctx, sqlstr, confirmationCode)
	closeSpan := startSQLSpan(ctx, "EmailConfirmationByConfirmationCode", sqlstr, confirmationCode)
	defer closeSpan()
	ec := EmailConfirmation{
		_exists: true,
	}

	// run query
	err = db.QueryRowxContext(ctx, sqlstr, confirmationCode).Scan(&ec.Email, &ec.ConfirmationCode, &ec.ConfirmedAt, &ec.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &ec, nil
}

// GetEmailConfirmationByEmailConfirmationCode retrieves a row from 'email_confirmations' as a EmailConfirmation.
// Generated from index 'confirmation_code_idx'.
func GetEmailConfirmationByEmailConfirmationCode(ctx context.Context, db Queryer, email string, confirmationCode string) (*EmailConfirmation, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`email, confirmation_code, confirmed_at, created_at ` +
		`FROM email_confirmations ` +
		`WHERE email = ? AND confirmation_code = ?`

	// log and trace
	XOLog(ctx, sqlstr, email, confirmationCode)
	closeSpan := startSQLSpan(ctx, "EmailConfirmationByEmailConfirmationCode", sqlstr, email, confirmationCode)
	defer closeSpan()
	ec := EmailConfirmation{
		_exists: true,
	}

	// run query
	err = db.QueryRowxContext(ctx, sqlstr, email, confirmationCode).Scan(&ec.Email, &ec.ConfirmationCode, &ec.ConfirmedAt, &ec.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &ec, nil
}

// GetEmailConfirmationByEmail retrieves a row from 'email_confirmations' as a EmailConfirmation.
// Generated from index 'email_confirmations_email_pkey'.
func GetEmailConfirmationByEmail(ctx context.Context, db Queryer, email string) (*EmailConfirmation, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`email, confirmation_code, confirmed_at, created_at ` +
		`FROM email_confirmations ` +
		`WHERE email = ?`

	// log and trace
	XOLog(ctx, sqlstr, email)
	closeSpan := startSQLSpan(ctx, "EmailConfirmationByEmail", sqlstr, email)
	defer closeSpan()
	ec := EmailConfirmation{
		_exists: true,
	}

	// run query
	err = db.QueryRowxContext(ctx, sqlstr, email).Scan(&ec.Email, &ec.ConfirmationCode, &ec.ConfirmedAt, &ec.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &ec, nil
}
