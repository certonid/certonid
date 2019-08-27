package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Client store aws info
type Client struct {
	Session aws.Session
}

// New init aws client session
func New(region string) (*Client, error) {
	sessionOptions := session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}

	if region != "" {
		sessionOptions.Config = aws.Config{
			Region: region,
		}
	}

	sess := session.Must(session.NewSessionWithOptions(sessionOptions))

	return &Client{
		Session: sess,
	}, nil
}
