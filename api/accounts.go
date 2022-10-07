package api

import (
	"database/sql"
	"net/http"

	db "backend_masterclass/db/sqlc"

	"github.com/gin-gonic/gin"
)

//Gin uses a validator under the hood to make sure that request bodies from the users are valid
//Using the "binding:required" tag is one of those forms of validation. This means the current field
//is required.
type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=NGN USD"`
}

func (server *Server) createAccount(ctx *gin.Context) {

	var request CreateAccountRequest
	err := ctx.ShouldBindJSON(&request)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner:    request.Owner,
		Currency: request.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   account,
	})

}

//Since the id is a Uri parameter i.e the id is in the path like accounts/:id, we have
//to do a different kind of binding
type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var request GetAccountRequest
	err := ctx.ShouldBindUri(&request)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, request.ID)
	if err != nil {

		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   account,
	})
}