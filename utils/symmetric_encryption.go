package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
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
		return "", err
	}

	cphr, err := aes.NewCipher(symmetricKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(cphr)
	if err != nil {
		return "", err
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	encryptedData := gcm.Seal(nonce, nonce, data, nil)

	return base64.StdEncoding.EncodeToString(encryptedData), nil
}

// SymmetricDecrypt decrypt value by CERTONID_SYMMETRIC_KEY
func SymmetricDecrypt(val string) ([]byte, error) {
	symmetricKey, err := getSymmetricKey()
	if err != nil {
		return []byte{}, err
	}

	data, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		return []byte{}, err
	}

	cphr, err := aes.NewCipher(symmetricKey)
	if err != nil {
		return []byte{}, err
	}

	gcm, err := cipher.NewGCM(cphr)
	if err != nil {
		return []byte{}, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return []byte{}, err
	}

	nonce, data := data[:nonceSize], data[nonceSize:]

	return gcm.Open(nil, nonce, data, nil)
}
