package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saintbyte/home-ctrl/internal/auth"
	"github.com/saintbyte/home-ctrl/internal/config"
)

// Server represents the HTTP server
type Server struct {
	config *config.Config
	auth   *auth.Auth
	router *gin.Engine
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config, authService *auth.Auth) *Server {
	return &Server{
		config: cfg,
		auth:   authService,
		router: gin.Default(),
	}
}

// SetupRoutes sets up the HTTP routes
func (s *Server) SetupRoutes() {
	// Public routes (no authentication required)
	s.setupPublicRoutes()
	
	// Auth routes
	s.setupAuthRoutes()
	
	// Protected routes (require authentication)
	s.setupProtectedRoutes()

	// 404 handler
	s.router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Not found",
			"path":  c.Request.URL.Path,
		})
	})
}

// setupPublicRoutes sets up routes that don't require authentication
func (s *Server) setupPublicRoutes() {
	// Health check endpoint
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Service is running",
		})
	})

	// API version endpoint (versioned)
	s.router.GET("/api/v1/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version": "0.1.0",
			"name":    "home-ctrl",
		})
	})
}

// setupAuthRoutes sets up authentication-related routes
func (s *Server) setupAuthRoutes() {
	authGroup := s.router.Group("/api/v1/auth")
	{
		// Login endpoint
		authGroup.POST("/login", s.auth.LoginHandler())
		
		// Logout endpoint (requires authentication)
		authGroup.POST("/logout", s.auth.AuthMiddleware(), s.auth.LogoutHandler())
	}
}

// setupProtectedRoutes sets up routes that require authentication
func (s *Server) setupProtectedRoutes() {
	protectedGroup := s.router.Group("/api/v1")
	protectedGroup.Use(s.auth.AuthMiddleware())
	{
		// Example protected endpoint
		protectedGroup.GET("/example", func(c *gin.Context) {
			username, _ := c.Get("username")
			c.JSON(http.StatusOK, gin.H{
				"message": "Hello from protected endpoint!",
				"user":    username,
				"config": gin.H{
					"host": s.config.Server.Host,
					"port": s.config.Server.Port,
				},
			})
		})
		
		// User info endpoint
		protectedGroup.GET("/me", func(c *gin.Context) {
			username, _ := c.Get("username")
			c.JSON(http.StatusOK, gin.H{
				"username": username,
				"message": "You are authenticated!",
			})
		})
	}
}

// Run starts the HTTP server
func (s *Server) Run() error {
	address := s.config.GetServerAddress()
	fmt.Printf("Starting server on %s\n", address)
	
	return s.router.Run(address)
}

// GetRouter returns the gin router (for testing)
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}