package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type TokenType string

const (
	TokenTypeAccess  TokenType = "pfa-access"
	TokenTypeRefresh TokenType = "pfa-refresh"
)

var ErrNoAuthHeaderIncluded = errors.New("not auth header included in request")

func HashPassword(password string) (string, error) {
	dat, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID int32, tokenSecret string, expiresIn time.Duration, tokenType TokenType) (string, error) {
	signingKey := []byte(tokenSecret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    string(tokenType),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   fmt.Sprintf("%d", userID),
	})
	return token.SignedString(signingKey)
}

func RefreshToken(tokenString, tokenSecret string) (string, error) { // takes in a refresh token, and if still valid, creates and returns a new access token.
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		return "", err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return "", err
	}
	if issuer != string(TokenTypeRefresh) {
		return "", errors.New("invalid issuer")
	}

	userIDInt, err := strconv.Atoi(userIDString)
	if err != nil {
		return "", err
	}
	userID := int32(userIDInt)

	newToken, err := MakeJWT(
		userID,
		tokenSecret,
		time.Hour,
		TokenTypeAccess,
	)
	if err != nil {
		return "", err
	}

	return newToken, nil
}

func ValidateUser(r *http.Request, jwtSecret string) (string, error) {
	token, err := GetBearerToken(r.Header)
	if err != nil {
		return "", err
	}

	userID, err := ValidateJWT(token, jwtSecret)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}
	return splitAuth[1], nil
}

func ValidateJWT(tokenString, tokenSecret string) (string, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		return "", err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return "", err
	}
	if issuer != string(TokenTypeAccess) && issuer != string(TokenTypeRefresh) {
		return "", errors.New("invalid access token, potentially expired")
	}

	return userIDString, nil
}
