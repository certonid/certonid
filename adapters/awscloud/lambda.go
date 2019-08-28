package awscloud

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
)

// LambdaClient store aws info
type LambdaClient struct {
	Client *lambda.Lambda
}

// LambdaClient return kms client
func (client *Client) LambdaClient(region string) *LambdaClient {
	return &LambdaClient{
		Client: lambda.New(client.Session, &aws.Config{Region: aws.String(region)}),
	}
}
