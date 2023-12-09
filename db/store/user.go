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

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	var i User
	row, err := q.callStoredFunction(ctx, "create_user",
		arg.Username,
		arg.HashedPassword,
		arg.Fullname,
		arg.Email)
	if err != nil {
		return User{}, err
	}
	err = row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.Fullname,
		&i.Email,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

func (q *Queries) GetUser(ctx context.Context, username string) (User, error) {
	row, err := q.callStoredFunction(ctx, "get_user", username)
	if err != nil {
		return User{}, err
	}
	var i User
	err = row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.Fullname,
		&i.Email,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

type UpdateUserParams struct {
	HashedPassword    sql.NullString `json:"hashed_password"`
	PasswordChangedAt sql.NullTime   `json:"password_changed_at"`
	Fullname          sql.NullString `json:"fullname"`
	Email             sql.NullString `json:"email"`
	Username          string         `json:"username"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row, err := q.callStoredFunction(ctx, "update_user",
		arg.Username,
		arg.HashedPassword,
		arg.PasswordChangedAt,
		arg.Fullname,
		arg.Email)
	if err != nil {
		return User{}, err
	}
	var i User
	err = row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.Fullname,
		&i.Email,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}
