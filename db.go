package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

type MySqlDB struct {
	Db *sql.DB
}

type TransactionStorage interface {
	AddNewTransaction(*Transaction) error
	// DeleteAccount(int) error
	// UpdateAccount(*Account) error
	// GetAccounts() ([]*Account, error)
	// GetAccountByID(int) (*Account, error)
	// GetAccountByNumber(int) (*Account, error)
}

func ConnectToDB() (*MySqlDB, error) {
	// Connect to the database
	db_connection := os.Getenv("DB_CONNECTION")
	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")
	db_database := os.Getenv("DB_DATABASE")
	db_username := os.Getenv("DB_USERNAME")
	db_password := os.Getenv("DB_PASSWORD")

	cfg := mysql.Config{
		User:                 db_username,
		Passwd:               db_password,
		Net:                  "tcp",
		Addr:                 db_host + ":" + db_port,
		DBName:               db_database,
		AllowNativePasswords: true,
	}

	db, err := sql.Open(db_connection, cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
		return nil, err
	}

	log.Println("DB Connected!")

	return &MySqlDB{Db: db}, nil
}

func (s *MySqlDB) Init() error {
	return s.createTransactionTable()
}

func (s *MySqlDB) createTransactionTable() error {
	query := `create table if not exists transactions (
		id serial primary key,
		transaction_type varchar(100) null,
		transaction_date timestamp,
		transaction_amount DECIMAL(20, 2),
		transaction_id varchar(100),
		transaction_name varchar(150) null,
		transaction_memo varchar(150) null,
		created_at timestamp,
		transaction_type_id int null,
		bank_name varchar(100)
	)`

	_, err := s.Db.Exec(query)
	return err
}

func (s *MySqlDB) AddNewTransaction(transaction *Transaction) error {

	// Check if the transaction already exists
	existsQuery := `select count(*) from transactions where transaction_id = ?`
	var count int
	err := s.Db.QueryRow(existsQuery, transaction.TransactionID).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		// Transaction already exists, do not add it
		log.Println("Transaction already exists")
		return nil
	}

	// Add new transaction
	query := `insert into transactions (
		transaction_type,
		transaction_date,
		transaction_amount,
		transaction_id,
		transaction_name,
		transaction_memo,
		created_at,
		transaction_type_id,
		bank_name
	) values (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = s.Db.Exec(
		query,
		transaction.TransactionType,
		transaction.TransactionDate,
		transaction.TransactionAmount.FloatString(2),
		transaction.TransactionID,
		transaction.TransactionName,
		transaction.TransactionMemo,
		transaction.CreatedAt,
		transaction.TransactionTypeID,
		transaction.BankName,
	)

	if err != nil {
		return err
	}

	return nil
}
