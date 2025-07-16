package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)


type PostgresStore struct {
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