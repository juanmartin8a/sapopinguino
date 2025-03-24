package main

import (
	"context"
	"encoding/json"
	"log"
	aiutils "sapopinguino/internal/ai"
	awsutils "sapopinguino/internal/aws"
	"sapopinguino/internal/config"
	dbutils "sapopinguino/internal/db"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/openai/openai-go"
)

func init() {
    log.Printf("ws endpoint: %s", config.C.Websockets.Endpoint)
    awsutils.ConfigAWS(&config.C.Websockets.Endpoint)

	config.ReadConfig(config.ReadConfigOption{})

	aiutils.ConfigOpenAI()

    dbutils.ConfigDB()
}

func handler(ctx context.Context, event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
    connectionID := event.RequestContext.ConnectionID

    tokens := []*aiutils.Token{}

    log.Println(event.Body)
    log.Println(event.RequestContext.Stage)

    // tokenStreamChannel := aiutils.ChatCompletion(ctx, openai.ChatModelGPT4o, aiutils.SystemRoleContent, event.Body)
    tokenStreamChannel := aiutils.ChatCompletion(ctx, openai.ChatModelGPT4o, aiutils.SystemRoleContent, `{
        "input_language": "English",
        "target_language": "Spanish",
        "input": "abc, easy as do re mi, or as simple as 123, abc 123 baby you and me girl"
    }`)

    for res := range tokenStreamChannel {
        log.Println("hi")
		if res.Error != nil {
			log.Printf("\nError encountered: %v\n", res.Error)
            _, err := awsutils.APIGatewayClient.PostToConnection(ctx, &apigatewaymanagementapi.PostToConnectionInput{
                ConnectionId: &connectionID,
                Data:         []byte("<error:/>"),
            })
            if err != nil {
                log.Println("Error sending error token to client: %s", err)
                awsutils.HandleDeleteConnection(ctx, &connectionID, "sending \"<error:/>\" in PostConnection")
            }
			break
		}

        log.Println(res.Response.Type)

        tokens = append(tokens, res.Response)

        jsonData, err := json.Marshal(res.Response)
	    if err != nil {
		    log.Println("Error marshaling JSON:", err)
            break
	    }

        _, err = awsutils.APIGatewayClient.PostToConnection(ctx, &apigatewaymanagementapi.PostToConnectionInput{
			ConnectionId: &connectionID,
			Data:         jsonData,
		})
        if err != nil {
		    log.Println("Error sending token to client: %s", err)
            awsutils.HandleDeleteConnection(ctx, &connectionID, "sending token in PostConnection")
            break
        }
	}
    _, err := awsutils.APIGatewayClient.PostToConnection(ctx, &apigatewaymanagementapi.PostToConnectionInput{
        ConnectionId: &connectionID,
        Data:         []byte("<end:)>"),
    })
    if err != nil {
        log.Println("Error sending <end:)> thingy to client: %s", err)
        awsutils.HandleDeleteConnection(ctx, &connectionID, "sending \"<end:/>\" in PostConnection")
    }

	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "\"Hello from Lambda!\"",
	}

	return response, nil
}


func main() {
    lambda.Start(handler)
}
