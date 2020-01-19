package cmd

import (
	"fmt"
	"time"

	"github.com/certonid/certonid/adapters/awscloud"
	"github.com/certonid/certonid/kmsauth"
	log "github.com/sirupsen/logrus"
)

// GenerateAwsKMSAuthToken return kmsauth token
func GenerateAwsKMSAuthToken(kmsAuthKeyID, kmsAuthServiceID, kmsAuthTokenValidUntil, awsProfile, awsRegion string) (string, error) {
	validUntil, err := time.ParseDuration(kmsAuthTokenValidUntil)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"value": kmsAuthTokenValidUntil,
		}).Error("Invalid KMSAuth ValidUntil value")
		return "", fmt.Errorf("Invalid KMSAuth ValidUntil value: %w", err)
	}

	kmsClient := awscloud.New(awsProfile).KmsClient(awsRegion)

	kmsauthContext := &kmsauth.AuthContextV2{
		From:     genUsername,
		To:       kmsAuthServiceID,
		UserType: "user",
	}

	tg := kmsauth.NewTokenGenerator(
		kmsAuthKeyID,
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
