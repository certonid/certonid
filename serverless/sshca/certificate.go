package sshca

import (
	"io/ioutil"
	"strings"
	"time"

	"github.com/le0pard/certonid/adapters/awscloud"
	"github.com/le0pard/certonid/serverless/signer"
	"github.com/le0pard/certonid/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// CertificateRequest used for function arguments
type CertificateRequest struct {
	CertType   string `json:"cert_type"`
	Key        string `json:"key"`
	Username   string `json:"username"`
	Hostnames  string `json:"hostnames"`
	ValidUntil string `json:"valid_until"`
}

func getCAPassphrase() ([]byte, error) {
	var (
		err        error
		passphrase []byte
	)

	encryptedPassphrase := viper.GetString("ca.passphrase.content")

	switch strings.ToLower(viper.GetString("ca.passphrase.encryption")) {
	case "aws_kms":
		var region string

		if viper.IsSet("ca.passphrase.region") {
			region = viper.GetString("ca.passphrase.region")
		}
		kmsClient := awscloud.New().KmsClient(region)
		passphrase, err = kmsClient.KmsDecryptText(encryptedPassphrase)
	default: // symmetric
		passphrase, err = utils.SymmetricDecrypt(encryptedPassphrase)
	}

	return passphrase, err
}

func getCAFromStorage() ([]byte, error) {
	var (
		err      error
		certData []byte
	)
	switch strings.ToLower(viper.GetString("ca.storage")) {
	case "aws_s3":
		// empty
	default: // file
		certData, err = ioutil.ReadFile(viper.GetString("ca.path"))
	}

	return certData, err
}

// GenerateCetrificate main function to get user of host cert
func GenerateCetrificate(req *CertificateRequest) (string, error) {
	var (
		err        error
		certData   []byte
		passphrase []byte
		validUntil time.Duration
	)

	validUntil, err = time.ParseDuration(req.ValidUntil)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"value": req.ValidUntil,
		}).Error("Invalid ValidUntil value")
		return "", err
	}

	certData, err = getCAFromStorage()
	if err != nil {
		log.WithFields(log.Fields{
			"error":    err,
			"filepath": viper.GetString("ca.path"),
		}).Error("Error to read CA file")
		return "", err
	}

	passphrase, err = getCAPassphrase()
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"encryption": viper.GetString("ca.passphrase.encryption"),
		}).Error("Error to decrypt passphrase for CA key")

		return "", err
	}

	certSigner, err := signer.New(certData, passphrase)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error to parse CA key")
		return "", err
	}

	cert, err := certSigner.SignKey(&signer.SignRequest{
		CertType:   req.CertType,
		Key:        req.Key,
		Username:   req.Username,
		Hostnames:  req.Hostnames,
		ValidUntil: time.Now().UTC().Add(validUntil),
	})
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error to sign user key")
		return "", err
	}

	return cert, nil
}
