package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/saintbyte/home-ctrl/internal/config"
	"github.com/saintbyte/home-ctrl/internal/database"
)

// Auth represents the authentication service
type Auth struct {
	config    *config.Config
	database  *database.Database
	users     map[string]string // username: password
	sessionTTL time.Duration
}

// NewAuth creates a new authentication service
func NewAuth(cfg *config.Config, db *database.Database) *Auth {
	return &Auth{
		config:    cfg,
		database:  db,
		users:     make(map[string]string),
		sessionTTL: 24 * time.Hour, // 24 hours by default
	}
}

// AddUser adds a user to the in-memory user store
func (a *Auth) AddUser(username, password string) {
	a.users[username] = password
}

// Authenticate authenticates a user and returns a session token
func (a *Auth) Authenticate(username, password string) (string, error) {
	// Check if user exists and password matches
	storedPassword, exists := a.users[username]
	if !exists || storedPassword != password {
		return "", fmt.Errorf("invalid username or password")
	}

	// Generate session token
	sessionID, err := generateRandomString(32)
	if err != nil {
		return "", fmt.Errorf("failed to generate session token: %w", err)
	}

	// Create session in database
	expiresAt := time.Now().Add(a.sessionTTL)
	_, err = a.database.CreateSession(sessionID, username, expiresAt)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	return sessionID, nil
}

// ValidateSession validates a session token
func (a *Auth) ValidateSession(sessionID string) (string, bool) {
	if sessionID == "" {
		return "", false
	}

	// Check if session is valid
	if !a.database.ValidateSession(sessionID) {
		return "", false
	}

	// Get session to return username
	session, err := a.database.GetSessionByID(sessionID)
	if err != nil || session == nil {
		return "", false
	}

	return session.Username, true
}

// ValidateAPIKey validates an API key
func (a *Auth) ValidateAPIKey(apiKey string) bool {
	if apiKey == "" {
		return false
	}

	return a.database.ValidateAPIKey(apiKey)
}

// generateRandomString generates a random string of given length
func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// AuthMiddleware is a Gin middleware for authentication
func (a *Auth) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for API key in header
		apiKey := c.GetHeader("X-API-Key")
		if apiKey != "" && a.ValidateAPIKey(apiKey) {
			c.Next()
			return
		}

		// Check for Authorization header (Bearer token)
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			// Extract token from "Bearer <token>" format
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				sessionID := authHeader[7:]
				if username, valid := a.ValidateSession(sessionID); valid {
					c.Set("username", username)
					c.Next()
					return
				}
			}
		}

		// If no valid authentication, return 401
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
			"message": "Authentication required",
		})
	}
}

// LoginHandler handles user login
func (a *Auth) LoginHandler() gin.HandlerFunc {
	type LoginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})
			return
		}

		// Authenticate user
		sessionID, err := a.Authenticate(req.Username, req.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid username or password",
			})
			return
		}

		// Return bearer token
		c.JSON(http.StatusOK, gin.H{
			"token":        sessionID,
			"token_type":   "bearer",
			"expires_in":   int(a.sessionTTL.Seconds()),
			"username":     req.Username,
			"message":      "Login successful",
		})
	}
}

// LogoutHandler handles user logout
func (a *Auth) LogoutHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			sessionID := authHeader[7:]
			
			// Delete session
			if err := a.database.DeleteSession(sessionID); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Internal Server Error",
					"message": "Failed to logout",
				})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Logout successful",
		})
	}
}