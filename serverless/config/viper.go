package config

import (
	"fmt"
	"os"

	"github.com/le0pard/certonid/utils"
	"github.com/spf13/viper"
)

// InitConfig initialize config for serverless function
func InitConfig() {
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix(utils.EnvPrefix)
	viper.AutomaticEnv()

	cfgFile, ok := utils.GetENV("CONFIG")
	// Don't forget to read config either from cfgFile or from home directory!
	if ok && cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("certonid")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Fatal error to read config file: %s", err)
		os.Exit(1)
	}
	// init logging system
	initLogging()

	// ca
	viper.SetDefault("ca.storage", "file")
	viper.SetDefault("ca.path", "ca.pem")
	viper.SetDefault("ca.passphrase.encryption", "symmetric")
	viper.SetDefault("ca.random_seed.source", "urandom")
	// user cert
	viper.SetDefault("certificates.user.max_valid_until", "24h")
	viper.SetDefault("certificates.user.additional_principals", []string{})
	viper.SetDefault("certificates.user.critical_options", []string{})
	viper.SetDefault("certificates.user.extensions", []string{})
	// host cert
	viper.SetDefault("certificates.host.max_valid_until", "24h")
	viper.SetDefault("certificates.host.additional_principals", []string{})
	viper.SetDefault("certificates.host.critical_options", []string{})
	viper.SetDefault("certificates.host.extensions", []string{})
}
