package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MeirionL/personal-finance-app/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateTransaction(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		TransactionTime string  `json:"transaction_time"`
		TransactionType string  `json:"transaction_type"`
		Amount          float32 `json:"amount"`
		PreBalance      float32 `json:"pre_balance"`
		PostBalance     float32 `json:"post_balance"`
		NewAccount      bool    `json:"new_account"`
		AccountName     string  `json:"account_name"`
		AccountNumber   string  `json:"account_number"`
		SortCode        string  `json:"sort_code"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	layout := "2006-01-02T15:04:05Z"

	t, err := time.Parse(layout, params.TransactionTime)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	transaction, err := cfg.DB.CreateTransaction(r.Context(), database.CreateTransactionParams{
		ID:              uuid.New(),
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
		TransactionTime: t,
		Type:            params.TransactionType,
		Amount:          params.Amount,
		PreBalance:      params.PreBalance,
		PostBalance:     params.PostBalance,
		NewAccount:      params.NewAccount,
		Name:            params.AccountName,
		AccountNumber:   params.AccountNumber,
		SortCode:        params.SortCode,
		UserID:          user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't create transaction")
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseTransactionToTransaction(transaction))
}

func (cfg *apiConfig) handlerGetTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, err := cfg.DB.GetTransactions(r.Context())
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get transactions: %v", err))
		return
	}

	respondWithJSON(w, 201, databaseTransactionsToTransactions(transactions))
}
