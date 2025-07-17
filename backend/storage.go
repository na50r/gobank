package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Storage interface {
	GetAccountByID(id int) (*Account, error)
	CreateAccount(a *Account) error
	UpdateAccount(a *Account) error
	DeleteAccountByNumber(id int) error
	GetAccounts() ([]*Account, error)
	GetAccountByNumber(number int) (*Account, error)
	CreateRefreshToken(rt *RefreshToken) error
	GetRefreshToken(id int) (*RefreshToken, error)
	GetAccountByRefreshToken(token string) (*Account, error)
	UpdateRefreshToken(rt *RefreshToken) error
	AddImage(image []byte, name string) error
	GetImage(accountID string) ([]byte, error)
	GetImages() ([]*Image, error)
	NewImageForAccount(accountNumber int) string
	AddElement(element *Element) error
	GetElement(a, b string) (*string, error)
	Init() error

}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	acc := new(Account)
	err := rows.Scan(
		&acc.ID,
		&acc.FirstName,
		&acc.LastName,
		&acc.ImageName,
		&acc.Number,
		&acc.EncryptedPassword,
		&acc.Balance,
		&acc.CreatedAt,
	)
	return acc, err
}

func scanIntoRefreshToken(rows *sql.Rows) (*RefreshToken, error) {
	rt := new(RefreshToken)
	err := rows.Scan(
		&rt.ID,
		&rt.AccountID,
		&rt.Token,
	)
	return rt, err
}

func scanIntoImage(rows *sql.Rows) (*Image, error) {
	image := new(Image)
	err := rows.Scan(
		&image.Name,
		&image.Image,
	)
	return image, err
}

func scanIntoElement(rows *sql.Rows) (*Element, error) {
	element := new(Element)
	err := rows.Scan(
		&element.A,
		&element.B,
		&element.Result,
	)
	return element, err
}
