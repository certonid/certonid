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

const (
	timeSkew     = 5 * time.Minute // to protect against time-skew issues we potentially generate a certificate timeSkew duration
	userCertType = "user"
	hostCertType = "host"
)

// KeySigner does the work of signing a ssh public key with the CA key.
type KeySigner struct {
	ca ssh.Signer
}

// SignRequest pass information for sign
type SignRequest struct {
	CertType   string    `json:"cert_type"`
	Key        string    `json:"key"`
	Username   string    `json:"username"`
	Hostnames  string    `json:"hostnames"`
	ValidUntil time.Time `json:"valid_until"`
}

func setPrincipals(cert *ssh.Certificate, req *SignRequest) {
	if req.CertType == hostCertType {
		hosts := strings.Split(req.Hostnames, ",")
		for i := range hosts {
			hosts[i] = strings.TrimSpace(hosts[i])
		}
		cert.ValidPrincipals = hosts
	} else {
		cert.ValidPrincipals = []string{req.Username}
	}
	configKey := fmt.Sprintf("certificates.%s.additional_principals", req.CertType)
	cert.ValidPrincipals = append(cert.ValidPrincipals, viper.GetStringSlice(configKey)...)
}

func setCriticalOptions(cert *ssh.Certificate, req *SignRequest) {
	cert.CriticalOptions = make(map[string]string)

	configKey := fmt.Sprintf("certificates.%s.critical_options", req.CertType)
	if len(viper.GetStringSlice(configKey)) > 0 {
		for _, perm := range viper.GetStringSlice(configKey) {
			if strings.Contains(perm, "=") || strings.Contains(perm, " ") {
				var opt []string
				if strings.Contains(perm, "=") {
					opt = strings.Split(perm, "=")
				} else {
					opt = strings.Split(perm, " ")
				}

				cert.CriticalOptions[strings.TrimSpace(opt[0])] = strings.TrimSpace(opt[1])
			}
		}
	}
}

func setExtensions(cert *ssh.Certificate, req *SignRequest) {
	cert.Extensions = make(map[string]string)
	extensions := defaultEntensions

	configKey := fmt.Sprintf("certificates.%s.extensions", req.CertType)
	if len(viper.GetStringSlice(configKey)) > 0 {
		extensions = viper.GetStringSlice(configKey)
	}

	for _, perm := range extensions {
		cert.Extensions[perm] = ""
	}
}

// signPublicKey returns a signed ssh certificate.
func (s *KeySigner) signPublicKey(req *SignRequest) (*ssh.Certificate, error) {
	if req.CertType != hostCertType && req.CertType != userCertType {
		req.CertType = userCertType // be sure we have at least user key
	}
	// parse public key
	pubkey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(req.Key))
	if err != nil {
		return nil, err
	}
	// check duration
	maxTTLConfigKey := fmt.Sprintf("certificates.%s.max_valid_until", req.CertType)
	maxKeyDuration, err := time.ParseDuration(viper.GetString(maxTTLConfigKey))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Invalid TTL for cert in config. Switched to TTL 24h")
		maxKeyDuration = time.Duration(24) * time.Hour
	}
	expires := time.Now().UTC().Add(maxKeyDuration)
	if req.ValidUntil.After(expires) || req.ValidUntil.Before(time.Now().UTC()) {
		req.ValidUntil = expires
	}
	// check cert type
	var certType uint32 = ssh.UserCert
	if req.CertType == hostCertType {
		certType = ssh.HostCert
	}
	// init certificate
	cert := &ssh.Certificate{
		CertType:    certType,
		Key:         pubkey,
		KeyId:       fmt.Sprintf("%s_%d", req.Username, time.Now().UTC().Unix()),
		ValidAfter:  uint64(time.Now().UTC().Add(-1 * timeSkew).Unix()),
		ValidBefore: uint64(req.ValidUntil.Unix()),
	}
	// principals
	setPrincipals(cert, req)
	// critical options
	setCriticalOptions(cert, req)
	// extensions
	setExtensions(cert, req)
	// sign client key
	if err := cert.SignCert(rand.Reader, s.ca); err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"ID":          cert.KeyId,
		"principals":  cert.ValidPrincipals,
		"fingerprint": ssh.FingerprintSHA256(pubkey),
		"valid until": time.Unix(int64(cert.ValidBefore), 0).UTC(),
	}).Info("Successfully issued cert")

	return cert, nil
}

// SignKey sign user key and return certificate as string
func (s *KeySigner) SignKey(req *SignRequest) (string, error) {
	cert, err := s.signPublicKey(req)
	if err != nil {
		return "", err
	}

	marshaled := ssh.MarshalAuthorizedKey(cert)
	return string(marshaled[:len(marshaled)-1]), nil
}

// New unseal CA key by passphrase
func New(pemBytes, passPhrase []byte) (*KeySigner, error) {
	key, err := sshkeys.ParseEncryptedPrivateKey(pemBytes, passPhrase)
	if err != nil {
		return nil, err
	}
	return &KeySigner{
		ca: key,
	}, nil
}
