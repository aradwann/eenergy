package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/aradwann/eenergy/util"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..", "app")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDB, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testStore = NewStore(testDB)
	os.Exit(m.Run())
}
