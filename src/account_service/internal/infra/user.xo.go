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

// User represents a row from 'users'.
type User struct {
	ID        string    `db:"id"`         // id
	CreatedAt time.Time `db:"created_at"` // created_at
	UpdatedAt time.Time `db:"updated_at"` // updated_at

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the User exists in the database.
func (u *User) Exists() bool {
	return u._exists
}

// GetAllUsers gets all Users
func GetAllUsers(ctx context.Context, db Queryer) ([]*User, error) {
	// sql query
	const sqlstr = `SELECT ` +
		`id, created_at, updated_at ` +
		`FROM users`

	// log and trace
	XOLog(ctx, sqlstr)
	closeSpan := startSQLSpan(ctx, "GetAllUsers", sqlstr)
	defer closeSpan()

	var us []*User
	rows, err := db.QueryContext(ctx, sqlstr)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		u := User{_exists: true}
		if err := rows.Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		us = append(us, &u)
	}
	return us, nil
}

// GetUser gets a User by primary key
func GetUser(ctx context.Context, db Queryer, key string) (*User, error) {
	// sql query
	const sqlstr = `SELECT ` +
		`id, created_at, updated_at ` +
		`FROM users ` +
		`WHERE id = ?`

	// log and trace
	XOLog(ctx, sqlstr, key)
	closeSpan := startSQLSpan(ctx, "GetUser", sqlstr, key)
	defer closeSpan()

	u := User{_exists: true}
	err := db.QueryRowxContext(ctx, sqlstr, key).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetUsers gets User list by primary keys
func GetUsers(ctx context.Context, db Queryer, keys []string) ([]*User, error) {
	// sql query
	sqlstr, args, err := sqlx.In(`SELECT `+
		`id, created_at, updated_at `+
		`FROM users `+
		`WHERE id IN (?)`, keys)
	if err != nil {
		return nil, err
	}

	// log and trace
	XOLog(ctx, sqlstr, args)
	closeSpan := startSQLSpan(ctx, "GetUsers", sqlstr, args)
	defer closeSpan()

	rows, err := db.QueryContext(ctx, sqlstr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// load results
	var res []*User
	for rows.Next() {
		u := User{
			_exists: true,
		}

		// scan
		err = rows.Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &u)
	}

	return res, nil
}

func QueryUser(ctx context.Context, q sqlx.QueryerContext, sqlstr string, args ...interface{}) (*User, error) {
	// log and trace
	XOLog(ctx, sqlstr, args)
	closeSpan := startSQLSpan(ctx, "QueryUser", sqlstr, args)
	defer closeSpan()

	var dest User
	err := sqlx.GetContext(ctx, q, &dest, sqlstr, args...)
	return &dest, err
}

func QueryUsers(ctx context.Context, q sqlx.QueryerContext, sqlstr string, args ...interface{}) ([]*User, error) {
	// log and trace
	XOLog(ctx, sqlstr, args)
	closeSpan := startSQLSpan(ctx, "QueryUsers", sqlstr, args)
	defer closeSpan()

	var dest []*User
	err := sqlx.SelectContext(ctx, q, &dest, sqlstr, args...)
	return dest, err
}

// Deleted provides information if the User has been deleted from the database.
func (u *User) Deleted() bool {
	return u._deleted
}

// Insert inserts the User to the database.
func (u *User) Insert(ctx context.Context, db Execer) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	// if already exist, bail
	if u._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key must be provided
	const sqlstr = `INSERT INTO users (` +
		`id, created_at, updated_at` +
		`) VALUES (` +
		`?, ?, ?` +
		`)`

	// log and trace
	XOLog(ctx, sqlstr, u.ID, u.CreatedAt, u.UpdatedAt)
	closeSpan := startSQLSpan(ctx, "User_Insert", sqlstr, u.ID, u.CreatedAt, u.UpdatedAt)
	defer closeSpan()

	// run query
	_, err := db.ExecContext(ctx, sqlstr, u.ID, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return err
	}

	// set existence
	u._exists = true

	return nil
}

// Update updates the User in the database.
func (u *User) Update(ctx context.Context, db Execer) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	// if doesn't exist, bail
	if !u._exists {
		return errors.New("update failed: does not exist")
	}
	// if deleted, bail
	if u._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE users SET ` +
		`created_at = ?, updated_at = ?` +
		` WHERE id = ?`

	// log and trace
	XOLog(ctx, sqlstr, u.CreatedAt, u.UpdatedAt, u.ID)
	closeSpan := startSQLSpan(ctx, "User_Update", sqlstr, u.CreatedAt, u.UpdatedAt, u.ID)
	defer closeSpan()

	// run query
	_, err := db.ExecContext(ctx, sqlstr, u.CreatedAt, u.UpdatedAt, u.ID)
	return err
}

// Delete deletes the User from the database.
func (u *User) Delete(ctx context.Context, db Execer) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	// if doesn't exist, bail
	if !u._exists {
		return nil
	}

	// if deleted, bail
	if u._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM users WHERE id = ?`

	// log and trace
	XOLog(ctx, sqlstr, u.ID)
	closeSpan := startSQLSpan(ctx, "{ .Name }}_Delete", sqlstr, u.ID)
	defer closeSpan()

	// run query
	_, err := db.ExecContext(ctx, sqlstr, u.ID)
	if err != nil {
		return err
	}

	// set deleted
	u._deleted = true

	return nil
}

// InsertOrUpdate inserts or updates the User to the database.
func (u *User) InsertOrUpdate(ctx context.Context, db Executor) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	_, err := GetUser(ctx, db, u.ID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == sql.ErrNoRows {
		return u.Insert(ctx, db)
	} else {
		u._exists = true
		return u.Update(ctx, db)
	}
}

// InsertOrUpdate inserts or updates the User to the database.
func (u *User) InsertIfNotExist(ctx context.Context, db Executor) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	_, err := GetUser(ctx, db, u.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return u.Insert(ctx, db)
		}
		return err
	}

	return nil
}

// GetUserByID retrieves a row from 'users' as a User.
// Generated from index 'users_id_pkey'.
func GetUserByID(ctx context.Context, db Queryer, id string) (*User, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, created_at, updated_at ` +
		`FROM users ` +
		`WHERE id = ?`

	// log and trace
	XOLog(ctx, sqlstr, id)
	closeSpan := startSQLSpan(ctx, "UserByID", sqlstr, id)
	defer closeSpan()
	u := User{
		_exists: true,
	}

	// run query
	err = db.QueryRowxContext(ctx, sqlstr, id).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
