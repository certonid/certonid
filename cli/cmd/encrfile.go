package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/certonid/certonid/adapters/awscloud"
	"github.com/certonid/certonid/utils"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	encrfileType          string
	encrfileAwsKmsKeyID   string
	encrfileAwsKmsProfile string
	encrfileAwsKmsRegion  string

	encrfileCmd = &cobra.Command{
		Use:   "encrfile [OPTIONS] FILEPATH",
		Short: "Encrypt file",
		Long:  `Encrypt file with symmetric or kms encryption`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var (
				err     error
				encText string
			)

			origFilepath, err := homedir.Expand(args[0])
			if err != nil {
				er(fmt.Errorf("Error to expand home dir: %w", err))
			}

			encrFilepath := fmt.Sprintf("%s.enc", origFilepath)

			fileBytes, err := os.ReadFile(origFilepath)
			if err != nil {
				er(fmt.Errorf("Error to read file %s: %w", origFilepath, err))
			}

			switch strings.ToLower(encrfileType) {
			case "aws_kms":
				awsclient, err := awscloud.New(encrfileAwsKmsProfile)
				if err != nil {
					er(err)
				}

				kmsClient := awsclient.KmsClient(encrfileAwsKmsRegion)
				encText, err = kmsClient.KmsEncryptText(encrfileAwsKmsKeyID, fileBytes)
			default: // symmetric
				encText, err = utils.SymmetricEncrypt(fileBytes)
			}

			if err != nil {
				er(err)
			}

			f, err := os.OpenFile(encrFilepath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
			if err != nil {
					er(fmt.Errorf("Error creating file (it may already exist): %w", err))
			}
			_, err = f.Write([]byte(encText))
			f.Close() // Ensure the file is closed
			if err != nil {
					er(fmt.Errorf("Error writing to file %s: %w", encrFilepath, err))
			}

			log.Info().
				Str("original", origFilepath).
				Str("encrypted", encrFilepath).
				Msg("Successfully encrypted file")
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
