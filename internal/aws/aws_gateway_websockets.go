package awsutils

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
)

func HandleDeleteConnection(ctx context.Context, connectionID *string, afterErrorMessage string) {
    _, err := APIGatewayClient.DeleteConnection(ctx, &apigatewaymanagementapi.DeleteConnectionInput{
        ConnectionId: connectionID,
    })
    if err != nil {
        log.Println("Error while deleting connection after %s failed: %s", afterErrorMessage, err)
    }
}
