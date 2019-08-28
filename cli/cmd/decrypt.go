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
	decryptType         string
	decryptAwsKmsRegion string

	decryptCmd = &cobra.Command{
		Use:   "decrypt [OPTIONS] TEXT",
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

			switch strings.ToLower(decryptType) {
			case "aws_kms":
				kmsClient := awscloud.New().KmsClient(decryptAwsKmsRegion)
				text, err = kmsClient.KmsDecryptText(args[0])
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
	rootCmd.AddCommand(decryptCmd)
	decryptCmd.Flags().StringVarP(&decryptType, "type", "t", "symmetric", "Decryption type (symmetric, aws_kms, gcloud_kms)")
	viper.BindPFlag("type", decryptCmd.PersistentFlags().Lookup("type"))
	viper.SetDefault("type", "symmetric")
}
