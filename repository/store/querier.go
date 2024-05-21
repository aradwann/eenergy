package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	// User queries
	// Session queries
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	// Email queries
	CreateVerifyEmail(ctx context.Context, arg CreateVerifyEmail) (VerifyEmail, error)
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
