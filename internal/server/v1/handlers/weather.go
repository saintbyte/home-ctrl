package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saintbyte/home-ctrl/internal/weather"
)

type WeatherHandler struct {
	weather *weather.YandexWeather
}

func NewWeatherHandler(w *weather.YandexWeather) *WeatherHandler {
	return &WeatherHandler{weather: w}
}

func (h *WeatherHandler) SetupRoutes(router *gin.RouterGroup) {
	router.GET("/yandex-weather", h.getYandexWeather)
}

func (h *WeatherHandler) getYandexWeather(c *gin.Context) {
	if h.weather == nil || !h.weather.Enabled() {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "Yandex weather is not configured",
		})
		return
	}
	c.JSON(http.StatusOK, h.weather.Get())
}
