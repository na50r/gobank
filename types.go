package main

import (
	"math/rand"
	"time"
	"golang.org/x/crypto/bcrypt"
)

type Transaction struct {
	Recipient int       `json:"recipient"`
	Amount    float64       `json:"amount"`
}

type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Number    int       `json:"number"`
	EncryptedPassword string `json:"-"`
	Balance   float64       `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

type UpdateAccountRequest struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Balance   *int    `json:"balance"`
	Number    *int    `json:"number"`
}

type LoginRequest struct {
	Number int `json:"number"`
	Password string `json:"password"`
}

func NewAccount(firstName, lastName, password string) (*Account, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:        rand.Intn(10000),
		FirstName: firstName,
		LastName:  lastName,
		Number:    rand.Intn(1000000),
		EncryptedPassword: string(encpw),
		Balance:   0,
		CreatedAt: time.Now().UTC(),
	}, nil
}

