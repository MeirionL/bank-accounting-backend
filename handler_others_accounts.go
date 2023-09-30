package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/MeirionL/personal-finance-app/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetOthersAccounts(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		AccountID string `json:"account_id"`
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

	accountNumber := r.URL.Query().Get("account_number")
	sortCode := r.URL.Query().Get("sort_code")

	if accountNumber != "" && sortCode != "" {
		cfg.handlerGetOthersAccountByDetails(w, r, accountNumber, sortCode, accountID)
		return
	}

	accounts, err := cfg.DB.GetOthersAccounts(r.Context(), accountID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get others accounts: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, databaseOthersAccountsToOthersAccounts(accounts))
}

func (cfg *apiConfig) handlerGetOthersAccountByDetails(w http.ResponseWriter, r *http.Request, accountNumber, sortCode string, accID uuid.UUID) {
	account, err := cfg.DB.GetOthersAccountByDetails(r.Context(), database.GetOthersAccountByDetailsParams{
		AccountID:     accID,
		AccountNumber: accountNumber,
		SortCode:      sortCode,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get others account")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseOthersAccountToOthersAccount(account))
}

func (cfg *apiConfig) handlerGetOthersAccountByID(w http.ResponseWriter, r *http.Request, accountIDString string) {
	accountIDInt, err := strconv.Atoi(accountIDString)
	if err != nil {
		respondWithError(w, 400, "Couldn't convert account id string to int")
	}
	accountID := int32(accountIDInt)
	account, err := cfg.DB.GetOthersAccountByID(r.Context(), database.GetOthersAccountByIDParams{
		ID: accountID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get others account")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseOthersAccountToOthersAccount(account))
}

func (cfg *apiConfig) createOthersAccount(w http.ResponseWriter, r *http.Request, accName, accNumber, sortCode string, accID uuid.UUID) {
	err := cfg.DB.CreateOthersAccount(r.Context(), database.CreateOthersAccountParams{
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		AccountName:   accName,
		AccountNumber: accNumber,
		SortCode:      sortCode,
		AccountID:     accID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't create others account")
		return
	}
}

func (cfg *apiConfig) checkOthersAccountName(w http.ResponseWriter, r *http.Request, accName, accNumber, sortCode string, accID uuid.UUID) {
	account, err := cfg.DB.GetOthersAccountByDetails(r.Context(), database.GetOthersAccountByDetailsParams{
		AccountID:     accID,
		AccountNumber: accNumber,
		SortCode:      sortCode,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get others account")
		return
	}

	if account.AccountName != accName {
		err = cfg.DB.UpdateOthersAccountName(r.Context(), database.UpdateOthersAccountNameParams{
			AccountName:   accName,
			AccountNumber: accNumber,
			SortCode:      sortCode,
			AccountID:     account.AccountID,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to update others account name")
			return
		}
	}
}
