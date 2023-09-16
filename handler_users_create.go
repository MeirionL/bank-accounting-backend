package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/MeirionL/personal-finance-app/internal/auth"
	"github.com/MeirionL/personal-finance-app/internal/database"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
		Name:           params.Name,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseUserToUser(user))
}
