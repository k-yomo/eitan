// Package infra contains the types for schema 'eitandb'.
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

// QuizRoom represents a row from 'quiz_rooms'.
type QuizRoom struct {
	ID        string    `db:"id"`         // id
	CreatedAt time.Time `db:"created_at"` // created_at
	UpdatedAt time.Time `db:"updated_at"` // updated_at

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the QuizRoom exists in the database.
func (qr *QuizRoom) Exists() bool {
	return qr._exists
}

// GetAllQuizRooms gets all QuizRooms
func GetAllQuizRooms(ctx context.Context, db Queryer) ([]*QuizRoom, error) {
	// sql query
	const sqlstr = `SELECT ` +
		`id, created_at, updated_at ` +
		`FROM quiz_rooms`

	// log and trace
	XOLog(ctx, sqlstr)
	closeSpan := startSQLSpan(ctx, "GetAllQuizRooms", sqlstr)
	defer closeSpan()

	var qrs []*QuizRoom
	rows, err := db.QueryContext(ctx, sqlstr)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		qr := QuizRoom{_exists: true}
		if err := rows.Scan(&qr.ID, &qr.CreatedAt, &qr.UpdatedAt); err != nil {
			return nil, err
		}
		qrs = append(qrs, &qr)
	}
	return qrs, nil
}

// GetQuizRoom gets a QuizRoom by primary key
func GetQuizRoom(ctx context.Context, db Queryer, key string) (*QuizRoom, error) {
	// sql query
	const sqlstr = `SELECT ` +
		`id, created_at, updated_at ` +
		`FROM quiz_rooms ` +
		`WHERE id = ?`

	// log and trace
	XOLog(ctx, sqlstr, key)
	closeSpan := startSQLSpan(ctx, "GetQuizRoom", sqlstr, key)
	defer closeSpan()

	qr := QuizRoom{_exists: true}
	err := db.QueryRowxContext(ctx, sqlstr, key).Scan(&qr.ID, &qr.CreatedAt, &qr.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &qr, nil
}

// GetQuizRooms gets QuizRoom list by primary keys
func GetQuizRooms(ctx context.Context, db Queryer, keys []string) ([]*QuizRoom, error) {
	// sql query
	sqlstr, args, err := sqlx.In(`SELECT `+
		`id, created_at, updated_at `+
		`FROM quiz_rooms `+
		`WHERE id IN (?)`, keys)
	if err != nil {
		return nil, err
	}

	// log and trace
	XOLog(ctx, sqlstr, args)
	closeSpan := startSQLSpan(ctx, "GetQuizRooms", sqlstr, args)
	defer closeSpan()

	rows, err := db.QueryContext(ctx, sqlstr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// load results
	var res []*QuizRoom
	for rows.Next() {
		qr := QuizRoom{
			_exists: true,
		}

		// scan
		err = rows.Scan(&qr.ID, &qr.CreatedAt, &qr.UpdatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &qr)
	}

	return res, nil
}

func QueryQuizRoom(ctx context.Context, q sqlx.QueryerContext, sqlstr string, args ...interface{}) (*QuizRoom, error) {
	// log and trace
	XOLog(ctx, sqlstr, args)
	closeSpan := startSQLSpan(ctx, "QueryQuizRoom", sqlstr, args)
	defer closeSpan()

	var dest QuizRoom
	err := sqlx.GetContext(ctx, q, &dest, sqlstr, args...)
	return &dest, err
}

func QueryQuizRooms(ctx context.Context, q sqlx.QueryerContext, sqlstr string, args ...interface{}) ([]*QuizRoom, error) {
	// log and trace
	XOLog(ctx, sqlstr, args)
	closeSpan := startSQLSpan(ctx, "QueryQuizRooms", sqlstr, args)
	defer closeSpan()

	var dest []*QuizRoom
	err := sqlx.SelectContext(ctx, q, &dest, sqlstr, args...)
	return dest, err
}

// Deleted provides information if the QuizRoom has been deleted from the database.
func (qr *QuizRoom) Deleted() bool {
	return qr._deleted
}

// Insert inserts the QuizRoom to the database.
func (qr *QuizRoom) Insert(ctx context.Context, db Execer) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	// if already exist, bail
	if qr._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key must be provided
	const sqlstr = `INSERT INTO quiz_rooms (` +
		`id, created_at, updated_at` +
		`) VALUES (` +
		`?, ?, ?` +
		`)`

	// log and trace
	XOLog(ctx, sqlstr, qr.ID, qr.CreatedAt, qr.UpdatedAt)
	closeSpan := startSQLSpan(ctx, "QuizRoom_Insert", sqlstr, qr.ID, qr.CreatedAt, qr.UpdatedAt)
	defer closeSpan()

	// run query
	_, err := db.ExecContext(ctx, sqlstr, qr.ID, qr.CreatedAt, qr.UpdatedAt)
	if err != nil {
		return err
	}

	// set existence
	qr._exists = true

	return nil
}

// Update updates the QuizRoom in the database.
func (qr *QuizRoom) Update(ctx context.Context, db Execer) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	// if doesn't exist, bail
	if !qr._exists {
		return errors.New("update failed: does not exist")
	}
	// if deleted, bail
	if qr._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE quiz_rooms SET ` +
		`created_at = ?, updated_at = ?` +
		` WHERE id = ?`

	// log and trace
	XOLog(ctx, sqlstr, qr.CreatedAt, qr.UpdatedAt, qr.ID)
	closeSpan := startSQLSpan(ctx, "QuizRoom_Update", sqlstr, qr.CreatedAt, qr.UpdatedAt, qr.ID)
	defer closeSpan()

	// run query
	_, err := db.ExecContext(ctx, sqlstr, qr.CreatedAt, qr.UpdatedAt, qr.ID)
	return err
}

// Delete deletes the QuizRoom from the database.
func (qr *QuizRoom) Delete(ctx context.Context, db Execer) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	// if doesn't exist, bail
	if !qr._exists {
		return nil
	}

	// if deleted, bail
	if qr._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM quiz_rooms WHERE id = ?`

	// log and trace
	XOLog(ctx, sqlstr, qr.ID)
	closeSpan := startSQLSpan(ctx, "{ .Name }}_Delete", sqlstr, qr.ID)
	defer closeSpan()

	// run query
	_, err := db.ExecContext(ctx, sqlstr, qr.ID)
	if err != nil {
		return err
	}

	// set deleted
	qr._deleted = true

	return nil
}

// InsertOrUpdate inserts or updates the QuizRoom to the database.
func (qr *QuizRoom) InsertOrUpdate(ctx context.Context, db Executor) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	_, err := GetQuizRoom(ctx, db, qr.ID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == sql.ErrNoRows {
		return qr.Insert(ctx, db)
	} else {
		qr._exists = true
		return qr.Update(ctx, db)
	}
}

// InsertOrUpdate inserts or updates the QuizRoom to the database.
func (qr *QuizRoom) InsertIfNotExist(ctx context.Context, db Executor) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	_, err := GetQuizRoom(ctx, db, qr.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return qr.Insert(ctx, db)
		}
		return err
	}

	return nil
}

// GetQuizRoomByID retrieves a row from 'quiz_rooms' as a QuizRoom.
// Generated from index 'quiz_rooms_id_pkey'.
func GetQuizRoomByID(ctx context.Context, db Queryer, id string) (*QuizRoom, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, created_at, updated_at ` +
		`FROM quiz_rooms ` +
		`WHERE id = ?`

	// log and trace
	XOLog(ctx, sqlstr, id)
	closeSpan := startSQLSpan(ctx, "QuizRoomByID", sqlstr, id)
	defer closeSpan()
	qr := QuizRoom{
		_exists: true,
	}

	// run query
	err = db.QueryRowxContext(ctx, sqlstr, id).Scan(&qr.ID, &qr.CreatedAt, &qr.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &qr, nil
}
