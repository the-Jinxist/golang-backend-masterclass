package api

import (
	db "backend_masterclass/db/sqlc"
	"backend_masterclass/token"
	"backend_masterclass/util"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

//This struct [Server] will serve all our HTTP requests for our banking services
type Server struct {
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
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

	//Here, we register the custom validator, we created for validating currency here
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validcurrency", validCurrency)
	}

	server.serveRouter()
	return server, nil
}

func (server *Server) serveRouter() {

	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.POST("/users/login", server.loginUser)

	//The above routes /accounts and /users/login don't need any authorization, so we create the endpoints that need
	//an authorization token after registering those endpoints

	authGroups := router.Group("/", authMiddleWare(server.tokenMaker))

	authGroups.GET("/account/:id", server.getAccount)
	authGroups.GET("/accounts", server.listAccounts)
	authGroups.POST("/transfers", server.transferMoney)

	authGroups.GET("/user", server.getUser)
	authGroups.POST("/user", server.createUser)

	server.router = router
}

//Startes the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
