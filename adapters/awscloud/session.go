package awscloud

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

// Client store aws info
type Client struct {
	Config aws.Config
}

// New init aws client session
func New(profile string) (*Client, error) {
	var (
		cfg aws.Config
		err error
	)

	if profile != "" {
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithSharedConfigProfile(profile),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(context.TODO())
	}

	if err != nil {
		return nil, fmt.Errorf("failed loading aws config: %w", err)
	}

	return &Client{
		Config: cfg,
	}, nil
}
