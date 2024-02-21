package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/marcelofeitoza/weather-web-service-go/models"
	"github.com/marcelofeitoza/weather-web-service-go/server"
	"github.com/redis/go-redis/v9"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	BaseURL         = "http://0.0.0.0:8080"
	DefaultCity     = "Montana"
	NonExistentCity = "NonExistentCity"
	RedirectStatus  = 308
	SuccessStatus   = 200
	NotFoundStatus  = 404
	RequestTimes    = 50
	MeanLimit       = 50
)

var dbConnStr = os.Getenv("DATABASE_URL")
var db, _ = sqlx.Open("postgres", dbConnStr)

var rdbConnStr = os.Getenv("REDIS_URL")
var rdb = redis.NewClient(&redis.Options{
	Addr:     rdbConnStr,
	Password: "",
	DB:       0,
})

func TestPerformance(t *testing.T) {
	var totalRequestTime time.Duration

	gin.SetMode(gin.TestMode)

	router := gin.Default()

	for i := 0; i < RequestTimes; i++ {
		start := time.Now()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, BaseURL+"/weather/"+DefaultCity, nil)
		router.ServeHTTP(w, req)

		totalRequestTime += time.Since(start)
	}

	meanRequestTime := totalRequestTime.Milliseconds() / RequestTimes
	if meanRequestTime > MeanLimit {
		t.Fatalf("Mean request time is too high: %d", meanRequestTime)
		return
	}

	t.Logf("Mean request time: %d", meanRequestTime)

	//router.GET("/weather/:city", handler.Weather)
	//
	//req, _ := http.NewRequest(http.MethodGet, "/weather/São Paulo", nil)
	//resp := httptest.NewRecorder()
	//
	//router.ServeHTTP(resp, req)
	//
	//if resp.Code != http.StatusOK {
	//	t.Fatalf("Expected status %d, got %d", http.StatusOK, resp.Code)
	//}
}

func TestWeatherHandlerReturnsNotFoundForInvalidCity(t *testing.T) {
	handler := &server.Handler{Handler: &models.Handler{DB: db, RDB: rdb}}

	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.GET("/weather/:city", handler.Weather)

	req, _ := http.NewRequest(http.MethodGet, "/weather/InvalidCity", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusNotFound {
		t.Fatalf("Expected status %d, got %d", http.StatusNotFound, resp.Code)
	}
}

func TestWeatherHandlerReturnsBadRequestWhenNoCityProvided(t *testing.T) {
	handler := &server.Handler{Handler: &models.Handler{DB: db, RDB: rdb}}

	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.GET("/weather/:city", handler.Weather)

	req, _ := http.NewRequest(http.MethodGet, "/weather/", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("Expected status %d, got %d", http.StatusBadRequest, resp.Code)
	}
}
