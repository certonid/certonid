package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/certonid/certonid/adapters/awscloud"
	"github.com/certonid/certonid/utils"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	decrfileType          string
	decrfileAwsKmsRegion  string
	decrfileAwsKmsProfile string

	decrfileCmd = &cobra.Command{
		Use:   "decrfile [OPTIONS] FILEPATH",
		Short: "Decrypt file",
		Long:  `Decrypt file with symmetric or kms encryption`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var (
				err     error
				results []byte
			)

			encrFilepath, err := homedir.Expand(args[0])
			if err != nil {
				er(fmt.Errorf("Error to expand home dir: %w", err))
			}

			origFilepath := fmt.Sprintf("%s.orig", encrFilepath)

			fileBytes, err := ioutil.ReadFile(encrFilepath)
			if err != nil {
				er(fmt.Errorf("Error to read file %s: %w", encrFilepath, err))
			}

			fileContent := string(fileBytes)

			switch strings.ToLower(decrfileType) {
			case "aws_kms":
				kmsClient := awscloud.New(decrfileAwsKmsProfile).KmsClient(decrfileAwsKmsRegion)
				results, err = kmsClient.KmsDecryptText(fileContent)
			default: // symmetric
				results, err = utils.SymmetricDecrypt(fileContent)
			}

			if err != nil {
				er(err)
			}

			err = ioutil.WriteFile(origFilepath, results, 0600)
			if err != nil {
				er(fmt.Errorf("Error to write file %s: %w", origFilepath, err))
			}

			log.WithFields(log.Fields{
				"original":  origFilepath,
				"encrypted": encrFilepath,
			}).Info("Successfully decrypted file")
		},
	}
)

func init() {
	rootCmd.AddCommand(decrfileCmd)
	decrfileCmd.Flags().StringVarP(&decrfileType, "type", "t", "symmetric", "Decryption type (symmetric, aws_kms, gcloud_kms)")
	decrfileCmd.Flags().StringVar(&decrfileAwsKmsProfile, "aws-kms-profile", "", "AWS KMS Profile")
	decrfileCmd.Flags().StringVar(&decrfileAwsKmsRegion, "aws-kms-region", "", "AWS KMS Region")
}
