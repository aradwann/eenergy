package common

import (
	"context"
	"database/sql"
	"fmt"
)

// TransactionManager handles the lifecycle of a database transaction.
type TransactionManager struct {
	db *sql.DB
}

// NewTransactionManager creates a new TransactionManager.
func NewTransactionManager(db *sql.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

// ExecTx executes a function within a database transaction.
func (tm *TransactionManager) ExecTx(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := tm.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
