package cmd

import (
	"fmt"
	"time"

	"github.com/le0pard/certonid/adapters/awscloud"
	"github.com/le0pard/certonid/kmsauth"
	log "github.com/sirupsen/logrus"
)

// GenerateKMSAuthToken return kmsauth token
func GenerateKMSAuthToken() (string, error) {
	validUntil, err := time.ParseDuration(genKMSAuthTokenValidUntil)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"value": genKMSAuthTokenValidUntil,
		}).Error("Invalid KMSAuth ValidUntil value")
		return "", fmt.Errorf("Invalid KMSAuth ValidUntil value: %w", err)
	}

	kmsClient := awscloud.New(genAwsProfile).KmsClient(genAwsRegion)

	kmsauthContext := &kmsauth.AuthContextV2{
		From:     genUsername,
		To:       genKMSAuthServiceID,
		UserType: "user",
	}

	tg := kmsauth.NewTokenGenerator(
		genKMSAuthKeyID,
		kmsauth.TokenVersion2,
		validUntil,
		genKMSAuthCachePath,
		kmsauthContext,
		kmsClient,
	)

	encryptedToken, err := tg.GetEncryptedToken()

	if err != nil {
		return "", fmt.Errorf("Error to generate kmsauth token: %w", err)
	}

	return encryptedToken.String(), nil
}
