package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MeirionL/boing-block/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetTransactions(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		AccountID string `json:"account_id"`
	}

	userID, ok := r.Context().Value(userIDKey).(int32)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "couldn't get userID from context")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("couldn't decode parameters: %v", err))
		return
	}

	accountIDString := params.AccountID
	accountID, err := uuid.Parse(accountIDString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("unable to parse UUID: %v", err))
		return
	}

	accountUserID, err := cfg.DB.GetUserIDByAccountID(r.Context(), accountID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("couldn't get user id: %v", err))
		return
	}

	if userID != accountUserID {
		respondWithError(w, http.StatusForbidden, fmt.Sprintf("action is not authorized for user: %v", userID))
		return
	}

	transacIDString := r.URL.Query().Get("transaction_id")
	if transacIDString != "" {
		transacID, err := uuid.Parse(transacIDString)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("unable to parse UUID: %v", err))
			return
		}

		cfg.handlerGetTransactionByID(w, r, accountID, transacID)
		return
	}

	transacType := r.URL.Query().Get("transaction_type")
	if transacType != "" {
		cfg.handlerGetTransactionsByType(w, r, accountID, transacType)
		return
	}

	accountNumber := r.URL.Query().Get("account_number")
	sortCode := r.URL.Query().Get("sort_code")
	if accountNumber != "" && sortCode != "" {
		cfg.handlerGetTransactionsByAccount(w, r, accountID, accountNumber, sortCode)
		return
	}

	limitString := r.URL.Query().Get("limit")
	if limitString != "" {
		limitInt, err := strconv.Atoi(limitString)
		limit := int32(limitInt)
		if err != nil {
			respondWithError(w, 400, "couldn't convert limit string to int")
		}
		cfg.handlerGetTransactionsWithLimit(w, r, accountID, limit)
		return
	}

	othersAccount := r.URL.Query().Get("others_account_id")
	if othersAccount != "" {
		cfg.getTransactionsByOthersAccount(w, r, accountID, othersAccount)
		return
	}

	transactions, err := cfg.DB.GetTransactions(r.Context(), accountID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("couldn't get transactions: %v", err))
		return
	}

	respondWithJSON(w, 201, databaseTransactionsToTransactions(transactions))
}

func (cfg *apiConfig) handlerGetTransactionByID(w http.ResponseWriter, r *http.Request, accountID, transacID uuid.UUID) {
	transaction, err := cfg.DB.GetTransactionByID(r.Context(), database.GetTransactionByIDParams{
		AccountID: accountID,
		ID:        transacID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("couldn't get transaction by id: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseTransactionToTransaction(transaction))
}

func (cfg *apiConfig) handlerGetTransactionsByType(w http.ResponseWriter, r *http.Request, accountID uuid.UUID, transacType string) {
	transactions, err := cfg.DB.GetTransactionsByType(r.Context(), database.GetTransactionsByTypeParams{
		AccountID: accountID,
		Type:      transacType,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("couldn't get transactions by type: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseTransactionsToTransactions(transactions))
}

func (cfg *apiConfig) handlerGetTransactionsByAccount(w http.ResponseWriter, r *http.Request, accountID uuid.UUID, accountNumber, sortCode string) {
	transactions, err := cfg.DB.GetTransactionsByAccount(r.Context(), database.GetTransactionsByAccountParams{
		AccountID:     accountID,
		AccountNumber: accountNumber,
		SortCode:      sortCode,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("couldn't get transactions by account: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseTransactionsToTransactions(transactions))
}

func (cfg *apiConfig) handlerGetTransactionsWithLimit(w http.ResponseWriter, r *http.Request, accountID uuid.UUID, limit int32) {
	transactions, err := cfg.DB.GetTransactionsWithLimit(r.Context(), database.GetTransactionsWithLimitParams{
		AccountID: accountID,
		Limit:     limit,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("couldn't get transactions: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseTransactionsToTransactions(transactions))
}

func (cfg *apiConfig) getTransactionsByOthersAccount(w http.ResponseWriter, r *http.Request, accountID uuid.UUID, othersAccountIDString string) {
	othersAccountIDInt, err := strconv.Atoi(othersAccountIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("couldn't convert others account id string to int: %v", err))
	}
	othersAccountID := int32(othersAccountIDInt)
	account, err := cfg.DB.GetOthersAccountByID(r.Context(), database.GetOthersAccountByIDParams{
		AccountID: accountID,
		ID:        othersAccountID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("couldn't get others account by id: %v", err))
		return
	}

	transactions, err := cfg.DB.GetTransactionsByOthersAccount(r.Context(), database.GetTransactionsByOthersAccountParams{
		AccountID:     accountID,
		AccountNumber: account.AccountNumber,
		SortCode:      account.SortCode,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("couldn't get transactions by others account: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseTransactionsToTransactions(transactions))
}
