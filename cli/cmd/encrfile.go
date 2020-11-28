package cmd

import (
	"fmt"
	"io/ioutil"
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

			fileBytes, err := ioutil.ReadFile(origFilepath)
			if err != nil {
				er(fmt.Errorf("Error to read file %s: %w", origFilepath, err))
			}

			switch strings.ToLower(encrfileType) {
			case "aws_kms":
				kmsClient := awscloud.New(encrfileAwsKmsProfile).KmsClient(encrfileAwsKmsRegion)
				encText, err = kmsClient.KmsEncryptText(encrfileAwsKmsKeyID, fileBytes)
			default: // symmetric
				encText, err = utils.SymmetricEncrypt(fileBytes)
			}

			if err != nil {
				er(err)
			}

			err = ioutil.WriteFile(encrFilepath, []byte(encText), 0600)
			if err != nil {
				er(fmt.Errorf("Error to write file %s: %w", encrFilepath, err))
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
