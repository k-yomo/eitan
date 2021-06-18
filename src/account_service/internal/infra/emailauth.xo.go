// Package infra contains the types for schema 'accountdb'.
package infra

// GENERATED BY XO. DO NOT EDIT.

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/k-yomo/eitan/src/pkg/tx"
)

// EmailAuth represents a row from 'email_auth'.
type EmailAuth struct {
	UserID         string    `db:"user_id"`         // user_id
	Email          string    `db:"email"`           // email
	PasswordDigest string    `db:"password_digest"` // password_digest
	CreatedAt      time.Time `db:"created_at"`      // created_at
	UpdatedAt      time.Time `db:"updated_at"`      // updated_at

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the EmailAuth exists in the database.
func (ea *EmailAuth) Exists() bool {
	return ea._exists
}

// GetAllEmailAuths gets all EmailAuths
func GetAllEmailAuths(ctx context.Context, db Queryer) ([]*EmailAuth, error) {
	// sql query
	const sqlstr = `SELECT ` +
		`user_id, email, password_digest, created_at, updated_at ` +
		`FROM email_auth`

	// log and trace
	XOLog(ctx, sqlstr)
	closeSpan := startSQLSpan(ctx, "GetAllEmailAuths", sqlstr)
	defer closeSpan()

	var eas []*EmailAuth
	rows, err := db.QueryContext(ctx, sqlstr)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		ea := EmailAuth{_exists: true}
		if err := rows.Scan(&ea.UserID, &ea.Email, &ea.PasswordDigest, &ea.CreatedAt, &ea.UpdatedAt); err != nil {
			return nil, err
		}
		eas = append(eas, &ea)
	}
	return eas, nil
}

// GetEmailAuth gets a EmailAuth by primary key
func GetEmailAuth(ctx context.Context, db Queryer, key string) (*EmailAuth, error) {
	// sql query
	const sqlstr = `SELECT ` +
		`user_id, email, password_digest, created_at, updated_at ` +
		`FROM email_auth ` +
		`WHERE user_id = ?`

	// log and trace
	XOLog(ctx, sqlstr, key)
	closeSpan := startSQLSpan(ctx, "GetEmailAuth", sqlstr, key)
	defer closeSpan()

	ea := EmailAuth{_exists: true}
	err := db.QueryRowxContext(ctx, sqlstr, key).Scan(&ea.UserID, &ea.Email, &ea.PasswordDigest, &ea.CreatedAt, &ea.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &ea, nil
}

// GetEmailAuths gets EmailAuth list by primary keys
func GetEmailAuths(ctx context.Context, db Queryer, keys []string) ([]*EmailAuth, error) {
	// sql query
	sqlstr, args, err := sqlx.In(`SELECT `+
		`user_id, email, password_digest, created_at, updated_at `+
		`FROM email_auth `+
		`WHERE user_id IN (?)`, keys)
	if err != nil {
		return nil, err
	}

	// log and trace
	XOLog(ctx, sqlstr, args)
	closeSpan := startSQLSpan(ctx, "GetEmailAuths", sqlstr, args)
	defer closeSpan()

	rows, err := db.QueryContext(ctx, sqlstr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// load results
	var res []*EmailAuth
	for rows.Next() {
		ea := EmailAuth{
			_exists: true,
		}

		// scan
		err = rows.Scan(&ea.UserID, &ea.Email, &ea.PasswordDigest, &ea.CreatedAt, &ea.UpdatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &ea)
	}

	return res, nil
}

func QueryEmailAuth(ctx context.Context, q sqlx.QueryerContext, sqlstr string, args ...interface{}) (*EmailAuth, error) {
	// log and trace
	XOLog(ctx, sqlstr, args)
	closeSpan := startSQLSpan(ctx, "QueryEmailAuth", sqlstr, args)
	defer closeSpan()

	var dest EmailAuth
	err := sqlx.GetContext(ctx, q, &dest, sqlstr, args...)
	return &dest, err
}

func QueryEmailAuths(ctx context.Context, q sqlx.QueryerContext, sqlstr string, args ...interface{}) ([]*EmailAuth, error) {
	// log and trace
	XOLog(ctx, sqlstr, args)
	closeSpan := startSQLSpan(ctx, "QueryEmailAuths", sqlstr, args)
	defer closeSpan()

	var dest []*EmailAuth
	err := sqlx.SelectContext(ctx, q, &dest, sqlstr, args...)
	return dest, err
}

// Deleted provides information if the EmailAuth has been deleted from the database.
func (ea *EmailAuth) Deleted() bool {
	return ea._deleted
}

// Insert inserts the EmailAuth to the database.
func (ea *EmailAuth) Insert(ctx context.Context, db Execer) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	// if already exist, bail
	if ea._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key must be provided
	const sqlstr = `INSERT INTO email_auth (` +
		`user_id, email, password_digest, created_at, updated_at` +
		`) VALUES (` +
		`?, ?, ?, ?, ?` +
		`)`

	// log and trace
	XOLog(ctx, sqlstr, ea.UserID, ea.Email, ea.PasswordDigest, ea.CreatedAt, ea.UpdatedAt)
	closeSpan := startSQLSpan(ctx, "EmailAuth_Insert", sqlstr, ea.UserID, ea.Email, ea.PasswordDigest, ea.CreatedAt, ea.UpdatedAt)
	defer closeSpan()

	// run query
	_, err := db.ExecContext(ctx, sqlstr, ea.UserID, ea.Email, ea.PasswordDigest, ea.CreatedAt, ea.UpdatedAt)
	if err != nil {
		return err
	}

	// set existence
	ea._exists = true

	return nil
}

// Update updates the EmailAuth in the database.
func (ea *EmailAuth) Update(ctx context.Context, db Execer) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	// if doesn't exist, bail
	if !ea._exists {
		return errors.New("update failed: does not exist")
	}
	// if deleted, bail
	if ea._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE email_auth SET ` +
		`email = ?, password_digest = ?, created_at = ?, updated_at = ?` +
		` WHERE user_id = ?`

	// log and trace
	XOLog(ctx, sqlstr, ea.Email, ea.PasswordDigest, ea.CreatedAt, ea.UpdatedAt, ea.UserID)
	closeSpan := startSQLSpan(ctx, "EmailAuth_Update", sqlstr, ea.Email, ea.PasswordDigest, ea.CreatedAt, ea.UpdatedAt, ea.UserID)
	defer closeSpan()

	// run query
	_, err := db.ExecContext(ctx, sqlstr, ea.Email, ea.PasswordDigest, ea.CreatedAt, ea.UpdatedAt, ea.UserID)
	return err
}

// Delete deletes the EmailAuth from the database.
func (ea *EmailAuth) Delete(ctx context.Context, db Execer) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	// if doesn't exist, bail
	if !ea._exists {
		return nil
	}

	// if deleted, bail
	if ea._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM email_auth WHERE user_id = ?`

	// log and trace
	XOLog(ctx, sqlstr, ea.UserID)
	closeSpan := startSQLSpan(ctx, "{ .Name }}_Delete", sqlstr, ea.UserID)
	defer closeSpan()

	// run query
	_, err := db.ExecContext(ctx, sqlstr, ea.UserID)
	if err != nil {
		return err
	}

	// set deleted
	ea._deleted = true

	return nil
}

// InsertOrUpdate inserts or updates the EmailAuth to the database.
func (ea *EmailAuth) InsertOrUpdate(ctx context.Context, db Executor) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	_, err := GetEmailAuth(ctx, db, ea.UserID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == sql.ErrNoRows {
		return ea.Insert(ctx, db)
	} else {
		ea._exists = true
		return ea.Update(ctx, db)
	}
}

// InsertOrUpdate inserts or updates the EmailAuth to the database.
func (ea *EmailAuth) InsertIfNotExist(ctx context.Context, db Executor) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	_, err := GetEmailAuth(ctx, db, ea.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ea.Insert(ctx, db)
		}
		return err
	}

	return nil
}

// User returns the User associated with the EmailAuth's UserID (user_id).
//
// Generated from foreign key 'email_auth_ibfk_1'.
func (ea *EmailAuth) User(ctx context.Context, db Executor) (*User, error) {
	return GetUserByID(ctx, db, ea.UserID)
}

// GetEmailAuthByEmail retrieves a row from 'email_auth' as a EmailAuth.
// Generated from index 'email'.
func GetEmailAuthByEmail(ctx context.Context, db Queryer, email string) (*EmailAuth, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`user_id, email, password_digest, created_at, updated_at ` +
		`FROM email_auth ` +
		`WHERE email = ?`

	// log and trace
	XOLog(ctx, sqlstr, email)
	closeSpan := startSQLSpan(ctx, "EmailAuthByEmail", sqlstr, email)
	defer closeSpan()
	ea := EmailAuth{
		_exists: true,
	}

	// run query
	err = db.QueryRowxContext(ctx, sqlstr, email).Scan(&ea.UserID, &ea.Email, &ea.PasswordDigest, &ea.CreatedAt, &ea.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &ea, nil
}

// GetEmailAuthByUserID retrieves a row from 'email_auth' as a EmailAuth.
// Generated from index 'email_auth_user_id_pkey'.
func GetEmailAuthByUserID(ctx context.Context, db Queryer, userID string) (*EmailAuth, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`user_id, email, password_digest, created_at, updated_at ` +
		`FROM email_auth ` +
		`WHERE user_id = ?`

	// log and trace
	XOLog(ctx, sqlstr, userID)
	closeSpan := startSQLSpan(ctx, "EmailAuthByUserID", sqlstr, userID)
	defer closeSpan()
	ea := EmailAuth{
		_exists: true,
	}

	// run query
	err = db.QueryRowxContext(ctx, sqlstr, userID).Scan(&ea.UserID, &ea.Email, &ea.PasswordDigest, &ea.CreatedAt, &ea.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &ea, nil
}
