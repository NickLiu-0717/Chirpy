package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestExpiredToken(t *testing.T) {
	// Setup
	userID := uuid.New()
	secret := "testSecret"

	// Exercise - create token
	token, err := MakeJWT(userID, secret)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	// Wait for token to expire
	time.Sleep(time.Millisecond * 2)

	// Verify - token should be invalid
	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Error("Expected error for expired token, got nil")
	}
}

func TestJWT(t *testing.T) {
	userID := uuid.New()
	secretKey := "secret key"

	tokenstr, err := MakeJWT(userID, secretKey)
	if err != nil {
		t.Fatalf("Error creating token: %v", err)
	}

	getID, err := ValidateJWT(tokenstr, secretKey)
	if err != nil {
		t.Fatalf("Error validating token: %v", err)
	}
	if getID != userID {
		t.Errorf("wrong ID: got %v, want %v", getID, userID)
	}
}

func TestBearerToken(t *testing.T) {
	tokenString := "Mystring"
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+tokenString)

	gotToken, err := GetBearerToken(req.Header)
	if err != nil {
		t.Fatalf("Failed to get bearer token: %v", err)
	}
	if gotToken != tokenString {
		t.Fatalf("Wrong token string: got %v, want %v", gotToken, tokenString)
	}

	req.Header.Del("Authorization")
	_, err = GetBearerToken(req.Header)
	if err == nil {
		t.Error("Expected error for no token, got nil")
	}

}
