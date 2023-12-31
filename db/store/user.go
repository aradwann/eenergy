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

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	var user User
	row := q.callStoredFunction(ctx, "create_user",
		arg.Username,
		arg.HashedPassword,
		arg.FullName,
		arg.Email)

	err := scanUserFromRow(row, &user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (q *Queries) GetUser(ctx context.Context, username string) (User, error) {
	var user User
	row := q.callStoredFunction(ctx, "get_user", username)
	err := scanUserFromRow(row, &user)
	if err != nil {
		return user, err
	}
	return user, nil
}

type UpdateUserParams struct {
	HashedPassword    sql.NullString `json:"hashed_password"`
	PasswordChangedAt sql.NullTime   `json:"password_changed_at"`
	FullName          sql.NullString `json:"fullname"`
	Email             sql.NullString `json:"email"`
	Username          string         `json:"username"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	var user User
	params := []interface{}{
		arg.Username,
		arg.HashedPassword,
		arg.PasswordChangedAt,
		arg.FullName,
		arg.Email,
	}
	row := q.callStoredFunction(ctx, "update_user", params...)

	// Execute the stored procedure and scan the results into the variables
	err := scanUserFromRow(row, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
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
