package cmd

import (
	"strings"

	"github.com/certonid/certonid/adapters/awscloud"
	"github.com/certonid/certonid/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	decrstrType          string
	decrstrAwsKmsRegion  string
	decrstrAwsKmsProfile string

	decrstrCmd = &cobra.Command{
		Use:   "decrstr [OPTIONS] TEXT",
		Short: "Decrypt text",
		Long:  `Decrypt text with symmetric or kms encryption`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var (
				err  error
				text []byte
			)

			switch strings.ToLower(decrstrType) {
			case "aws_kms":
				kmsClient := awscloud.New(decrstrAwsKmsProfile).KmsClient(decrstrAwsKmsRegion)
				text, err = kmsClient.KmsDecryptText(args[0])
			default: // symmetric
				text, err = utils.SymmetricDecrypt(args[0])
			}

			if err != nil {
				er(err)
			}

			log.Info().
				Str("text", string(text)).
				Msg("Successfully decrypted")
		},
	}
)

func init() {
	rootCmd.AddCommand(decrstrCmd)
	decrstrCmd.Flags().StringVarP(&decrstrType, "type", "t", "symmetric", "Decryption type (symmetric, aws_kms, gcloud_kms)")
	decrstrCmd.Flags().StringVar(&decrstrAwsKmsProfile, "aws-kms-profile", "", "AWS KMS Profile")
	decrstrCmd.Flags().StringVar(&decrstrAwsKmsRegion, "aws-kms-region", "", "AWS KMS Region")
}
