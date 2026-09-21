package server

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/saintbyte/home-ctrl/internal/auth"
	"github.com/saintbyte/home-ctrl/internal/config"
	"github.com/saintbyte/home-ctrl/internal/database"
	"github.com/saintbyte/home-ctrl/internal/scheduler"
	"github.com/saintbyte/home-ctrl/internal/server/v1"
	"github.com/saintbyte/home-ctrl/internal/weather"
)

// Server represents the HTTP server
type Server struct {
	config   *config.Config
	auth     *auth.Auth
	v1Router *v1.Router
	router   *gin.Engine
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config, authService *auth.Auth, db *database.Database, sched *scheduler.Scheduler, weatherService *weather.YandexWeather) *Server {
	return &Server{
		config:   cfg,
		auth:     authService,
		v1Router: v1.NewRouter(cfg, authService, db, sched, weatherService),
		router:   gin.Default(),
	}
}

// SetupRoutes sets up the HTTP routes
func (s *Server) SetupRoutes() {
	// Configure CORS to allow all origins
	s.router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * 3600, // 12 hours
	}))

	// HSTS header is only sent on HTTPS responses (RFC 6797). The main router
	// serves HTTPS when TLS is enabled, so enable the middleware then.
	if s.config.Server.TLS {
		s.router.Use(func(c *gin.Context) {
			c.Header("Strict-Transport-Security", "max-age=31536000")
			c.Next()
		})
	}

	// Setup v1 routes - use the main router instead of v1Router
	s.v1Router.SetupRoutesOn(s.router)

	// Health check endpoint (not versioned)
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Service is running",
		})
	})

	// Serve static files from public directory in root
	// Files from public will be accessible directly (e.g., /style.css, /script.js)
	// Use StaticFile for root to serve index.html
	s.router.StaticFile("/", "./public/index.html")

	// Serve other static files from public directory
	// This middleware checks if a file exists and serves it, otherwise passes to next handler
	s.router.Use(func(c *gin.Context) {
		path := c.Request.URL.Path
		// Skip root path, health check, and API routes
		if path == "/" || path == "/health" ||
			path == "/api" ||
			(len(path) > 5 && path[:5] == "/api/") {
			c.Next()
			return
		}

		// Check if file exists in public directory
		filePath := filepath.Join("./public", path)
		if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
			c.File(filePath)
			c.Abort()
			return
		}

		c.Next()
	})

	// Fallback: serve index.html for all other routes (for SPA support)
	s.router.NoRoute(func(c *gin.Context) {
		// Skip health check and API routes
		if c.Request.URL.Path == "/health" ||
			c.Request.URL.Path == "/api" ||
			(len(c.Request.URL.Path) > 5 && c.Request.URL.Path[:5] == "/api/") {
			c.Next()
			return
		}
		c.File("./public/index.html")
	})
}

// Run starts the HTTP(S) server
func (s *Server) Run() error {
	address := s.config.GetServerAddress()

	if !s.config.Server.TLS {
		fmt.Printf("Starting server on http://%s\n", address)
		return s.router.Run(address)
	}

	if s.config.Server.TLSCert == "" || s.config.Server.TLSKey == "" {
		return fmt.Errorf("tls enabled but tls_cert or tls_key is not set")
	}

	errCh := make(chan error, 2)

	// When http_port is set, start a plain HTTP listener that serves only the
	// bootstrap endpoints (main page and root certificate download) and
	// redirects everything else to HTTPS. This is what makes HSTS bootstrap
	// work: the browser first trusts the CA and then sees an HTTPS response
	// with the Strict-Transport-Security header.
	if s.config.Server.HTTPPort > 0 && s.config.Server.HTTPPort != s.config.Server.Port {
		httpAddr := net.JoinHostPort(s.config.Server.Host, strconv.Itoa(s.config.Server.HTTPPort))
		httpServer := &http.Server{
			Addr:    httpAddr,
			Handler: s.buildBootstrapRouter(),
		}
		fmt.Printf("Starting HTTP bootstrap server on http://%s\n", httpAddr)
		go func() {
			if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				errCh <- fmt.Errorf("http bootstrap server: %w", err)
			}
		}()
	}

	fmt.Printf("Starting server on https://%s\n", address)
	go func() {
		if err := s.router.RunTLS(address, s.config.Server.TLSCert, s.config.Server.TLSKey); err != nil {
			errCh <- err
		}
	}()

	return <-errCh
}

// buildBootstrapRouter builds the plain HTTP router that is only used for
// HSTS bootstrap: it serves the main page and the root certificate download
// and redirects every other path to HTTPS.
func (s *Server) buildBootstrapRouter() *gin.Engine {
	router := gin.Default()

	// Main page over HTTP (bootstrap)
	router.GET("/", func(c *gin.Context) {
		c.File("./public/index.html")
	})

	// Root certificate download over HTTP (bootstrap)
	// Reuse the v1 handler so the endpoint behaves identically to HTTPS.
	certHandler := v1.NewCertificateHandler(s.config)
	certHandler.SetupRoutes(router.Group("/api/v1"))

	// Everything else goes to HTTPS
	router.NoRoute(func(c *gin.Context) {
		s.redirectToHTTPS(c)
	})

	return router
}

// redirectToHTTPS issues a permanent redirect to the same path on the HTTPS
// port, preserving the hostname the client used.
func (s *Server) redirectToHTTPS(c *gin.Context) {
	host := c.Request.Host
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	} else if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		host = strings.Trim(host, "[]")
	}

	tlsPort := s.config.Server.Port
	if tlsPort != 443 {
		host = net.JoinHostPort(host, strconv.Itoa(tlsPort))
	}

	c.Redirect(http.StatusMovedPermanently, "https://"+host+c.Request.URL.RequestURI())
}

// GetRouter returns the gin router (for testing)
func (s *Server) GetRouter() *gin.Engine {
	return s.v1Router.GetRouter()
}
