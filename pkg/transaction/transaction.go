package transaction

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(tx *sqlx.Tx) error) error
}

type TxManager struct {
	db *sqlx.DB
}

func NewTxManager(
	db *sqlx.DB,
) TransactionManager {
	return &TxManager{
		db: db,
	}
}

func (tm *TxManager) WithTransaction(ctx context.Context, fn func(tx *sqlx.Tx) error) error {
	tx, err := tm.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	if err = fn(tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
