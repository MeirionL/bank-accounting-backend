package auth

import (
	"net/http"
	"testing"
	"time"
)

func TestHashPassword(t *testing.T) {
	password := "password"

	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Errorf("HashPassword failed: %v", err)
	}

	if len(hashedPassword) == 0 {
		t.Error("hashed password is empty")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	hash, _ := HashPassword("foo")
	err := CheckPasswordHash("foo", hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash failed: %v", err)
	}
}

func TestMakeJWT(t *testing.T) {
	userID := int32(123)
	tokenSecret := "my-secret-key"
	expiresIn := time.Minute * 30

	tokenString, err := MakeJWT(userID, tokenSecret, expiresIn, TokenTypeRefresh)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	userIDString, err := ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}

	if userIDString != "123" {
		t.Fatalf("expected: %v got: %v", "123", userIDString)
	}
}

func TestRefreshToken(t *testing.T) {
	tokenSecret := "qRcrjPQ989usFYG/O2fcBXASlGNisXI4v0+9bWNwgXCNXjYPKww4d6j93faBB8poxi9K7QzkWLSC8UGtvgNzOw=="

	tokenString, err := MakeJWT(int32(1), tokenSecret, (time.Minute * 30), TokenTypeRefresh)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	accessToken, err := RefreshToken(tokenString, tokenSecret)
	if err != nil {
		t.Fatalf("RefreshToken failed: %v", err)
	}
	if len(accessToken) != 167 {
		t.Fatal("Invalid access token")
	}
}

func TestGetBearerToken(t *testing.T) {
	token := "12345"
	header := make(http.Header)
	header.Add("Authorization", "Bearer "+token)

	bearerToken, err := GetBearerToken(header)
	if err != nil {
		t.Fatalf("GetBearerToken failed: %v", err)
	}

	if token != bearerToken {
		t.Fatalf("expected: %v got: %v", token, bearerToken)
	}
}
