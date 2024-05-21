package account

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aradwann/eenergy/repository/postgres/common"
)

type Account struct {
	ID        int64     `json:"id"`
	Owner     string    `json:"owner"`
	Balance   int64     `json:"balance"`
	Unit      string    `json:"unit"`
	CreatedAt time.Time `json:"created_at"`
}

type AccountRepository interface {
	GetAccount(ctx context.Context, id int64) (*Account, error)
	AddAccountBalance(ctx context.Context, arg AddAccountBalanceParams) (*Account, error)
	CreateAccount(ctx context.Context, arg CreateAccountParams) (*Account, error)
	DeleteAccount(ctx context.Context, id int64) error
	ListAccounts(ctx context.Context, arg ListAccountsParams) ([]*Account, error)
	UpdateAccount(ctx context.Context, arg UpdateAccountParams) (*Account, error)
}

type accountRepository struct {
	db *sql.DB
}

// https://github.com/uber-go/guide/blob/master/style.md#verify-interface-compliance
var _ AccountRepository = (*accountRepository)(nil)

func NewAccountRepository(db *sql.DB) AccountRepository {
	return &accountRepository{db: db}
}

type AddAccountBalanceParams struct {
	Amount int64 `json:"amount"`
	ID     int64 `json:"id"`
}

func (r *accountRepository) AddAccountBalance(ctx context.Context, arg AddAccountBalanceParams) (*Account, error) {
	acc := &Account{}
	row := common.CallStoredFunction(ctx, r.db, "add_account_balance",
		arg.Amount,
		arg.ID,
	)

	if err := scanAccount(row, acc); err != nil {
		return nil, fmt.Errorf("failed to scan account: %w", err)
	}

	return acc, nil
}

type CreateAccountParams struct {
	Owner   string `json:"owner"`
	Balance int64  `json:"balance"`
	Unit    string `json:"unit"`
}

func (r *accountRepository) CreateAccount(ctx context.Context, arg CreateAccountParams) (*Account, error) {
	acc := &Account{}
	row := common.CallStoredFunction(ctx, r.db, "create_account",
		arg.Owner,
		arg.Balance,
		arg.Unit,
	)

	if err := scanAccount(row, acc); err != nil {
		return nil, fmt.Errorf("failed to scan account: %w", err)
	}

	return acc, nil
}

func (r *accountRepository) DeleteAccount(ctx context.Context, id int64) error {
	var result bool

	row := common.CallStoredFunction(ctx, r.db, "delete_account",
		id,
	)
	err := row.Scan(&result)
	// TODO: handle logging
	if err != nil {
		fmt.Printf("Failed to execute function: %v", err)
	} else if result {
		fmt.Println("Account deleted successfully.")
	} else {
		fmt.Println("No account was deleted.")
	}

	return row.Err()
}

func (r *accountRepository) GetAccount(ctx context.Context, id int64) (*Account, error) {
	acc := &Account{}
	row := common.CallStoredFunction(ctx, r.db, "get_account", id)

	if err := scanAccount(row, acc); err != nil {
		return nil, fmt.Errorf("failed to scan account: %w", err)
	}

	return acc, nil
}

type ListAccountsParams struct {
	Owner  string `json:"owner"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

func (r *accountRepository) ListAccounts(ctx context.Context, arg ListAccountsParams) ([]*Account, error) {
	rows, err := common.CallStoredFunctionRows(ctx, r.db, "list_accounts", arg.Owner, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Account
	for rows.Next() {
		acc := &Account{}
		if err := scanAccount(rows, acc); err != nil {
			return nil, fmt.Errorf("failed to scan account: %w", err)
		}
		items = append(items, acc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

type UpdateAccountParams struct {
	ID      int64 `json:"id"`
	Balance int64 `json:"balance"`
}

func (r *accountRepository) UpdateAccount(ctx context.Context, arg UpdateAccountParams) (*Account, error) {
	acc := &Account{}
	row := common.CallStoredFunction(ctx, r.db, "update_account", arg.ID, arg.Balance)
	if err := scanAccount(row, acc); err != nil {
		return nil, fmt.Errorf("failed to scan account: %w", err)
	}
	return acc, nil
}

// scanAccount abstracts the row scanning logic for an Account. Works with *sql.Row and *sql.Rows.
func scanAccount(scanner interface {
	Scan(dest ...interface{}) error
}, acc *Account) error {
	if scanner == nil {
		return fmt.Errorf("row is nil")
	}
	if acc == nil {
		return fmt.Errorf("account is nil")
	}

	err := scanner.Scan(&acc.ID, &acc.Owner, &acc.Balance, &acc.Unit, &acc.CreatedAt)
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
