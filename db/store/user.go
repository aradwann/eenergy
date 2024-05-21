package db

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"log/slog"
)

type CreateUserParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	FullName       string `json:"fullname"`
	Email          string `json:"email"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	var user User
	slog.Info("CreateUser", slog.String("username", arg.Username))
	row := q.callStoredFunction(ctx, "create_user",
		arg.Username,
		arg.HashedPassword,
		arg.FullName,
		arg.Email)

	err := scanUserFromRow(row, &user)
	if err != nil {
		slog.Error("error scaning the created user", slog.String("error message", err.Error()))
		return user, err
	}
	return user, nil
}

func (q *Queries) GetUser(ctx context.Context, user_id int64) (User, error) {
	var user User
	row := q.callStoredFunction(ctx, "get_user", user_id)
	err := scanUserFromRow(row, &user)
	if err != nil {
		return user, err
	}
	return user, nil
}

type UpdateUserParams struct {
	ID                int64          `json:"id"`
	HashedPassword    sql.NullString `json:"hashed_password"`
	PasswordChangedAt sql.NullTime   `json:"password_changed_at"`
	FullName          sql.NullString `json:"fullname"`
	Email             sql.NullString `json:"email"`
	Username          sql.NullString `json:"username"`
	IsEmailVerified   sql.NullBool   `json:"is_email_verified"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	var user User
	params := []interface{}{
		arg.ID,
		arg.Username,
		arg.HashedPassword,
		arg.PasswordChangedAt,
		arg.FullName,
		arg.Email,
		arg.IsEmailVerified,
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
		&user.ID,
		&user.Username,
		&user.HashedPassword,
		&user.FullName,
		&user.Email,
		&user.PasswordChangedAt,
		&user.CreatedAt,
		&user.IsEmailVerified,
		&user.RoleID,
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
