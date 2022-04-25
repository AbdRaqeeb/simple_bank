package api

import (
	db "github.com/AbdRaqeeb/simple_bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serves http requests for the banking service
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
func NewServer(store db.Store) *Server {
	server := &Server{store: store}

	router := gin.Default()

	// endpoints
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	server.router = router
	return server
}

// Start runs the http server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// errorResponse --> format error into json format
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
