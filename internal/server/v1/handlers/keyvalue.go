package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saintbyte/home-ctrl/internal/database"
	"github.com/saintbyte/home-ctrl/internal/database/models"
)

// KeyValueHandler handles key-value storage operations
type KeyValueHandler struct {
	db *database.Database
}

// NewKeyValueHandler creates a new key-value handler
func NewKeyValueHandler(db *database.Database) *KeyValueHandler {
	return &KeyValueHandler{db: db}
}

// SetupRoutes sets up key-value related routes
func (h *KeyValueHandler) SetupRoutes(router *gin.RouterGroup) {
	keyValueGroup := router.Group("/keyvalue")
	{
		keyValueGroup.POST("", h.createKeyValue)
		keyValueGroup.GET("/:key", h.getKeyValue)
		keyValueGroup.PUT("/:key", h.updateKeyValue)
		keyValueGroup.PATCH("/:key/status", h.updateKeyValueStatus)
		keyValueGroup.PATCH("/:key/hidden", h.updateKeyValueHidden)
		keyValueGroup.DELETE("/:key", h.deleteKeyValue)
		keyValueGroup.GET("", h.listKeyValues)
		keyValueGroup.GET("/:key/status", h.checkKeyValueStatus)
		keyValueGroup.GET("/:key/exists", h.checkKeyValueExists)
	}
}

// createKeyValue handles POST /keyvalue
func (h *KeyValueHandler) createKeyValue(c *gin.Context) {
	type request struct {
		Key   string `json:"key" binding:"required"`
		Value string `json:"value" binding:"required"`
	}

	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	kv, err := h.db.CreateKeyValue(req.Key, req.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to create key-value pair",
		})
		return
	}

	c.JSON(http.StatusCreated, kv)
}

// getKeyValue handles GET /keyvalue/:key
func (h *KeyValueHandler) getKeyValue(c *gin.Context) {
	key := c.Param("key")

	kv, err := h.db.GetKeyValue(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to get key-value pair",
		})
		return
	}

	if kv == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "Key not found",
		})
		return
	}

	c.JSON(http.StatusOK, kv)
}

// updateKeyValue handles PUT /keyvalue/:key
func (h *KeyValueHandler) updateKeyValue(c *gin.Context) {
	key := c.Param("key")

	type request struct {
		Value string `json:"value" binding:"required"`
	}

	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	kv, err := h.db.UpdateKeyValue(key, req.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to update key-value pair",
		})
		return
	}

	c.JSON(http.StatusOK, kv)
}

// updateKeyValueStatus handles PATCH /keyvalue/:key/status
func (h *KeyValueHandler) updateKeyValueStatus(c *gin.Context) {
	key := c.Param("key")

	type request struct {
		Status string `json:"status" binding:"required,oneof=unread read archived"`
	}

	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	kv, err := h.db.UpdateKeyValueStatus(key, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to update key-value status",
		})
		return
	}

	c.JSON(http.StatusOK, kv)
}

// updateKeyValueHidden handles PATCH /keyvalue/:key/hidden
func (h *KeyValueHandler) updateKeyValueHidden(c *gin.Context) {
	key := c.Param("key")

	type request struct {
		Hidden bool `json:"hidden"`
	}

	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	kv, err := h.db.UpdateKeyValueHidden(key, req.Hidden)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to update key-value hidden flag",
		})
		return
	}

	c.JSON(http.StatusOK, kv)
}

// deleteKeyValue handles DELETE /keyvalue/:key
func (h *KeyValueHandler) deleteKeyValue(c *gin.Context) {
	key := c.Param("key")

	if err := h.db.DeleteKeyValue(key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to delete key-value pair",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Key-value pair deleted successfully",
	})
}

// listKeyValues handles GET /keyvalue
func (h *KeyValueHandler) listKeyValues(c *gin.Context) {
	includeHidden := c.DefaultQuery("include_hidden", "false") == "true"

	keyValues, err := h.db.ListKeyValues(includeHidden)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to list key-value pairs",
		})
		return
	}

	c.JSON(http.StatusOK, keyValues)
}

// checkKeyValueStatus handles GET /keyvalue/:key/status
func (h *KeyValueHandler) checkKeyValueStatus(c *gin.Context) {
	key := c.Param("key")

	status, exists, err := h.db.CheckKeyValueStatus(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to check key-value status",
		})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "Key not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"key":    key,
		"status": status,
		"exists": exists,
	})
}

// checkKeyValueExists handles GET /keyvalue/:key/exists
func (h *KeyValueHandler) checkKeyValueExists(c *gin.Context) {
	key := c.Param("key")

	exists, err := h.db.CheckKeyValueExists(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to check key-value existence",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"key":    key,
		"exists": exists,
	})
}