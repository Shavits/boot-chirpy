package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	secret := "mysecretkey"
	expiresIn := 2 * time.Hour

	token, err := MakeJWT(uuid.MustParse(userID), secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	returnedUserID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}

	if returnedUserID.String() != userID {
		t.Errorf("Expected userID %s, got %s", userID, returnedUserID.String())
	}
}

func TestJWT_ExpiredToken(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	secret := "mysecretkey"
	expiresIn := -1 * time.Hour // Token already expired

	token, err := MakeJWT(uuid.MustParse(userID), secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Fatalf("Expected error for expired token, got none")
	}
}

func TestJWT_InvalidToken(t *testing.T) {
	secret := "mysecretkey"
	invalidToken := "invalid.token.string"

	_, err := ValidateJWT(invalidToken, secret)
	if err == nil {
		t.Fatalf("Expected error for invalid token, got none")
	}
}

func TestJWT_WrongSecret(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	secret := "mysecretkey"
	wrongSecret := "wrongsecretkey"
	expiresIn := 2 * time.Hour

	token, err := MakeJWT(uuid.MustParse(userID), secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Fatalf("Expected error for wrong secret, got none")
	}
}