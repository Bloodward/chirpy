package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	got, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT returned error: %v", err)
	}
	if got != userID {
		t.Errorf("expected user ID %v, got %v", userID, got)
	}
}

func TestValidateJWT_Expired(t *testing.T) {
	token, err := MakeJWT(uuid.New(), "test-secret", -time.Second)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	_, err = ValidateJWT(token, "test-secret")
	if !errors.Is(err, jwt.ErrTokenExpired) {
		t.Errorf("expected ErrTokenExpired, got %v", err)
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	token, err := MakeJWT(uuid.New(), "right-secret", time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT returned errror: %v", err)
	}

	_, err = ValidateJWT(token, "wrong-secret")
	if !errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		t.Errorf("expected ErrTokenSignatureInvalid, got %v", err)
	}
}

func TestVaslidateJWT_Malformed(t *testing.T) {
	_, err := ValidateJWT("not.a.jwt", "test-secret")
	if err == nil {
		t.Error("expected error for malformed token, got nil")
	}
}
