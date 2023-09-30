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

func (cfg *apiConfig) handlerCreateAccount(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		AccountName   string  `json:"account_name"`
		Balance       float32 `json:"balance"`
		AccountNumber string  `json:"account_number"`
		SortCode      string  `json:"sort_code"`
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

	account, err := cfg.DB.CreateAccount(r.Context(), database.CreateAccountParams{
		ID:            uuid.New(),
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
		AccountName:   params.AccountName,
		Balance:       params.Balance,
		AccountNumber: params.AccountNumber,
		SortCode:      params.SortCode,
		UserID:        userID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't create account")
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseAccountToAccount(account))
}

func (cfg *apiConfig) handlerUpdateAccount(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ID          string  `json:"id"`
		AccountName string  `json:"account_name"`
		Balance     float32 `json:"balance"`
	}
	type response struct {
		Account
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

	accountID, err := uuid.Parse(params.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Unable to parse UUID: %v", err))
		return
	}

	account, err := cfg.DB.UpdateAccount(r.Context(), database.UpdateAccountParams{
		ID:          accountID,
		UserID:      userID,
		UpdatedAt:   time.Now(),
		AccountName: params.AccountName,
		Balance:     params.Balance,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't update account: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Account: Account{
			ID:            account.ID,
			CreatedAt:     account.CreatedAt,
			UpdatedAt:     account.UpdatedAt,
			AccountName:   account.AccountName,
			Balance:       account.Balance,
			AccountNumber: account.AccountNumber,
			SortCode:      account.SortCode,
			UserID:        account.UserID,
		},
	})
}

func (cfg *apiConfig) handlerDeleteAccount(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int32)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get userID from context")
		return
	}

	accountIDString := chi.URLParam(r, "accountID")
	accountID, err := uuid.Parse(accountIDString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Unable to parse UUID: %v", err))
		return
	}

	err = cfg.DB.DeleteAccount(r.Context(), database.DeleteAccountParams{
		ID:     accountID,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't delete account: %v", err))
		return
	}

	respondWithJSON(w, 200, struct{}{})
}
