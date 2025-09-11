package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	aiutils "sapopinguino/internal/ai"
	awsutils "sapopinguino/internal/aws"
	"sapopinguino/internal/config"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/openai/openai-go/v2"
)

func init() {
	awsutils.ConfigAWS()

	config.ReadConfig(config.ReadConfigOption{})

	awsutils.ConfigAWSGateway(&config.C.Websocket.Endpoint)

	aiutils.ConfigOpenAI()
}

func handler(ctx context.Context, event events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	connectionID := event.RequestContext.ConnectionID

	// tokens := []*aiutils.Token{}

	bodyBytes := []byte(event.Body)

	var bodyS awsutils.Body

	err := json.Unmarshal(bodyBytes, &bodyS)
	if err != nil {
		error := fmt.Errorf("Failed to unmarshal request's body: %v", err)
		log.Println(error)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `"Internal server error :/"`,
		}, error
	}

	tokenStreamChannel := aiutils.StreamResponse(ctx, openai.ChatModelGPT5Mini, aiutils.DeveloperPrompt, bodyS.Message)

	for res := range tokenStreamChannel {
		if res.Error != nil {
			log.Printf("Error while streaming LLM's response: %v", res.Error)
			_, err = awsutils.APIGatewayClient.PostToConnection(ctx, &apigatewaymanagementapi.PostToConnectionInput{
				ConnectionId: &connectionID,
				Data:         []byte("<error:/>"),
			})
			if err != nil {
				log.Printf("Error sending error token to client: %v", err)
				awsutils.HandleDeleteConnection(ctx, &connectionID, "sending \"<error:/>\" in PostConnection")
			}
			break
		}

		// tokens = append(tokens, res.Response)

		var jsonData []byte
		jsonData, err = json.Marshal(res.Response)
		if err != nil {
			log.Printf("Error marshaling JSON: %v", err)
			_, err := awsutils.APIGatewayClient.PostToConnection(ctx, &apigatewaymanagementapi.PostToConnectionInput{
				ConnectionId: &connectionID,
				Data:         []byte("<error:/>"),
			})
			if err != nil {
				log.Printf("Error sending error token to client: %v", err)
				awsutils.HandleDeleteConnection(ctx, &connectionID, "sending \"<error:/>\" in PostConnection")
			}
			break
		}

		_, err = awsutils.APIGatewayClient.PostToConnection(ctx, &apigatewaymanagementapi.PostToConnectionInput{
			ConnectionId: &connectionID,
			Data:         jsonData,
		})
		if err != nil {
			if strings.Contains(err.Error(), "410") {
				log.Printf("Client disconnected: %v", err)
				break
			} else {
				log.Printf("Error sending token to client: %v", err)
				awsutils.HandleDeleteConnection(ctx, &connectionID, "sending token in PostConnection")
				break
			}
		}
	}

	if err == nil {
		_, err = awsutils.APIGatewayClient.PostToConnection(ctx, &apigatewaymanagementapi.PostToConnectionInput{
			ConnectionId: &connectionID,
			Data:         []byte("<end:)>"),
		})
		if err != nil {
			log.Printf("Error sending <end:)> thingy to client: %v", err)
			awsutils.HandleDeleteConnection(ctx, &connectionID, "sending \"<end:/>\" in PostConnection")
		}
	}

	// Add data to DB

	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       `"SIIUUUUU! :D"`,
	}

	return response, nil
}

func main() {
	lambda.Start(handler)
}
