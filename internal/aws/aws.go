package awsutils

import (
	"context"
	"log"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
)

var (
	// SSMClient *ssm.Client
	// KMSClient *kms.Client
	APIGatewayClient *apigatewaymanagementapi.Client
)

func ConfigAWS() {
	_, err := awsConfig.LoadDefaultConfig(
		context.TODO(),
	)
	if err != nil {
		log.Fatalf("Error while loading the AWS config: %s", err)
	}

	// SSMClient = ssm.NewFromConfig(cfg)
	// KMSClient = kms.NewFromConfig(cfg)
}

func ConfigAWSGateway(websocketsEndpoint *string) {
	cfg, err := awsConfig.LoadDefaultConfig(
		context.TODO(),
	)
	if err != nil {
		log.Fatalf("Error while loading the AWS config: %s", err)
	}

	APIGatewayClient = apigatewaymanagementapi.NewFromConfig(cfg, func(o *apigatewaymanagementapi.Options) {
		o.BaseEndpoint = websocketsEndpoint
	})
}
