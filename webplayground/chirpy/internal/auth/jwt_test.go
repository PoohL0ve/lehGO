package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	// Table-driven tests to verify token generation.
	testSuite := []struct {
		name    string
		userID  string
		secret  string
		wantErr bool
	}{
		{
			name:    "Valid Token",
			userID:  "123e4567-e89b-12d3-a456-426614174000",
			secret:  "super-secret-key",
			wantErr: false,
		},
		{
			name:    "Empty Secret Key",
			userID:  "123e4567-e89b-12d3-a456-426614174000",
			secret:  "",
			wantErr: false,
		},
	}

	for _, ts := range testSuite {
		t.Run(ts.name, func(t *testing.T) {
			id, err := uuid.Parse(ts.userID)
			if err != nil {
				t.Fatalf("Could not parse user ID: %v", err)
			}

			// Generate a token valid for 1 hour.
			token, err := MakeJWT(id, ts.secret, 1*time.Hour)
			if (err != nil) != ts.wantErr {
				t.Errorf("Wanted error = %v, got error = %v", ts.wantErr, err)
				return
			}

			if !ts.wantErr && token == "" {
				t.Errorf("Expected a non-empty token string, got empty string")
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	validUserID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	secret := "super-secret-key"

	// Create a valid token that expires in 1 hour.
	validToken, err := MakeJWT(validUserID, secret, 1*time.Hour)
	if err != nil {
		t.Fatalf("Setup failed: could not create valid test token: %v", err)
	}

	// Create an EXPIRED token by setting the expiration duration to negative 1 hour (-1 * time.Hour).
	expiredToken, err := MakeJWT(validUserID, secret, -1*time.Hour)
	if err != nil {
		t.Fatalf("Setup failed: could not create expired test token: %v", err)
	}

	tests := []struct {
		name     string
		token    string
		secret   string
		wantUUID string
		wantErr  bool
	}{
		{
			name:     "Valid token with correct secret",
			token:    validToken,
			secret:   secret,
			wantUUID: validUserID.String(),
			wantErr:  false,
		},
		{
			name:     "Invalid token due to incorrect secret",
			token:    validToken,
			secret:   "wrong-secret",
			wantUUID: uuid.Nil.String(),
			wantErr:  true,
		},
		{
			name:     "Expired token should fail validation",
			token:    expiredToken,
			secret:   secret,
			wantUUID: uuid.Nil.String(),
			wantErr:  true,
		},
		{
			name:     "Malformed token string",
			token:    "not.a.valid.jwt.token",
			secret:   secret,
			wantUUID: uuid.Nil.String(),
			wantErr:  true,
		},
		{
			name:     "Empty token string",
			token:    "",
			secret:   secret,
			wantUUID: uuid.Nil.String(),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			returnedUUID, err := ValidateJWT(tt.token, tt.secret)

			// Check whether an error occurred when we expected one (or vice versa).
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}

			// If no error was expected, verify that the returned UUID matches what we expected.
			if !tt.wantErr && returnedUUID.String() != tt.wantUUID {
				t.Errorf("ValidateJWT() returned userID = %q, want %q", returnedUUID, tt.wantUUID)
			}
		})
	}
	// Run tests: go test -v -run "TestMakeJWT|TestValidateJWT".
}
