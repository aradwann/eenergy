package user

import (
	"context"
	"database/sql"

	"github.com/aradwann/eenergy/entities"
)

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
