package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/aradwann/eenergy/util"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	config, err := util.LoadConfig(".", "app")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}
	db, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	runDBMigrations(db, config.MigrationsURL)

	var greeting string
	err = db.QueryRow("select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(greeting)
}

func runDBMigrations(db *sql.DB, migrationsURL string) {

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

}
