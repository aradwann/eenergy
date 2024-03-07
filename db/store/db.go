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

// type storedProcedureParams struct {
// 	InParams  []interface{}
// 	OutParams []interface{}
// }

// func (q *Queries) callStoredProcedure(ctx context.Context, procedureName string, params storedProcedureParams) *sql.Row {
// 	sqlStatement := fmt.Sprintf(`CALL %s(%s)`, procedureName, generateParamPlaceholders(len(params.InParams)))

//		return q.db.QueryRowContext(
//			ctx,
//			sqlStatement,
//			params.InParams...,
//		)
//	}
func (q *Queries) callStoredFunction(ctx context.Context, functionName string, params ...interface{}) *sql.Row {
	// Assuming generateParamPlaceholders generates the placeholders for parameters
	placeholders := generateParamPlaceholders(len(params))

	// Use SELECT statement to call the stored function
	sqlStatement := fmt.Sprintf(`SELECT * FROM %s(%s)`, functionName, placeholders)

	return q.db.QueryRowContext(ctx, sqlStatement, params...)
}

func (q *Queries) callStoredFunctionRows(ctx context.Context, functionName string, params ...interface{}) (*sql.Rows, error) {
	// Assuming generateParamPlaceholders generates the placeholders for parameters
	placeholders := generateParamPlaceholders(len(params))

	// Use SELECT statement to call the stored function
	sqlStatement := fmt.Sprintf(`SELECT * FROM %s(%s)`, functionName, placeholders)

	return q.db.QueryContext(ctx, sqlStatement, params...)
}

func generateParamPlaceholders(count int) string {
	placeholders := make([]string, count)
	for i := 1; i <= count; i++ {
		placeholders[i-1] = fmt.Sprintf("$%d", i)
	}
	return strings.Join(placeholders, ", ")
}
