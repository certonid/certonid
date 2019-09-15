package awscloud

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
)

// LambdaClient store aws info
type LambdaClient struct {
	Client *lambda.Lambda
}

// LambdaInvoke allow to encrypt text by AWS KMS
func (cl *LambdaClient) LambdaInvoke(funcName string, payload []byte) ([]byte, error) {
	result, err := cl.Client.Invoke(&lambda.InvokeInput{
		FunctionName: aws.String(funcName),
		Payload:      payload,
	})

	if err != nil {
		return []byte{}, fmt.Errorf("Error to invoke AWS Lambda function: %w", err)
	}

	return result.Payload, nil
}

// LambdaClient return kms client
func (client *Client) LambdaClient(region string) *LambdaClient {
	awsConfig := aws.Config{}

	if region != "" {
		awsConfig.Region = aws.String(region)
	}

	return &LambdaClient{
		Client: lambda.New(client.Session, &awsConfig),
	}
}
