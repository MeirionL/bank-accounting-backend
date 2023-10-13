package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MeirionL/personal-finance-app/internal/auth"
)

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ID       int32  `json:"id"`
		Password string `json:"password"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("couldn't decode parameters: %v", err))
		return
	}

	user, err := cfg.DB.GetUserByID(r.Context(), params.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("couldn't get user: %v", err))
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("invalid password: %v", err))
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour, auth.TokenTypeAccess)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("couldn't create access JWT: %v", err))
		return
	}

	refreshToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour*24*60, auth.TokenTypeRefresh)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("couldn't create refresh JWT: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:   user.ID,
			Name: user.Name,
		},
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}
