package main

import (
	"net/http"
)

func (apiCfg *apiConfig) handlerTransactionGet(w http.ResponseWriter, r *http.Request) {
	// transacIDString := chi.URLParam(r, "id")
	// transacID, err := strconv.Atoi(transacIDString)
	// if err != nil {
	// 	respondWithError(w, http.StatusBadRequest, "Invalid transaction ID")
	// 	return
	// }

	// dbTransac, err := cfg.DB.GetTransaction(transacID)

	// respondWithJSON(w, http.StatusOK, Transaction{
	// 	ID:             dbTransac.ID,
	// 	PaymentAmount:  dbTransac.PaymentAmount,
	// 	PostBalance:    dbTransac.CurrentAccount.Balance,
	// 	PartnerAccount: dbTransac.PartnerAccount.HolderName,
	// 	TimeOf:         dbTransac.TimeOf,
	// 	Message:        dbTransac.Message,
	// })
}

func (apiCfg *apiConfig) handlerTransactionsRetrieve(w http.ResponseWriter, r *http.Request) {
	// 	dbTransacs, err := cfg.DB.GetTransactions()
	// 	if err != nil {
	// 		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve transactions")
	// 		return
	// 	}

	// transacs := []Transaction{}
	//
	//	for _, dbTransac := range dbTransacs {
	//		transacs = append(transacs, Transaction{
	//			ID:             dbTransac.ID,
	//			PaymentAmount:  dbTransac.PaymentAmount,
	//			PostBalance:    dbTransac.CurrentAccount.Balance,
	//			PartnerAccount: dbTransac.PartnerAccount.HolderName,
	//			TimeOf:         dbTransac.TimeOf,
	//			Message:        dbTransac.Message,
	//		})
	//	}
}

func (apiCfg *apiConfig) handlerLastTransactionGet(w http.ResponseWriter, r *http.Request) {
	// 	dbTransac, err := cfg.DB.GetLastTransaction()
	// 	if err != nil {
	// 		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve transaction")
	// 		return
	// 	}

	//	respondWithJSON(w, http.StatusOK, Transaction{
	//		ID:             dbTransac.ID,
	//		PaymentAmount:  dbTransac.PaymentAmount,
	//		PostBalance:    dbTransac.CurrentAccount.Balance,
	//		PartnerAccount: dbTransac.PartnerAccount.HolderName,
	//		TimeOf:         dbTransac.TimeOf,
	//		Message:        dbTransac.Message,
	//	})
}
