package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	testSuite := []struct {
		name     string
		password string
		want     bool
	}{
		{name: "Valid Password", password: "S0m3Text!!", want: false},
		{name: "Long Password", password: "verylongpasswordthatshouldstillwork123!", want: false},
		{name: "Empty Password", password: "", want: false},
	}

	for _, ts := range testSuite {
		t.Run(ts.name, func(t *testing.T) {
			hash, err := HashPassword(ts.password)
			if err != nil {
				t.Fatalf("Unexpected error with hashing function: %v", err)
			}

			if (hash == "") != ts.want {
				t.Errorf("Wanted hash empty=%v, got: %q", ts.want, hash)
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testing123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to prepare test hash: %v", err)
	}

	testSuite := []struct {
		name      string
		password  string
		hash      string
		wantMatch bool
	}{
		{name: "Correct Password", password: "testing123", hash: hash, wantMatch: true},
		{name: "Wrong Password", password: "wrongpass", hash: hash, wantMatch: false},
		{name: "Empty Password", password: "", hash: hash, wantMatch: false},
	}

	for _, ts := range testSuite {
		t.Run(ts.name, func(t *testing.T) {
			match, err := CheckPasswordHash(ts.password, ts.hash)
			if err != nil {
				t.Fatalf("Unexpected error checking password hash: %v", err)
			}

			if match != ts.wantMatch {
				t.Errorf("Match mismatch for %s: wanted %v, got %v", ts.name, ts.wantMatch, match)
			}
		})
	}
}
