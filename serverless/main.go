package main

import (
	"io/ioutil"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/le0pard/certonid/serverless/config"
	"github.com/le0pard/certonid/serverless/signer"
	log "github.com/sirupsen/logrus"
	// "github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/s3"
)

type SignEvent struct {
	Username string `json:"username"`
	Key      string `json:"key"`
}

type SignResponse struct {
	Cert string `json:"cert"`
}

func init() {
	config.InitConfig()
}

func LambdaHandler(event SignEvent) (SignResponse, error) {
	var err error

	certData, err := ioutil.ReadFile("./ca")
	if err != nil {
		log.Error("Error to read certificate:", err)
		return SignResponse{}, err
	}
	certSigner, err := signer.New(certData, []byte("password"))
	if err != nil {
		log.Error("Error to init signer:", err)
		return SignResponse{}, err
	}
	cert, err := certSigner.SignKey(&signer.SignRequest{
		Key:        event.Key,
		Username:   event.Username,
		ValidUntil: time.Now().UTC().Add(12 * time.Hour),
	})
	if err != nil {
		log.Error("Error to SignUserKey:", err)
		return SignResponse{}, err
	}
	return SignResponse{
		Cert: cert,
	}, nil
}

func main() {
	lambda.Start(LambdaHandler)
}
