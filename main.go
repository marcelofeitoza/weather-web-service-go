package main

import (
	"github.com/marcelofeitoza/weather-web-service-go/server"
	"log"
)

func main() {
	srv, err := server.NewServer()
	if err != nil {
		log.Fatalf("Error starting server: %s", err.Error())
		return
	}

	srv.Gin.Run()
}
