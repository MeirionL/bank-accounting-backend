package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MeirionL/personal-finance-app/internal/auth"
	"github.com/MeirionL/personal-finance-app/internal/database"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}

	subject, err := auth.ValidateUser(r, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't validate user: %v", err))
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	userIDint, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't parse user ID: %v", err))
		return
	}
	userID := int32(userIDint)

	user, err := cfg.DB.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userID,
		Name:           params.Name,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			Name: user.Name,
			ID:   user.ID,
		},
	})
}

func (cfg *apiConfig) handlerDeleteUser(w http.ResponseWriter, r *http.Request) {
	subject, err := auth.ValidateUser(r, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't validate user: %v", err))
	}

	userIDint, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't parse user ID: %v", err))
		return
	}
	userID := int32(userIDint)

	err = cfg.DB.DeleteUser(r.Context(), userID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't delete user: %v", err))
		return
	}
	respondWithJSON(w, 200, struct{}{})
}
