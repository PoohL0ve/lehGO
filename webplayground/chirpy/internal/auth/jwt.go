package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// MakeJWT creates a signed JSON Web Token string for a given user ID.
// It accepts the user's UUID, a secret key string used for signing, and an expiration duration.
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	// Secret keys must be passed as a byte slice to the HMAC signing function.
	signingKey := []byte(tokenSecret)

	// RegisteredClaims contains standard claims defined by the JWT spec (RFC 7519).
	// Issuer: who created the token.
	// IssuedAt: creation time (UTC).
	// ExpiresAt: exact timestamp when the token becomes invalid.
	// Subject: the entity the token represents (here, our user's string UUID).
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	}

	// Create a new token struct specifying HMAC-SHA256 (HS256) as the signing algorithm.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret byte slice and return the complete, encoded JWT string.
	return token.SignedString(signingKey)
}

// ValidateJWT verifies the signature of a token string using the secret key,
// checks that it hasn't expired, and returns the user's UUID extracted from the subject claim.
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	// Destination struct where the decoded token claims will be stored.
	claimsStruct := &jwt.RegisteredClaims{}

	// ParseWithClaims parses the raw token string, populates claimsStruct,
	// and invokes the key-Func callback to verify the signing algorithm and return the secret key.
	token, err := jwt.ParseWithClaims(
		tokenString,
		claimsStruct,
		func(token *jwt.Token) (interface{}, error) {
			// CRITICAL SECURITY CHECK: Ensure the token's header specifies HMAC.
			// This prevents attackers from supplying an unsigned token ("alg": "none").
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// Return the secret key as a byte slice to verify the signature.
			return []byte(tokenSecret), nil
		},
	)

	// Return any parsing error (e.g., malformed token, signature mismatch, or expired token).
	if err != nil {
		return uuid.Nil, err
	}

	// Double-check if the library marked the token as valid.
	if !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	// Extract the subject claim ("sub") which contains our user ID string.
	userIDString, err := claimsStruct.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	// Convert the user ID string back into a Go uuid.UUID.
	return uuid.Parse(userIDString)
}
