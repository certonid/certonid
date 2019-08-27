package aws

import (
	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

// KmsDecryptText allow to decrypt text from AWS KMS
func KmsDecryptText(text string) ([]byte, error) {
	blob, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return []byte{}, err
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create KMS service client
	svc := kms.New(sess)

	result, err := svc.Decrypt(&kms.DecryptInput{CiphertextBlob: blob})

	if err != nil {
		return []byte{}, err
	}

	return result.Plaintext, nil
}
