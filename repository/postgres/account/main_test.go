package account

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/aradwann/eenergy/util"
)

var testAccRepo AccountRepository

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../..", ".env")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDB, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testAccRepo = NewAccountRepository(testDB)
	os.Exit(m.Run())
}
