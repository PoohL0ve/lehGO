package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name      string
		header    http.Header
		wantToken string
		wantErr   bool
	}{
		{
			name: "Valid Bearer token",
			header: map[string][]string{
				"Authorization": {"Bearer my-token-123"},
			},
			wantToken: "my-token-123",
			wantErr:   false,
		},
		{
			name:      "Missing Authorization header",
			header:    map[string][]string{},
			wantToken: "",
			wantErr:   true,
		},
		{
			name: "Missing Bearer prefix",
			header: map[string][]string{
				"Authorization": {"my-token-only"},
			},
			wantToken: "",
			wantErr:   true,
		},
		{
			name: "Bearer with spaces trimmed",
			header: map[string][]string{
				"Authorization": {"Bearer   spaced-token   "},
			},
			wantToken: "spaced-token",
			wantErr:   false,
		},
		{
			name: "Empty Bearer token",
			header: map[string][]string{
				"Authorization": {"Bearer "},
			},
			wantToken: "",
			wantErr:   false, // Empty string after trimming is technically valid
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetBearerToken(tt.header)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && token != tt.wantToken {
				t.Errorf("GetBearerToken() = %q, want %q", token, tt.wantToken)
			}
		})
	}
}
