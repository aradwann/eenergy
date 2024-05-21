package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
)

// DBTX defines the interface for database transactions.
type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

// New creates a new instance of Queries using the provided database transaction.
func New(db DBTX) *Queries {
	return &Queries{db: db}
}

// Queries provides methods for executing database queries.
type Queries struct {
	db DBTX
}

// callStoredFunction executes a stored function and returns a single row result.
func (q *Queries) callStoredFunction(ctx context.Context, functionName string, params ...interface{}) *sql.Row {
	placeholders := generateParamPlaceholders(len(params))

	// Construct the SQL statement to call the stored function.
	sqlStatement := fmt.Sprintf(`SELECT * FROM %s(%s)`, functionName, placeholders)
	slog.Info("PostgreSQL function called",
		slog.String("function name", functionName),
		// slog.Any("params", params), // TODO: filter out sensitive info
		slog.String("SQL statement", sqlStatement),
	)

	return q.db.QueryRowContext(ctx, sqlStatement, params...)
}

// callStoredFunctionRows executes a stored function and returns multiple rows result.
func (q *Queries) callStoredFunctionRows(ctx context.Context, functionName string, params ...interface{}) (*sql.Rows, error) {
	placeholders := generateParamPlaceholders(len(params))

	// Construct the SQL statement to call the stored function.
	sqlStatement := fmt.Sprintf(`SELECT * FROM %s(%s)`, functionName, placeholders)

	return q.db.QueryContext(ctx, sqlStatement, params...)
}

// generateParamPlaceholders generates placeholders for SQL parameters.
func generateParamPlaceholders(count int) string {
	placeholders := make([]string, count)
	for i := 1; i <= count; i++ {
		placeholders[i-1] = fmt.Sprintf("$%d", i)
	}
	return strings.Join(placeholders, ", ")
}
