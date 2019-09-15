package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var (
	genAwsLambdaProfile  string
	genAwsLambdaRegion   string
	genAwsLambdaFuncName string
	genCertCertName      string
	genSkipCertCache     bool
	genCertRunner        string
	genCertPath          string
	genCertType          string
	genPublicKeyPath     string
	genUsername          string
	genHostnames         string
	genValidUntil        string
	genAddToSSHAgent     string

	gencertCmd = &cobra.Command{
		Use:   "gencert [OPTIONS] [KEY NAME]",
		Short: "Generate user or host certificate",
		Long:  `Generate user or host sertificate by involke serverless function`,
		Run: func(cmd *cobra.Command, args []string) {
			var (
				certBytes     []byte
				serverlessErr error
			)

			if len(args) > 0 && args[0] != "" {
				genCertCertName = args[0]
			}

			genValidateOptions()

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

			isFresh, cachedCert := genIsCertValidInCache()

			if !genSkipCertCache && isFresh {
				genPostScripts(cachedCert)
				// exit from program
				os.Exit(0)
			}

			log.WithFields(log.Fields{
				"runner":      genCertRunner,
				"public key":  genPublicKeyPath,
				"certificate": genCertPath,
			}).Info("Signing public key")

			publicKeyData, err := ioutil.ReadFile(genPublicKeyPath)

			if err != nil {
				er(fmt.Errorf("Error to read public key: %w", err))
			}

			switch strings.ToLower(genCertRunner) {
			case "gcloud":
				// TODO
			default: // aws
				certBytes, serverlessErr = genCertFromAws(publicKeyData)
			}

			if serverlessErr != nil {
				er(serverlessErr)
			}

			err = genStoreCertAtFile(genCertPath, certBytes)

			if err != nil {
				er(err)
			}

			cert, err := genParseCertificate(certBytes)

			if err != nil {
				er(err)
			}

			log.WithFields(log.Fields{
				"public_key":  genPublicKeyPath,
				"certificate": genCertPath,
				"valid until": time.Unix(int64(cert.ValidBefore), 0).UTC(),
			}).Info("Certificate generated and stored")

			genPostScripts(cert)
		},
	}
)

func genPostScripts(cert *ssh.Certificate) {
	if len(genAddToSSHAgent) > 0 {
		genAddCertToAgent(cert)
	}
}

func init() {
	rootCmd.AddCommand(gencertCmd)
	gencertCmd.Flags().StringVarP(&genCertRunner, "runner", "r", "", "Serverless runner (aws)")
	gencertCmd.Flags().StringVarP(&genCertType, "type", "t", "user", "Certificate type (user, host)")
	gencertCmd.Flags().StringVarP(&genPublicKeyPath, "public-key-path", "p", "", "Path to public file, which will used for certificate")
	gencertCmd.Flags().StringVarP(&genCertPath, "certificate-path", "o", "", "Path to cerrtificate file")
	gencertCmd.Flags().StringVarP(&genUsername, "username", "u", "", "Username for certificate")
	gencertCmd.Flags().StringVar(&genAwsLambdaProfile, "aws-lambda-profile", "", "AWS Lambda Profile")
	gencertCmd.Flags().StringVar(&genAwsLambdaRegion, "aws-lambda-region", "", "AWS Lambda Region")
	gencertCmd.Flags().StringVar(&genAwsLambdaFuncName, "aws-lambda-func-name", "", "AWS Lambda Function name")
	gencertCmd.Flags().BoolVar(&genSkipCertCache, "skip-cache", false, "Skip certificate in cache and run serverless function")
	gencertCmd.Flags().StringVar(&genCertCertName, "key-name", "", "Certificate name")
	gencertCmd.Flags().StringVar(&genHostnames, "hostnames", "", "Hostnames for certificate (use comma as divider)")
	gencertCmd.Flags().StringVar(&genValidUntil, "valid-until", "", "TTL for certificate")
	gencertCmd.Flags().StringVar(&genAddToSSHAgent, "add-to-ssh-agent", "", "Private key path, which will added with certificate to ssh-agent")
}
