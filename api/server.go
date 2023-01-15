package api

import (
	db "bank/db/sqlc"

	"github.com/gin-gonic/gin"
)

// Server serve Http request for our services
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// Create new Http server and setup routing
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)

	server.router = router
	return server
}

// Run http server on the specific address
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
