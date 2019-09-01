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
		var (
			profile string
			region  string
		)

		if viper.IsSet("ca.passphrase.profile") {
			profile = viper.GetString("ca.passphrase.profile")
		}
		if viper.IsSet("ca.passphrase.region") {
			region = viper.GetString("ca.passphrase.region")
		}
		kmsClient := awscloud.New(profile).KmsClient(region)
		passphrase, err = kmsClient.KmsDecryptText(encryptedPassphrase)
	default: // symmetric
		passphrase, err = utils.SymmetricDecrypt(encryptedPassphrase)
	}

	return passphrase, err
}

func decryptCAContent(data []byte) ([]byte, error) {
	var (
		decryptedErr     error
		decryptedContent []byte
	)

	// file is not encrypted
	if !viper.IsSet("ca.encrypted") {
		return data, nil
	}

	encryptedContent := string(data)

	switch strings.ToLower(viper.GetString("ca.encrypted.encryption")) {
	case "aws_kms":
		var (
			profile string
			region  string
		)

		if viper.IsSet("ca.encrypted.profile") {
			profile = viper.GetString("ca.encrypted.profile")
		}
		if viper.IsSet("ca.encrypted.region") {
			region = viper.GetString("ca.encrypted.region")
		}
		kmsClient := awscloud.New(profile).KmsClient(region)
		decryptedContent, decryptedErr = kmsClient.KmsDecryptText(encryptedContent)
	default: // symmetric
		decryptedContent, decryptedErr = utils.SymmetricDecrypt(encryptedContent)
	}

	return decryptedContent, decryptedErr
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

	if err != nil {
		return []byte{}, err
	}

	return decryptCAContent(certData)
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
