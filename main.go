package main

import (
	"database/sql"
	"log"

	"github.com/Agentic_Bank_Server/api"
	db "github.com/Agentic_Bank_Server/db/sqlc"

	_ "github.com/lib/pq"
)

const (
	dbDriver    = "postgres"
	dbSource    = "postgres://root:secret@localhost:5432/agentic_bank_server?sslmode=disable"
	portAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("cannot connect to the database ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	if err := server.Start(portAddress); err != nil {
		log.Fatal("cannot start server on port :8080")
	}
}
