package cmd

import (
	"fmt"
	"time"

	"github.com/certonid/certonid/adapters/awscloud"
	"github.com/certonid/certonid/kmsauth"
	"github.com/rs/zerolog/log"
)

// GenerateAwsKMSAuthToken return kmsauth token
func GenerateAwsKMSAuthToken(kmsAuthKeyID, kmsAuthServiceID, kmsAuthTokenValidUntil, awsProfile, awsRegion string, skipCache bool) (string, error) {
	validUntil, err := time.ParseDuration(kmsAuthTokenValidUntil)
	if err != nil {
		log.Error().
			Err(err).
			Str("value", kmsAuthTokenValidUntil).
			Msg("Invalid KMSAuth ValidUntil value")
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

	encryptedToken, err := tg.GetEncryptedToken(skipCache)

	if err != nil {
		return "", fmt.Errorf("Error to generate kmsauth token: %w", err)
	}

	return encryptedToken.String(), nil
}
