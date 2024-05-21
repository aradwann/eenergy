package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aradwann/eenergy/repository/postgres/common"
)

type Transfer struct {
	ID            int64 `json:"id"`
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	// it must be positive
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
type TransferRepository interface {
	CreateTransfer(ctx context.Context, arg CreateTransferParams) (*Transfer, error)
	GetTransfer(ctx context.Context, id int64) (*Transfer, error)
	ListTransfers(ctx context.Context, arg ListTransfersParams) ([]*Transfer, error)
}

type transferRepository struct {
	db *sql.DB
}

// https://github.com/uber-go/guide/blob/master/style.md#verify-interface-compliance
var _ TransferRepository = (*transferRepository)(nil)

func NewTransferRepository(db *sql.DB) TransferRepository {
	return &transferRepository{db: db}
}

type CreateTransferParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

func (r *transferRepository) CreateTransfer(ctx context.Context, arg CreateTransferParams) (*Transfer, error) {
	t := &Transfer{}
	row := common.CallStoredFunction(ctx, r.db, "create_transfer",
		arg.FromAccountID,
		arg.ToAccountID,
		arg.Amount,
	)
	if err := scanTransfer(row, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (r *transferRepository) GetTransfer(ctx context.Context, id int64) (*Transfer, error) {
	t := &Transfer{}

	row := common.CallStoredFunction(ctx, r.db, "get_transfer", id)
	if err := scanTransfer(row, t); err != nil {
		return nil, err
	}
	return t, nil
}

type ListTransfersParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Limit         int32 `json:"limit"`
	Offset        int32 `json:"offset"`
}

func (r *transferRepository) ListTransfers(ctx context.Context, arg ListTransfersParams) ([]*Transfer, error) {
	rows, err := common.CallStoredFunctionRows(ctx, r.db, "list_transfers",
		arg.FromAccountID,
		arg.ToAccountID,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []*Transfer{}
	for rows.Next() {
		t := &Transfer{}
		if err := scanTransfer(rows, t); err != nil {
			return nil, err
		}
		items = append(items, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

// scanTransfer abstracts the row scanning logic for an Transfer. Works with *sql.Row and *sql.Rows.
func scanTransfer(scanner interface {
	Scan(dest ...interface{}) error
}, t *Transfer) error {
	if scanner == nil {
		return fmt.Errorf("row is nil")
	}
	if t == nil {
		return fmt.Errorf("transfer is nil")
	}

	err := scanner.Scan(
		&t.ID,
		&t.FromAccountID,
		&t.ToAccountID,
		&t.Amount,
		&t.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Handle the case where no rows are found.
			log.Println("No rows were returned.")
			return common.ErrRecordNotFound
		}
		// Log and return other scan-related errors.
		log.Printf("Error scanning row: %s", err)
		return err
	}
	return nil
}
