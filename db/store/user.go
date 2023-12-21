package db

import (
	"context"
	"database/sql"
)

type CreateUserParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	Fullname       string `json:"fullname"`
	Email          string `json:"email"`
}

// CreateUser calls the create_user stored procedure
func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	var user User
	interf := []interface{}{
		arg.Username,
		arg.HashedPassword,
		arg.Fullname,
		arg.Email,
		&user.Username,
		&user.HashedPassword,
		&user.Fullname,
		&user.Email,
		&user.PasswordChangedAt,
		&user.CreatedAt,
	}
	row, err := q.callStoredProcedure(ctx, "create_user", interf...)
	if err != nil {
		return User{}, err
	}
	// Execute the stored procedure and scan the results into the variables
	err = row.Scan(
		&user.Username,
		&user.HashedPassword,
		&user.Fullname,
		&user.Email,
		&user.PasswordChangedAt,
		&user.CreatedAt,
	)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (q *Queries) GetUser(ctx context.Context, username string) (User, error) {
	var user User
	interf := []interface{}{
		username,
		&user.Username,
		&user.HashedPassword,
		&user.Fullname,
		&user.Email,
		&user.PasswordChangedAt,
		&user.CreatedAt,
	}
	row, err := q.callStoredProcedure(ctx, "get_user", interf...)
	if err != nil {
		return User{}, err
	}
	err = row.Scan(
		&user.Username,
		&user.HashedPassword,
		&user.Fullname,
		&user.Email,
		&user.PasswordChangedAt,
		&user.CreatedAt,
	)
	return user, err
}

type UpdateUserParams struct {
	HashedPassword    sql.NullString `json:"hashed_password"`
	PasswordChangedAt sql.NullTime   `json:"password_changed_at"`
	Fullname          sql.NullString `json:"fullname"`
	Email             sql.NullString `json:"email"`
	Username          string         `json:"username"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	var user User
	interf := []interface{}{
		arg.Username,
		arg.HashedPassword,
		arg.PasswordChangedAt,
		arg.Fullname,
		arg.Email,
		&user.Username,
		&user.HashedPassword,
		&user.Fullname,
		&user.Email,
		&user.PasswordChangedAt,
		&user.CreatedAt,
	}
	row, err := q.callStoredProcedure(ctx, "update_user", interf...)
	if err != nil {
		return User{}, err
	}
	// Execute the stored procedure and scan the results into the variables
	err = row.Scan(
		&user.Username,
		&user.HashedPassword,
		&user.Fullname,
		&user.Email,
		&user.PasswordChangedAt,
		&user.CreatedAt,
	)

	if err != nil {
		return User{}, err
	}

	return user, nil
}
