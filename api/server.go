package api

import (
	db "github.com/Agentic_Bank_Server/db/sqlc"
	"github.com/Agentic_Bank_Server/token"
	"github.com/Agentic_Bank_Server/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests fro our banking service.
type Server struct {
	config     utils.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(config utils.Config, store db.Store) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker(config.SymmetricKey)

	if err != nil {
		return nil, err
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
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
	router.POST("/account", server.createAccount)
	router.GET("/account/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)
	router.POST("/transfer", server.createTransfer)
	router.POST("/login", server.loginUser)

	server.router = router
}

// Start runs the HTTP server on the given address
func (server *Server) Start(address string) error {
	err := server.router.Run(address)
	return err
}

func errResponsne(err error) gin.H {
	return gin.H{"error": err.Error()}
}
