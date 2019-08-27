package cmd

import (
	"strings"

	"github.com/le0pard/certonid/adapters/awscloud"
	"github.com/le0pard/certonid/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	decrType string

	decryptTextCmd = &cobra.Command{
		Use:   "decrypttext [text]",
		Short: "Decrypt text",
		Long:  `Decrypt text with symmetric or kms encryption`,
		Run: func(cmd *cobra.Command, args []string) {
			var (
				err  error
				text []byte
			)

			if len(args) < 1 {
				er("provide text for decryption")
			}

			switch strings.ToLower(decrType) {
			case "aws_kms":
				awsClient := awscloud.New("")
				text, err = awsClient.KmsDecryptText(args[0])
			default: // symmetric
				text, err = utils.SymmetricDecrypt(args[0])
			}

			if err != nil {
				er(err)
			}

			log.WithFields(log.Fields{
				"text": string(text),
			}).Info("Successfully decrypted")
		},
	}
)

func init() {
	rootCmd.AddCommand(decryptTextCmd)
	decryptTextCmd.Flags().StringVarP(&decrType, "type", "t", "symmetric", "Decryption type (symmetric, aws_kms, gcloud_kms)")
	viper.BindPFlag("type", decryptTextCmd.PersistentFlags().Lookup("type"))
	viper.SetDefault("type", "symmetric")
}
