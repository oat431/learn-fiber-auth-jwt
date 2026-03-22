package utils

import (
	"testing"

	"github.com/google/uuid"
)

func TestGenerateAccessToken(t *testing.T) {
	authID := uuid.New()
	token, err := GenerateAccessToken(authID)
	if err != nil {
		t.Fatalf("failed to generate access token: %v", err)
	}

	if token == "" {
		t.Fatal("generated token is empty")
	}

	parsedToken, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}

	claims, ok := parsedToken.Claims.(*JWTClaims)
	if !ok || !parsedToken.Valid {
		t.Fatal("invalid token or claims")
	}

	if claims.AuthID != authID {
		t.Errorf("expected AuthID %v, got %v", authID, claims.AuthID)
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	token, err := GenerateRefreshToken()
	if err != nil {
		t.Fatalf("failed to generate refresh token: %v", err)
	}

	if token == "" {
		t.Fatal("generated refresh token is empty")
	}
}
