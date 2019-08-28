package awscloud

import (
	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
)

// KMSClient store aws info
type KMSClient struct {
	Client *kms.KMS
}

// KmsEncryptText allow to encrypt text by AWS KMS
func (cl *KMSClient) KmsEncryptText(keyId string, text []byte) (string, error) {
	result, err := cl.Client.Encrypt(&kms.EncryptInput{
		KeyId:     aws.String(keyId),
		Plaintext: text,
	})

	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(result.CiphertextBlob), nil
}

// KmsDecryptText allow to decrypt text by AWS KMS
func (cl *KMSClient) KmsDecryptText(text string) ([]byte, error) {
	blob, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return []byte{}, err
	}

	result, err := cl.Client.Decrypt(&kms.DecryptInput{CiphertextBlob: blob})

	if err != nil {
		return []byte{}, err
	}

	return result.Plaintext, nil
}

// KmsClient return kms client
func (client *Client) KmsClient(region string) *KMSClient {
	return &KMSClient{
		Client: kms.New(client.Session, &aws.Config{Region: aws.String(region)}),
	}
}
