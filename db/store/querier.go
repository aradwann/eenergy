package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	GetUser(ctx context.Context, username string) (User, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
