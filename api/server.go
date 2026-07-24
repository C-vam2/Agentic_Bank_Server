package api

import (
	db "github.com/Agentic_Bank_Server/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests fro our banking service.
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(store *db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()

	router.POST("/account", server.createAccount)
	router.GET("account/:id", server.getAccount)

	server.router = router
	return server
}

// Start runs the HTTP server on the given address
func (server *Server) Start(address string) error {
	err := server.router.Run(address)
	return err
}

func errResponsne(err error) gin.H {
	return gin.H{"error": err.Error()}
}
