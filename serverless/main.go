package main

import (
	"io/ioutil"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/le0pard/certonid/serverless/config"
	"github.com/le0pard/certonid/serverless/signer"
	log "github.com/sirupsen/logrus"
)

// SignEvent used for function arguments
type SignEvent struct {
	Username string `json:"username"`
	Key      string `json:"key"`
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

// main function
func main() {
	lambda.Start(LambdaHandler)
}
