package gapi

import (
	db "backend_masterclass/db/sqlc"
	"backend_masterclass/pb"
	"backend_masterclass/token"
	"backend_masterclass/util"
	"fmt"
)

//This server will serve gRPC requests for our banking service

//We added pb.UnimplementedSimpleBankServer to enable forward compatibility. This means the server can accept calls to CreateUser and LoginUser before they are
//implemented
type Server struct {
	pb.UnimplementedSimpleBankServer
	store      db.Store
	tokenMaker token.Maker
	config     util.Config
}

func NewServer(config util.Config, store db.Store) (*Server, error) {

	tokenMaker, err := token.NewJwtMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store:      store,
		config:     config,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
