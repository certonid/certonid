package cmd

import (
	"strings"

	"github.com/certonid/certonid/adapters/awscloud"
	"github.com/certonid/certonid/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	encrstrType          string
	encrstrAwsKmsKeyID   string
	encrstrAwsKmsProfile string
	encrstrAwsKmsRegion  string

	encrstrCmd = &cobra.Command{
		Use:   "encrstr [OPTIONS] TEXT",
		Short: "Encrypt text",
		Long:  `Encrypt text with symmetric or kms encryption`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var (
				err     error
				encText string
			)

			switch strings.ToLower(encrstrType) {
			case "aws_kms":
				kmsClient := awscloud.New(encrstrAwsKmsProfile).KmsClient(encrstrAwsKmsRegion)
				encText, err = kmsClient.KmsEncryptText(encrstrAwsKmsKeyID, []byte(args[0]))
			default: // symmetric
				encText, err = utils.SymmetricEncrypt([]byte(args[0]))
			}

			if err != nil {
				er(err)
			}

			log.Info().
				Str("text", encText).
				Msg("Successfully encrypted")
		},
	}
)

func init() {
	rootCmd.AddCommand(encrstrCmd)
	encrstrCmd.Flags().StringVarP(&encrstrType, "type", "t", "symmetric", "Encryption type (symmetric, aws_kms, gcp_kms)")
	encrstrCmd.Flags().StringVar(&encrstrAwsKmsKeyID, "aws-kms-key-id", "", "AWS KMS Key ID")
	encrstrCmd.Flags().StringVar(&encrstrAwsKmsProfile, "aws-kms-profile", "", "AWS KMS Profile")
	encrstrCmd.Flags().StringVar(&encrstrAwsKmsRegion, "aws-kms-region", "", "AWS KMS Region")
}
