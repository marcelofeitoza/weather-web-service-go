package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/marcelofeitoza/weather-web-service-go/models"
	"github.com/redis/go-redis/v9"
	"net/http"
	"os"
)

func NewServer() (*models.Server, error) {
	r := gin.New()
	//_ = r.SetTrustedProxies([]string{"127.0.0.1", "0.0.0.0", "192.168.0.1", "192.168.0.2"})

	dbConnStr := os.Getenv("DATABASE_URL")
	fmt.Println(dbConnStr)
	db, err := sqlx.Open("postgres", dbConnStr)
	if err != nil {
		return nil, err
	}

	rdbConnStr := os.Getenv("REDIS_URL")
	rdb := redis.NewClient(&redis.Options{
		Addr:     rdbConnStr,
		Password: "",
		DB:       0,
	})

	setupRoutes(r, db, rdb)

	return &models.Server{Gin: r, DB: db, RDB: rdb}, nil
}

func setupRoutes(r *gin.Engine, db *sqlx.DB, rdb *redis.Client) {
	handler := &Handler{Handler: &models.Handler{DB: db, RDB: rdb}}

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/health")
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, weather!"})
	})

	r.GET("/weather/:city", handler.Weather)
}
