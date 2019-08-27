package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// InitConfig initialize config for serverless function
func InitConfig() {
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix(EnvPrefix)
	viper.AutomaticEnv()

	cfgFile, ok := GetENV("config")
	// Don't forget to read config either from cfgFile or from home directory!
	if ok && cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("certonid-serverless")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Fatal error to read config file: %s", err)
		os.Exit(1)
	}
	// init logging system
	initLogging()

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
