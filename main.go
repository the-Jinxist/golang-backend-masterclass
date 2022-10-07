package main

import (
	"backend_masterclass/api"
	db "backend_masterclass/db/sqlc"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

const (
	//We used the lib/pq library to use the correct postgres driver
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal(err.Error())
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("Cannot start server")
	}

}
