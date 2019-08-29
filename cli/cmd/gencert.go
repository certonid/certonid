package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/le0pard/certonid/adapters/awscloud"
	"github.com/le0pard/certonid/proto"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	genAwsLambdaProfile  string
	genAwsLambdaRegion   string
	genAwsLambdaFuncName string
	genCertCertName      string
	genCertRunner        string
	genCertPath          string
	genCertType          string
	genPublicKeyPath     string
	genUsername          string
	genHostnames         string
	genValidUntil        string

	gencertCmd = &cobra.Command{
		Use:   "gencert [OPTIONS] [KEY NAME]",
		Short: "Generate user or host certificate",
		Long:  `Generate user or host sertificate by involke serverless function`,
		Run: func(cmd *cobra.Command, args []string) {
			var (
				serverlessErr error
			)

			if len(args) > 0 && args[0] != "" {
				genCertCertName = args[0]
			}

			validateOptions()

			if len(genPublicKeyPath) == 0 {
				er("You need to provide public key for key sign")
			}

			if len(genCertPath) == 0 {
				er("You need to provide certificate path")
			}

			if genCertType == "host" && len(genHostnames) == 0 {
				er("You need to hostnames for certificate")
			} else if len(genUsername) == 0 {
				er("You need to username for certificate")
			}

			publicKeyData, err := ioutil.ReadFile(genPublicKeyPath)

			if err != nil {
				er(err)
			}

			switch strings.ToLower(genCertRunner) {
			case "gcloud":
				// TODO
			default: // aws
				serverlessErr = genCertFromAws(publicKeyData)
			}

			if serverlessErr != nil {
				er(serverlessErr)
			}
		},
	}
)

func validateOptions() {
	var (
		keyPrefix    string
		hasConfigKey bool
	)

	if len(genCertCertName) > 0 {
		keyPrefix = fmt.Sprintf("keys.%s", genCertCertName)
		hasConfigKey = len(keyPrefix) > 0 && viper.IsSet(keyPrefix)
	}

	if len(genPublicKeyPath) == 0 && hasConfigKey {
		genPublicKeyPath = viper.GetString(fmt.Sprintf("%s.public_key_path", keyPrefix))
	}

	resolvedPublicKeyPath, err := homedir.Expand(genPublicKeyPath)
	if err != nil {
		er(err)
	}

	genPublicKeyPath = resolvedPublicKeyPath

	if len(genCertRunner) == 0 && hasConfigKey {
		genCertRunner = viper.GetString(fmt.Sprintf("%s.runner", keyPrefix))
	}

	if len(genCertPath) == 0 && hasConfigKey {
		genCertPath = viper.GetString(fmt.Sprintf("%s.certificate_path", keyPrefix))

		if len(genCertPath) == 0 && viper.IsSet("cache_keys_path") {
			certFilename := []string{
				genCertCertName,
				".cert",
			}

			certFilePath, err := homedir.Expand(filepath.Join(viper.GetString("cache_keys_path"), strings.Join(certFilename, "")))
			if err != nil {
				er(err)
			}

			genCertPath = certFilePath
		}
	}

	if len(genUsername) == 0 && hasConfigKey {
		genUsername = viper.GetString(fmt.Sprintf("%s.username", keyPrefix))
	}

	if len(genHostnames) == 0 && hasConfigKey {
		genHostnames = viper.GetString(fmt.Sprintf("%s.hostnames", keyPrefix))
	}

	if len(genValidUntil) == 0 && hasConfigKey {
		genValidUntil = viper.GetString(fmt.Sprintf("%s.valid_until", keyPrefix))
	}

	// aws
	if len(genAwsLambdaProfile) == 0 && hasConfigKey {
		genAwsLambdaProfile = viper.GetString(fmt.Sprintf("%s.aws.profile", keyPrefix))
	}
	if len(genAwsLambdaRegion) == 0 && hasConfigKey {
		genAwsLambdaRegion = viper.GetString(fmt.Sprintf("%s.aws.region", keyPrefix))
	}
	if len(genAwsLambdaFuncName) == 0 && hasConfigKey {
		genAwsLambdaFuncName = viper.GetString(fmt.Sprintf("%s.aws.function_name", keyPrefix))
	}

	genCertType = strings.ToLower(genCertType)
}

func storeCertAtFile(filename, cert string) error {
	os.MkdirAll(filepath.Dir(filename), os.ModePerm)
	return ioutil.WriteFile(filename, []byte(cert), 0600)
}

func genCertFromAws(keyData []byte) error {
	if len(genAwsLambdaFuncName) == 0 {
		return errors.New("You need to provide AWS Lambda function name")
	}

	awsSignRequest, err := json.Marshal(proto.AwsSignEvent{
		CertType:   genCertType,
		Key:        string(keyData),
		Username:   genUsername,
		Hostnames:  genHostnames,
		ValidUntil: genValidUntil,
	})

	if err != nil {
		return err
	}

	lambdaClient := awscloud.New(genAwsLambdaProfile).LambdaClient(genAwsLambdaRegion)

	invokePayload, err := lambdaClient.LambdaInvoke(genAwsLambdaFuncName, awsSignRequest)

	if err != nil {
		return err
	}

	var resp proto.AwsSignResponse

	err = json.Unmarshal(invokePayload, &resp)

	if err != nil {
		return err
	}

	if len(resp.Cert) == 0 {
		log.WithFields(log.Fields{
			"response": string(invokePayload),
		}).Error("Error to execute serverless function")
		return errors.New("Function not return cert in result")
	}

	err = storeCertAtFile(genCertPath, resp.Cert)

	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"public_key":  genPublicKeyPath,
		"certificate": genCertPath,
	}).Info("Certificate generated and stored")

	return nil
}

func init() {
	rootCmd.AddCommand(gencertCmd)
	gencertCmd.Flags().StringVar(&genAwsLambdaProfile, "aws-lambda-profile", "", "AWS Lambda Profile")
	gencertCmd.Flags().StringVar(&genAwsLambdaRegion, "aws-lambda-region", "", "AWS Lambda Region")
	gencertCmd.Flags().StringVar(&genAwsLambdaFuncName, "aws-lambda-func-name", "", "AWS Lambda Function name")
	gencertCmd.Flags().StringVarP(&genCertCertName, "key-name", "n", "", "Certificate name")
	gencertCmd.Flags().StringVarP(&genCertRunner, "runner", "r", "", "Serverless runner (aws)")
	gencertCmd.Flags().StringVarP(&genCertType, "type", "t", "user", "Certificate type (user, host)")
	gencertCmd.Flags().StringVarP(&genPublicKeyPath, "public-key-path", "p", "", "Path to public file, which will used for certificate")
	gencertCmd.Flags().StringVarP(&genCertPath, "certificate-path", "o", "", "Path to cerrtificate file, which will saved output")
	gencertCmd.Flags().StringVarP(&genUsername, "username", "u", "", "Username for certificate")
	gencertCmd.Flags().StringVarP(&genHostnames, "hostnames", "s", "", "Hostnames for certificate (use , for division)")
	gencertCmd.Flags().StringVarP(&genValidUntil, "valid-until", "l", "24h", "TTL for certificate")
}
