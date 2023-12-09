package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db DBTX
}

// callStoredFunction executes a stored function with parameters.
func (q *Queries) callStoredFunction(ctx context.Context, functionName string, params ...interface{}) (*sql.Row, error) {
	placeholders := make([]string, len(params))
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	sqlStatement := fmt.Sprintf(`SELECT * FROM %s(%s)`, functionName, strings.Join(placeholders, ", "))

	return q.db.QueryRowContext(ctx, sqlStatement, params...), nil
}
