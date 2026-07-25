package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/Agentic_Bank_Server/utils"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var dbConn *sql.DB

func TestMain(m *testing.M) {
	var err error
	config, err := utils.LoadConfig("../../")

	if err != nil {
		log.Fatal("cannot load configs ", err)
		return
	}

	dbConn, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to the database", err)
	}

	testQueries = New(dbConn)

	os.Exit(m.Run())
}
