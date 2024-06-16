package postgres

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	migrate "github.com/aradwann/eenergy/repository/postgres/migrate"
	"github.com/aradwann/eenergy/util"
)

// InitDatabase initializes the database connection and performs migrations.
func InitDatabase(config util.Config) *sql.DB {
	// Initialize database connection.
	dbConn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		handleError("Unable to connect to database", err)
	}
	// Run database migrations.
	migrate.RunDBMigrations(dbConn, config.MigrationsURL)

	return dbConn
}

// handleError logs an error message and exits the program with status code 1.
func handleError(message string, err error) {
	slog.Error(fmt.Sprintf("%s: %v", message, err))
	os.Exit(1)
}
