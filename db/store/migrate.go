package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func RunDBMigrations(db *sql.DB, migrationsURL string) {

	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		// log.Fatal().Msg("cannot create postgres driver")
		fmt.Printf("cannot create postgres driver %s", err)
	}
	migration, err := migrate.NewWithDatabaseInstance(
		migrationsURL,
		"eenergy", driver)
	if err != nil {
		// log.Fatal().Msg("cannot create new migrate instance")
		fmt.Printf("cannot create new migrate instance %s", err)
	}
	migration.Up()
	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		// log.Fatal().Msg("failed to run migrate up")
		fmt.Printf("failed to run migrate up %s", err)

	}

	// log.Info().Msg("DB migrated successfully")
	fmt.Println("DB migrated successfully")

	// Run unversioned migrations
	err = runUnversionedMigrations(db, "./db/migrations/procs")
	if err != nil {
		fmt.Println("Error applying unversioned migrations:", err)
		os.Exit(1)
	}

	fmt.Println("Unversioned migrations applied successfully")

}

// Get a list of SQL files in the migration directory
func getSQLFiles(migrationDir string) ([]string, error) {
	var sqlFiles []string

	err := filepath.WalkDir(migrationDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Check if the file has a .sql extension
		if strings.HasSuffix(path, ".sql") {
			sqlFiles = append(sqlFiles, path)
		}

		return nil
	})

	return sqlFiles, err
}

func runUnversionedMigrations(db *sql.DB, migrationDir string) error {

	sqlFiles, err := getSQLFiles(migrationDir)

	if err != nil {
		return err
	}
	// Sort files to ensure execution order
	// Note: You may need a custom sorting logic if file names include version numbers
	// For simplicity, we assume alphabetical order here.
	// Sorting ensures that the files are executed in the correct order.
	sortFiles(sqlFiles)

	// Execute each SQL file
	for _, file := range sqlFiles {

		contents, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		// Execute the SQL content
		_, err = db.Exec(string(contents))
		if err != nil {
			return fmt.Errorf("error executing SQL file %s: %w", file, err)
		}

		fmt.Printf("Executed migration: %s\n", file)
	}

	return nil
}

// Simple alphabetical sorting function
func sortFiles(files []string) {
	sort.Slice(files, func(i, j int) bool {
		return strings.Compare(files[i], files[j]) < 0
	})
}
