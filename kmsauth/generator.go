package kmsauth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"sync"
	"time"

	"github.com/certonid/certonid/adapters/awscloud"

	"github.com/rs/zerolog/log"
)

// TokenGenerator generates a token
type TokenGenerator struct {
	// AuthKey the key_arn or alias to use for authentication
	AuthKey string
	// TokenVersion is a kmsauth token version
	TokenVersion TokenVersion
	// The token lifetime
	TokenLifetime time.Duration
	// A file to use as a cache
	TokenCacheFile string
	// An auth context
	AuthContext AuthContext

	// KMSClient for kms encryption
	KMSClient *awscloud.KMSClient
	// rw mutex
	mutex sync.RWMutex
}

// NewTokenGenerator returns a new token generator
func NewTokenGenerator(
	authKey string,
	tokenVersion TokenVersion,
	tokenLifetime time.Duration,
	tokenCacheFile string,
	authContext AuthContext,
	kmsClient *awscloud.KMSClient,
) *TokenGenerator {
	return &TokenGenerator{
		AuthKey:        authKey,
		TokenVersion:   tokenVersion,
		TokenLifetime:  tokenLifetime,
		TokenCacheFile: tokenCacheFile,
		AuthContext:    authContext,
		KMSClient:      kmsClient,
	}
}

// Validate validates the TokenGenerator
func (tg *TokenGenerator) Validate() error {
	if tg == nil {
		return errors.New("Nil token generator")
	}
	return tg.AuthContext.Validate()
}

// getCachedToken tries to fetch a token from the cache
func (tg *TokenGenerator) getCachedToken() (*Token, error) {
	// lock for reading
	tg.mutex.RLock()
	defer tg.mutex.RUnlock()

	_, err := os.Stat(tg.TokenCacheFile)
	if os.IsNotExist(err) {
		// token cache file does not exist
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("Error os.Stat token cache: %w", err)
	}
	cacheBytes, err := ioutil.ReadFile(tg.TokenCacheFile)
	if err != nil {
		return nil, fmt.Errorf("Could not open token cache file: %w", err)
	}

	tokenCache := &TokenCache{}
	err = json.Unmarshal(cacheBytes, tokenCache)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal token cache: %w", err)
	}
	// Compare token cache with current params
	ok := reflect.DeepEqual(tokenCache.AuthContext, tg.AuthContext.GetKMSContext())
	if !ok {
		log.Debug().Msg("Cached token invalid")
		return nil, nil
	}
	now := time.Now().UTC()
	// subtract timeSkew to account for clock skew
	notAfter := tokenCache.Token.NotAfter.Add(-1 * timeSkew)
	if now.After(notAfter) { // expired, need new token
		return nil, nil
	}
	return &tokenCache.Token, nil
}

// cacheToken caches a token
func (tg *TokenGenerator) cacheToken(tokenCache *TokenCache) error {
	// lock for writing
	tg.mutex.Lock()
	defer tg.mutex.Unlock()

	dir := path.Dir(tg.TokenCacheFile)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	data, err := json.Marshal(tokenCache)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(tg.TokenCacheFile, data, 0644)
	return err
}

// getToken gets a token
func (tg *TokenGenerator) getToken(skipCache bool) (*Token, error) {
	if !skipCache {
		token, err := tg.getCachedToken()
		if err != nil {
			return nil, err
		}
		// If we could not find a token then return a new one
		if token != nil {
			return token, err
		}
	}
	return NewToken(tg.TokenLifetime), nil
}

// GetEncryptedToken returns the encrypted kmsauth token
func (tg *TokenGenerator) GetEncryptedToken(skipCache bool) (*EncryptedToken, error) {
	token, err := tg.getToken(skipCache)
	if err != nil {
		return nil, err
	}

	tokenBytes, err := json.Marshal(token)
	if err != nil {
		return nil, err
	}

	encryptedData, err := tg.KMSClient.KmsEncrypt(
		tg.AuthKey,
		tokenBytes,
		tg.AuthContext.GetKMSContext(),
	)

	if err != nil {
		return nil, err
	}

	encryptedToken := EncryptedToken(base64.StdEncoding.EncodeToString(encryptedData))

	if !skipCache {
		tokenCache := &TokenCache{
			Token:          *token,
			EncryptedToken: encryptedToken,
			AuthContext:    tg.AuthContext.GetKMSContext(),
		}
		err = tg.cacheToken(tokenCache)
		if err != nil {
			return nil, err
		}
	}

	return &encryptedToken, err
}
