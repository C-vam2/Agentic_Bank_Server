package main

import (
	"database/sql"
	"log"

	"github.com/Agentic_Bank_Server/api"
	db "github.com/Agentic_Bank_Server/db/sqlc"
	"github.com/Agentic_Bank_Server/utils"

	_ "github.com/lib/pq"
)

func main() {
	config, err := utils.LoadConfig(".")

	if err != nil {
		log.Fatal("Unable to load configs ", err)
		return
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to the database ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	if err := server.Start(config.ServerAddress); err != nil {
		log.Fatal("cannot start server on port :8080")
	}
}
