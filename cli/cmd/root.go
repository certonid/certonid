package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "certonid",
	Short: "Certonid is a Serverless SSH Certificate Authority",
	Long:  `Serverless SSH Certificate Authority`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.SetVersionTemplate("Docker version {{.Version}}\n")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.certonid/config.yml)")
}

func initLogging() {
	log.SetOutput(os.Stdout)

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
	} else {
		log.SetLevel(log.InfoLevel)
	}
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

		viper.AddConfigPath(filepath.Join(home, ".certonid"))
		viper.AddConfigPath(".certonid")
		viper.SetConfigName("config")
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Config not found. Continue without it")
		} else {
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
