package config

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func initLogging() {
	if viper.IsSet("logger.format") && strings.ToLower(viper.GetString("logger.format")) == "json" {
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if viper.IsSet("logger.level") {
		var level, err = zerolog.ParseLevel(viper.GetString("logger.level"))
		if err == nil {
			zerolog.SetGlobalLevel(level)
		} else {
			log.Warn().
				Err(err).
				Str("level", viper.GetString("logger.level")).
				Msg("Invalid log level")
		}
		return
	}
}
