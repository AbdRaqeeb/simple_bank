package db

import (
	"database/sql"
	"github.com/AbdRaqeeb/simple_bank/util"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	var err error

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDb, err = sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		log.Fatal("can't connect to database", err)
	}

	testQueries = New(testDb)

	os.Exit(m.Run())
}
