package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/Franklynoble/bankapp/db/sqlc"
	"github.com/Franklynoble/bankapp/db/token"
	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"` //change this latter base on the currency
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {

	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	//store the valid result of from the account and  valid variable
	fromAccount, valid := server.validAccount(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return
	}
	// create account payload
	authPayload := ctx.MustGet(authorizationPayloadkey).(*token.Payload) // this returns an interface so convert  it to tokenpayload type
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("from account deos not belong to Authenticated user")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}
	_, valid = server.validAccount(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return
	}
	// use the validation to check for the accounts: from Account and To Account
	//if !server.validAccount(ctx, req.FromAccountID, req.Currency) {
	//return
	//}
	//if !server.validAccount(ctx, req.ToAccountID, req.Currency) {
	//	return
	//}

	// get args for new account, for first time user
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	//create the new account, if err return else create the new account
	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	// if no errors occour, the the account is successfully created
	ctx.JSON(http.StatusOK, result)
}

//to check the validation for currency in the two account
func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID) //get the account from the databse
	if err != nil {
		// this error is for when the account does not exist
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return account, false
	}
	if account.Currency != currency {
		err := fmt.Errorf("account[%d] currency mismatch: %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return account, false
	}
	return account, true

}
