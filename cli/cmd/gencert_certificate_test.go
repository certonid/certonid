package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

func generateTestCertificate(t *testing.T, validBefore uint64) []byte {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	require.NoError(t, err)

	cert := &ssh.Certificate{
		CertType:    ssh.UserCert,
		Key:         pub,
		KeyId:       "test-cert",
		ValidBefore: validBefore,
	}

	// Self-sign for testing parsing
	signer, err := ssh.NewSignerFromKey(privateKey)
	require.NoError(t, err)
	err = cert.SignCert(rand.Reader, signer)
	require.NoError(t, err)

	return ssh.MarshalAuthorizedKey(cert)
}

func TestGenParseCertificate(t *testing.T) {
	// Generate a valid certificate bytes
	validBefore := uint64(time.Now().Add(1 * time.Hour).Unix())
	certBytes := generateTestCertificate(t, validBefore)

	cert, err := genParseCertificate(certBytes)
	require.NoError(t, err)
	assert.NotNil(t, cert)
	assert.Equal(t, "test-cert", cert.KeyId)
	assert.Equal(t, validBefore, cert.ValidBefore)

	// Test with invalid bytes
	_, err = genParseCertificate([]byte("invalid ssh key data"))
	assert.Error(t, err)
}

func TestGenIsCertStillFresh(t *testing.T) {
	// Test valid cert (expires in 1 hour)
	validBefore := uint64(time.Now().Add(1 * time.Hour).Unix())
	validCert := &ssh.Certificate{ValidBefore: validBefore}
	assert.True(t, genIsCertStillFresh(validCert))

	// Test expired cert (expired 1 hour ago)
	expiredBefore := uint64(time.Now().Add(-1 * time.Hour).Unix())
	expiredCert := &ssh.Certificate{ValidBefore: expiredBefore}
	assert.False(t, genIsCertStillFresh(expiredCert))

	// Test time skew protection (expires in 2 minutes, which is < timeSkew of 5 min)
	skewBefore := uint64(time.Now().Add(2 * time.Minute).Unix())
	skewCert := &ssh.Certificate{ValidBefore: skewBefore}

	// Should be false because it is within the 5m skew window and deemed "too close to expired"
	assert.False(t, genIsCertStillFresh(skewCert))
}
