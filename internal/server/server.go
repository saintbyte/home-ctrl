package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/saintbyte/home-ctrl/internal/auth"
	"github.com/saintbyte/home-ctrl/internal/config"
	"github.com/saintbyte/home-ctrl/internal/server/v1"
)

// Server represents the HTTP server
type Server struct {
	config *config.Config
	auth   *auth.Auth
	v1Router *v1.Router
	router *gin.Engine
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config, authService *auth.Auth) *Server {
	return &Server{
		config: cfg,
		auth:   authService,
		v1Router: v1.NewRouter(cfg, authService),
		router: gin.Default(),
	}
}

// SetupRoutes sets up the HTTP routes
func (s *Server) SetupRoutes() {
	// Setup v1 routes
	s.v1Router.SetupRoutes()
	
	// Health check endpoint (not versioned)
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Service is running",
		})
	})
}

// Run starts the HTTP server
func (s *Server) Run() error {
	address := s.config.GetServerAddress()
	fmt.Printf("Starting server on %s\n", address)
	
	return s.v1Router.Run(address)
}

// GetRouter returns the gin router (for testing)
func (s *Server) GetRouter() *gin.Engine {
	return s.v1Router.GetRouter()
}