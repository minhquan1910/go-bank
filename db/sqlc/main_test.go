package sqlc

import (
	"bank/util"
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	config, err := util.LoadConfig("../../app.env")
	if err != nil {
		log.Fatal("cannot load config ", err)
	}
	testDB, err = sql.Open(config.DBDrive, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database ", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
