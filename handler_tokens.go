package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MeirionL/personal-finance-app/internal/auth"
	"github.com/MeirionL/personal-finance-app/internal/database"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int32)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "couldn't get userID from context")
		return
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("couldn't get JWT: %v", err))
		return
	}

	err = cfg.DB.CreateRevokedToken(r.Context(), database.CreateRevokedTokenParams{
		Token:     refreshToken,
		RevokedAt: time.Now(),
		UserID:    userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("couldn't revoke token: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}

func (cfg *apiConfig) handlerGetRevokedTokens(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int32)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "couldn't get userID from context")
		return
	}

	revokedTokens, err := cfg.DB.GetRevokedTokens(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("couldn't get revoked tokens: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, revokedTokens)
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	userID, ok := r.Context().Value(userIDKey).(int32)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "couldn't get userID from context")
		return
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("couldn't get JWT: %v", err))
		return
	}

	revokedTokens, err := cfg.DB.GetRevokedTokens(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("couldn't get revoked tokens: %v", err))
		return
	}
	for _, token := range revokedTokens {
		if token == refreshToken {
			respondWithError(w, http.StatusUnauthorized, "refresh token is revoked")
			return
		}
	}

	accessToken, err := auth.RefreshToken(refreshToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("couldn't refresh token: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}
