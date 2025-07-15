package awsutils

import (
	"context"
	"log"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

var (
	SSMClient *ssm.Client
	KMSClient *kms.Client
)

type Body struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

func ConfigAWS() {
	cfg, err := awsConfig.LoadDefaultConfig(
		context.TODO(),
	)
	if err != nil {
		log.Fatalf("Error while loading the AWS config: %s", err)
	}

	SSMClient = ssm.NewFromConfig(cfg)
	KMSClient = kms.NewFromConfig(cfg)
}

