package api

import (
	db "github.com/Franklynoble/bankapp/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

//NewServer create a new HTTP server and setup routing

func NewServer(store *db.Store) *Server {

}
