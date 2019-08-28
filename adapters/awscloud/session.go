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
func New(profile string) *Client {
	sessionOptions := session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}

	if profile != "" {
		sessionOptions.Profile = aws.String(profile)
	}

	return &Client{
		Session: session.Must(session.NewSessionWithOptions(sessionOptions)),
	}
}
