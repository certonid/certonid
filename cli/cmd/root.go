package cmd

import (
	"fmt"
	"os"

	"github.com/le0pard/certonid/cli/version"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
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
	Version:               fmt.Sprintf("%s, build %s", version.Version, version.GitCommit),
	SilenceUsage:          true,
	SilenceErrors:         true,
	DisableFlagsInUseLine: true,
	TraverseChildren:      true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return showHelp(cmd, args)
		}
		return fmt.Errorf("certonid: '%s' is not a certonid command.\nSee 'certonid --help'", args[0])
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.SetVersionTemplate("Certonid version {{.Version}}\n")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.certonid.yml)")
	rootCmd.PersistentFlags().BoolVar(&debugLogging, "debug", false, "use debug logging")
}

func initLogging() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)

	if debugLogging {
		log.SetLevel(log.DebugLevel)
		return
	}

	if viper.IsSet("logger.level") {
		var level, err = log.ParseLevel(viper.GetString("logger.level"))
		if err == nil {
			log.SetLevel(level)
		} else {
			log.WithFields(log.Fields{
				"error": err,
				"level": viper.GetString("logger.level"),
			}).Info("Invalid log level")
			log.SetLevel(log.InfoLevel)
		}
		return
	}

	log.SetLevel(log.InfoLevel)
}

func initConfig() {
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("certonid")
	viper.AutomaticEnv()
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
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Fatal error in config file")
			os.Exit(1)
		}

	}

	initLogging()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
