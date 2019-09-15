package kmsauth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/le0pard/certonid/adapters/awscloud"
	log "github.com/sirupsen/logrus"
)

// TokenValidator validates a token
type TokenValidator struct {
	// An auth context
	AuthContext AuthContext
	// TokenLifetime is the max lifetime we accept tokens to have
	TokenLifetime time.Duration
	// AuthKey the key_arn or alias to use for authentication
	AuthKey string
	// KMSClient for kms encryption
	KMSClient *awscloud.KMSClient
}

// NewTokenValidator returns a new token validator
func NewTokenValidator(
	authKey string,
	authContext AuthContext,
	tokenLifetime time.Duration,
	kmsClient *awscloud.KMSClient,
) *TokenValidator {
	return &TokenValidator{
		AuthKey:       authKey,
		AuthContext:   authContext,
		TokenLifetime: tokenLifetime,
		KMSClient:     kmsClient,
	}
}

// ValidateToken validates a token
func (tv *TokenValidator) ValidateToken(tokenb64 string) error {
	token, err := tv.decryptToken(tokenb64)
	if err != nil {
		return err
	}
	log.Info("Token")
	log.Info(token)
	return token.IsValid(tv.TokenLifetime)
}

// decryptToken decrypts a token
func (tv *TokenValidator) decryptToken(tokenb64 string) (*Token, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(tokenb64)
	if err != nil {
		return nil, err
	}
	plaintext, keyID, err := tv.KMSClient.KmsDecrypt(ciphertext, tv.AuthContext.GetKMSContext())
	if err != nil {
		return nil, err
	}
	if tv.AuthKey != keyID {
		return nil, fmt.Errorf("Invalid KMS key used %s", keyID)
	}
	token := &Token{}
	err = json.Unmarshal(plaintext, token)
	return token, err
}
