package main

import (
	"math/big"
	"time"
)

type Transaction struct {
	ID                int       `json:"id"`
	TransactionType   string    `json:"transactionType"`
	TransactionDate   time.Time `json:"transactionDate"`
	TransactionAmount *big.Rat  `json:"transactionAmount"`
	TransactionID     string    `json:"transactionID"`
	TransactionName   string    `json:"transactionName"`
	TransactionMemo   string    `json:"transactionMemo"`
	CreatedAt         time.Time `json:"createdAt"`
	TransactionTypeID int       `json:"transactionTypeID"`
	BankName          string    `json:"bankName"`
}

func NewTransaction(
	transactionType string,
	transactionDate string,
	transactionAmount *big.Rat,
	transactionID string,
	transactionName string,
	transactionMemo string,
) (*Transaction, error) {

	transactionDateParsed, err := time.Parse("2006-01-02", transactionDate)
	if err != nil {
		return nil, err
	}

	return &Transaction{
		TransactionType:   transactionType,
		TransactionDate:   transactionDateParsed,
		TransactionAmount: transactionAmount,
		TransactionID:     transactionID,
		TransactionName:   transactionName,
		TransactionMemo:   transactionMemo,
		BankName:          "FNB",
		CreatedAt:         time.Now().UTC(),
	}, nil
}
