package kmsauth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
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

// readCacheFile tries to fetch a token from the cache without locking.
// The caller is responsible for acquiring the necessary RLock or Lock.
func (tg *TokenGenerator) readCacheFile() (*TokenCache, error) {
	_, err := os.Stat(tg.TokenCacheFile)
	if os.IsNotExist(err) {
		// token cache file does not exist yet
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("Error os.Stat token cache: %w", err)
	}
	cacheBytes, err := os.ReadFile(tg.TokenCacheFile)
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
		log.Debug().Msg("Cached token context invalid")
		return nil, nil
	}

	now := time.Now().UTC()
	// subtract timeSkew to account for clock skew
	notAfter := tokenCache.Token.NotAfter.Add(-1 * timeSkew)
	if now.After(notAfter) { // expired, need new token
		return nil, nil
	}

	return tokenCache, nil
}

// GetEncryptedToken returns the encrypted kmsauth token safely handling concurrency
func (tg *TokenGenerator) GetEncryptedToken(skipCache bool) (*EncryptedToken, error) {
	if !skipCache {
		// First pass: Acquire read lock to check if a valid token already exists
		tg.mutex.RLock()
		cached, err := tg.readCacheFile()
		tg.mutex.RUnlock()

		if err != nil {
			return nil, err
		}
		// If we found a valid cache, return the ALREADY encrypted token!
		if cached != nil && cached.EncryptedToken != "" {
			return &cached.EncryptedToken, nil
		}
	}

	// Wait for the write lock. If multiple concurrent requests arrived,
	// they will queue up here.
	tg.mutex.Lock()
	defer tg.mutex.Unlock()

	// Double-Check: Once we have the lock, check the cache one more time.
	// Another goroutine might have just generated and cached the token while we were waiting!
	if !skipCache {
		cached, err := tg.readCacheFile()
		if err != nil {
			return nil, err
		}
		if cached != nil && cached.EncryptedToken != "" {
			return &cached.EncryptedToken, nil
		}
	}

	// Generate a fresh token
	token := NewToken(tg.TokenLifetime)

	tokenBytes, err := json.Marshal(token)
	if err != nil {
		return nil, err
	}

	// Only 1 goroutine makes the AWS KMS Network call
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
			EncryptedToken: encryptedToken, // Store the encrypted token for future requests
			AuthContext:    tg.AuthContext.GetKMSContext(),
		}

		// Write the cache out securely
		dir := path.Dir(tg.TokenCacheFile)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, err
		}

		data, err := json.Marshal(tokenCache)
		if err != nil {
			return nil, err
		}

		tmpFilename := tg.TokenCacheFile + ".tmp"
		f, err := os.OpenFile(tmpFilename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
		if err != nil {
			return nil, err
		}
		_, err = f.Write(data)
		f.Close()
		if err != nil {
			return nil, err
		}

		if err := os.Rename(tmpFilename, tg.TokenCacheFile); err != nil {
			return nil, err
		}
	}

	return &encryptedToken, nil
}
