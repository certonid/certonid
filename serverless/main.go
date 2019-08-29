package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/le0pard/certonid/proto"
	"github.com/le0pard/certonid/serverless/config"
	"github.com/le0pard/certonid/serverless/sshca"
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
		CertType:   event.CertType,
		Key:        event.Key,
		Username:   event.Username,
		Hostnames:  event.Hostnames,
		ValidUntil: event.ValidUntil,
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
