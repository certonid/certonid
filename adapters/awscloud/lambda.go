package awscloud

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
)

// LambdaClient store aws info
type LambdaClient struct {
	Client *lambda.Lambda
}

// LambdaInvoke allow to call AWS Lambda
func (cl *LambdaClient) LambdaInvoke(funcName string, payload []byte, timeout int) ([]byte, error) {
	ctx, done := context.WithTimeout(context.Background(), timeout*time.Second)
	defer done()

	result, err := cl.Client.InvokeWithContext(ctx, &lambda.InvokeInput{
		FunctionName: aws.String(funcName),
		Payload:      payload,
	})

	if err != nil {
		return []byte{}, fmt.Errorf("Error to invoke AWS Lambda function: %w", err)
	}

	return result.Payload, nil
}

// LambdaClient return AWS Lambda client
func (client *Client) LambdaClient(region string) *LambdaClient {
	awsConfig := aws.Config{}

	if region != "" {
		awsConfig.Region = aws.String(region)
	}

	return &LambdaClient{
		Client: lambda.New(client.Session, &awsConfig),
	}
}
