package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomBytes(t *testing.T) {
	bytes, err := GenerateRandomBytes(32)
	assert.NoError(t, err)
	assert.Len(t, bytes, 32)
}

func TestGenerateRandomString(t *testing.T) {
	str, err := GenerateRandomString(32)
	assert.NoError(t, err)
	assert.Len(t, str, 32)
}
