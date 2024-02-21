package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/marcelofeitoza/weather-web-service-go/models"
	"github.com/marcelofeitoza/weather-web-service-go/services"
	"net/http"
	"time"
)

type Handler struct {
	*models.Handler
}

func (h *Handler) Weather(c *gin.Context) {
	city := c.Param("city")

	start := time.Now()

	latlong, err := services.GetLatLong(city, h.DB, h.RDB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	weather, err := services.GetWeather(*latlong, h.RDB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	forecasts := make([]models.Forecast, len(weather.Hourly.Time))
	for i := range weather.Hourly.Time {
		forecasts[i] = models.Forecast{
			Date:        weather.Hourly.Time[i],
			Temperature: weather.Hourly.Temperature2m[i],
		}
	}

	duration := time.Since(start)

	display := models.WeatherDisplay{
		Duration:  duration.String(),
		City:      city,
		Forecasts: forecasts,
	}

	c.JSON(http.StatusOK, display)

	fmt.Printf("Weather request to %s took: %s\n", city, duration)
}
