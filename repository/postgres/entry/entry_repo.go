package entry

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aradwann/eenergy/repository/postgres/common"
)

type Entry struct {
	ID        int64 `json:"id"`
	AccountID int64 `json:"account_id"`
	// can be negative or positive
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type EntryRepository interface {
	CreateEntry(ctx context.Context, arg CreateEntryParams) (*Entry, error)
	GetEntry(ctx context.Context, id int64) (*Entry, error)
	ListEntries(ctx context.Context, arg ListEntriesParams) ([]*Entry, error)
}

type entryRepository struct {
	db *sql.DB
}

// https://github.com/uber-go/guide/blob/master/style.md#verify-interface-compliance
var _ EntryRepository = (*entryRepository)(nil)

func NewEntryRepository(db *sql.DB) EntryRepository {
	return &entryRepository{db: db}
}

type CreateEntryParams struct {
	AccountID int64 `json:"account_id"`
	Amount    int64 `json:"amount"`
}

func (r *entryRepository) CreateEntry(ctx context.Context, arg CreateEntryParams) (*Entry, error) {
	e := &Entry{}
	row := common.CallStoredFunction(ctx, r.db, "create_entry",
		arg.AccountID,
		arg.Amount,
	)
	err := scanEntry(row, e)
	return e, err
}

func (r *entryRepository) GetEntry(ctx context.Context, id int64) (*Entry, error) {
	e := &Entry{}
	row := common.CallStoredFunction(ctx, r.db, "get_entry", id)
	err := scanEntry(row, e)
	return e, err
}

type ListEntriesParams struct {
	AccountID int64 `json:"account_id"`
	Limit     int32 `json:"limit"`
	Offset    int32 `json:"offset"`
}

func (r *entryRepository) ListEntries(ctx context.Context, arg ListEntriesParams) ([]*Entry, error) {
	rows, err := common.CallStoredFunctionRows(ctx, r.db, "list_entries",
		arg.AccountID,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []*Entry{}
	for rows.Next() {
		e := &Entry{}
		if err := scanEntry(rows, e); err != nil {
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
	if scanner == nil {
		return fmt.Errorf("row is nil")
	}
	if t == nil {
		return fmt.Errorf("account is nil")
	}
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
			return common.ErrRecordNotFound
		}
		// Log and return other scan-related errors.
		log.Printf("Error scanning row: %s", err)
		return err
	}
	return nil
}
