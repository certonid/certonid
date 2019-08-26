package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// InitConfig initialize config for serverless function
func InitConfig() {
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("certonid")
	viper.AutomaticEnv()
	var cfgFile string = viper.GetString("config")
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
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
	// defaults
	viper.SetDefault("certificates.max_valid_until", "24h")
	viper.SetDefault("certificates.additional_principals", []string{})
	viper.SetDefault("certificates.critical_options", []string{})
	viper.SetDefault("certificates.extensions", []string{})
}
