package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/certonid/certonid/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var (
	genAwsProfile             string
	genAwsRegion              string
	genAwsFuncName            string
	genCertCertName           string
	genSkipCertCache          bool
	genCertRunner             string
	genCertTimeout            int
	genCertPath               string
	genCertType               string
	genPublicKeyPath          string
	genUsername               string
	genHostnames              string
	genValidUntil             string
	genAddToSSHAgent          string
	genKMSAuthCachePath       string
	genKMSAuthKeyID           string
	genKMSAuthServiceID       string
	genKMSAuthTokenValidUntil string

	genFailoverVariants []FailoverSchema

	gencertCmd = &cobra.Command{
		Use:   "gencert [OPTIONS] [KEY NAME]",
		Short: "Generate user or host certificate",
		Long:  `Generate user or host sertificate by involke serverless function`,
		Run: func(cmd *cobra.Command, args []string) {
			var (
				kmsauthToken  string
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

			if genCertType == utils.HostCertType && len(genHostnames) == 0 {
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

			log.Info().
				Str("runner", genCertRunner).
				Str("public_key", genPublicKeyPath).
				Str("certificate", genCertPath).
				Msg("Signing public key")

			publicKeyData, err := ioutil.ReadFile(genPublicKeyPath)

			if err != nil {
				er(fmt.Errorf("Error to read public key: %w", err))
			}

			switch strings.ToLower(genCertRunner) {
			case "gcloud":
				// TODO
			default: // aws
				// kmsauth for aws
				if genCertType != utils.HostCertType && len(genKMSAuthKeyID) != 0 && len(genKMSAuthServiceID) != 0 {
					kmsauthToken, err = GenerateAwsKMSAuthToken(
						genKMSAuthKeyID,
						genKMSAuthServiceID,
						genKMSAuthTokenValidUntil,
						genAwsProfile,
						genAwsRegion,
						false,
					)
					if err != nil {
						er(err)
					}
				}

				certBytes, serverlessErr = genCertFromAws(
					genAwsProfile,
					genAwsRegion,
					genAwsFuncName,
					publicKeyData,
					kmsauthToken,
					genCertTimeout,
				)

				if serverlessErr != nil {
					if len(genFailoverVariants) > 0 {
						log.Warn().
							Err(serverlessErr).
							Msg("Error to generate certificate. Switching to failover")

						certBytes, serverlessErr = genCertAWSFailover(publicKeyData)

						if serverlessErr != nil {
							er(serverlessErr)
						}
					} else {
						er(serverlessErr)
					}
				}

			}

			err = genStoreCertAtFile(genCertPath, certBytes)

			if err != nil {
				er(err)
			}

			cert, err := genParseCertificate(certBytes)

			if err != nil {
				er(err)
			}

			log.Info().
				Str("public_key", genPublicKeyPath).
				Str("certificate", genCertPath).
				Time("valid_until", time.Unix(int64(cert.ValidBefore), 0).UTC()).
				Msg("Certificate generated and stored")

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
	gencertCmd.Flags().StringVarP(&genCertRunner, "runner", "r", "", "Serverless runner (aws, gcloud, etc)")
	gencertCmd.Flags().IntVar(&genCertTimeout, "timeout", 0, "Timeout for serverless function (seconds)")
	gencertCmd.Flags().StringVarP(&genCertType, "type", "t", "user", "Certificate type (user, host)")
	gencertCmd.Flags().StringVarP(&genPublicKeyPath, "public-key-path", "p", "", "Path to public file, which will used for certificate")
	gencertCmd.Flags().StringVarP(&genCertPath, "certificate-path", "o", "", "Path to cerrtificate file")
	gencertCmd.Flags().StringVarP(&genUsername, "username", "u", "", "Username for certificate")
	gencertCmd.Flags().BoolVar(&genSkipCertCache, "skip-cache", false, "Skip certificate in cache and run serverless function")
	gencertCmd.Flags().StringVar(&genCertCertName, "key-name", "", "Certificate name")
	gencertCmd.Flags().StringVar(&genHostnames, "hostnames", "", "Hostnames for certificate (use comma as divider)")
	gencertCmd.Flags().StringVar(&genValidUntil, "valid-until", "", "TTL for certificate")
	gencertCmd.Flags().StringVar(&genAddToSSHAgent, "add-to-ssh-agent", "", "Private key path, which will added with certificate to ssh-agent")
	// aws
	gencertCmd.Flags().StringVar(&genAwsProfile, "aws-lambda-profile", "", "AWS Lambda Profile")
	gencertCmd.Flags().StringVar(&genAwsRegion, "aws-lambda-region", "", "AWS Lambda Region")
	gencertCmd.Flags().StringVar(&genAwsFuncName, "aws-lambda-func-name", "", "AWS Lambda Function name")
	// kmsauth
	gencertCmd.Flags().StringVar(&genKMSAuthCachePath, "kmsauth-cache-path", "", "Path to KMSAuth cached token")
	gencertCmd.Flags().StringVar(&genKMSAuthKeyID, "kmsauth-key-id", "", "KMSAuth key ID")
	gencertCmd.Flags().StringVar(&genKMSAuthServiceID, "kmsauth-service-id", "", "KMSAuth service ID")
	gencertCmd.Flags().StringVar(&genKMSAuthTokenValidUntil, "kmsauth-token-ttl", "", "KMSAuth token TTL")
}
