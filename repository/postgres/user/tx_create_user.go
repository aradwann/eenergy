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
	slog.Info("CreateUser", slog.String("username", arg.Username))
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

// CreateUserTxParams represents the parameters for creating a user transaction.
type CreateUserTxParams struct {
	CreateUserParams // Embedding CreateUserParams for reuse
	AfterCreate      func(user *entities.User) error
}

// CreateUserTxResult represents the result of creating a user transaction.
type CreateUserTxResult struct {
	User *entities.User
}

// CreateUserTx executes a transaction for creating a user.
func (r *userRepository) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	r.logger.Info("Create User Transaction")
	var result CreateUserTxResult

	err := r.transaction.ExecTx(ctx, func(tx *sql.Tx) error {
		var err error

		// Create user using provided parameters.
		result.User, err = r.createUser(ctx, tx, arg.CreateUserParams)
		if err != nil {
			return err
		}
		// Execute user-defined callback after creating the user.
		if arg.AfterCreate != nil {
			return arg.AfterCreate(result.User)
		}

		return nil
	})

	return result, err
}
