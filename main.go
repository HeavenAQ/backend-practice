package main

import (
	"database/sql"
	"log"

	"github.com/HeavenAQ/simple-bank/api"
	db "github.com/HeavenAQ/simple-bank/db/sqlc"
	"github.com/HeavenAQ/simple-bank/utils"
	_ "github.com/lib/pq"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("Failed to load configurations")
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Failed to connect to the database")
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)
	if err = server.StartServer(config.ServerAddress); err != nil {
		log.Fatal("Failed to start server", err)
	}
}
