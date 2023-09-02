package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type Account struct {
	ID             int       `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	AccountName    string    `json:"account_name"`
	Password       string    `json:"-"`
	AccountDetails struct {
		AccountNumber string `json:"account_number"`
		SortCode      string `json:"sort_code"`
	}
}

func (cfg *apiConfig) handlerCreateAccount(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		AccountName    string `json:"account_name"`
		Password       string `json:"password"`
		AccountDetails struct {
			AccountNumber string `json:"account_number"`
			SortCode      string `json:"sort_code"`
		}
	}
	type response struct {
		Account
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	// account, err := cfg.DB.CreateAccount(r.Context(), database.CreateAccountParams{
	// 	ID:             uuid.New(),
	// 	CreatedAt:      time.Now().UTC(),
	// 	UpdatedAt:      time.Now().UTC(),
	// 	AccountName:    params.AccountName,
	// 	Password:       params.Password,
	// 	AccountDetails: params.AccountDetails,
	// })
	// if err != nil {
	// 	respondWithError(w, http.StatusBadRequest, "Couldn't create account")
	// 	return
	// }

	// respondWithJSON(w, http.StatusCreated, response{
	// 	Account: Account{
	// 		ID:             account.ID,
	// 		AccountDetails: account.AccountDetails,
	// 	},
	// })
}
