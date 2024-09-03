package main

import (
	"fmt"
	"log"
	"net/http"
)



const serverPort = "1025"

type AppConfig struct{}



func main() {
	app := AppConfig{}

	log.Printf("Broker Service startup on web port %s\n", serverPort)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", serverPort),
		Handler: app.newRouter(),
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Panicf("Failed to start server: %v", err)
	}
}
