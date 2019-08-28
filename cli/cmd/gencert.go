package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/le0pard/certonid/adapters/awscloud"
	"github.com/le0pard/certonid/proto"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	genAwsLambdaProfile string
	genAwsLambdaRegion  string
	genCertType         string
	genPublicKey        string
	genUsername         string
	genHostnames        string
	genValidUntil       string

	gencertCmd = &cobra.Command{
		Use:   "gencert [OPTIONS] [KEY NAME]",
		Short: "Generate user or host certificate",
		Long:  `Generate user or host sertificate by involke serverless function`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				er("provide key name")
			}

			publicKeyData, err := ioutil.ReadFile(genPublicKey)

			if err != nil {
				er(err)
			}

			awsSignRequest, err := json.Marshal(proto.AwsSignEvent{
				CertType:   genCertType,
				Key:        string(publicKeyData),
				Username:   genUsername,
				Hostnames:  genHostnames,
				ValidUntil: genValidUntil,
			})

			if err != nil {
				er(err)
			}

			lambdaClient := awscloud.New(genAwsLambdaProfile).LambdaClient(genAwsLambdaRegion)

			invokePayload, err := lambdaClient.LambdaInvoke("BressFunction", awsSignRequest)

			if err != nil {
				er(err)
			}

			var resp proto.AwsSignResponse

			err = json.Unmarshal(invokePayload, &resp)

			if err != nil {
				er(err)
			}

			if len(resp.Cert) == 0 {
				log.WithFields(log.Fields{
					"response": string(invokePayload),
				}).Error("Error to execute serverless function")
				os.Exit(0)
			}

			log.WithFields(log.Fields{
				"body": resp.Cert,
			}).Info("Result")

		},
	}
)

func init() {
	rootCmd.AddCommand(gencertCmd)
	gencertCmd.Flags().StringVarP(&genAwsLambdaProfile, "aws-lambda-profile", "", "", "AWS Lambda Profile")
	gencertCmd.Flags().StringVarP(&genAwsLambdaRegion, "aws-lambda-region", "", "", "AWS Lambda Region")
	gencertCmd.Flags().StringVarP(&genCertType, "type", "t", "user", "Certificate type (user, host)")
	gencertCmd.Flags().StringVarP(&genPublicKey, "public-file", "p", "", "Path to public file, which will used for certificate")
	gencertCmd.MarkFlagRequired("public-file")
	gencertCmd.Flags().StringVarP(&genUsername, "username", "u", "", "Username for certificate")
	gencertCmd.Flags().StringVarP(&genHostnames, "hostnames", "n", "", "Hostnames for certificate (use , for division)")
	gencertCmd.Flags().StringVarP(&genValidUntil, "valid-until", "l", "24h", "TTL for certificate")
}
