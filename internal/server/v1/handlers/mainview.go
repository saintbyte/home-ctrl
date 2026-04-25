package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saintbyte/home-ctrl/internal/config"
)

type MainViewHandler struct {
	config *config.Config
}

func NewMainViewHandler(cfg *config.Config) *MainViewHandler {
	return &MainViewHandler{config: cfg}
}

func (h *MainViewHandler) SetupRoutes(router *gin.RouterGroup) {
	router.GET("/mainview", h.getMainView)
}

func (h *MainViewHandler) getMainView(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"widgets": h.config.MainView.Widgets,
	})
}
