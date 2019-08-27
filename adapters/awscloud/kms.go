package awscloud

import (
	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
)

// KmsEncryptText allow to encrypt text by AWS KMS
func (client *Client) KmsEncryptText(keyId string, text []byte) (string, error) {
	// Create KMS service client
	svc := kms.New(client.Session)

	result, err := svc.Encrypt(&kms.EncryptInput{
		KeyId:     aws.String(keyId),
		Plaintext: text,
	})

	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(result.CiphertextBlob), nil
}

// KmsDecryptText allow to decrypt text by AWS KMS
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
