package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Search for the authorisation header and returns its token string
func GetBearerToken(headers http.Header) (string, error) {
	// Get the authorisation header
	authheader := headers.Get("Authorization")
	if authheader == "" {
		return "", errors.New("The header is empty")
	}

	// Check that the authheader has "Bearer "
	isBearer := strings.HasPrefix(authheader, "Bearer ")
	if !isBearer {
		return "", errors.New("Missing prefix from authorisation")
	}

	// If all is well, trim prefix and white space
	removeBearer := strings.TrimPrefix(authheader, "Bearer ")
	cleanAuthToken := strings.TrimSpace(removeBearer)
	if cleanAuthToken == "" {
		return "", errors.New("Token is empty")
	}

	return cleanAuthToken, nil
}
