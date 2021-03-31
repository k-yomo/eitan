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

// Account represents a row from 'accounts'.
type Account struct {
	ID           string         `db:"id"`             // id
	Provider     string         `db:"provider"`       // provider
	Email        string         `db:"email"`          // email
	DisplayName  string         `db:"display_name"`   // display_name
	ScreenImgURL sql.NullString `db:"screen_img_url"` // screen_img_url
	CreatedAt    time.Time      `db:"created_at"`     // created_at
	UpdatedAt    time.Time      `db:"updated_at"`     // updated_at

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the Account exists in the database.
func (a *Account) Exists() bool {
	return a._exists
}

// GetAllAccounts gets all Accounts
func GetAllAccounts(ctx context.Context, db Queryer) ([]*Account, error) {
	// sql query
	const sqlstr = `SELECT ` +
		`id, provider, email, display_name, screen_img_url, created_at, updated_at ` +
		`FROM accounts`

	// log and trace
	XOLog(ctx, sqlstr)
	closeSpan := startSQLSpan(ctx, "GetAllAccounts", sqlstr)
	defer closeSpan()

	var as []*Account
	rows, err := db.QueryContext(ctx, sqlstr)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		a := Account{_exists: true}
		if err := rows.Scan(&a.ID, &a.Provider, &a.Email, &a.DisplayName, &a.ScreenImgURL, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		as = append(as, &a)
	}
	return as, nil
}

// GetAccount gets a Account by primary key
func GetAccount(ctx context.Context, db Queryer, key string) (*Account, error) {
	// sql query
	const sqlstr = `SELECT ` +
		`id, provider, email, display_name, screen_img_url, created_at, updated_at ` +
		`FROM accounts ` +
		`WHERE id = ?`

	// log and trace
	XOLog(ctx, sqlstr, key)
	startSQLSpan(ctx, "GetAccount", sqlstr, key)
	a := Account{_exists: true}
	err := db.QueryRowxContext(ctx, sqlstr, key).Scan(&a.ID, &a.Provider, &a.Email, &a.DisplayName, &a.ScreenImgURL, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// GetAccounts gets Account list by primary keys
func GetAccounts(ctx context.Context, db Queryer, keys []string) ([]*Account, error) {
	// sql query
	sqlstr, args, err := sqlx.In(`SELECT `+
		`id, provider, email, display_name, screen_img_url, created_at, updated_at `+
		`FROM accounts `+
		`WHERE id IN (?)`, keys)
	if err != nil {
		return nil, err
	}

	// log and trace
	XOLog(ctx, sqlstr, args)
	startSQLSpan(ctx, "GetAccounts", sqlstr, args)

	rows, err := db.QueryContext(ctx, sqlstr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// load results
	var res []*Account
	for rows.Next() {
		a := Account{
			_exists: true,
		}

		// scan
		err = rows.Scan(&a.ID, &a.Provider, &a.Email, &a.DisplayName, &a.ScreenImgURL, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &a)
	}

	return res, nil
}

// Deleted provides information if the Account has been deleted from the database.
func (a *Account) Deleted() bool {
	return a._deleted
}

// Insert inserts the Account to the database.
func (a *Account) Insert(ctx context.Context, db Execer) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	// if already exist, bail
	if a._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key must be provided
	const sqlstr = `INSERT INTO accounts (` +
		`id, provider, email, display_name, screen_img_url, created_at, updated_at` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?` +
		`)`

	// log and trace
	XOLog(ctx, sqlstr, a.ID, a.Provider, a.Email, a.DisplayName, a.ScreenImgURL, a.CreatedAt, a.UpdatedAt)
	closeSpan := startSQLSpan(ctx, "Account_Insert", sqlstr, a.ID, a.Provider, a.Email, a.DisplayName, a.ScreenImgURL, a.CreatedAt, a.UpdatedAt)
	defer closeSpan()

	// run query
	_, err := db.ExecContext(ctx, sqlstr, a.ID, a.Provider, a.Email, a.DisplayName, a.ScreenImgURL, a.CreatedAt, a.UpdatedAt)
	if err != nil {
		return err
	}

	// set existence
	a._exists = true

	return nil
}

// Update updates the Account in the database.
func (a *Account) Update(ctx context.Context, db Execer) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	// if doesn't exist, bail
	if !a._exists {
		return errors.New("update failed: does not exist")
	}
	// if deleted, bail
	if a._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE accounts SET ` +
		`provider = ?, email = ?, display_name = ?, screen_img_url = ?, created_at = ?, updated_at = ?` +
		` WHERE id = ?`

	// log and trace
	XOLog(ctx, sqlstr, a.Provider, a.Email, a.DisplayName, a.ScreenImgURL, a.CreatedAt, a.UpdatedAt, a.ID)
	closeSpan := startSQLSpan(ctx, "Account_Update", sqlstr, a.Provider, a.Email, a.DisplayName, a.ScreenImgURL, a.CreatedAt, a.UpdatedAt, a.ID)
	defer closeSpan()

	// run query
	_, err := db.ExecContext(ctx, sqlstr, a.Provider, a.Email, a.DisplayName, a.ScreenImgURL, a.CreatedAt, a.UpdatedAt, a.ID)
	return err
}

// Delete deletes the Account from the database.
func (a *Account) Delete(ctx context.Context, db Execer) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	// if doesn't exist, bail
	if !a._exists {
		return nil
	}

	// if deleted, bail
	if a._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM accounts WHERE id = ?`

	// log and trace
	XOLog(ctx, sqlstr, a.ID)
	closeSpan := startSQLSpan(ctx, "{ .Name }}_Delete", sqlstr, a.ID)
	defer closeSpan()

	// run query
	_, err := db.ExecContext(ctx, sqlstr, a.ID)
	if err != nil {
		return err
	}

	// set deleted
	a._deleted = true

	return nil
}

// InsertOrUpdate inserts or updates the Account to the database.
func (a *Account) InsertOrUpdate(ctx context.Context, db Executor) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	_, err := GetAccount(ctx, db, a.ID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == sql.ErrNoRows {
		return a.Insert(ctx, db)
	} else {
		a._exists = true
		return a.Update(ctx, db)
	}
}

// InsertOrUpdate inserts or updates the Account to the database.
func (a *Account) InsertIfNotExist(ctx context.Context, db Executor) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	_, err := GetAccount(ctx, db, a.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return a.Insert(ctx, db)
		}
		return err
	}

	return nil
}

// AccountByID retrieves a row from 'accounts' as a Account.
// Generated from index 'accounts_id_pkey'.
func AccountByID(ctx context.Context, db Queryer, id string) (*Account, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, provider, email, display_name, screen_img_url, created_at, updated_at ` +
		`FROM accounts ` +
		`WHERE id = ?`

	// log and trace
	XOLog(ctx, sqlstr, id)
	closeSpan := startSQLSpan(ctx, "AccountByID", sqlstr, id)
	defer closeSpan()
	a := Account{
		_exists: true,
	}

	// run query
	err = db.QueryRowxContext(ctx, sqlstr, id).Scan(&a.ID, &a.Provider, &a.Email, &a.DisplayName, &a.ScreenImgURL, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

// AccountByEmail retrieves a row from 'accounts' as a Account.
// Generated from index 'email'.
func AccountByEmail(ctx context.Context, db Queryer, email string) (*Account, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, provider, email, display_name, screen_img_url, created_at, updated_at ` +
		`FROM accounts ` +
		`WHERE email = ?`

	// log and trace
	XOLog(ctx, sqlstr, email)
	closeSpan := startSQLSpan(ctx, "AccountByEmail", sqlstr, email)
	defer closeSpan()
	a := Account{
		_exists: true,
	}

	// run query
	err = db.QueryRowxContext(ctx, sqlstr, email).Scan(&a.ID, &a.Provider, &a.Email, &a.DisplayName, &a.ScreenImgURL, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

// AccountsByEmail retrieves a row from 'accounts' as a Account.
// Generated from index 'email_idx'.
func AccountsByEmail(ctx context.Context, db Queryer, email string) ([]*Account, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, provider, email, display_name, screen_img_url, created_at, updated_at ` +
		`FROM accounts ` +
		`WHERE email = ?`

	// log and trace
	XOLog(ctx, sqlstr, email)
	closeSpan := startSQLSpan(ctx, "AccountsByEmail", sqlstr, email)
	defer closeSpan()
	// run query
	rows, err := db.QueryContext(ctx, sqlstr, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// load results
	var res []*Account
	for rows.Next() {
		a := Account{
			_exists: true,
		}

		// scan
		err = rows.Scan(&a.ID, &a.Provider, &a.Email, &a.DisplayName, &a.ScreenImgURL, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &a)
	}

	return res, nil
}
