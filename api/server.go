package api

import (
	db "backend_masterclass/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

//This struct [Server] will serve all our HTTP requests for our banking services
type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()

	//Here, we register the custom validator, we created for validating currency here
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validcurrency", validCurrency)
	}

	//Add routes to the router in a bit
	router.POST("/accounts", server.createAccount)
	router.GET("/account/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.POST("/transfers", server.transferMoney)

	router.GET("/user", server.getUser)
	router.POST("/user", server.createUser)

	server.router = router
	return server
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
