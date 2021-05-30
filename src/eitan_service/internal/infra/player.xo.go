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

// Player represents a row from 'players'.
type Player struct {
	ID        string    `db:"id"`         // id
	UserID    string    `db:"user_id"`    // user_id
	CreatedAt time.Time `db:"created_at"` // created_at
	UpdatedAt time.Time `db:"updated_at"` // updated_at

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the Player exists in the database.
func (p *Player) Exists() bool {
	return p._exists
}

// GetAllPlayers gets all Players
func GetAllPlayers(ctx context.Context, db Queryer) ([]*Player, error) {
	// sql query
	const sqlstr = `SELECT ` +
		`id, user_id, created_at, updated_at ` +
		`FROM players`

	// log and trace
	XOLog(ctx, sqlstr)
	closeSpan := startSQLSpan(ctx, "GetAllPlayers", sqlstr)
	defer closeSpan()

	var ps []*Player
	rows, err := db.QueryContext(ctx, sqlstr)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		p := Player{_exists: true}
		if err := rows.Scan(&p.ID, &p.UserID, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		ps = append(ps, &p)
	}
	return ps, nil
}

// GetPlayer gets a Player by primary key
func GetPlayer(ctx context.Context, db Queryer, key string) (*Player, error) {
	// sql query
	const sqlstr = `SELECT ` +
		`id, user_id, created_at, updated_at ` +
		`FROM players ` +
		`WHERE id = ?`

	// log and trace
	XOLog(ctx, sqlstr, key)
	closeSpan := startSQLSpan(ctx, "GetPlayer", sqlstr, key)
	defer closeSpan()

	p := Player{_exists: true}
	err := db.QueryRowxContext(ctx, sqlstr, key).Scan(&p.ID, &p.UserID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// GetPlayers gets Player list by primary keys
func GetPlayers(ctx context.Context, db Queryer, keys []string) ([]*Player, error) {
	// sql query
	sqlstr, args, err := sqlx.In(`SELECT `+
		`id, user_id, created_at, updated_at `+
		`FROM players `+
		`WHERE id IN (?)`, keys)
	if err != nil {
		return nil, err
	}

	// log and trace
	XOLog(ctx, sqlstr, args)
	closeSpan := startSQLSpan(ctx, "GetPlayers", sqlstr, args)
	defer closeSpan()

	rows, err := db.QueryContext(ctx, sqlstr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// load results
	var res []*Player
	for rows.Next() {
		p := Player{
			_exists: true,
		}

		// scan
		err = rows.Scan(&p.ID, &p.UserID, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &p)
	}

	return res, nil
}

func QueryPlayer(ctx context.Context, q sqlx.QueryerContext, sqlstr string, args ...interface{}) (*Player, error) {
	// log and trace
	XOLog(ctx, sqlstr, args)
	closeSpan := startSQLSpan(ctx, "QueryPlayer", sqlstr, args)
	defer closeSpan()

	var dest Player
	err := sqlx.GetContext(ctx, q, &dest, sqlstr, args...)
	return &dest, err
}

func QueryPlayers(ctx context.Context, q sqlx.QueryerContext, sqlstr string, args ...interface{}) ([]*Player, error) {
	// log and trace
	XOLog(ctx, sqlstr, args)
	closeSpan := startSQLSpan(ctx, "QueryPlayers", sqlstr, args)
	defer closeSpan()

	var dest []*Player
	err := sqlx.SelectContext(ctx, q, &dest, sqlstr, args...)
	return dest, err
}

// Deleted provides information if the Player has been deleted from the database.
func (p *Player) Deleted() bool {
	return p._deleted
}

// Insert inserts the Player to the database.
func (p *Player) Insert(ctx context.Context, db Execer) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	// if already exist, bail
	if p._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key must be provided
	const sqlstr = `INSERT INTO players (` +
		`id, user_id, created_at, updated_at` +
		`) VALUES (` +
		`?, ?, ?, ?` +
		`)`

	// log and trace
	XOLog(ctx, sqlstr, p.ID, p.UserID, p.CreatedAt, p.UpdatedAt)
	closeSpan := startSQLSpan(ctx, "Player_Insert", sqlstr, p.ID, p.UserID, p.CreatedAt, p.UpdatedAt)
	defer closeSpan()

	// run query
	_, err := db.ExecContext(ctx, sqlstr, p.ID, p.UserID, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return err
	}

	// set existence
	p._exists = true

	return nil
}

// Update updates the Player in the database.
func (p *Player) Update(ctx context.Context, db Execer) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	// if doesn't exist, bail
	if !p._exists {
		return errors.New("update failed: does not exist")
	}
	// if deleted, bail
	if p._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE players SET ` +
		`user_id = ?, created_at = ?, updated_at = ?` +
		` WHERE id = ?`

	// log and trace
	XOLog(ctx, sqlstr, p.UserID, p.CreatedAt, p.UpdatedAt, p.ID)
	closeSpan := startSQLSpan(ctx, "Player_Update", sqlstr, p.UserID, p.CreatedAt, p.UpdatedAt, p.ID)
	defer closeSpan()

	// run query
	_, err := db.ExecContext(ctx, sqlstr, p.UserID, p.CreatedAt, p.UpdatedAt, p.ID)
	return err
}

// Delete deletes the Player from the database.
func (p *Player) Delete(ctx context.Context, db Execer) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	// if doesn't exist, bail
	if !p._exists {
		return nil
	}

	// if deleted, bail
	if p._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM players WHERE id = ?`

	// log and trace
	XOLog(ctx, sqlstr, p.ID)
	closeSpan := startSQLSpan(ctx, "{ .Name }}_Delete", sqlstr, p.ID)
	defer closeSpan()

	// run query
	_, err := db.ExecContext(ctx, sqlstr, p.ID)
	if err != nil {
		return err
	}

	// set deleted
	p._deleted = true

	return nil
}

// InsertOrUpdate inserts or updates the Player to the database.
func (p *Player) InsertOrUpdate(ctx context.Context, db Executor) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	_, err := GetPlayer(ctx, db, p.ID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == sql.ErrNoRows {
		return p.Insert(ctx, db)
	} else {
		p._exists = true
		return p.Update(ctx, db)
	}
}

// InsertOrUpdate inserts or updates the Player to the database.
func (p *Player) InsertIfNotExist(ctx context.Context, db Executor) error {
	if t, ok := tx.GetTx(ctx); ok {
		db = t
	}
	_, err := GetPlayer(ctx, db, p.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return p.Insert(ctx, db)
		}
		return err
	}

	return nil
}

// GetPlayerByID retrieves a row from 'players' as a Player.
// Generated from index 'players_id_pkey'.
func GetPlayerByID(ctx context.Context, db Queryer, id string) (*Player, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, user_id, created_at, updated_at ` +
		`FROM players ` +
		`WHERE id = ?`

	// log and trace
	XOLog(ctx, sqlstr, id)
	closeSpan := startSQLSpan(ctx, "PlayerByID", sqlstr, id)
	defer closeSpan()
	p := Player{
		_exists: true,
	}

	// run query
	err = db.QueryRowxContext(ctx, sqlstr, id).Scan(&p.ID, &p.UserID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// GetPlayerByUserID retrieves a row from 'players' as a Player.
// Generated from index 'user_id'.
func GetPlayerByUserID(ctx context.Context, db Queryer, userID string) (*Player, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, user_id, created_at, updated_at ` +
		`FROM players ` +
		`WHERE user_id = ?`

	// log and trace
	XOLog(ctx, sqlstr, userID)
	closeSpan := startSQLSpan(ctx, "PlayerByUserID", sqlstr, userID)
	defer closeSpan()
	p := Player{
		_exists: true,
	}

	// run query
	err = db.QueryRowxContext(ctx, sqlstr, userID).Scan(&p.ID, &p.UserID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// GetPlayersByUserID retrieves a row from 'players' as a Player.
// Generated from index 'user_id_idx'.
func GetPlayersByUserID(ctx context.Context, db Queryer, userID string) ([]*Player, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, user_id, created_at, updated_at ` +
		`FROM players ` +
		`WHERE user_id = ?`

	// log and trace
	XOLog(ctx, sqlstr, userID)
	closeSpan := startSQLSpan(ctx, "PlayersByUserID", sqlstr, userID)
	defer closeSpan()
	// run query
	rows, err := db.QueryContext(ctx, sqlstr, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// load results
	var res []*Player
	for rows.Next() {
		p := Player{
			_exists: true,
		}

		// scan
		err = rows.Scan(&p.ID, &p.UserID, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &p)
	}

	return res, nil
}
