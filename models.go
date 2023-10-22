// The following functions are required to apply JSON tags to the parameters of returned types from sqlc generated code, so that they
// may be used in JSON responses.

package main

import (
	"time"

	"github.com/MeirionL/boing-block/internal/database"
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

type Account struct {
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	AccountName   string    `json:"account_name"`
	Balance       float32   `json:"balance"`
	AccountNumber string    `json:"account_number"`
	SortCode      string    `json:"sort_code"`
	UserID        int32     `json:"user_id"`
}

func databaseAccountToAccount(dbAccount database.Account) Account {
	return Account{
		ID:            dbAccount.ID,
		CreatedAt:     dbAccount.CreatedAt,
		UpdatedAt:     dbAccount.UpdatedAt,
		AccountName:   dbAccount.AccountName,
		Balance:       dbAccount.Balance,
		AccountNumber: dbAccount.AccountNumber,
		SortCode:      dbAccount.SortCode,
		UserID:        dbAccount.UserID,
	}
}

func databaseAccountsToAccounts(dbAccounts []database.Account) []Account {
	accounts := []Account{}
	for _, dbAccount := range dbAccounts {
		accounts = append(accounts, databaseAccountToAccount(dbAccount))
	}
	return accounts
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
	AccountID       uuid.UUID `json:"account_id"`
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
		AccountID:       dbTransaction.AccountID,
	}
}

func databaseTransactionsToTransactions(dbTransactions []database.Transaction) []Transaction {
	transactions := []Transaction{}
	for _, dbTransaction := range dbTransactions {
		transactions = append(transactions, databaseTransactionToTransaction(dbTransaction))
	}
	return transactions
}

type OthersAccount struct {
	ID            int32     `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updates_at"`
	AccountName   string    `json:"account_name"`
	AccountNumber string    `json:"account_number"`
	SortCode      string    `json:"sort_code"`
}

func databaseOthersAccountToOthersAccount(dbOthersAccount database.OthersAccount) OthersAccount {
	return OthersAccount{
		ID:            dbOthersAccount.ID,
		CreatedAt:     dbOthersAccount.CreatedAt,
		UpdatedAt:     dbOthersAccount.UpdatedAt,
		AccountName:   dbOthersAccount.AccountName,
		AccountNumber: dbOthersAccount.AccountNumber,
		SortCode:      dbOthersAccount.SortCode,
	}
}

func databaseOthersAccountsToOthersAccounts(dbOthersAccounts []database.OthersAccount) []OthersAccount {
	accounts := []OthersAccount{}
	for _, dbOthersAccount := range dbOthersAccounts {
		accounts = append(accounts, databaseOthersAccountToOthersAccount(dbOthersAccount))
	}
	return accounts
}
