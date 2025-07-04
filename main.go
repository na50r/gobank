package main

import (
	"flag"
	"fmt"
	"log"
)

func seedAccount(store Storage, fname, lname, pw string) *Account {
	acc, err := NewAccount(fname, lname, pw)
	if err != nil {
		log.Fatal(err)
	}

	if err := store.CreateAccount(acc); err != nil {
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

	server := NewAPIServer(":3000", store)
	server.Run()
}
