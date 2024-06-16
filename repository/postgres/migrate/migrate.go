package migrate

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// RunDBMigrations runs database migrations using the provided database connection and migrations URL.
func RunDBMigrations(db *sql.DB, migrationsURL string) {
	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		slog.Error("cannot create postgres driver", slog.String("error", err.Error()))
		os.Exit(1) // Exit on error
	}

	migration, err := migrate.NewWithDatabaseInstance(migrationsURL, "eenergy", driver)
	if err != nil {
		slog.Error("cannot create new migrate instance", slog.String("error", err.Error()))
		os.Exit(1) // Exit on error
	}

	err = migration.Up()
	if err != nil && err != migrate.ErrNoChange {
		slog.Error("failed to run migrate up", slog.String("error", err.Error()))
		os.Exit(1) // Exit on error
	}

	slog.Info("DB migrated successfully")

	// Concatenate the migrations URL with "/functions"
	functionsMigrationsURL := filepath.Join(migrationsURL, "functions")

	// Run unversioned migrations
	err = runUnversionedMigrations(db, functionsMigrationsURL)
	if err != nil {
		slog.Error("Error applying unversioned migrations:", slog.String("error", err.Error()))
		os.Exit(1) // Exit on error
	}

	slog.Info("Unversioned migrations applied successfully")
}

// Get a list of SQL files in the migration directory
func getSQLFiles(migrationDir string) ([]string, error) {
	var sqlFiles []string

	err := filepath.Walk(migrationDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // Immediately return errors for short-circuiting
		}

		// Only process regular files with .sql extension
		if info.Mode().IsRegular() && strings.HasSuffix(path, ".sql") {
			sqlFiles = append(sqlFiles, path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return sqlFiles, nil
}

func runUnversionedMigrations(db *sql.DB, migrationDir string) error {
	sqlFiles, err := getSQLFiles(migrationDir)
	if err != nil {
		return err
	}

	// Execute each SQL file
	for _, file := range sqlFiles {
		// log.Printf("Executing SQL file: %s", file)

		contents, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("error reading SQL file %s: %w", file, err)
		}

		// Execute the SQL content
		_, err = db.Exec(string(contents))
		if err != nil {
			return fmt.Errorf("error executing SQL file %s: %w", file, err)
		}

		// log.Printf("Finished executing SQL file: %s", file)
	}

	return nil
}

// Simple alphabetical sorting function
// func sortFiles(files []string) {
// 	sort.Slice(files, func(i, j int) bool {
// 		return strings.Compare(files[i], files[j]) < 0
// 	})
// }
