package main

import (
	"fmt"
	"log"
	"net/http"
)



const serverPort = "1025"

type AppConfig struct{}



func main() {
	log.Printf("Broker service startup on web port %s\n", serverPort)
	
	app := AppConfig{}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", serverPort),
		Handler: app.newRouter(),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
