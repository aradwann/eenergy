package db

// // CreateUserTxParams represents the parameters for creating a user transaction.
// type CreateUserTxParams struct {
// 	CreateUserParams // Embedding CreateUserParams for reuse
// 	AfterCreate      func(user User) error
// }

// // CreateUserTxResult represents the result of creating a user transaction.
// type CreateUserTxResult struct {
// 	User User
// }

// // CreateUserTx executes a transaction for creating a user.
// func (r *userRepository) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
// 	slog.Info("Create User Transaction")
// 	var result CreateUserTxResult

// 	err := store.execTx(ctx, func(q *Queries) error {
// 		var err error

// 		// Create user using provided parameters.
// 		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
// 		if err != nil {
// 			return err
// 		}

// 		// Execute user-defined callback after creating the user.
// 		return arg.AfterCreate(result.User)
// 	})

// 	return result, err
// }
