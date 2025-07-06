package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/na50r/gobank/backend/sse"
	"log"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

var JWT_SECRET = os.Getenv("JWT_SECRET")

func seedAccount(store Storage, fname, lname, pw string) *Account {
	acc, err := NewAccount(fname, lname, pw)
	acc.Balance = float64(1000000)
	if err != nil {
		log.Fatal(err)
	}

	if err := store.CreateAccount(acc); err != nil {
		log.Fatal(err)
	}

	acc, err = store.GetAccountByNumber(acc.Number)
	if err != nil {
		log.Fatal(err)
	}
	rt, err := NewRefreshToken(acc)
	if err != nil {
		log.Fatal(err)
	}
	if err := store.CreateRefreshToken(rt); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Created account:", acc.Number)
	return acc
}

func seedAccounts(s Storage) {
	seedAccount(s, "Andrew", "Levitt", "test1")
	seedAccount(s, "John", "Doe", "test2")
}

func main() {
	seed := flag.Bool("seed", false, "seed the database")
	flag.Parse()

	store, err := NewSQLiteStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	//./bin/gobank --seed
	if *seed {
		fmt.Println("Seeding the database...")
		seedAccounts(store)
	}

	server := NewAPIServer(":3000", store, sse.NewServer())
	server.Run()
}
