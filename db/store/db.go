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

// callStoredProcedure executes a stored procedure with parameters.
func (q *Queries) callStoredProcedure(ctx context.Context, procedureName string, params ...interface{}) (*sql.Row, error) {
	// Construct placeholders for the parameters
	placeholders := make([]string, len(params))
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	// Construct the SQL statement for calling the stored procedure
	sqlStatement := fmt.Sprintf(`CALL %s(%s)`, procedureName, strings.Join(placeholders, ", "))

	// Execute the stored procedure and return the result
	return q.db.QueryRowContext(ctx, sqlStatement, params...), nil
}
