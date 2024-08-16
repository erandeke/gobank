package storage

import (
	"database/sql"
	"gobank/types"
)

type Store interface {
	CreateNewAccount(*types.Account) error
	DeleteAccount(id interface{}) error
	UpdateAccount(*types.Account) (*types.Account, error)
	GetAccounts() ([]*types.Account, error)
	GetACcountsById(id int) (*types.Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.CreateAccountTable()
}

func (s *PostgresStore) CreateAccountTable() error {

	query := `create table if not exists account (
           id serial primary key,
		   first_name  varchar(100),
		   last_name varchar(100),
		   number serial,
		   encrypted_password varchar(100),
		   balance serial,
		   created_at timestamp

	)`

	_, err := s.db.Exec(query)
	return err

}

//create account

func (s *PostgresStore) CreateAccount(acc *types.Account) error {
	query := `insert into account 
	(first_name, last_name, number, encrypted_password, balance, created_at)
	values ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.Query(
		query,
		acc.FirstName,
		acc.LastName,
		acc.Number,
		acc.EncryptedPassword,
		acc.Balance,
		acc.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

//update account

func (s *PostgresStore) UpdateAccount(acc *types.Account) error {

	return nil

}

//delete account

func (s *PostgresStore) DeleteAccount(id int) error {
	query := `delete from account where id=$1`
	_, err := s.db.Query(query, id)
	return err

}

// getaccount
func (s *PostgresStore) GetAccounts() ([]*types.Account, error) {
	query := `select * from accounts`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	accounts := []*types.Account{}
	for rows.Next() {
		account, err := scanIntoNextRow(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil

}

func scanIntoNextRow(rows *sql.Rows) (*types.Account, error) {
	account := new(types.Account)
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.EncryptedPassword,
		&account.Balance,
		&account.CreatedAt,
	)

	return account, err

}

//getaccountbyId
