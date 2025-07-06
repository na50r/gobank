package main

import (
	"database/sql"
	"fmt"

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
	Init() error
}

type PostgresStore struct {
	db *sql.DB
}

type SQLiteStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	//To avoid conflicts with dockerized postgres, make sure to create it with:
	//docker run --name postgres -e POSTGRES_PASSWORD=gobank -p 5433:5432 -d postgres
	//Map port 5432 of the container to 5433 of the host
	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable port=5433"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("Connected to the Postgres database successfully.")
	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Init() error {
	s.createAccountTable()
	return nil
}

func (s *PostgresStore) createAccountTable() error {
	query := `create table if not exists account (
		id integer primary key,
		first_name varchar(100),
		last_name varchar(100),
		number serial,
		password varchar(100),
		balance serial,
		created_at timestamp
		)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `insert into account 
	(first_name, last_name, number, password, balance, created_at)
	values ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.Exec(
		query,
		acc.FirstName,
		acc.LastName,
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

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	rows, err := s.db.Query("select * from account where id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account %d not found", id)
}

func (s *PostgresStore) UpdateAccount(acc *Account) error {
	query := `UPDATE account SET
        first_name = $1,
        last_name = $2,
        balance = $3
		WHERE id = $4`
	_, err := s.db.Exec(
		query,
		acc.FirstName,
		acc.LastName,
		acc.Balance,
		acc.ID,
	)
	return err
}

func (s *PostgresStore) DeleteAccountByNumber(id int) error {
	_, err := s.db.Query("delete from account where number = $1", id)
	return err
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
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
	return accounts, nil
}

func NewSQLiteStore() (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", "./bank.db")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	_, err = db.Exec("PRAGMA busy_timeout = 5000")
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to the SQLite database successfully.")
	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) createAccountTable() error {
	query := `create table if not exists account (
		id integer primary key,
		first_name text,
		last_name text,
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

func (s *SQLiteStore) Init() error {
	if err := s.createAccountTable(); err != nil {
		return err
	}
	if err := s.createRefreshTokenTable(); err != nil {
		return err
	}
	return nil
	
}

func (s * SQLiteStore) GetRefreshToken(accountID int) (*RefreshToken, error) {
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
	(id, first_name, last_name, number, password, balance, created_at)
	values (?, ?, ?, ?, ?, ?, ?)`

	_, err := s.db.Exec(
		query,
		acc.ID,
		acc.FirstName,
		acc.LastName,
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

func (s * SQLiteStore) UpdateRefreshToken(rt *RefreshToken) error {
	query := `UPDATE refresh_token SET
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
	query := `UPDATE account SET
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

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	fmt.Println("Fetching account")
	acc := new(Account)
	err := rows.Scan(
		&acc.ID,
		&acc.FirstName,
		&acc.LastName,
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
