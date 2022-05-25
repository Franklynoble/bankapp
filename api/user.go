package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	db "github.com/Franklynoble/bankapp/db/sqlc"
	"github.com/Franklynoble/bankapp/db/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	Fullname string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

//convert the db.user object to user response
func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	//hash the password before storing
	hashedPassword, err := util.HashedPassword(req.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	// get args for new account, for first time user
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.Fullname,
		Email:          req.Email,
	}
	//arg = db.CreateUserParams{}
	//create the new account, if err return else create the new account
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok { //finding type of db errors using
			fmt.Printf("pqerr.Code.Name(): %v\n", pqerr.Code.Name())
			switch pqerr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	//call the newUserResponse func to create the response object
	rsp := newUserResponse(user)
	// if no errors occour, the the account is successfully created
	ctx.JSON(http.StatusOK, rsp)
}

//to login the user
type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}
type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

//login user function
func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	user, err := server.store.GetUser(ctx, req.Username) // ge the user from the database
	if err != nil {
		//can't get user or user does not exist
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		// or display the internal server error
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	//if password, check the password if it is valid
	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}
	// if password valid no create accessToken for a specific user
	accessToken, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuraion,
	)
	//if error occours, write to client
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
	}

	//if all good, pass in the created accessToken
	resp := loginUserResponse{
		AccessToken: accessToken, //pass the created accessToken
		User:        newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, resp)

}
