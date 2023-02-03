package main

import (
	"bank/api"
	db "bank/db/sqlc"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	connection, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(connection)
	server := api.NewServer(store)

	err = server.StartServer(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
