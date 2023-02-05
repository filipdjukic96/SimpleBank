package main

import (
	"bank/api"
	db "bank/db/sqlc"
	"bank/db/util"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config")
	}

	connection, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(connection)
	server := api.NewServer(store)

	err = server.StartServer(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
