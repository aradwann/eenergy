package db

// import (
// 	"context"
// 	"fmt"
// )

// // TransferTxParams contains the input parameters of the transfer transaction
// type TransferTxParams struct {
// 	FromAccountID int64 `json:"from_account_id"`
// 	ToAccountID   int64 `json:"to_account_id"`
// 	Amount        int64 `json:"amount"`
// }

// // TransferTxResult is the result of the tranfer transaction
// type TransferTxResult struct {
// 	Transfer    Transfer `json:"transfer"`
// 	FromAccount Account  `json:"from_account"`
// 	ToAccount   Account  `json:"to_account"`
// 	FromEntry   Entry    `json:"from_entry"`
// 	ToEntry     Entry    `json:"to_entry"`
// }

// // TransferTx performs a money/energy unit transfer from one account to another.
// // It creates a new transfer record, adds new account entries, and updates account balances in a single database transaction.
// func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
// 	var result TransferTxResult

// 	err := store.execTx(ctx, func(q *Queries) error {
// 		var err error

// 		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))
// 		if err != nil {
// 			return fmt.Errorf("failed to create transfer: %w", err)
// 		}

// 		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
// 			AccountID: arg.FromAccountID,
// 			Amount:    -arg.Amount,
// 		})
// 		if err != nil {
// 			return fmt.Errorf("failed to create entry for from account: %w", err)
// 		}

// 		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
// 			AccountID: arg.ToAccountID,
// 			Amount:    arg.Amount,
// 		})
// 		if err != nil {
// 			return fmt.Errorf("failed to create entry for to account: %w", err)
// 		}

// 		// Helper function to update account balance.
// 		updateBalance := func(accountID int64, amount int64) (Account, error) {
// 			return q.AddAccountBalance(ctx, AddAccountBalanceParams{
// 				ID:     accountID,
// 				Amount: amount,
// 			})
// 		}

// 		// Ensure account balances are updated in a consistent order to avoid deadlocks.
// 		accounts := []int64{arg.FromAccountID, arg.ToAccountID}
// 		amounts := []int64{-arg.Amount, arg.Amount}
// 		if arg.FromAccountID > arg.ToAccountID {
// 			accounts[0], accounts[1] = accounts[1], accounts[0]
// 			amounts[0], amounts[1] = amounts[1], amounts[0]
// 		}

// 		for i, accountID := range accounts {
// 			account, err := updateBalance(accountID, amounts[i])
// 			if err != nil {
// 				return fmt.Errorf("failed to update balance for account %d: %w", accountID, err)
// 			}
// 			if i == 0 {
// 				result.FromAccount = account
// 			} else {
// 				result.ToAccount = account
// 			}
// 		}

// 		return nil
// 	})

// 	return result, err
// }
