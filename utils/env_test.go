package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetENV(t *testing.T) {
	os.Setenv("CERTONID_TEST_VAR", "test_value")
	defer os.Unsetenv("CERTONID_TEST_VAR")

	val, ok := GetENV("TEST_VAR")
	assert.True(t, ok)
	assert.Equal(t, "test_value", val)

	valMissing, okMissing := GetENV("MISSING_VAR")
	assert.False(t, okMissing)
	assert.Equal(t, "", valMissing)
}
