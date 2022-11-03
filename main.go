package main

import (
	"backend_masterclass/api"
	db "backend_masterclass/db/sqlc"
	"backend_masterclass/util"
	"database/sql"
	"log"

	_ "github.com/golang-jwt/jwt/v4"
	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/google/uuid"
	_ "github.com/lib/pq"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err.Error())
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err.Error())
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot start server")
	}

}
