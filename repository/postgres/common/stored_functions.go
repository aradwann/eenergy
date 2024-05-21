package common

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
)

// callStoredFunction executes a stored function and returns a single row result.
func CallStoredFunction(ctx context.Context, db *sql.DB, functionName string, params ...interface{}) *sql.Row {
	placeholders := generateParamPlaceholders(len(params))

	// Construct the SQL statement to call the stored function.
	sqlStatement := fmt.Sprintf(`SELECT * FROM %s(%s)`, functionName, placeholders)
	slog.Info("PostgreSQL function called",
		slog.String("function name", functionName),
		// slog.Any("params", params), // TODO: filter out sensitive info
		slog.String("SQL statement", sqlStatement),
	)

	return db.QueryRowContext(ctx, sqlStatement, params...)
}

// callStoredFunctionRows executes a stored function and returns multiple rows result.
func CallStoredFunctionRows(ctx context.Context, db *sql.DB, functionName string, params ...interface{}) (*sql.Rows, error) {
	placeholders := generateParamPlaceholders(len(params))

	// Construct the SQL statement to call the stored function.
	sqlStatement := fmt.Sprintf(`SELECT * FROM %s(%s)`, functionName, placeholders)

	return db.QueryContext(ctx, sqlStatement, params...)
}

// generateParamPlaceholders generates placeholders for SQL parameters.
func generateParamPlaceholders(count int) string {
	placeholders := make([]string, count)
	for i := 1; i <= count; i++ {
		placeholders[i-1] = fmt.Sprintf("$%d", i)
	}
	return strings.Join(placeholders, ", ")
}
