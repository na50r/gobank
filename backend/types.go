package main

import (
	jwt "github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"time"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

type Transaction struct {
	Recipient int     `json:"recipient"`
	Amount    float64 `json:"amount"`
}

type Account struct {
	ID                int       `json:"id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	ImageName         string    `json:"image_name"`
	Number            int       `json:"number"`
	EncryptedPassword string    `json:"-"`
	Balance           float64   `json:"balance"`
	CreatedAt         time.Time `json:"created_at"`
}

type RefreshToken struct {
	ID        int    `json:"id"`
	AccountID int    `json:"account_id"`
	Token     string `json:"token"`
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
	Number   int    `json:"number"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type Image struct {
	Name  string `json:"name"`
	Image []byte `json:"image"`
}

type ImageRequest struct {
	Name  string `json:"name"`
	Image []byte `json:"image"`
}

type ImageResponse struct {
	Name  string `json:"name"`
	Image []byte `json:"image"`
}

type Element struct {
	A string `json:"a"`
	B string `json:"b"`
	Result string `json:"result"`
}

type ElementRequest struct {
	A string `json:"a"`
	B string `json:"b"`
}

type ElementResponse struct {
	Result string `json:"result"`
}

func NewAccount(firstName, lastName, password string) (*Account, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	imageName := "default.png"
	return &Account{
		ID:                rand.Intn(10000),
		FirstName:         firstName,
		LastName:          lastName,
		ImageName:         imageName,
		Number:            rand.Intn(1000000),
		EncryptedPassword: string(encpw),
		Balance:           5,
		CreatedAt:         time.Now().UTC(),
	}, nil
}

func NewRefreshToken(account *Account) (*RefreshToken, error) {
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	claims := refreshToken.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	rt, err := refreshToken.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return nil, err
	}

	return &RefreshToken{
		ID:        rand.Intn(10000),
		AccountID: account.ID,
		Token:     rt,
	}, nil
}
