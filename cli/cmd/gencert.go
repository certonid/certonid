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

type lambdaAwsResponseError struct {
	Message string `json:"message"`
}

type lambdaAwsResponseBody struct {
	Result string                 `json:"result"`
	Data   proto.AwsSignResponse  `json:"data"`
	Error  lambdaAwsResponseError `json:"error"`
}

type lambdaAwsResponseHeaders struct {
	ContentType string `json:"Content-Type"`
}

type lambdaAwsResponse struct {
	StatusCode int                      `json:"statusCode"`
	Headers    lambdaAwsResponseHeaders `json:"headers"`
	Body       lambdaAwsResponseBody    `json:"body"`
}

var (
	genAwsLambdaRegion string
	genCertType        string
	genPublicKey       string
	genUsername        string
	genHostnames       string
	genValidUntil      string

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

			lambdaClient := awscloud.New().LambdaClient(genAwsLambdaRegion)

			invokePayload, err := lambdaClient.LambdaInvoke("BressFunction", awsSignRequest)

			if err != nil {
				er(err)
			}

			var resp lambdaAwsResponse

			err = json.Unmarshal(invokePayload, &resp)

			if err != nil {
				er(err)
			}

			if resp.StatusCode != 200 {
				log.WithFields(log.Fields{
					"code": resp.StatusCode,
				}).Error("Error to execute serverless function")
				os.Exit(0)
			}

			if resp.Body.Result == "failure" {
				log.WithFields(log.Fields{
					"body": resp.Body,
				}).Error("Failed to get results")
				os.Exit(0)
			}

			log.WithFields(log.Fields{
				"body": resp.Body.Data,
			}).Info("Result")

		},
	}
)

func init() {
	rootCmd.AddCommand(gencertCmd)
	gencertCmd.Flags().StringVarP(&genAwsLambdaRegion, "aws-kms-region", "", "", "AWS Lambda Region")
	gencertCmd.Flags().StringVarP(&genCertType, "type", "t", "user", "Certificate type (user, host)")
	gencertCmd.Flags().StringVarP(&genPublicKey, "public-file", "p", "", "Path to public file, which will used for certificate")
	gencertCmd.MarkFlagRequired("public-file")
	gencertCmd.Flags().StringVarP(&genUsername, "username", "u", "", "Username for certificate")
	gencertCmd.Flags().StringVarP(&genHostnames, "hostnames", "h", "", "Hostnames for certificate (use , for division)")
	gencertCmd.Flags().StringVarP(&genValidUntil, "valid-until", "l", "24h", "TTL for certificate")
}
