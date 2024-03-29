package main

import (
	"database/sql"
	"fmt"
	"os"

	mig "github.com/aradwann/eenergy/db/store"
	"github.com/aradwann/eenergy/util"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	config, err := util.LoadConfig(".", ".env")
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

	mig.RunDBMigrations(db, config.MigrationsURL)

}
