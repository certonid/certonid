package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/certonid/certonid/proto"
	"github.com/certonid/certonid/serverless/config"
	"github.com/certonid/certonid/serverless/sshca"
)

// init function
func init() {
	config.InitConfig()
}

// LambdaHandler used to handle lambda call
func LambdaHandler(event proto.AwsSignEvent) (proto.AwsSignResponse, error) {
	var (
		err  error
		cert string
	)

	cert, err = sshca.GenerateCetrificate(&sshca.CertificateRequest{
		CertType:     event.CertType,
		Key:          event.Key,
		Username:     event.Username,
		Hostnames:    event.Hostnames,
		ValidUntil:   event.ValidUntil,
		KMSAuthToken: event.KMSAuthToken,
	})
	if err != nil {
		return proto.AwsSignResponse{}, err
	}

	return proto.AwsSignResponse{
		Cert: cert,
	}, nil
}

// main function
func main() {
	lambda.Start(LambdaHandler)
}
