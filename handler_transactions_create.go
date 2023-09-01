package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func (apiCfg *apiConfig) handlerTransactionsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		CurrentAccount struct {
			HolderName    string  `json:"holder_name"`
			Balance       float64 `json:"balance"`
			AccountNumber string  `json:"accountNumber"`
			SortCode      string  `json:"sortCode"`
		}
		PartnerAccount struct {
			HolderName    string `json:"holder_name"`
			AccountNumber string `json:"accountNumber"`
			SortCode      string `json:"sortCode"`
		}
		TimeOf struct {
			Time time.Time `json:"time"`
			Date time.Time `json:"data"`
		}
		PaymentAmount float64 `json:"payment_amount"`
		IsOutgoing    bool    `json:"is_outgoing"`
		Message       string  `json:"message"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	transactionID := 1
	preBalance := params.CurrentAccount.Balance
	params.CurrentAccount.Balance -= params.PaymentAmount

	newPartner := true

	respondWithJSON(w, http.StatusCreated, transaction{
		ID:             transactionID,
		PreBalance:     preBalance,
		PaymentAmount:  params.PaymentAmount,
		PostBalance:    params.CurrentAccount.Balance,
		PartnerAccount: params.PartnerAccount.HolderName,
		AccountNumber:  params.PartnerAccount.AccountNumber,
		SortCode:       params.PartnerAccount.SortCode,
		NewPartner:     newPartner,
		TimeOf:         params.TimeOf,
		Message:        params.Message,
	})
}
