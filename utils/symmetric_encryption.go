package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

func getSymmetricKey() ([]byte, error) {
	symmetricKeyStr, ok := GetENV("SYMMETRIC_KEY")

	if !ok {
		return []byte{}, errors.New("CERTONID_SYMMETRIC_KEY not set")
	}

	return []byte(symmetricKeyStr), nil
}

// SymmetricEncrypt encrypt value by CERTONID_SYMMETRIC_KEY
func SymmetricEncrypt(data []byte) (string, error) {
	symmetricKey, err := getSymmetricKey()
	if err != nil {
		return "", fmt.Errorf("Not found symmetric key: %w", err)
	}

	cphr, err := aes.NewCipher(symmetricKey)
	if err != nil {
		return "", fmt.Errorf("Error to init aes cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(cphr)
	if err != nil {
		return "", fmt.Errorf("Error to init gcm: %w", err)
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("Error to populate nonce with random sequence: %w", err)
	}

	encryptedData := gcm.Seal(nonce, nonce, data, nil)

	return base64.StdEncoding.EncodeToString(encryptedData), nil
}

// SymmetricDecrypt decrypt value by CERTONID_SYMMETRIC_KEY
func SymmetricDecrypt(val string) ([]byte, error) {
	symmetricKey, err := getSymmetricKey()
	if err != nil {
		return []byte{}, fmt.Errorf("Not found symmetric key: %w", err)
	}

	data, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		return []byte{}, fmt.Errorf("Error to decode base64 encrypted value: %w", err)
	}

	cphr, err := aes.NewCipher(symmetricKey)
	if err != nil {
		return []byte{}, fmt.Errorf("Error to init aes cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(cphr)
	if err != nil {
		return []byte{}, fmt.Errorf("Error to init gcm: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return []byte{}, fmt.Errorf("Invalid nonce size for decrypt: %w", err)
	}

	nonce, data := data[:nonceSize], data[nonceSize:]

	return gcm.Open(nil, nonce, data, nil)
}
