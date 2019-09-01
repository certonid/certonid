package utils

import (
	"crypto/rand"

	"github.com/sethvargo/go-password/password"
)

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(n int) (string, error) {
	// Generate a password that is n characters long with 10 digits, 20 symbols,
	// allowing upper and lower case letters, disallowing repeat characters.
	return password.Generate(n, 10, 20, false, false)
}
