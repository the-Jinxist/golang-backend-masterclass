package main

import (
	// "backend_masterclass/api"

	db "backend_masterclass/db/sqlc"
	"backend_masterclass/gapi"
	"backend_masterclass/pb"
	"backend_masterclass/util"
	"database/sql"
	"log"
	"net"

	_ "github.com/golang-jwt/jwt/v4"
	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/google/uuid"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	reflection "google.golang.org/grpc/reflection"
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
	// runHTTPServer(config, store)
	runGRPCServer(config, store)

}

// func runHTTPServer(config util.Config, store db.Store) {
// 	server, err := api.NewServer(config, store)

// 	if err != nil {
// 		log.Fatal("Cannot create server")
// 	}

// 	err = server.Start(config.ServerAddress)
// 	if err != nil {
// 		log.Fatal("Cannot start server")
// 	}
// }

func runGRPCServer(config util.Config, store db.Store) {
	log.Println("GRPC getting ready..")
	server, err := gapi.NewServer(config, store)

	log.Println("new server created")

	if err != nil {
		log.Fatal("Cannot create server")
	}

	grpcServer := grpc.NewServer()
	log.Println("new gRPC server created")

	reflection.Register(grpcServer)
	pb.RegisterSimpleBankServer(grpcServer, server)
	log.Println("new gRPC server registered")

	listener, err := net.Listen("tcp", config.ServerAddress)
	log.Println("listener done listener")
	if err != nil {
		log.Fatal("Cannot create listener")
	}

	log.Printf("server started at %s", listener.Addr())
	err1 := grpcServer.Serve(listener)

	if err1 != nil {
		log.Fatalf("Error occurred while running sever: %s", err1.Error())
	}

}
