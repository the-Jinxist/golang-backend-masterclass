package main

import (
	// "backend_masterclass/api"

	db "backend_masterclass/db/sqlc"
	"backend_masterclass/gapi"
	"backend_masterclass/pb"
	"backend_masterclass/util"
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	_ "github.com/golang-jwt/jwt/v4"
	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	reflection "google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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
	go runGRPCGatewayServer(config, store)
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

	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatalf("Cannot create server: %s", err.Error())
	}

	grpcServer := grpc.NewServer()

	reflection.Register(grpcServer)
	pb.RegisterSimpleBankServer(grpcServer, server)

	listener, err := net.Listen("tcp", config.ServerAddress)
	if err != nil {
		log.Fatalf("Cannot create listener: %s", err.Error())
	}

	log.Printf("server started at %s", listener.Addr())
	err1 := grpcServer.Serve(listener)
	if err1 != nil {
		log.Fatalf("Error occurred while running sever: %s", err1.Error())
	}

}

//We're using in-process translation method because createUser and loginUser services are unary operations
func runGRPCGatewayServer(config util.Config, store db.Store) {

	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatalf("Cannot create gateway server: %s", err.Error())
	}

	//These option code was mainly used to make sure the field names in our json output
	//follow snake_case_convention
	jsonOption := runtime.WithMarshalerOption(
		runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		})

	//This code creates the new gateway
	grpcMux := runtime.NewServeMux(
		jsonOption,
	)

	//This creates a context object with a cancel. Cancelling a context is a way to prevent the object from doing
	//unnecessary work
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatalf("Cannot register gateway server: %s", err.Error())
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		log.Fatalf("Cannot create gateway listener: %s", err.Error())
	}

	log.Printf("grpcGateway server started at %s", listener.Addr())
	err1 := http.Serve(listener, mux)
	if err1 != nil {
		log.Fatalf("Error occurred while running sever: %s", err1.Error())
	}

}
