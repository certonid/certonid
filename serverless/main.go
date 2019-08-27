package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/le0pard/certonid/serverless/config"
	"github.com/le0pard/certonid/serverless/sshca"
)

// SignEvent used for function arguments
type SignEvent struct {
	CertType  string `json:"cert_type"`
	Key       string `json:"key"`
	Username  string `json:"username"`
	Hostnames string `json:"hostnames"`
}

// SignResponse used for function response
type SignResponse struct {
	Cert string `json:"cert"`
}

// init function
func init() {
	config.InitConfig()
}

// LambdaHandler used to handle lambda call
func LambdaHandler(event SignEvent) (SignResponse, error) {
	var (
		err  error
		cert string
	)

	cert, err = sshca.GenerateCetrificate(&sshca.CertificateRequest{
		CertType:  event.CertType,
		Key:       event.Key,
		Username:  event.Username,
		Hostnames: event.Hostnames,
	})
	if err != nil {
		return SignResponse{}, err
	}

	return SignResponse{
		Cert: cert,
	}, nil
}

// main function
func main() {
	lambda.Start(LambdaHandler)
}
