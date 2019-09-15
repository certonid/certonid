package kmsauth

import (
	"fmt"
	"strings"
	"time"

	"errors"

	log "github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"
)

// ------------- AuthContext --------------

// AuthContext is a kms encryption context used to ascertain a user's identiy
type AuthContext interface {
	Validate() error
	GetUsername() string
	GetKMSContext() map[string]*string
}

// AuthContextV1 is a kms encryption context used to ascertain a user's identiy
type AuthContextV1 struct {
	From string `json:"from" validate:"required"`
	To   string `json:"to" validate:"required"`
}

// Validate validates
func (ac *AuthContextV1) Validate() error {
	if ac == nil {
		return errors.New("NilAuthContext")
	}
	v := validator.New()
	return v.Struct(ac)
}

// GetUsername returns a username
func (ac *AuthContextV1) GetUsername() string {
	return ac.From
}

// GetKMSContext gets the kms context
func (ac *AuthContextV1) GetKMSContext() map[string]*string {
	return map[string]*string{
		"from": &ac.From,
		"to":   &ac.To,
	}
}

// AuthContextV2 is a kms encryption context used to ascertain a user's identiy
type AuthContextV2 struct {
	From     string `json:"from" validate:"required"`
	To       string `json:"to" validate:"required"`
	UserType string `json:"user_type" validate:"required"`
}

// Validate validates
func (ac *AuthContextV2) Validate() error {
	if ac == nil {
		return errors.New("NilAuthContext")
	}
	v := validator.New()
	return v.Struct(ac)
}

// GetUsername returns a username
func (ac *AuthContextV2) GetUsername() string {
	return fmt.Sprintf("%d/%s/%s", TokenVersion2, ac.UserType, ac.From)
}

// GetKMSContext gets the kms context
func (ac *AuthContextV2) GetKMSContext() map[string]*string {
	context := map[string]*string{
		"from":      &ac.From,
		"to":        &ac.To,
		"user_type": &ac.UserType,
	}

	return context
}

// ------------- Token --------------

// TokenTime is a custom time formatter
type TokenTime struct {
	time.Time
}

// MarshalJSON marshals into json
func (t *TokenTime) MarshalJSON() ([]byte, error) {
	formatted := t.Time.Format(time.RFC3339Nano)
	stamp := fmt.Sprintf("\"%s\"", formatted)
	return []byte(stamp), nil
}

// UnmarshalJSON unmarshals
func (t *TokenTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	// Unmarshal gives us back a string that looks like "<some_time>". Get rid of the quotes
	s = strings.Trim(s, "\"")
	parsed, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return err
	}
	t = &TokenTime{parsed}
	return nil
}

// Token is a kmsauth token
type Token struct {
	NotBefore TokenTime `json:"not_before"`
	NotAfter  TokenTime `json:"not_after"`
}

// IsValid returns an error if token is invalid, nil if valid
func (t *Token) IsValid(tokenLifetime time.Duration) error {
	now := time.Now().UTC()
	delta := t.NotAfter.Sub(t.NotBefore.Time)
	if delta > tokenLifetime {
		return errors.New("Token issued for longer than Tokenlifetime")
	}

	log.WithFields(log.Fields{
		"now":          now,
		"before":       t.NotBefore.Time,
		"before_check": now.Before(t.NotBefore.Time),
		"after":        t.NotAfter.Time,
		"after_check":  now.After(t.NotAfter.Time),
	}).Info("IsValid")

	if now.Before(t.NotBefore.Time) || now.After(t.NotAfter.Time) {
		return errors.New("Invalid time validity for token")
	}
	return nil
}

// NewToken generates a new token
func NewToken(tokenLifetime time.Duration) *Token {
	now := time.Now().UTC()
	// Start the notBefore x time in the past to avoid clock skew
	notBefore := now.Add(-1 * timeSkew)
	// Set the notAfter x time in the future but account for timeSkew
	notAfter := now.Add(tokenLifetime - timeSkew)
	return &Token{
		NotBefore: TokenTime{notBefore},
		NotAfter:  TokenTime{notAfter},
	}
}

// ------------- EncryptedToken --------------

// EncryptedToken is a b64 kms encrypted token
type EncryptedToken string

//  String satisfies the stringer interface
func (e EncryptedToken) String() string {
	return string(e)
}

// ------------- TokenCache --------------

// TokenCache is a cached token, consists of a token and an encryptedToken
type TokenCache struct {
	Token          Token              `json:"token,omitempty"`
	EncryptedToken EncryptedToken     `json:"encrypted_token,omitempty"`
	AuthContext    map[string]*string `json:"auth_context,omitempty"`
}
