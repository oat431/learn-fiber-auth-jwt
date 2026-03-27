package utils

import (
	"encoding/hex"
	"testing"
)

func TestGenerateVerifyToken(t *testing.T) {
	token := GenerateVerifyToken()

	if token == "" {
		t.Fatal("generated verify token is empty")
	}

	// 32 random bytes encode to 64 hex characters
	const expectedLen = 64
	if len(token) != expectedLen {
		t.Errorf("expected token length %d, got %d", expectedLen, len(token))
	}

	if _, err := hex.DecodeString(token); err != nil {
		t.Errorf("token is not valid hex encoding: %v", err)
	}
}

func TestGenerateVerifyTokenUniqueness(t *testing.T) {
	token1 := GenerateVerifyToken()
	token2 := GenerateVerifyToken()

	if token1 == token2 {
		t.Error("two generated tokens should not be equal")
	}
}
