package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MeirionL/personal-finance-app/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateTransaction(w http.ResponseWriter, r *http.Request) {
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
		AccountID       string  `json:"account_id"`
	}

	userID, ok := r.Context().Value(userIDKey).(int32)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get userID from context")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	accountID, err := uuid.Parse(params.AccountID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("couldn't parse account id: %v", err))
		return
	}

	accountUserID, err := cfg.DB.GetUserIDByAccountID(r.Context(), accountID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("couldn't get user id: %v", err))
		return
	}

	if userID != accountUserID {
		respondWithError(w, http.StatusForbidden, "action is not authorized for user")
		return
	}

	layout := "2006-01-02T15:04:05Z"
	transacTime, err := time.Parse(layout, params.TransactionTime)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("couldn't parse transaction time: %v", err))
		return
	}

	transaction, err := cfg.DB.CreateTransaction(r.Context(), database.CreateTransactionParams{
		ID:              uuid.New(),
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
		TransactionTime: transacTime,
		Type:            params.TransactionType,
		Amount:          params.Amount,
		PreBalance:      params.PreBalance,
		PostBalance:     params.PostBalance,
		NewAccount:      params.NewAccount,
		Name:            params.AccountName,
		AccountNumber:   params.AccountNumber,
		SortCode:        params.SortCode,
		AccountID:       accountID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn't create transaction: %v", err))
		return
	}

	if params.NewAccount {
		cfg.createOthersAccount(w, r, params.AccountName, params.AccountNumber, params.SortCode, accountID)
	}

	if !params.NewAccount {
		cfg.checkOthersAccountName(w, r, params.AccountName, params.AccountNumber, params.SortCode, accountID)
	}

	respondWithJSON(w, http.StatusCreated, databaseTransactionToTransaction(transaction))
}

func (cfg *apiConfig) handlerDeleteTransaction(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		AccountID string `json:"account_id"`
	}

	userID, ok := r.Context().Value(userIDKey).(int32)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get userID from context")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	accountID, err := uuid.Parse(params.AccountID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("couldn't parse account id: %v", err))
	}

	accountUserID, err := cfg.DB.GetUserIDByAccountID(r.Context(), accountID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("couldn't get user id: %v", err))
		return
	}

	if userID != accountUserID {
		respondWithError(w, http.StatusForbidden, "action is not authorized for user")
		return
	}

	transacIDString := chi.URLParam(r, "transactionID")
	transacID, err := uuid.Parse(transacIDString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Unable to parse UUID: %v", err))
		return
	}

	err = cfg.DB.DeleteTransaction(r.Context(), database.DeleteTransactionParams{
		AccountID: accountID,
		ID:        transacID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't delete transaction: %v", err))
		return
	}

	respondWithJSON(w, 200, struct{}{})
}
