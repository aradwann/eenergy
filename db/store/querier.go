package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	// User queries
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetUser(ctx context.Context, username string) (User, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
	// Session queries
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	// Email queries
	CreateVerifyEmail(ctx context.Context, arg CreateVerifyEmail) (VerifyEmail, error)
	// Account queries
	AddAccountBalance(ctx context.Context, arg AddAccountBalanceParams) (Account, error)
	CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error)
	DeleteAccount(ctx context.Context, id int64) error
	GetAccount(ctx context.Context, id int64) (Account, error)
	GetAccountForUpdate(ctx context.Context, id int64) (Account, error)
	ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error)
	UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error)
	// Transfer queries
	CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error)
	GetTransfer(ctx context.Context, id int64) (Transfer, error)
	ListTransfers(ctx context.Context, arg ListTransfersParams) ([]Transfer, error)
	// Entry queries
	CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error)
	GetEntry(ctx context.Context, id int64) (Entry, error)
	ListEntries(ctx context.Context, arg ListEntriesParams) ([]Entry, error)
}

var _ Querier = (*Queries)(nil)
