package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore() (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", "./store.db")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Release lock in 5 seconds
	// Reference: https://stackoverflow.com/questions/66909180/increase-the-lock-timeout-with-sqlite-and-what-is-the-default-values
	_, err = db.Exec("PRAGMA busy_timeout = 5000")
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to the SQLite database successfully.")
	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) createElementTable() error {
	query := `create table if not exists element (
		a text,
		b text,
		result text,
		unique(a, b)
		)`
	_, err := s.db.Exec(query)
	return err
}

func (s *SQLiteStore) createAccountTable() error {
	query := `create table if not exists account (
		id integer primary key,
		first_name text,
		last_name text,
		image_name text,
		number integer,
		password text,
		balance real,
		created_at datetime
		)`
	_, err := s.db.Exec(query)
	return err
}

func (s *SQLiteStore) createRefreshTokenTable() error {
	query := `create table if not exists refresh_token (
		id integer primary key autoincrement,
		account_id integer,
		token text
		)`
	_, err := s.db.Exec(query)
	return err
}

func (s *SQLiteStore) createImageTable() error {
	query := `create table if not exists image (
		name text primary key,
		image blob
		)`
	_, err := s.db.Exec(query)
	return err
}

func (s *SQLiteStore) Init() error {
	if err := s.createAccountTable(); err != nil {
		return err
	}
	if err := s.createRefreshTokenTable(); err != nil {
		return err
	}
	if err := s.createImageTable(); err != nil {
		return err
	}
	if err := s.createElementTable(); err != nil {
		return err
	}
	return nil

}

func (s *SQLiteStore) GetRefreshToken(accountID int) (*RefreshToken, error) {
	rows, err := s.db.Query("select * from refresh_token where account_id = ?", accountID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		defer rows.Close()
		return scanIntoRefreshToken(rows)
	}
	return nil, fmt.Errorf("refresh token for account %d not found", accountID)
}

func (s *SQLiteStore) GetAccountByRefreshToken(token string) (*Account, error) {
	rows, err := s.db.Query("select * from refresh_token where token = ?", token)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		defer rows.Close()
		rt, err := scanIntoRefreshToken(rows)
		if err != nil {
			return nil, err
		}
		return s.GetAccountByID(rt.AccountID)
	}
	return nil, fmt.Errorf("refresh token %s not found", token[:10])
}

func (s *SQLiteStore) CreateAccount(acc *Account) error {
	query := `insert into account 
	(id, first_name, last_name, image_name, number, password, balance, created_at)
	values (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.db.Exec(
		query,
		acc.ID,
		acc.FirstName,
		acc.LastName,
		acc.ImageName,
		acc.Number,
		acc.EncryptedPassword,
		acc.Balance,
		acc.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLiteStore) CreateRefreshToken(rt *RefreshToken) error {
	query := `insert into refresh_token 
	(id, account_id, token)
	values (?, ?, ?)`
	_, err := s.db.Exec(
		query,
		rt.ID,
		rt.AccountID,
		rt.Token,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLiteStore) UpdateRefreshToken(rt *RefreshToken) error {
	query := `update refresh_token set
        token = ?
		WHERE account_id = ?`
	_, err := s.db.Exec(
		query,
		rt.Token,
		rt.AccountID,
	)
	return err
}

func (s *SQLiteStore) GetAccountByNumber(number int) (*Account, error) {
	rows, err := s.db.Query("select * from account where number = ?", number)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		defer rows.Close()
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account %d not found", number)
}

func (s *SQLiteStore) GetAccountByID(id int) (*Account, error) {
	rows, err := s.db.Query("select * from account where id = ?", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		defer rows.Close()
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account %d not found", id)
}

func (s *SQLiteStore) DeleteAccountByNumber(id int) error {
	_, err := s.db.Query("delete from account where number = ?", id)
	return err
}

func (s *SQLiteStore) UpdateAccount(acc *Account) error {
	query := `update account set
        first_name = ?,
        last_name = ?,
        balance = ?
        WHERE id = ?`

	_, err := s.db.Exec(
		query,
		acc.FirstName,
		acc.LastName,
		acc.Balance,
		acc.ID,
	)
	return err
}

func (s *SQLiteStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("select * from account")
	if err != nil {
		return nil, err
	}
	accounts := []*Account{}
	for rows.Next() {
		acc, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)
	}
	defer rows.Close()
	return accounts, nil
}

func (s *SQLiteStore) AddImage(image []byte, name string) error {
	_, err := s.db.Exec(
		"insert or replace into image (name, image) values (?, ?)",
		name,
		image,
	)
	return err
}

func (s *SQLiteStore) AddElement(element *Element) error {
	a := element.A
	b := element.B
	sorted := a < b
	if !sorted {
		a, b = b, a
	}
	_, err := s.db.Exec(
		"insert or ignore into element (a, b, result) values (?, ?, ?)",
		a,
		b,
		element.Result,
	)
	return err
}

func (s *SQLiteStore) GetElement(a, b string) (*string, error) {
	a = strings.ToLower(a)
	b = strings.ToLower(b)
	sorted := a < b
	if !sorted {
		a, b = b, a   
	}
	var result string
	err := s.db.QueryRow("SELECT result FROM element WHERE a = ? AND b = ?", a, b).Scan(&result)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("element for %s and %s not found", a, b)
	} else if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *SQLiteStore) GetImage(name string) ([]byte, error) {
	fmt.Println("Fetching image")
	rows, err := s.db.Query("select image from image where name = ?", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var image []byte
		err := rows.Scan(&image)
		if err != nil {
			return nil, err
		}
		return image, nil
	}
	return nil, fmt.Errorf("image for account %s not found", name)
}

func (s *SQLiteStore) GetImages() ([]*Image, error) {
	rows, err := s.db.Query("select * from image")
	if err != nil {
		return nil, err
	}
	images := []*Image{}
	for rows.Next() {
		img, err := scanIntoImage(rows)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, nil
}


func (s *SQLiteStore) NewImageForAccount(accountNumber int) string {
	images, err := s.GetImages()
	if err != nil {
		return err.Error()
	}
	size := len(images)
	hash := accountNumber % size
	image := images[hash]
	return image.Name
}