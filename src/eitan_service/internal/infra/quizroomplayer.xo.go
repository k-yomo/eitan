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

// QuizRoomPlayer represents a row from 'quiz_room_players'.
type QuizRoomPlayer struct {
	QuizRoomID string    `db:"quiz_room_id"` // quiz_room_id
	PlayerID   string    `db:"player_id"`    // player_id
	CreatedAt  time.Time `db:"created_at"`   // created_at

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the QuizRoomPlayer exists in the database.
func (qrp *QuizRoomPlayer) Exists() bool {
	return qrp._exists
}

// GetAllQuizRoomPlayers gets all QuizRoomPlayers
func GetAllQuizRoomPlayers(ctx context.Context, db Queryer) ([]*QuizRoomPlayer, error) {
	// sql query
	const sqlstr = `SELECT ` +
		`quiz_room_id, player_id, created_at ` +
		`FROM quiz_room_players`

	// log and trace
	XOLog(ctx, sqlstr)
	closeSpan := startSQLSpan(ctx, "GetAllQuizRoomPlayers", sqlstr)
	defer closeSpan()

	var qrps []*QuizRoomPlayer
	rows, err := db.QueryContext(ctx, sqlstr)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		qrp := QuizRoomPlayer{_exists: true}
		if err := rows.Scan(&qrp.QuizRoomID, &qrp.PlayerID, &qrp.CreatedAt); err != nil {
			return nil, err
		}
		qrps = append(qrps, &qrp)
	}
	return qrps, nil
}

// GetQuizRoomPlayer gets a QuizRoomPlayer by primary key
func GetQuizRoomPlayer(ctx context.Context, db Queryer, key string) (*QuizRoomPlayer, error) {
	// sql query
	const sqlstr = `SELECT ` +
		`quiz_room_id, player_id, created_at ` +
		`FROM quiz_room_players ` +
		`WHERE player_id = ?`

	// log and trace
	XOLog(ctx, sqlstr, key)
	closeSpan := startSQLSpan(ctx, "GetQuizRoomPlayer", sqlstr, key)
	defer closeSpan()

	qrp := QuizRoomPlayer{_exists: true}
	err := db.QueryRowxContext(ctx, sqlstr, key).Scan(&qrp.QuizRoomID, &qrp.PlayerID, &qrp.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &qrp, nil
}

// GetQuizRoomPlayers gets QuizRoomPlayer list by primary keys
func GetQuizRoomPlayers(ctx context.Context, db Queryer, keys []string) ([]*QuizRoomPlayer, error) {
	// sql query
	sqlstr, args, err := sqlx.In(`SELECT `+
		`quiz_room_id, player_id, created_at `+
		`FROM quiz_room_players `+
		`WHERE player_id IN (?)`, keys)
	if err != nil {
		return nil, err
	}

	// log and trace
	XOLog(ctx, sqlstr, args)
	closeSpan := startSQLSpan(ctx, "GetQuizRoomPlayers", sqlstr, args)
	defer closeSpan()

	rows, err := db.QueryContext(ctx, sqlstr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// load results
	var res []*QuizRoomPlayer
	for rows.Next() {
		qrp := QuizRoomPlayer{
			_exists: true,
		}

		// scan
		err = rows.Scan(&qrp.QuizRoomID, &qrp.PlayerID, &qrp.CreatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &qrp)
	}

	return res, nil
}

func QueryQuizRoomPlayer(ctx context.Context, q sqlx.QueryerContext, sqlstr string, args ...interface{}) (*QuizRoomPlayer, error) {
	// log and trace
	XOLog(ctx, sqlstr, args)
	closeSpan := startSQLSpan(ctx, "QueryQuizRoomPlayer", sqlstr, args)
	defer closeSpan()

	var dest QuizRoomPlayer
	err := sqlx.GetContext(ctx, q, &dest, sqlstr, args...)
	return &dest, err
}

func QueryQuizRoomPlayers(ctx context.Context, q sqlx.QueryerContext, sqlstr string, args ...interface{}) ([]*QuizRoomPlayer, error) {
	// log and trace
	XOLog(ctx, sqlstr, args)
	closeSpan := startSQLSpan(ctx, "QueryQuizRoomPlayers", sqlstr, args)
	defer closeSpan()

	var dest []*QuizRoomPlayer
	err := sqlx.SelectContext(ctx, q, &dest, sqlstr, args...)
	return dest, err
}

// Deleted provides information if the QuizRoomPlayer has been deleted from the database.
func (qrp *QuizRoomPlayer) Deleted() bool {
	return qrp._deleted
}

// Insert inserts the QuizRoomPlayer to the database.
func (qrp *QuizRoomPlayer) Insert(ctx context.Context, db Execer) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	// if already exist, bail
	if qrp._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key must be provided
	const sqlstr = `INSERT INTO quiz_room_players (` +
		`quiz_room_id, player_id, created_at` +
		`) VALUES (` +
		`?, ?, ?` +
		`)`

	// log and trace
	XOLog(ctx, sqlstr, qrp.QuizRoomID, qrp.PlayerID, qrp.CreatedAt)
	closeSpan := startSQLSpan(ctx, "QuizRoomPlayer_Insert", sqlstr, qrp.QuizRoomID, qrp.PlayerID, qrp.CreatedAt)
	defer closeSpan()

	// run query
	_, err := db.ExecContext(ctx, sqlstr, qrp.QuizRoomID, qrp.PlayerID, qrp.CreatedAt)
	if err != nil {
		return err
	}

	// set existence
	qrp._exists = true

	return nil
}

// Update updates the QuizRoomPlayer in the database.
func (qrp *QuizRoomPlayer) Update(ctx context.Context, db Execer) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	// if doesn't exist, bail
	if !qrp._exists {
		return errors.New("update failed: does not exist")
	}
	// if deleted, bail
	if qrp._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query with composite primary key
	const sqlstr = `UPDATE quiz_room_players SET ` +
		`created_at = ?` +
		` WHERE quiz_room_id = ? AND player_id = ?`

	// log and trace
	XOLog(ctx, sqlstr, qrp.CreatedAt, qrp.QuizRoomID, qrp.PlayerID)
	closeSpan := startSQLSpan(ctx, "QuizRoomPlayer_Update", sqlstr, qrp.CreatedAt, qrp.QuizRoomID, qrp.PlayerID)
	defer closeSpan()

	// run query
	_, err := db.ExecContext(ctx, sqlstr, qrp.CreatedAt, qrp.QuizRoomID, qrp.PlayerID)
	return err
}

// Delete deletes the QuizRoomPlayer from the database.
func (qrp *QuizRoomPlayer) Delete(ctx context.Context, db Execer) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	// if doesn't exist, bail
	if !qrp._exists {
		return nil
	}

	// if deleted, bail
	if qrp._deleted {
		return nil
	}

	// sql query with composite primary key
	const sqlstr = `DELETE FROM quiz_room_players WHERE quiz_room_id = ? AND player_id = ?`

	// log and trace
	XOLog(ctx, sqlstr, qrp.QuizRoomID, qrp.PlayerID)
	closeSpan := startSQLSpan(ctx, "{ .Name }}Delete", sqlstr, qrp.QuizRoomID, qrp.PlayerID)
	defer closeSpan()

	// run query
	_, err := db.ExecContext(ctx, sqlstr, qrp.QuizRoomID, qrp.PlayerID)
	if err != nil {
		return err
	}

	// set deleted
	qrp._deleted = true

	return nil
}

// InsertOrUpdate inserts or updates the QuizRoomPlayer to the database.
func (qrp *QuizRoomPlayer) InsertOrUpdate(ctx context.Context, db Executor) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	_, err := GetQuizRoomPlayer(ctx, db, qrp.PlayerID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == sql.ErrNoRows {
		return qrp.Insert(ctx, db)
	} else {
		qrp._exists = true
		return qrp.Update(ctx, db)
	}
}

// InsertOrUpdate inserts or updates the QuizRoomPlayer to the database.
func (qrp *QuizRoomPlayer) InsertIfNotExist(ctx context.Context, db Executor) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	_, err := GetQuizRoomPlayer(ctx, db, qrp.PlayerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return qrp.Insert(ctx, db)
		}
		return err
	}

	return nil
}

// QuizRoom returns the QuizRoom associated with the QuizRoomPlayer's QuizRoomID (quiz_room_id).
//
// Generated from foreign key 'quiz_room_players_ibfk_1'.
func (qrp *QuizRoomPlayer) QuizRoom(ctx context.Context, db Executor) (*QuizRoom, error) {
	return GetQuizRoomByID(ctx, db, qrp.QuizRoomID)
}

// Player returns the Player associated with the QuizRoomPlayer's PlayerID (player_id).
//
// Generated from foreign key 'quiz_room_players_ibfk_2'.
func (qrp *QuizRoomPlayer) Player(ctx context.Context, db Executor) (*Player, error) {
	return GetPlayerByID(ctx, db, qrp.PlayerID)
}

// GetQuizRoomPlayersByPlayerID retrieves a row from 'quiz_room_players' as a QuizRoomPlayer.
// Generated from index 'player_id'.
func GetQuizRoomPlayersByPlayerID(ctx context.Context, db Queryer, playerID string) ([]*QuizRoomPlayer, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`quiz_room_id, player_id, created_at ` +
		`FROM quiz_room_players ` +
		`WHERE player_id = ?`

	// log and trace
	XOLog(ctx, sqlstr, playerID)
	closeSpan := startSQLSpan(ctx, "QuizRoomPlayersByPlayerID", sqlstr, playerID)
	defer closeSpan()
	// run query
	rows, err := db.QueryContext(ctx, sqlstr, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// load results
	var res []*QuizRoomPlayer
	for rows.Next() {
		qrp := QuizRoomPlayer{
			_exists: true,
		}

		// scan
		err = rows.Scan(&qrp.QuizRoomID, &qrp.PlayerID, &qrp.CreatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &qrp)
	}

	return res, nil
}

// GetQuizRoomPlayerByPlayerID retrieves a row from 'quiz_room_players' as a QuizRoomPlayer.
// Generated from index 'quiz_room_players_player_id_pkey'.
func GetQuizRoomPlayerByPlayerID(ctx context.Context, db Queryer, playerID string) (*QuizRoomPlayer, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`quiz_room_id, player_id, created_at ` +
		`FROM quiz_room_players ` +
		`WHERE player_id = ?`

	// log and trace
	XOLog(ctx, sqlstr, playerID)
	closeSpan := startSQLSpan(ctx, "QuizRoomPlayerByPlayerID", sqlstr, playerID)
	defer closeSpan()
	qrp := QuizRoomPlayer{
		_exists: true,
	}

	// run query
	err = db.QueryRowxContext(ctx, sqlstr, playerID).Scan(&qrp.QuizRoomID, &qrp.PlayerID, &qrp.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &qrp, nil
}
