package tx

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Manager interface {
	RunInTx(ctx context.Context, f func(ctx context.Context) error) error
}

type txCtxKey struct {}

type dbTxManager struct {
	db *sqlx.DB
}

func NewManager(db *sqlx.DB) Manager {
	return &dbTxManager{db: db}
}

func (t *dbTxManager) RunInTx(ctx context.Context, f func(ctx context.Context) error) error {
	// if transaction is already started
	if _, ok := GetTx(ctx); ok {
		return f(ctx)
	}

	tx, err := t.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	if err := f(context.WithValue(ctx, txCtxKey{}, tx)); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	return tx.Commit()
}

func GetTx(ctx context.Context) (*sqlx.Tx, bool) {
	extractedTx := ctx.Value(txCtxKey{})
	tx, ok := extractedTx.(*sqlx.Tx)
	return tx, ok
}
