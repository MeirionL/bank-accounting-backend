package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MeirionL/personal-finance-app/internal/auth"
)

type contextKey string

const userIDKey contextKey = "userID"

func (cfg *apiConfig) middlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIDString, err := auth.ValidateUser(r, cfg.jwtSecret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("unauthorized user: %v", err))
			return
		}

		userIDInt, err := strconv.Atoi(userIDString)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("couldn't convert userID to int: %v", err))
			return
		}
		userID := int32(userIDInt)

		r = r.WithContext(context.WithValue(r.Context(), userIDKey, userID))
		next.ServeHTTP(w, r)
	})
}
