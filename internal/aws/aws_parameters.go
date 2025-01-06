package awsutils

import (
	"context"
	"encoding/base64"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func GetSecretString(
	name string,
) ([]byte, error) {
	withDecryption := true

	parameterOutput, err := SSMClient.GetParameter(
		context.TODO(),
		&ssm.GetParameterInput{
			Name:           &name,
			WithDecryption: &withDecryption,
		},
	)
	if err != nil {
		log.Println(
			"Error while getting parameter: ",
			err,
		)
		return nil, err
	}

	decoded, err := base64.StdEncoding.DecodeString(
		*parameterOutput.Parameter.Value,
	)
	if err != nil {
		log.Println(
			"Error while decoding parameter output: ",
			err,
		)
		return nil, err
	}

	decrypted, err := KMSClient.Decrypt(
		context.TODO(),
		&kms.DecryptInput{
			CiphertextBlob: decoded,
		},
	)
	if err != nil {
		log.Println(
			"Error while decrypting decoded parameter output: ",
			err,
		)
		return nil, err
	}

	return decrypted.Plaintext, nil
}
