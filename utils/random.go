package utils

import (
	"crypto/rand"
	"fmt"

	passGenerator "github.com/sethvargo/go-password/password"
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
	gen, err := passGenerator.NewGenerator(&passGenerator.GeneratorInput{
		Symbols: "~@#&*()_+-=?,.",
	})

	if err != nil {
		return "", fmt.Errorf("Error to init random generator: %w", err)
	}
	// Generate a password that is n characters long with 5 digits, 3 symbols,
	// allowing upper and lower case letters, disallowing repeat characters.
	return gen.Generate(n, 5, 3, false, false)
}
