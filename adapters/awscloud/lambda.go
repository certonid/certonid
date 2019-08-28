package awscloud

import (
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
		return []byte{}, err
	}

	return result.Payload, nil
}

// LambdaClient return kms client
func (client *Client) LambdaClient(region string) *LambdaClient {
	return &LambdaClient{
		Client: lambda.New(client.Session, &aws.Config{Region: aws.String(region)}),
	}
}
