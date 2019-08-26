package config

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func initLogging() {
	log.SetOutput(os.Stdout)

	if viper.IsSet("logger.level") {
		var level, err = log.ParseLevel(viper.GetString("logger.level"))
		if err == nil {
			log.SetLevel(level)
		} else {
			log.Error("Invalid log level:", err)
			log.SetLevel(log.InfoLevel)
		}
	} else {
		log.SetLevel(log.InfoLevel)
	}
}
