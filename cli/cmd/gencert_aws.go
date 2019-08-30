package cmd

import (
	"encoding/json"
	"errors"

	"github.com/le0pard/certonid/adapters/awscloud"
	"github.com/le0pard/certonid/proto"
	log "github.com/sirupsen/logrus"
)

func genCertFromAws(keyData []byte) ([]byte, error) {
	if len(genAwsLambdaFuncName) == 0 {
		return []byte{}, errors.New("You need to provide AWS Lambda function name")
	}

	awsSignRequest, err := json.Marshal(proto.AwsSignEvent{
		CertType:   genCertType,
		Key:        string(keyData),
		Username:   genUsername,
		Hostnames:  genHostnames,
		ValidUntil: genValidUntil,
	})

	if err != nil {
		return []byte{}, err
	}

	lambdaClient := awscloud.New(genAwsLambdaProfile).LambdaClient(genAwsLambdaRegion)

	invokePayload, err := lambdaClient.LambdaInvoke(genAwsLambdaFuncName, awsSignRequest)

	if err != nil {
		return []byte{}, err
	}

	var resp proto.AwsSignResponse

	err = json.Unmarshal(invokePayload, &resp)

	if err != nil {
		return []byte{}, err
	}

	if len(resp.Cert) == 0 {
		log.WithFields(log.Fields{
			"response": string(invokePayload),
		}).Error("Error to execute serverless function")
		return []byte{}, errors.New("Function not return cert in result")
	}

	return []byte(resp.Cert), err
}
