package cmd

import (
	"strings"

	"github.com/le0pard/certonid/adapters/awscloud"
	"github.com/le0pard/certonid/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	encrfileType          string
	encrfileAwsKmsKeyID   string
	encrfileAwsKmsProfile string
	encrfileAwsKmsRegion  string

	encrfileCmd = &cobra.Command{
		Use:   "encrfile [OPTIONS] TEXT",
		Short: "Encrypt file",
		Long:  `Encrypt file with symmetric or kms encryption`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var (
				err     error
				encText string
			)

			switch strings.ToLower(encrfileType) {
			case "aws_kms":
				kmsClient := awscloud.New(encrfileAwsKmsProfile).KmsClient(encrfileAwsKmsRegion)
				encText, err = kmsClient.KmsEncryptText(encrfileAwsKmsKeyID, []byte(args[0]))
			default: // symmetric
				encText, err = utils.SymmetricEncrypt([]byte(args[0]))
			}

			if err != nil {
				er(err)
			}

			log.WithFields(log.Fields{
				"text": encText,
			}).Info("Successfully encrypted")
		},
	}
)

func init() {
	rootCmd.AddCommand(encrfileCmd)
	encrfileCmd.Flags().StringVarP(&encrfileType, "type", "t", "symmetric", "Encryption type (symmetric, aws_kms, gcloud_kms)")
	encrfileCmd.Flags().StringVar(&encrfileAwsKmsKeyID, "aws-kms-key-id", "", "AWS KMS Key ID")
	encrfileCmd.Flags().StringVar(&encrfileAwsKmsProfile, "aws-kms-profile", "", "AWS KMS Profile")
	encrfileCmd.Flags().StringVar(&encrfileAwsKmsRegion, "aws-kms-region", "", "AWS KMS Region")
}
