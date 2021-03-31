package tx

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type txKey string

const key txKey = "txKey"

type Manager interface {
	RunInTx(ctx context.Context, f func(ctx context.Context) error) error
}

type txManagerImpl struct {
	db *sqlx.DB
}

func NewManager(db *sqlx.DB) Manager {
	return &txManagerImpl{db: db}
}

func (t *txManagerImpl) RunInTx(ctx context.Context, f func(ctx context.Context) error) error {
	// if transaction is already started
	if _, ok := GetTx(ctx); ok {
		return f(ctx)
	}

	tx, err := t.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	if err := f(context.WithValue(ctx, key, tx)); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	return tx.Commit()
}

func GetTx(ctx context.Context) (*sqlx.Tx, bool) {
	extractedTx := ctx.Value(key)
	tx, ok := extractedTx.(*sqlx.Tx)
	return tx, ok
}
