package models

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	Gin *gin.Engine
	DB  *sqlx.DB
	RDB *redis.Client
}

type Handler struct {
	DB  *sqlx.DB
	RDB *redis.Client
}

type Forecast struct {
	Date        string  `json:"date"`
	Temperature float64 `json:"temperature"`
}

type WeatherDisplay struct {
	Duration  string     `json:"duration"`
	City      string     `json:"city"`
	Forecasts []Forecast `json:"forecasts"`
}

type WeatherResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`
	Hourly    struct {
		Time          []string  `json:"time"`
		Temperature2m []float64 `json:"temperature_2m"`
	} `json:"hourly"`
}

type GeoResponse struct {
	Results []LatLong `json:"results"`
}

type LatLong struct {
	Latitude  float64 `json:"Latitude"`
	Longitude float64 `json:"Longitude"`
}
