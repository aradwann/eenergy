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

type StoredProcedureParams struct {
	InParams  []interface{}
	OutParams []interface{}
}

func (q *Queries) callStoredProcedure(ctx context.Context, procedureName string, params StoredProcedureParams) *sql.Row {
	sqlStatement := fmt.Sprintf(`CALL %s(%s)`, procedureName, generateParamPlaceholders(len(params.InParams)))

	return q.db.QueryRowContext(
		ctx,
		sqlStatement,
		params.InParams...,
	)
}

func generateParamPlaceholders(count int) string {
	placeholders := make([]string, count)
	for i := 1; i <= count; i++ {
		placeholders[i-1] = fmt.Sprintf("$%d", i)
	}
	return strings.Join(placeholders, ", ")
}
