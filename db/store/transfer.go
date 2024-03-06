package db

import (
	"context"
	"database/sql"
	"errors"
	"log"
)

type CreateTransferParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

func (q *Queries) CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error) {
	var t Transfer
	row := q.callStoredFunction(ctx, "create_transfer",
		arg.FromAccountID,
		arg.ToAccountID,
		arg.Amount,
	)
	err := scanTransfer(row, &t)
	if err != nil {
		return t, err
	}

	return t, err
}

func (q *Queries) GetTransfer(ctx context.Context, id int64) (Transfer, error) {
	var t Transfer

	row := q.callStoredFunction(ctx, "get_transfer", id)
	err := row.Scan(
		&t.ID,
		&t.FromAccountID,
		&t.ToAccountID,
		&t.Amount,
		&t.CreatedAt,
	)
	return t, err
}

type ListTransfersParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Limit         int32 `json:"limit"`
	Offset        int32 `json:"offset"`
}

func (q *Queries) ListTransfers(ctx context.Context, arg ListTransfersParams) ([]Transfer, error) {
	rows, err := q.callStoredFunctionRows(ctx, "list_transfers",
		arg.FromAccountID,
		arg.ToAccountID,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []Transfer{}
	for rows.Next() {
		var t Transfer
		if err := scanTransfer(rows, &t); err != nil {
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
			return ErrRecordNotFound
		}
		// Log and return other scan-related errors.
		log.Printf("Error scanning row: %s", err)
		return err
	}
	return nil
}
