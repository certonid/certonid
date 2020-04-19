package sshca

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/certonid/certonid/adapters/awscloud"
	"github.com/certonid/certonid/kmsauth"
	"github.com/certonid/certonid/serverless/signer"
	"github.com/certonid/certonid/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// CertificateRequest used for function arguments
type CertificateRequest struct {
	CertType     string `json:"cert_type"`
	Key          string `json:"key"`
	Username     string `json:"username"`
	Hostnames    string `json:"hostnames"`
	ValidUntil   string `json:"valid_until"`
	KMSAuthToken string `json:"kmsauth_token"`
}

func getCAPassphrase() ([]byte, error) {
	var (
		err        error
		passphrase []byte
	)

	encryptedPassphrase := viper.GetString("ca.passphrase.content")

	log.WithFields(log.Fields{
		"encryptedPassphrase": encryptedPassphrase,
	}).Debug("Decrypting encrypted passphrase")

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
	case "gcp_kms":
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

	log.WithFields(log.Fields{
		"type":                viper.GetString("ca.passphrase.encryption"),
		"encryptedPassphrase": encryptedPassphrase,
		"error":               err,
	}).Debug("Decrypted encrypted passphrase")

	return passphrase, err
}

func decryptCAContent(data []byte) ([]byte, error) {
	var (
		decryptedErr     error
		decryptedContent []byte
	)

	// file is not encrypted
	if !viper.IsSet("ca.encrypted") {
		log.Debug("CA key is not encrypted")
		return data, nil
	}

	encryptedContent := string(data)

	log.WithFields(log.Fields{
		"encryptedContent": encryptedContent,
	}).Debug("Decrypting encrypted CA key")

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

	log.WithFields(log.Fields{
		"type":             viper.GetString("ca.encrypted.encryption"),
		"encryptedContent": encryptedContent,
		"error":            decryptedErr,
	}).Debug("Decrypted encrypted CA key")

	return decryptedContent, decryptedErr
}

func getCAFromStorage() ([]byte, error) {
	var (
		err      error
		certData []byte
	)

	log.WithFields(log.Fields{
		"type": viper.GetString("ca.storage"),
	}).Debug("Reading CA file content")

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

func validateKMSAuthToken(token, username string) error {
	var (
		region string
	)

	validUntil, err := time.ParseDuration(viper.GetString("kmsauth.max_valid_until"))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"value": validUntil,
		}).Error("Invalid KMSAuth ValidUntil value")
		return fmt.Errorf("Invalid KMSAuth ValidUntil value: %w", err)
	}

	log.WithFields(log.Fields{
		"validUntil": validUntil,
	}).Debug("Validate KMSAuth TTL")

	if viper.IsSet("kmsauth.region") {
		region = viper.GetString("kmsauth.region")
	}

	kmsClient := awscloud.New("").KmsClient(region)

	kmsauthContext := &kmsauth.AuthContextV2{
		From:     username,
		To:       viper.GetString("kmsauth.service_id"),
		UserType: "user",
	}

	log.WithFields(log.Fields{
		"kmsauthContext": kmsauthContext,
	}).Debug("KMSAuth context")

	tv := kmsauth.NewTokenValidator(
		viper.GetString("kmsauth.key_id"),
		kmsauthContext,
		validUntil,
		kmsClient,
	)

	return tv.ValidateToken(token)
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
		return "", fmt.Errorf("Invalid ValidUntil value: %w", err)
	}

	log.WithFields(log.Fields{
		"validUntil": validUntil,
	}).Debug("Get TTL information for certificate")

	certData, err = getCAFromStorage()
	if err != nil {
		log.WithFields(log.Fields{
			"error":    err,
			"filepath": viper.GetString("ca.path"),
		}).Error("Error to read CA file")
		return "", fmt.Errorf("Error to read CA file from storage: %w", err)
	}

	passphrase, err = getCAPassphrase()
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"encryption": viper.GetString("ca.passphrase.encryption"),
		}).Error("Error to decrypt passphrase for CA key")

		return "", fmt.Errorf("Error to decrypt passphrase for CA key: %w", err)
	}

	if req.CertType != utils.HostCertType && viper.IsSet("kmsauth.key_id") && viper.IsSet("kmsauth.service_id") && viper.IsSet("kmsauth.region") {
		if len(req.KMSAuthToken) == 0 {
			return "", fmt.Errorf("Need to provide KMSAuth token to get certificate")
		}
		err = validateKMSAuthToken(req.KMSAuthToken, req.Username)
		if err != nil {
			return "", fmt.Errorf("Error to validate kmsauth token: %w", err)
		}
	}

	certSigner, err := signer.New(certData, passphrase)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error to parse CA key")
		return "", fmt.Errorf("Error to parse CA key: %w", err)
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
		return "", fmt.Errorf("Error to sign user key: %w", err)
	}

	return cert, nil
}
