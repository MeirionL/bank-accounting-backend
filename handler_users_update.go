package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("couldn't hash password: %v", err))
		return
	}

	duplicateUsers, err := cfg.DB.GetUsersByDetails(r.Context(), database.GetUsersByDetailsParams{
		Name:           params.Name,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("couldn't get users by details: %v", err))
		return
	}

	if len(duplicateUsers) >= 1 {
		respondWithError(w, http.StatusBadRequest, "users with entered parameters already exist")
		return
	}

	user, err := cfg.DB.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userID,
		UpdatedAt:      time.Now(),
		Name:           params.Name,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("couldn't update user: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			UpdatedAt: user.UpdatedAt,
			Name:      user.Name,
		},
	})
}

func (cfg *apiConfig) handlerDeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int32)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "couldn't get userID from context")
		return
	}

	err := cfg.DB.DeleteUser(r.Context(), userID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("couldn't delete user: %v", err))
		return
	}
	respondWithJSON(w, 200, struct{}{})
}
