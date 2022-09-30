package signer

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/certonid/certonid/adapters/awscloud"
	"github.com/certonid/certonid/utils"
	"github.com/rs/zerolog/log"
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
	timeSkew = time.Duration(5) * time.Minute // to protect against time-skew issues we potentially generate a certificate timeSkew duration
)

type sshAlgorithmSigner struct {
	algorithm string
	signer    ssh.AlgorithmSigner
}

func (s *sshAlgorithmSigner) PublicKey() ssh.PublicKey {
	return s.signer.PublicKey()
}

func (s *sshAlgorithmSigner) Sign(rand io.Reader, data []byte) (*ssh.Signature, error) {
	return s.signer.SignWithAlgorithm(rand, data, s.algorithm)
}

func newAlgorithmSignerFromSigner(signer ssh.Signer, algorithm string) (ssh.Signer, error) {
	algorithmSigner, ok := signer.(ssh.AlgorithmSigner)
	if !ok {
		return nil, errors.New("unable to cast to ssh.AlgorithmSigner")
	}
	s := sshAlgorithmSigner{
		signer:    algorithmSigner,
		algorithm: algorithm,
	}
	return &s, nil
}

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
	if req.CertType == utils.HostCertType {
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

func certRandomReader() io.Reader {
	switch strings.ToLower(viper.GetString("certificates.random_seed.source")) {
	case "aws_kms":
		var (
			profile string
			region  string
		)

		if viper.IsSet("certificates.random_seed.profile") {
			profile = viper.GetString("certificates.random_seed.profile")
		} else if viper.IsSet("ca.passphrase.profile") {
			profile = viper.GetString("ca.passphrase.profile")
		}
		if viper.IsSet("certificates.random_seed.region") {
			region = viper.GetString("certificates.random_seed.region")
		} else if viper.IsSet("ca.passphrase.region") {
			region = viper.GetString("ca.passphrase.region")
		}

		return awscloud.New(profile).KmsClient(region)
	default: // urandom
		return rand.Reader
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
	if req.CertType != utils.HostCertType && req.CertType != utils.UserCertType {
		req.CertType = utils.UserCertType // be sure we have at least user key
	}
	// parse public key
	pubkey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(req.Key))
	if err != nil {
		return nil, fmt.Errorf("Error to parse public key: %w", err)
	}
	// check duration
	maxTTLConfigKey := fmt.Sprintf("certificates.%s.max_valid_until", req.CertType)
	maxKeyDuration, err := time.ParseDuration(viper.GetString(maxTTLConfigKey))
	if err != nil {
		log.Warn().
			Err(err).
			Str("value", viper.GetString(maxTTLConfigKey)).
			Msg("Invalid TTL for cert in config. Switched to TTL 24h")
		maxKeyDuration = time.Duration(24) * time.Hour
	}
	expires := time.Now().UTC().Add(maxKeyDuration)
	if req.ValidUntil.After(expires) || req.ValidUntil.Before(time.Now().UTC()) {
		req.ValidUntil = expires
	}
	// check cert type
	var certType uint32 = ssh.UserCert
	if req.CertType == utils.HostCertType {
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
	// use rsa-sha2-512 for sign keys
	sshAlgorithmSigner, err := newAlgorithmSignerFromSigner(s.ca, ssh.SigAlgoRSASHA2512)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Error to initialize rsa-sha2-512 signer")
		return nil, fmt.Errorf("Error to initialize rsa-sha2-512 signer: %w", err)
	}
	// sign client key
	if err := cert.SignCert(certRandomReader(), sshAlgorithmSigner); err != nil {
		return nil, fmt.Errorf("Error sign public key: %w", err)
	}

	log.Info().
		Str("ID", cert.KeyId).
		Str("principals", strings.Join(cert.ValidPrincipals, ",")).
		Str("fingerprint", ssh.FingerprintSHA256(pubkey)).
		Time("valid_until", time.Unix(int64(cert.ValidBefore), 0).UTC()).
		Msg("Successfully issued cert")

	return cert, nil
}

// SignKey sign user key and return certificate as string
func (s *KeySigner) SignKey(req *SignRequest) (string, error) {
	cert, err := s.signPublicKey(req)
	if err != nil {
		return "", err
	}

	marshaled := ssh.MarshalAuthorizedKey(cert)
	// Strip trailing newline
	return string(marshaled[:len(marshaled)-1]), nil
}

// New unseal CA key by passphrase
func New(pemBytes, passPhrase []byte) (*KeySigner, error) {
	key, err := ssh.ParsePrivateKeyWithPassphrase(pemBytes, passPhrase)
	if err != nil {
		return nil, fmt.Errorf("Error parse private key: %w", err)
	}
	return &KeySigner{
		ca: key,
	}, nil
}
