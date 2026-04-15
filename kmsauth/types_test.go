package kmsauth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthContextV1(t *testing.T) {
	ctx := &AuthContextV1{
		From: "user1",
		To:   "service1",
	}
	err := ctx.Validate()
	require.NoError(t, err)

	assert.Equal(t, "user1", ctx.GetUsername())
	kmsCtx := ctx.GetKMSContext()
	assert.Equal(t, "user1", kmsCtx["from"])
	assert.Equal(t, "service1", kmsCtx["to"])
}

func TestAuthContextV2(t *testing.T) {
	ctx := &AuthContextV2{
		From:     "user2",
		To:       "service2",
		UserType: "user",
	}
	err := ctx.Validate()
	require.NoError(t, err)

	assert.Equal(t, "2/user/user2", ctx.GetUsername())
	kmsCtx := ctx.GetKMSContext()
	assert.Equal(t, "user2", kmsCtx["from"])
	assert.Equal(t, "service2", kmsCtx["to"])
	assert.Equal(t, "user", kmsCtx["user_type"])
}

func TestTokenIsValid(t *testing.T) {
	lifetime := 1 * time.Hour
	token := NewToken(lifetime)

	// Should be valid right now
	err := token.IsValid(lifetime)
	require.NoError(t, err)

	// Test expiration bounds
	expiredToken := NewToken(lifetime)
	expiredToken.NotAfter = TokenTime{time.Now().Add(-1 * time.Minute)}
	err = expiredToken.IsValid(lifetime)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid time validity")
}
