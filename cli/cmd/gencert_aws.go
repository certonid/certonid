package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/le0pard/certonid/adapters/awscloud"
	"github.com/le0pard/certonid/proto"
	log "github.com/sirupsen/logrus"
)

func genCertFromAws(keyData []byte, kmsauthToken string) ([]byte, error) {
	if len(genAwsFuncName) == 0 {
		return []byte{}, errors.New("You need to provide AWS Lambda function name")
	}

	awsSignRequest, err := json.Marshal(proto.AwsSignEvent{
		CertType:     genCertType,
		Key:          string(keyData),
		Username:     genUsername,
		Hostnames:    genHostnames,
		ValidUntil:   genValidUntil,
		KMSAuthToken: kmsauthToken,
	})

	if err != nil {
		return []byte{}, fmt.Errorf("Error to marshal data in json: %w", err)
	}

	lambdaClient := awscloud.New(genAwsProfile).LambdaClient(genAwsRegion)

	invokePayload, err := lambdaClient.LambdaInvoke(genAwsFuncName, awsSignRequest)

	if err != nil {
		return []byte{}, err
	}

	var resp proto.AwsSignResponse

	err = json.Unmarshal(invokePayload, &resp)

	if err != nil {
		return []byte{}, fmt.Errorf("Error to unmarshal data from json: %w", err)
	}

	if len(resp.Cert) == 0 {
		log.WithFields(log.Fields{
			"response": string(invokePayload),
		}).Error("Error to execute serverless function")
		return []byte{}, errors.New("Function not return cert in result")
	}

	return []byte(resp.Cert), nil
}
