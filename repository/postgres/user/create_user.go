package user

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/aradwann/eenergy/entities"
	"github.com/aradwann/eenergy/repository/postgres/common"
)

type CreateUserParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	FullName       string `json:"fullname"`
	Email          string `json:"email"`
}

func (r *userRepository) createUser(ctx context.Context, tx *sql.Tx, arg CreateUserParams) (*entities.User, error) {
	user := &entities.User{}
	r.logger.Info("CreateUser", slog.String("username", arg.Username))
	row := common.CallStoredFunctionTx(ctx, tx, "create_user",
		arg.Username,
		arg.HashedPassword,
		arg.FullName,
		arg.Email)

	err := scanUserFromRow(row, user)
	if err != nil {
		slog.Error("error scaning the created user", slog.String("error message", err.Error()))
		return user, err
	}
	return user, nil
}
