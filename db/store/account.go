package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type AddAccountBalanceParams struct {
	Amount int64 `json:"amount"`
	ID     int64 `json:"id"`
}

func (q *Queries) AddAccountBalance(ctx context.Context, arg AddAccountBalanceParams) (Account, error) {
	var acc Account
	row := q.callStoredFunction(ctx, "add_account_balace",
		arg.Amount,
		arg.ID,
	)

	err := scanAccount(row, &acc)

	return acc, err
}

type CreateAccountParams struct {
	Owner   string `json:"owner"`
	Balance int64  `json:"balance"`
	Unit    string `json:"unit"`
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error) {
	var acc Account
	row := q.callStoredFunction(ctx, "create_account",
		arg.Owner,
		arg.Balance,
		arg.Unit,
	)

	err := scanAccount(row, &acc)

	return acc, err
}

func (q *Queries) DeleteAccount(ctx context.Context, id int64) error {
	var result bool

	row := q.callStoredFunction(ctx, "delete_account",
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

func (q *Queries) GetAccount(ctx context.Context, id int64) (Account, error) {
	var acc Account
	row := q.callStoredFunction(ctx, "get_account",
		id,
	)

	err := scanAccount(row, &acc)

	return acc, err
}

func (q *Queries) GetAccountForUpdate(ctx context.Context, id int64) (Account, error) {
	var acc Account
	row := q.callStoredFunction(ctx, "get_account_for_update", id)
	err := scanAccount(row, &acc)

	return acc, err
}

type ListAccountsParams struct {
	Owner  string `json:"owner"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

func (q *Queries) ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error) {
	rows, err := q.callStoredFunctionRows(ctx, "list_accounts", arg.Owner, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Account
	for rows.Next() {
		var acc Account
		if err := scanAccount(rows, &acc); err != nil {
			return nil, err
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

func (q *Queries) UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error) {
	row := q.callStoredFunction(ctx, "update_account", arg.ID, arg.Balance)
	var acc Account
	if err := scanAccount(row, &acc); err != nil {
		return Account{}, err
	}
	return acc, nil
}

// scanAccount abstracts the row scanning logic for an Account. Works with *sql.Row and *sql.Rows.
func scanAccount(scanner interface {
	Scan(dest ...interface{}) error
}, acc *Account) error {
	err := scanner.Scan(&acc.ID, &acc.Owner, &acc.Balance, &acc.Unit, &acc.CreatedAt)
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
