package awscloud

import (
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
)

// KMSClient store aws info
type KMSClient struct {
	Client *kms.KMS
}

// KmsEncrypt allow to encrypt data by AWS KMS
func (cl *KMSClient) KmsEncrypt(keyID string, ciphertextBlob []byte, encryptionContext map[string]*string) ([]byte, error) {
	result, err := cl.Client.Encrypt(&kms.EncryptInput{
		KeyId:             aws.String(keyID),
		Plaintext:         ciphertextBlob,
		EncryptionContext: encryptionContext,
	})

	if err != nil {
		return []byte{}, fmt.Errorf("Error in encrypt data by AWS KMS: %w", err)
	}

	return result.CiphertextBlob, nil
}

// KmsEncryptText allow to encrypt text by AWS KMS
func (cl *KMSClient) KmsEncryptText(keyID string, text []byte) (string, error) {
	result, err := cl.Client.Encrypt(&kms.EncryptInput{
		KeyId:     aws.String(keyID),
		Plaintext: text,
	})

	if err != nil {
		return "", fmt.Errorf("Error in encrypt data by AWS KMS: %w", err)
	}

	return base64.StdEncoding.EncodeToString(result.CiphertextBlob), nil
}

// KmsDecrypt allow to decrypt data AWS KMS
func (cl *KMSClient) KmsDecrypt(ciphertextBlob []byte, encryptionContext map[string]*string) ([]byte, string, error) {
	result, err := cl.Client.Decrypt(&kms.DecryptInput{
		CiphertextBlob:    ciphertextBlob,
		EncryptionContext: encryptionContext,
	})

	if err != nil {
		return []byte{}, "", fmt.Errorf("Error in decrypt data by AWS KMS: %w", err)
	}

	return result.Plaintext, *result.KeyId, nil
}

// KmsDecryptText allow to decrypt text by AWS KMS
func (cl *KMSClient) KmsDecryptText(text string) ([]byte, error) {
	blob, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return []byte{}, fmt.Errorf("Error in decode base64 encrypted data: %w", err)
	}

	result, err := cl.Client.Decrypt(&kms.DecryptInput{CiphertextBlob: blob})

	if err != nil {
		return []byte{}, fmt.Errorf("Error in decrypt data by AWS KMS: %w", err)
	}

	return result.Plaintext, nil
}

// Reader interface
func (cl *KMSClient) Reader(p []byte) (n int, err error) {
	input := &kms.GenerateRandomInput{
		NumberOfBytes: aws.Int64(len(p)),
	}

	result, err := cl.Client.GenerateRandom(input)
	if err != nil {
		n = 0
		return
	}

	copy(p, result.Plaintext)
	n = len(p)
	return
}

// KmsClient return kms client
func (client *Client) KmsClient(region string) *KMSClient {
	awsConfig := aws.Config{}

	if region != "" {
		awsConfig.Region = aws.String(region)
	}

	return &KMSClient{
		Client: kms.New(client.Session, &awsConfig),
	}
}
