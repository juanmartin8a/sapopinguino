package awsutils

import (
	"context"
	"fmt"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
)

var (
	// SSMClient *ssm.Client
	// KMSClient *kms.Client
	APIGatewayClient *apigatewaymanagementapi.Client
)

func ConfigAWS() error {
	_, err := awsConfig.LoadDefaultConfig(
		context.TODO(),
	)
	if err != nil {
		return fmt.Errorf("Error while loading the AWS config: %s", err)
	}

	// SSMClient = ssm.NewFromConfig(cfg)
	// KMSClient = kms.NewFromConfig(cfg)

	return nil
}

func ConfigAWSGateway(websocketsEndpoint *string) error {
	cfg, err := awsConfig.LoadDefaultConfig(
		context.TODO(),
	)
	if err != nil {
		return fmt.Errorf("Error while loading the AWS config: %s", err)
	}

	APIGatewayClient = apigatewaymanagementapi.NewFromConfig(cfg, func(o *apigatewaymanagementapi.Options) {
		o.BaseEndpoint = websocketsEndpoint
	})

	return nil
}
