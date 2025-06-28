package auth

import (
	"testing"
	"time"
	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret-key-32-bytes-long!!!"

	// Test 1: Create and validate a valid token
	t.Run("valid token", func(t *testing.T) {
		token, err := MakeJWT(userID, secret, time.Hour)
		if err != nil {
			t.Fatalf("failed to make JWT: %v", err)
		}

		gotUserID, err := ValidateJWT(token, secret)
		if err != nil {
			t.Fatalf("failed to validate JWT: %v", err)
		}

		if gotUserID != userID {
			t.Errorf("got userID %v, want %v", gotUserID, userID)
		}
	})

	// Test 2: Expired token should fail
	t.Run("expired token", func(t *testing.T) {
		token, err := MakeJWT(userID, secret, -time.Hour) // negative = already expired
		if err != nil {
			t.Fatalf("failed to make JWT: %v", err)
		}

		_, err = ValidateJWT(token, secret)
		if err == nil {
			t.Error("expected error for expired token, got nil")
		}
	})

	// Test 3: Wrong secret should fail
	t.Run("wrong secret", func(t *testing.T) {
		token, err := MakeJWT(userID, secret, time.Hour)
		if err != nil {
			t.Fatalf("failed to make JWT: %v", err)
		}

		wrongSecret := "wrong-secret-key!!!!!!!!!!!!!!!"
		_, err = ValidateJWT(token, wrongSecret)
		if err == nil {
			t.Error("expected error for wrong secret, got nil")
		}
	})

	// Test 4: Invalid token format
	t.Run("invalid token", func(t *testing.T) {
		_, err := ValidateJWT("not.a.valid.token", secret)
		if err == nil {
			t.Error("expected error for invalid token, got nil")
		}
	})
}