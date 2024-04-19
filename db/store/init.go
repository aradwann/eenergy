package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/aradwann/eenergy/util"
)

// InitDatabase initializes the database connection and performs migrations.
func InitDatabase(config util.Config) Store {
	// Initialize database connection.
	dbConn, err := initDatabaseConn(config)
	if err != nil {
		handleError("Unable to connect to database", err)
	}
	defer dbConn.Close()

	// Run database migrations.
	RunDBMigrations(dbConn, config.MigrationsURL)

	// Create store instance.
	store := newStore(dbConn)
	return store
}

// initDatabaseConn initializes the database connection.
func initDatabaseConn(config util.Config) (*sql.DB, error) {
	dbConn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		handleError("Unable to connect to database", err)
	}
	return dbConn, err
}

// handleError logs an error message and exits the program with status code 1.
func handleError(message string, err error) {
	slog.Error(fmt.Sprintf("%s: %v", message, err))
	os.Exit(1)
}
