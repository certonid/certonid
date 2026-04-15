package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSymmetricEncryptionLifecycle(t *testing.T) {
	// 32-byte key for AES-256
	testKey := "0123456789abcdef0123456789abcdef"
	os.Setenv("CERTONID_SYMMETRIC_KEY", testKey)
	defer os.Unsetenv("CERTONID_SYMMETRIC_KEY")

	plainText := []byte("secret_serverless_ca_data")

	// Test Encrypt
	encryptedData, err := SymmetricEncrypt(plainText)
	assert.NoError(t, err)
	assert.NotEmpty(t, encryptedData)

	// Test Decrypt
	decryptedData, err := SymmetricDecrypt(encryptedData)
	assert.NoError(t, err)
	assert.Equal(t, plainText, decryptedData)
}

func TestSymmetricEncryptionMissingKey(t *testing.T) {
	os.Unsetenv("CERTONID_SYMMETRIC_KEY")

	_, err := SymmetricEncrypt([]byte("data"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Not found symmetric key")
}
