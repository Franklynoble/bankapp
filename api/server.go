package api

import (
	"fmt"

	db "github.com/Franklynoble/bankapp/db/sqlc"
	"github.com/Franklynoble/bankapp/db/token"
	"github.com/Franklynoble/bankapp/db/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

//Server serves HTTP request for our banking service.
type Server struct {
	config     util.Config
	tokenMaker token.Maker
	store      db.Store
	router     *gin.Engine
}

//NewServer create a new HTTP server and setup all  routing

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokeMaker, err := token.NewJWTMaker(config.TokenSymmetricKey) // load this from  environment variable

	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %v", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokeMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)

	}
	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.POST("/users", server.createUser)

	router.POST("/user/login", server.loginUser)

	router.POST("/accounts", server.createAccount)

	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)
	router.PUT("/accounts", server.updateAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)

	router.POST("/transfers", server.createTransfer)

	server.router = router // pass the instance of the gin

}

// take an address input and return an error, Start runs HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
