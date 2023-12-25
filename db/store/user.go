package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type CreateUserParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	FullName       string `json:"fullname"`
	Email          string `json:"email"`
}

type UpdateUserParams struct {
	HashedPassword    sql.NullString `json:"hashed_password"`
	PasswordChangedAt sql.NullTime   `json:"password_changed_at"`
	FullName          sql.NullString `json:"fullname"`
	Email             sql.NullString `json:"email"`
	Username          string         `json:"username"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	user := User{}
	params := StoredProcedureParams{
		InParams:  []interface{}{arg.Username, arg.HashedPassword, arg.FullName, arg.Email, &user.PasswordChangedAt, &user.CreatedAt},
		OutParams: []interface{}{&user.Username, &user.HashedPassword, &user.FullName, &user.Email, &user.PasswordChangedAt, &user.CreatedAt},
	}

	row := q.callStoredProcedure(ctx, "create_user", params)
	err := scanUserFromRow(row, &user)
	return user, err
}

func (q *Queries) GetUser(ctx context.Context, username string) (User, error) {
	user := User{}
	params := StoredProcedureParams{
		InParams:  []interface{}{username, &user.Username, &user.HashedPassword, &user.FullName, &user.Email, &user.PasswordChangedAt, &user.CreatedAt},
		OutParams: []interface{}{&user.Username, &user.HashedPassword, &user.FullName, &user.Email, &user.PasswordChangedAt, &user.CreatedAt},
	}

	row := q.callStoredProcedure(ctx, "get_user", params)
	err := scanUserFromRow(row, &user)
	return user, err
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	user := User{}
	params := StoredProcedureParams{
		InParams:  []interface{}{arg.Username, arg.HashedPassword, arg.PasswordChangedAt, arg.FullName, arg.Email, &user.Username, &user.HashedPassword, &user.FullName, &user.Email, &user.PasswordChangedAt, &user.CreatedAt},
		OutParams: []interface{}{&user.Username, &user.HashedPassword, &user.FullName, &user.Email, &user.PasswordChangedAt, &user.CreatedAt},
	}

	row := q.callStoredProcedure(ctx, "update_user", params)
	err := scanUserFromRow(row, &user)
	return user, err
}

func scanUserFromRow(row *sql.Row, user *User) error {

	err := row.Scan(
		&user.Username,
		&user.HashedPassword,
		&user.FullName,
		&user.Email,
		&user.PasswordChangedAt,
		&user.CreatedAt,
	)
	// Check for errors after scanning
	if err := row.Err(); err != nil {
		// Handle row-related errors
		log.Fatal(err)
		return err
	}
	if err != nil {
		// Check for a specific error related to the scan
		if err == sql.ErrNoRows {
			fmt.Println("No rows were returned.")
			return err
		} else {
			// Handle other scan-related errors
			log.Fatal(err)
			return err
		}
	}

	return nil
}
