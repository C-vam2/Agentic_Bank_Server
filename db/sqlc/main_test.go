package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgres://root:secret@localhost:5432/agentic_bank_server?sslmode=disable"
)

var testQueries *Queries
var dbConn *sql.DB

func TestMain(m *testing.M) {
	var err error
	dbConn, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to the database", err)
	}

	testQueries = New(dbConn)

	os.Exit(m.Run())
}
