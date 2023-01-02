package main

import (
	// "backend_masterclass/api"

	db "backend_masterclass/db/sqlc"
	"backend_masterclass/gapi"
	"backend_masterclass/pb"
	"backend_masterclass/util"
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"

	_ "github.com/golang-jwt/jwt/v4"
	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	reflection "google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	_ "backend_masterclass/doc/statik"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msgf("cannot load config: ")
	}

	if config.Environment == "development" {
		//This code enables pretty logging, fancy stuff that I actually like
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	//We will run the DB migrations here
	runDBMigrations(config.MigrationURL, config.DBSource)

	store := db.NewStore(conn)
	// runHTTPServer(config, store)
	go runGRPCGatewayServer(config, store)
	runGRPCServer(config, store)

}

func runDBMigrations(migrationURL string, dbSourceString string) {
	migration, err := migrate.New(migrationURL, dbSourceString)
	if err != nil {
		log.Fatal().Msgf("cannot create new migrate instance: %s", err.Error())
	}

	err = migration.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal().Msgf("cannot run up migrations: %s", err.Error())
	}

	log.Info().Msgf("db migrated successfully")
}

// func runHTTPServer(config util.Config, store db.Store) {
// 	server, err := api.NewServer(config, store)

// 	if err != nil {
// 		log.Fatal().Msg("Cannot create server")
// 	}

// 	err = server.Start(config.ServerAddress)
// 	if err != nil {
// 		log.Fatal().Msg("Cannot start server")
// 	}
// }

func runGRPCServer(config util.Config, store db.Store) {

	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal().Msgf("Cannot create server: %s", err.Error())
	}

	//Here, wer're creating the unary interceptor
	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)

	//We're adding the interceptor to the server here, as well as creating a new grpc server
	grpcServer := grpc.NewServer(grpcLogger)

	reflection.Register(grpcServer)
	pb.RegisterSimpleBankServer(grpcServer, server)

	listener, err := net.Listen("tcp", config.ServerAddress)
	if err != nil {
		log.Fatal().Msgf("Cannot create listener: %s", err.Error())
	}

	log.Info().Msgf("server started at %s", listener.Addr())
	err1 := grpcServer.Serve(listener)
	if err1 != nil {
		log.Fatal().Msgf("Error occurred while running sever: %s", err1.Error())
	}

}

// We're using in-process translation method because createUser and loginUser services are unary operations
func runGRPCGatewayServer(config util.Config, store db.Store) {

	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal().Msgf("Cannot create gateway server: %s", err.Error())
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
		log.Fatal().Msgf("Cannot register gateway server: %s", err.Error())
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	//We're creating a file server and serving the front end files we copied into doc/swagger-ui
	// fs := http.FileServer(http.Dir("./doc/swagger-ui"))

	//We decided to use the statik library to serve our static front-end files inside our golang binary
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal().Msgf("Cannot create statik file system: %s", err.Error())
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		log.Fatal().Msgf("Cannot create gateway listener: %s", err.Error())
	}

	log.Info().Msgf("grpcGateway server started at %s", listener.Addr())
	err1 := http.Serve(listener, mux)
	if err1 != nil {
		log.Fatal().Msgf("Error occurred while running sever: %s", err1.Error())
	}

}
