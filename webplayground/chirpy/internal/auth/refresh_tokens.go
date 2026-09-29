package auth

import (
	"crypto/rand"
	"encoding/hex"
)

// Generates random 256-bit hex-encoded string to use as refresh token
func MakeRefreshToken() (string, error) {
	// Generate 32-bytes of random data
	dataByte := make([]byte, 32)
	_, err := rand.Read(dataByte)
	if err != nil {
		return "", err
	}

	// Convert to string and return
	return hex.EncodeToString(dataByte), nil
}
