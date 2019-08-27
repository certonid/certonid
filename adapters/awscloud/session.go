package awscloud

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

// Client store aws info
type Client struct {
	Session *session.Session
}

// New init aws client session
func New() *Client {
	sessionOptions := session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}

	return &Client{
		Session: session.Must(session.NewSessionWithOptions(sessionOptions)),
	}
}
