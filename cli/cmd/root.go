package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/certonid/certonid/cli/version"
	"github.com/certonid/certonid/utils"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile      string
	debugLogging bool
)

var rootCmd = &cobra.Command{
	Use:                   "certonid [OPTIONS] COMMAND [ARG...]",
	Short:                 "Certonid is a Serverless SSH Certificate Authority",
	Version:               fmt.Sprintf("%s, date %s, build %s", version.Version, version.BuildTime, version.GitCommit),
	SilenceUsage:          true,
	SilenceErrors:         true,
	DisableFlagsInUseLine: true,
	TraverseChildren:      true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return showHelp(cmd, args)
		}

		suggestions := cmd.SuggestionsFor(args[0])
		if len(suggestions) > 0 {
			return fmt.Errorf("ERROR: unknown command '%s' for '%s'. Did you mean '%s' ?", args[0], cmd.CalledAs(), strings.Join(suggestions, ", "))
		}
		return fmt.Errorf("ERROR: unknown command '%s' for '%s'", args[0], cmd.CalledAs())
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.SetVersionTemplate("Certonid version {{.Version}}\n")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.certonid.yml)")
	rootCmd.PersistentFlags().BoolVar(&debugLogging, "debug", false, "use debug logging")
}

func initLogging() {
	if viper.IsSet("logger.format") && strings.ToLower(viper.GetString("logger.format")) == "json" {
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debugLogging {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		return
	}

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

func initConfig() {
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix(utils.EnvPrefix)
	viper.SetEnvKeyReplacer(utils.EnvStrReplacer)

	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName(".certonid")
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Error().
				Err(err).
				Msg("Fatal error in config file")
			os.Exit(1)
		}

	}

	initLogging()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error().
			Err(err).
			Msg("Execute error")
		os.Exit(1)
	}
}
