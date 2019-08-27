package cmd

import (
	"strings"

	"github.com/le0pard/certonid/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	encType string

	encryptTextCmd = &cobra.Command{
		Use:   "encrypttext [text]",
		Short: "Encrypt text",
		Long:  `Encrypt text with symmetric or kms encryption`,
		Run: func(cmd *cobra.Command, args []string) {
			var (
				err     error
				encText string
			)

			if len(args) < 1 {
				er("provide text for encryption")
			}

			switch strings.ToLower(encType) {
			case "aws_kms":
				// empty
			default: // symmetric
				encText, err = utils.SymmetricEncrypt([]byte(args[0]))
			}

			if err != nil {
				er(err)
			}

			log.WithFields(log.Fields{
				"password": encText,
			}).Info("Successfully encrypted")
		},
	}
)

func init() {
	rootCmd.AddCommand(encryptTextCmd)
	encryptTextCmd.Flags().StringVarP(&encType, "type", "t", "symmetric", "Encryption type (symmetric, aws_kms, gcloud_kms)")
	viper.BindPFlag("type", encryptTextCmd.PersistentFlags().Lookup("type"))
	viper.SetDefault("type", "symmetric")
}
