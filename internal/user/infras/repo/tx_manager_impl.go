package repo

import (
	"context"
	"database/sql"
)

type TxManager struct {
	DB *sql.DB
}

func NewTxManager(db *sql.DB) *TxManager {
	return &TxManager{DB: db}
}

func (tm *TxManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := tm.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	ctxWithTx := context.WithValue(ctx, "tx", tx)

	if err := fn(ctxWithTx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
