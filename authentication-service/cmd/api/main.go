package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)



const serverPort = "1025"

type AppConfig struct {
	DB        *sql.DB
	DBModels  data.Models
}



func main() {
	log.Printf("Authentication sservice startup on web port %s\n", serverPort)

	connection := connectToDB()
	if connection == nil {
		log.Panic("Unable to establish a connection to Postgres!")
	}

	app := AppConfig{
		DB:        connection,
		DBModels:  data.NewDatabase(connection),
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", serverPort),
		Handler: app.newRouter(),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Check the connection by pinging the database
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	retryLimit := 10
	retryDelay := 2 * time.Second

	for attempts := 0; attempts <= retryLimit; attempts++ {
		connection, err := openDB(dsn)
		if err != nil {
			log.Printf("Postgres not yet ready... Attempt #%d\n", attempts+1)
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		if attempts == retryLimit {
			log.Printf("Reached retry limit. Last error: %v\n", err)
			return nil
		}

		log.Println("Backing off for two seconds...")
		time.Sleep(retryDelay)
	}

	return nil
}
