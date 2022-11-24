package api

import (
	db "backend_masterclass/db/sqlc"
	"backend_masterclass/token"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransferRequest struct {
	FromAccountId int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountId   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,validcurrency"`
}

func (server *Server) transferMoney(ctx *gin.Context) {

	var request TransferRequest
	err := ctx.ShouldBindJSON(&request)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, fromAccountIsValid := server.validateTransferCurrency(ctx, request.FromAccountId, request.Currency)
	_, toAccountIsValid := server.validateTransferCurrency(ctx, request.ToAccountId, request.Currency)

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Username != fromAccount.Owner {
		err := errors.New("users can only send money from their own accounts")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !fromAccountIsValid || !toAccountIsValid {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: request.FromAccountId,
		ToAccountID:   request.ToAccountId,
		Amount:        request.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   result,
	})

}

func (server *Server) validateTransferCurrency(ctx *gin.Context, accountId int64, currency string) (db.Accounts, bool) {
	account, err := server.store.GetAccount(ctx, accountId)
	if err != nil {

		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountId, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "failure",
			"data":   errorResponse(err),
		})

		return account, false
	}

	return account, true
}
