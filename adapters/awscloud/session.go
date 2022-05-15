package awscloud

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

// Client store aws info
type Client struct {
	Config *aws.Config
}

// New init aws client session
func New(profile string) *Client {
	var (
		cfg *aws.Config
		err error
	)

	if profile != "" {
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithSharedConfigProfile(profile),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(context.TODO())
	}
	// like session.Must in v1
	if err != nil {
		panic(fmt.Sprintf("failed loading config, %v", err))
	}

	return &Client{
		Config: cfg,
	}
}
