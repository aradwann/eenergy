package db

import (
	"database/sql"
)

// SQLStore provides all functions to execute db queries and transactions
type Store interface {
	Querier
}

// SQLStore provides all functions to execute db queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore create a new Store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}
