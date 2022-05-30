package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/Franklynoble/bankapp/db/sqlc"
	"github.com/Franklynoble/bankapp/db/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(ctx *gin.Context) {

	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	// create account payload
	authPayload := ctx.MustGet(authorizationPayloadkey).(*token.Payload) // this returns an interface so convert  it to tokenpayload type
	// get args for new account, for first time user
	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	//create the new account, if err return else create the new account
	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok { //finding type of db errors using
			fmt.Printf("pqerr.Code.Name(): %v\n", pqerr.Code.Name())
			switch pqerr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	// if no errors occour, the the account is successfully created
	ctx.JSON(http.StatusOK, account)
}

//using gin to bind id in the uri
type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// get single account
func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return

	}
	// create account payload
	authPayload := ctx.MustGet(authorizationPayloadkey).(*token.Payload) // this returns an interface so convert  it to tokenpayload type

	if account.Owner != authPayload.Username {
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}
	//account = db.Account{} for test
	ctx.JSON(http.StatusOK, account)

}

type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	// create account payload
	authPayload := ctx.MustGet(authorizationPayloadkey).(*token.Payload) // this returns an interface so convert  it to tokenpayload type
	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,                    // this would be the page size
		Offset: (req.PageID - 1) * req.PageSize, //number of records  the database should  skip
	}

	accounts, err := server.store.ListAccounts(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, accounts)

}

type updateAccountRequest struct {
	AccountID int64 `json:"id" binding:"required,min=1" `
	Balance   int64 `json:"balance" binding:"required,min=0"`
}

//update account balance
func (server *Server) updateAccount(ctx *gin.Context) {
	var req updateAccountRequest

	//use Bind JSON when using gin json binding
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fmt.Print("error printed")
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	arg := db.UpdateAccountParams{
		ID:      req.AccountID,
		Balance: req.Balance,
	}
	update, err := server.store.UpdateAccount(ctx, arg)

	if err != nil {
		fmt.Print("second  error")
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, update)

}

type deletAccountRequest struct {
	ID int64 `uri:"id" binding:"required"` //note when retriving you Must use caps for these fields
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var req deletAccountRequest

	//use shouldBindUri if request is bind
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	fmt.Print("retrived id: ", req.ID)
	//arg := req.ID

	err := server.store.DeleteAccount(ctx, req.ID)
	fmt.Print("retrived id: Again ", req.ID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	fmt.Print("deleted with id", req.ID)
	ctx.JSON(http.StatusOK, "Object Deleted")

}

func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
