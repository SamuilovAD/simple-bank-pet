package api

import (
	db "github.com/SamuilovAD/simple-bank-pet/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for my banking service
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewSever creates a new HTTP sever and setup routing
func NewSever(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
