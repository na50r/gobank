package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/na50r/gobank/backend/sse"
)

var JWT_SECRET string
var CLIENT string

func init() {
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found, continuing...")
    }
    JWT_SECRET = os.Getenv("JWT_SECRET")
    CLIENT = os.Getenv("CLIENT")
}


func main() {
	store, err := NewSQLiteStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	//Accounts for ports provided by hosting services
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "3000"
	}

	server := NewAPIServer(":"+PORT, store, sse.NewServer())
	server.Run()
}
