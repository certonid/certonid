package cmd

import (
	"github.com/certonid/certonid/utils"

	"github.com/rs/zerolog/log"
)

// FailoverKmsauthSchema used for faiover kmsauth
type FailoverKmsauthSchema struct {
	KeyID      string `mapstructure:"key_id"`
	ServiceID  string `mapstructure:"service_id"`
	Profile    string `mapstructure:"profile"`
	Region     string `mapstructure:"region"`
	ValidUntil string `mapstructure:"valid_until"`
}

// FailoverSchema used for faiover settings
type FailoverSchema struct {
	Profile      string                `mapstructure:"profile"`
	Region       string                `mapstructure:"region"`
	FunctionName string                `mapstructure:"function_name"`
	Kmsauth      FailoverKmsauthSchema `mapstructure:"kmsauth"`
}

func genCertAWSFailover(keyData []byte) ([]byte, error) {
	var (
		certBytes        []byte
		kmsauthToken     string
		err              error
		kmsAuthKeyID     string
		kmsAuthServiceID string
		kmsValidUntil    string
		awsProfile       string
		awsRegion        string
		awsFuncName      string
	)

	for _, failoverSettings := range genFailoverVariants {
		kmsAuthKeyID = failoverSettings.Kmsauth.KeyID
		if len(kmsAuthKeyID) == 0 {
			kmsAuthKeyID = genKMSAuthKeyID
		}
		kmsAuthServiceID = failoverSettings.Kmsauth.ServiceID
		if len(kmsAuthServiceID) == 0 {
			kmsAuthServiceID = genKMSAuthServiceID
		}
		kmsValidUntil = failoverSettings.Kmsauth.ValidUntil
		if len(kmsValidUntil) == 0 {
			kmsValidUntil = genKMSAuthTokenValidUntil
		}
		awsProfile = failoverSettings.Profile
		if len(awsProfile) == 0 {
			awsProfile = genAwsProfile
		}
		awsRegion = failoverSettings.Region
		if len(awsRegion) == 0 {
			awsRegion = genAwsRegion
		}
		awsFuncName = failoverSettings.FunctionName
		if len(awsFuncName) == 0 {
			awsFuncName = genAwsFuncName
		}

		if genCertType != utils.HostCertType && len(kmsAuthKeyID) != 0 && len(kmsAuthServiceID) != 0 {
			kmsauthToken, err = GenerateAwsKMSAuthToken(
				kmsAuthKeyID,
				kmsAuthServiceID,
				kmsValidUntil,
				awsProfile,
				awsRegion,
			)
			if err != nil {
				log.Warn().
					Str("kmsauth_key_id", kmsAuthKeyID).
					Str("kmsauth_service_id", kmsAuthServiceID).
					Str("valid_until", kmsValidUntil).
					Str("aws_profile", awsProfile).
					Str("aws_region", awsRegion).
					Msg("Error to generate kmsauth token on failover")
				continue
			}
		}

		certBytes, err = genCertFromAws(
			awsProfile,
			awsRegion,
			awsFuncName,
			keyData,
			kmsauthToken,
			genCertTimeout,
		)

		if err != nil {
			log.Warn().
				Err(err).
				Str("aws_profile", awsProfile).
				Str("aws_region", awsRegion).
				Msg("Error to generate certificate on failover")
			continue
		} else {
			log.Info().
				Str("aws_profile", awsProfile).
				Str("aws_region", awsRegion).
				Str("aws_function_name", awsFuncName).
				Msg("Failover successfully generate certificate")
			break
		}

	}

	return certBytes, err
}
