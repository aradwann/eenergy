package db

import (
	"context"
	"database/sql"
	"errors"
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

	err := scanAccountFromRow(row, &acc)
	if err != nil {
		return acc, err
	}

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

	err := scanAccountFromRow(row, &acc)
	if err != nil {
		return acc, err
	}

	return acc, err
}

func (q *Queries) DeleteAccount(ctx context.Context, id int64) error {
	row := q.callStoredFunction(ctx, "delete_account",
		id,
	)
	return row.Err()
}

func (q *Queries) GetAccount(ctx context.Context, id int64) (Account, error) {
	var acc Account
	row := q.callStoredFunction(ctx, "get_account",
		id,
	)

	err := scanAccountFromRow(row, &acc)
	if err != nil {
		return acc, err
	}

	return acc, err
}

func (q *Queries) GetAccountForUpdate(ctx context.Context, id int64) (Account, error) {
	row := q.db.QueryRow(ctx, getAccountForUpdate, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Unit,
		&i.CreatedAt,
	)
	return i, err
}

type ListAccountsParams struct {
	Owner  string `json:"owner"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

func (q *Queries) ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error) {

	rows, err := q.callStoredFunctionRows(ctx, "list_accounts",
		arg.Owner, arg.Limit, arg.Offset,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Account{}
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Balance,
			&i.Unit,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
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
	row := q.db.QueryRow(ctx, updateAccount, arg.ID, arg.Balance)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Unit,
		&i.CreatedAt,
	)
	return i, err
}

func scanAccountFromRow(row *sql.Row, acc *Account) error {
	err := row.Scan(
		&acc.ID,
		&acc.Owner,
		&acc.Balance,
		&acc.Unit,
		&acc.CreatedAt,
	)

	// Check for errors after scanning
	if err != nil {
		// Handle scan-related errors
		if errors.Is(err, ErrRecordNotFound) {
			// fmt.Println("No rows were returned.")
			return err
		} else {
			// Log and return other scan-related errors
			log.Printf("Error scanning row: %s", err)
			return err
		}
	}

	return nil
}
