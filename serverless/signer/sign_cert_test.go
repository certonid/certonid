package signer

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	"github.com/certonid/certonid/utils"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

func generateTestCAKey(t *testing.T) []byte {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	privDER := x509.MarshalPKCS1PrivateKey(privateKey)
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}
	return pem.EncodeToMemory(&privBlock)
}

func generateTestUserPubKey(t *testing.T) string {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	require.NoError(t, err)
	return string(ssh.MarshalAuthorizedKey(pub))
}

func TestSignKey_UserCert(t *testing.T) {
	viper.Set("certificates.user.max_valid_until", "24h")
	caPem := generateTestCAKey(t)
	userPubKey := generateTestUserPubKey(t)

	s, err := New(caPem, []byte{})
	require.NoError(t, err) // Stops execution immediately if err != nil

	req := &SignRequest{
		CertType:   utils.UserCertType,
		Key:        userPubKey,
		Username:   "leopard",
		ValidUntil: time.Now().UTC().Add(1 * time.Hour),
	}

	signedCert, err := s.SignKey(req)
	require.NoError(t, err)
	assert.Contains(t, signedCert, "ssh-rsa-cert-v01@openssh.com")
}

func TestSignKey_HostCert(t *testing.T) {
	viper.Set("certificates.host.max_valid_until", "24h")
	caPem := generateTestCAKey(t)
	hostPubKey := generateTestUserPubKey(t)

	s, err := New(caPem, []byte{})
	require.NoError(t, err) // Stops execution immediately if err != nil

	req := &SignRequest{
		CertType:   utils.HostCertType,
		Key:        hostPubKey,
		Hostnames:  "example.com, test.com",
		ValidUntil: time.Now().UTC().Add(1 * time.Hour),
	}

	certStr, err := s.SignKey(req)
	require.NoError(t, err)

	// Validate parsing back out
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(certStr))
	require.NoError(t, err)

	cert, ok := pubKey.(*ssh.Certificate)
	require.True(t, ok)
	assert.Equal(t, uint32(ssh.HostCert), cert.CertType)
	assert.Contains(t, cert.ValidPrincipals, "example.com")
	assert.Contains(t, cert.ValidPrincipals, "test.com")
}
