package main

import (
	"context"
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
			respondWithError(w, http.StatusUnauthorized, "Unauthorized user")
			return
		}

		userIDInt, err := strconv.Atoi(userIDString)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't convert userID to int")
			return
		}
		userID := int32(userIDInt)

		r = r.WithContext(context.WithValue(r.Context(), userIDKey, userID))
		next.ServeHTTP(w, r)
	})
}
