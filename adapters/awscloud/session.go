package awscloud

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Client store aws info
type Client struct {
	Session *session.Session
}

// New init aws client session
func New(region string) *Client {
	sessionOptions := session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}

	if region != "" {
		sessionOptions.Config = aws.Config{
			Region: aws.String(region),
		}
	}

	return &Client{
		Session: session.Must(session.NewSessionWithOptions(sessionOptions)),
	}
}
