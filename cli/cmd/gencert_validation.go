package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	genCertSufix         = "cert.pub"
	genKmsAuthFileSufix  = "kmsauth.json"
	genDefaultValidUntil = "30m" // default 30 min
)

func genValidateOptions() {
	var (
		keyPrefix    string
		hasConfigKey bool
	)

	if len(genCertCertName) > 0 {
		keyPrefix = fmt.Sprintf("certificates.%s", genCertCertName)
		hasConfigKey = len(keyPrefix) > 0 && viper.IsSet(keyPrefix)
	}

	if len(genPublicKeyPath) == 0 && hasConfigKey {
		genPublicKeyPath = viper.GetString(fmt.Sprintf("%s.public_key_path", keyPrefix))
	}

	resolvedPublicKeyPath, err := homedir.Expand(genPublicKeyPath)
	if err != nil {
		er(fmt.Errorf("Could not expand path: %w", err))
	}

	genPublicKeyPath = resolvedPublicKeyPath

	if len(genCertRunner) == 0 && hasConfigKey {
		genCertRunner = viper.GetString(fmt.Sprintf("%s.runner", keyPrefix))
	}

	if len(genCertPath) == 0 && hasConfigKey {
		genCertPath = viper.GetString(fmt.Sprintf("%s.certificate_path", keyPrefix))

		if len(genCertPath) != 0 {
			genCertPath, err = homedir.Expand(genCertPath)
			if err != nil {
				er(err)
			}
		} else if len(genCertPath) == 0 && viper.IsSet("cache_path") {
			certFilePath, err := homedir.Expand(filepath.Join(viper.GetString("cache_path"), fmt.Sprintf("%s-%s", genCertCertName, genCertSufix)))
			if err != nil {
				er(err)
			}

			genCertPath = certFilePath
		}
	}

	if len(genUsername) == 0 && hasConfigKey {
		genUsername = viper.GetString(fmt.Sprintf("%s.username", keyPrefix))
	}

	if len(genHostnames) == 0 && hasConfigKey {
		genHostnames = viper.GetString(fmt.Sprintf("%s.hostnames", keyPrefix))
	}

	if len(genValidUntil) == 0 && hasConfigKey {
		genValidUntil = viper.GetString(fmt.Sprintf("%s.valid_until", keyPrefix))

		if len(genValidUntil) == 0 {
			genValidUntil = genDefaultValidUntil
		}
	}

	if len(genAddToSSHAgent) == 0 && hasConfigKey {
		genAddToSSHAgent = viper.GetString(fmt.Sprintf("%s.add_to_ssh_agent", keyPrefix))
	}

	if !genSkipCertCache && hasConfigKey {
		genSkipCertCache = viper.GetBool(fmt.Sprintf("%s.skip_cache", keyPrefix))
	}

	// aws
	if len(genAwsProfile) == 0 && hasConfigKey {
		genAwsProfile = viper.GetString(fmt.Sprintf("%s.aws.profile", keyPrefix))
	}
	if len(genAwsRegion) == 0 && hasConfigKey {
		genAwsRegion = viper.GetString(fmt.Sprintf("%s.aws.region", keyPrefix))
	}
	if len(genAwsFuncName) == 0 && hasConfigKey {
		genAwsFuncName = viper.GetString(fmt.Sprintf("%s.aws.function_name", keyPrefix))
	}

	if len(genKMSAuthCachePath) == 0 && hasConfigKey {
		genKMSAuthCachePath = viper.GetString(fmt.Sprintf("%s.aws.kmsauth.cache_path", keyPrefix))

		if len(genKMSAuthCachePath) != 0 {
			genKMSAuthCachePath, err = homedir.Expand(genKMSAuthCachePath)
			if err != nil {
				er(err)
			}
		} else if len(genKMSAuthCachePath) == 0 && viper.IsSet("cache_path") {
			kmsauthCachePath, err := homedir.Expand(filepath.Join(viper.GetString("cache_path"), fmt.Sprintf("%s-%s", genCertCertName, genKmsAuthFileSufix)))
			if err != nil {
				er(err)
			}

			genKMSAuthCachePath = kmsauthCachePath
		}
	}

	if len(genKMSAuthKeyID) == 0 && hasConfigKey {
		genKMSAuthKeyID = viper.GetString(fmt.Sprintf("%s.aws.kmsauth.key_id", keyPrefix))
	}
	if len(genKMSAuthServiceID) == 0 && hasConfigKey {
		genKMSAuthServiceID = viper.GetString(fmt.Sprintf("%s.aws.kmsauth.service_id", keyPrefix))
	}
	if len(genKMSAuthTokenValidUntil) == 0 && hasConfigKey {
		genKMSAuthTokenValidUntil = viper.GetString(fmt.Sprintf("%s.aws.kmsauth.valid_until", keyPrefix))

		if len(genKMSAuthTokenValidUntil) == 0 {
			genKMSAuthTokenValidUntil = genDefaultValidUntil
		}
	}

	if hasConfigKey {
		genFailoverVariants = viper.GetStringMap(fmt.Sprintf("%s.failover", keyPrefix))
	}

	genCertType = strings.ToLower(genCertType)
}
