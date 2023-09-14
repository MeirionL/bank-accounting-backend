package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerGetBalance(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserID string `json:"user_id"`
	}
	type response struct {
		Balance int `json:"balance"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	// userID, err := strconv.Atoi(params.UserID)

	// user, err := cfg.DB.GetUser(userID)

	respondWithJSON(w, http.StatusOK, response{
		Balance: 24,
	})
}
