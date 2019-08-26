package signer

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"

	"github.com/ScaleFT/sshkeys"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

var (
	defaultEntensions = []string{
		"permit-X11-forwarding",
		"permit-agent-forwarding",
		"permit-port-forwarding",
		"permit-pty",
		"permit-user-rc",
	}
)

// KeySigner does the work of signing a ssh public key with the CA key.
type KeySigner struct {
	ca ssh.Signer
}

type SignRequest struct {
	Key        string    `json:"key"`
	Username   string    `json:"username"`
	ValidUntil time.Time `json:"valid_until"`
}

// SignUserKey returns a signed ssh certificate.
func (s *KeySigner) SignUserKey(req *SignRequest) (*ssh.Certificate, error) {
	pubkey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(req.Key))
	if err != nil {
		return nil, err
	}

	maxKeyDuration, err := time.ParseDuration(viper.GetString("certificates.max_valid_until"))
	if err != nil {
		maxKeyDuration = time.Duration(24) * time.Hour
	}
	expires := time.Now().UTC().Add(maxKeyDuration)
	if req.ValidUntil.After(expires) {
		req.ValidUntil = expires
	}
	cert := &ssh.Certificate{
		CertType:        ssh.UserCert,
		Key:             pubkey,
		KeyId:           fmt.Sprintf("%s_%d", req.Username, time.Now().UTC().Unix()),
		ValidAfter:      uint64(time.Now().UTC().Add(-5 * time.Minute).Unix()),
		ValidBefore:     uint64(req.ValidUntil.Unix()),
		ValidPrincipals: []string{req.Username},
	}
	cert.ValidPrincipals = append(cert.ValidPrincipals, viper.GetStringSlice("certificates.additional_principals")...)
	// critical options
	cert.CriticalOptions = make(map[string]string)
	if len(viper.GetStringSlice("certificates.critical_options")) > 0 {
		for _, perm := range viper.GetStringSlice("certificates.critical_options") {
			if strings.Contains(perm, "=") {
				opt := strings.Split(perm, "=")
				cert.CriticalOptions[strings.TrimSpace(opt[0])] = strings.TrimSpace(opt[1])
			} else {
				cert.CriticalOptions[perm] = ""
			}
		}
	}

	// extensions
	cert.Extensions = make(map[string]string)
	extensions := defaultEntensions
	if len(viper.GetStringSlice("certificates.extensions")) > 0 {
		extensions = viper.GetStringSlice("certificates.extensions")
	}

	for _, perm := range extensions {
		cert.Extensions[perm] = ""
	}
	// sign client key
	if err := cert.SignCert(rand.Reader, s.ca); err != nil {
		return nil, err
	}
	log.Info("Issued cert id: ", cert.KeyId, " principals: ", cert.ValidPrincipals, " fp: ", ssh.FingerprintSHA256(pubkey), " valid until ", time.Unix(int64(cert.ValidBefore), 0).UTC())
	return cert, nil
}

func New(pemBytes, passPhrase []byte) (*KeySigner, error) {
	key, err := sshkeys.ParseEncryptedPrivateKey(pemBytes, passPhrase)
	if err != nil {
		return nil, err
	}
	return &KeySigner{
		ca: key,
	}, nil
}
