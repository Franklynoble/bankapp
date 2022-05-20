package api

import (
	db "github.com/Franklynoble/bankapp/db/sqlc"
	"github.com/gin-gonic/gin"
)

//Server serves HTTP request for our banking service.
type Server struct {
	store  db.Store
	router *gin.Engine
}

//NewServer create a new HTTP server and setup all  routing

func NewServer(store db.Store) *Server {

	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)

	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)
	router.PUT("/accounts", server.updateAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)

	router.POST("/transfers", server.createTransfer)

	server.router = router // pass the instance of the gin
	return server
}

// take an address input and return an error, Start runs HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
