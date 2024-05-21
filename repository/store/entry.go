package db

import (
	"context"
	"database/sql"
	"errors"
	"log"
)

type CreateEntryParams struct {
	AccountID int64 `json:"account_id"`
	Amount    int64 `json:"amount"`
}

func (q *Queries) CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error) {
	var e Entry
	row := q.callStoredFunction(ctx, "create_entry",
		arg.AccountID,
		arg.Amount,
	)
	err := scanEntry(row, &e)
	return e, err
}

func (q *Queries) GetEntry(ctx context.Context, id int64) (Entry, error) {
	var e Entry
	row := q.callStoredFunction(ctx, "get_entry", id)
	err := scanEntry(row, &e)
	return e, err
}

type ListEntriesParams struct {
	AccountID int64 `json:"account_id"`
	Limit     int32 `json:"limit"`
	Offset    int32 `json:"offset"`
}

func (q *Queries) ListEntries(ctx context.Context, arg ListEntriesParams) ([]Entry, error) {
	rows, err := q.callStoredFunctionRows(ctx, "list_entries",
		arg.AccountID,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []Entry{}
	for rows.Next() {
		var e Entry
		if err := scanEntry(rows, &e); err != nil {
			return nil, err
		}
		items = append(items, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

// scanEntry abstracts the row scanning logic for an Entry. Works with *sql.Row and *sql.Rows.
func scanEntry(scanner interface {
	Scan(dest ...interface{}) error
}, t *Entry) error {
	err := scanner.Scan(
		&t.ID,
		&t.AccountID,
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
