package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var invokeCount = 0
var myObjects []*s3.Object

func init() {
	iniConfig()

	svc := s3.New(session.New())
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String("examplebucket"),
	}
	result, _ := svc.ListObjectsV2(input)
	myObjects = result.Contents
}

func LambdaHandler() (int, error) {
	invokeCount = invokeCount + 1
	log.Info(myObjects)
	return invokeCount, nil
}

func main() {
	lambda.Start(LambdaHandler)
}
