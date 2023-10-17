package main

import (
	"fmt"
	"net/http"

	"github.com/MeirionL/personal-finance-app/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetAccounts(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int32)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "couldn't get userID from context")
		return
	}

	accountNumber := r.URL.Query().Get("account_number")
	sortCode := r.URL.Query().Get("sort_code")

	if accountNumber != "" && sortCode != "" {
		cfg.handlerGetAccountByDetails(w, r, userID, accountNumber, sortCode)
		return
	}

	accountIDString := r.URL.Query().Get("id")

	if accountIDString != "" {
		accountID, err := uuid.Parse(accountIDString)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("unable to parse UUID: %v", err))
			return
		}
		cfg.handlerGetAccountByID(w, r, userID, accountID)
		return
	}

	accounts, err := cfg.DB.GetAccounts(r.Context(), userID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("couldn't get accounts: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, databaseAccountsToAccounts(accounts))
}

func (cfg *apiConfig) handlerGetAccountsBalances(w http.ResponseWriter, r *http.Request) {
	type Balance struct {
		AccountName string  `json:"account_name"`
		Balance     float32 `json:"balance"`
	}

	userID, ok := r.Context().Value(userIDKey).(int32)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "couldn't get userID from context")
		return
	}

	balances, err := cfg.DB.GetAccountsBalances(r.Context(), userID)
	if len(balances) == 0 {
		respondWithJSON(w, http.StatusOK, struct{}{})
		return
	} else if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("couldn't get accounts balances: %v", err))
		return
	}

	var finalBalances []Balance
	for _, balance := range balances {
		finalBalances = append(finalBalances, Balance{
			AccountName: balance.AccountName,
			Balance:     balance.Balance,
		})
	}

	respondWithJSON(w, http.StatusOK, finalBalances)
}

func (cfg *apiConfig) handlerGetAccountByDetails(w http.ResponseWriter, r *http.Request, userID int32, accountNumber, sortCode string) {
	account, err := cfg.DB.GetAccountByDetails(r.Context(), database.GetAccountByDetailsParams{
		AccountNumber: accountNumber,
		SortCode:      sortCode,
		UserID:        userID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("couldn't get account by details: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseAccountToAccount(account))
}

func (cfg *apiConfig) handlerGetAccountByID(w http.ResponseWriter, r *http.Request, userID int32, id uuid.UUID) {
	account, err := cfg.DB.GetAccountByID(r.Context(), database.GetAccountByIDParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("couldn't get account by ID: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseAccountToAccount(account))
}
