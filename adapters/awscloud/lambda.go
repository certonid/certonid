package awscloud

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

// LambdaClient store aws info
type LambdaClient struct {
	Client *lambda.Client
}

// LambdaInvoke allow to call AWS Lambda
func (cl *LambdaClient) LambdaInvoke(funcName string, payload []byte, timeout int) ([]byte, error) {
	ctx, done := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer done()

	result, err := cl.Client.Invoke(ctx, &lambda.InvokeInput{
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
	return &LambdaClient{
		Client: lambda.NewFromConfig(client.Config, func(o *lambda.Options) {
			if region != "" {
				o.Region = region
			}
		}),
	}
}
