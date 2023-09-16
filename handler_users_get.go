package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) handlerGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := cfg.DB.GetUsers(r.Context())
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get users: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, databaseUsersToUsers(users))
}

func (cfg *apiConfig) handlerGetUserByID(w http.ResponseWriter, r *http.Request) {
	userIDString := chi.URLParam(r, "userID")
	userIDint, err := strconv.Atoi(userIDString)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't parse user ID: %v", err))
		return
	}
	userID := int32(userIDint)

	user, err := cfg.DB.GetUserByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't delete user: %v", err))
		return
	}

	respondWithJSON(w, 200, databaseUserToUser(user))
}
