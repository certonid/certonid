package aws

import (
	"encoding/base64"

	"github.com/aws/aws-sdk-go/service/kms"
)

// KmsDecryptText allow to decrypt text from AWS KMS
func (client *Client) KmsDecryptText(text string) ([]byte, error) {
	blob, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return []byte{}, err
	}

	// Create KMS service client
	svc := kms.New(client.Session)

	result, err := svc.Decrypt(&kms.DecryptInput{CiphertextBlob: blob})

	if err != nil {
		return []byte{}, err
	}

	return result.Plaintext, nil
}
