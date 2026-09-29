package auth

import (
	"encoding/hex"
	"testing"
)

func TestMakeRefreshToken(t *testing.T) {
	tokenOne, err := MakeRefreshToken()
	if err != nil {
		t.Errorf("System failed to create token: %v", err)
	}
	tokenTwo, err := MakeRefreshToken()
	if err != nil {
		t.Errorf("System failed to create token: %v", err)
	}

	// Test for empty tokens
	if tokenOne == "" || tokenTwo == "" {
		t.Errorf("Token is empty")
	}

	// Test for randomness
	if tokenOne == tokenTwo {
		t.Errorf("Token is not random")
	}

	// Test length
	if len(tokenOne) != 64 || len(tokenTwo) != 64 {
		t.Errorf("Did not return 64 characters")
	}

	// Test for invalid chars
	_, err = hex.DecodeString(tokenOne)
	if err != nil {
		t.Errorf("tokenOne is not valid hex: %v", err)
	}
}
