package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/certonid/certonid/adapters/awscloud"
	"github.com/certonid/certonid/proto"
	"github.com/rs/zerolog/log"
)

func genCertFromAws(awsProfile, awsRegion, awsFuncName string, keyData []byte, kmsauthToken string, timeout int) ([]byte, error) {
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

	lambdaClient := awscloud.New(awsProfile).LambdaClient(awsRegion)

	invokePayload, err := lambdaClient.LambdaInvoke(awsFuncName, awsSignRequest, timeout)

	if err != nil {
		return []byte{}, err
	}

	var resp proto.AwsSignResponse

	err = json.Unmarshal(invokePayload, &resp)

	if err != nil {
		return []byte{}, fmt.Errorf("Error to unmarshal data from json: %w", err)
	}

	if len(resp.Cert) == 0 {
		log.Error().
			Str("response", string(invokePayload)).
			Msg("Error to execute serverless function")
		return []byte{}, errors.New("Function not return cert in result")
	}

	return []byte(resp.Cert), nil
}
