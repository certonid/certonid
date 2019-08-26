package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func iniConfig() {
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("certonid")
	viper.AutomaticEnv()
	var cfgFile string = viper.GetString("config")
	// Don't forget to read config either from cfgFile or from home directory!
	panic(cfgFile)
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("certonid-serverless")
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Error("Can't read config:", err)
		os.Exit(1)
	}
}
