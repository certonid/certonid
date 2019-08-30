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
		er(err)
	}

	genPublicKeyPath = resolvedPublicKeyPath

	if len(genCertRunner) == 0 && hasConfigKey {
		genCertRunner = viper.GetString(fmt.Sprintf("%s.runner", keyPrefix))
	}

	if len(genCertPath) == 0 && hasConfigKey {
		genCertPath = viper.GetString(fmt.Sprintf("%s.certificate_path", keyPrefix))

		if len(genCertPath) == 0 && viper.IsSet("cache_path") {
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

	if !genSkipCertCache && hasConfigKey {
		genSkipCertCache = viper.GetBool(fmt.Sprintf("%s.skip_cache", keyPrefix))
	}

	// aws
	if len(genAwsLambdaProfile) == 0 && hasConfigKey {
		genAwsLambdaProfile = viper.GetString(fmt.Sprintf("%s.aws.profile", keyPrefix))
	}
	if len(genAwsLambdaRegion) == 0 && hasConfigKey {
		genAwsLambdaRegion = viper.GetString(fmt.Sprintf("%s.aws.region", keyPrefix))
	}
	if len(genAwsLambdaFuncName) == 0 && hasConfigKey {
		genAwsLambdaFuncName = viper.GetString(fmt.Sprintf("%s.aws.function_name", keyPrefix))
	}

	genCertType = strings.ToLower(genCertType)
}
