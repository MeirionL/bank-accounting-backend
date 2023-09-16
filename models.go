package main

import (
	"time"

	"github.com/MeirionL/personal-finance-app/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        int32     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
	}
}

func databaseUsersToUsers(dbUsers []database.User) []User {
	users := []User{}
	for _, dbUser := range dbUsers {
		users = append(users, databaseUserToUser(dbUser))
	}
	return users
}

type Transaction struct {
	ID              uuid.UUID `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	TransactionTime time.Time `json:"transaction_time"`
	Type            string    `json:"type"`
	Amount          float32   `json:"amount"`
	PreBalance      float32   `json:"pre_balance"`
	PostBalance     float32   `json:"post_balance"`
	NewAccount      bool      `json:"new_account"`
	Name            string    `json:"name"`
	AccountNumber   string    `json:"account_number"`
	SortCode        string    `json:"sort_code"`
	UserID          int32     `json:"user_id"`
}

func databaseTransactionToTransaction(dbTransaction database.Transaction) Transaction {
	return Transaction{
		ID:              dbTransaction.ID,
		CreatedAt:       dbTransaction.CreatedAt,
		UpdatedAt:       dbTransaction.UpdatedAt,
		TransactionTime: dbTransaction.TransactionTime,
		Type:            dbTransaction.Type,
		Amount:          dbTransaction.Amount,
		PreBalance:      dbTransaction.PreBalance,
		PostBalance:     dbTransaction.PostBalance,
		NewAccount:      dbTransaction.NewAccount,
		Name:            dbTransaction.Name,
		AccountNumber:   dbTransaction.AccountNumber,
		SortCode:        dbTransaction.SortCode,
		UserID:          dbTransaction.UserID,
	}
}

func databaseTransactionsToTransactions(dbTransactions []database.Transaction) []Transaction {
	transactions := []Transaction{}
	for _, dbTransaction := range dbTransactions {
		transactions = append(transactions, databaseTransactionToTransaction(dbTransaction))
	}
	return transactions
}
