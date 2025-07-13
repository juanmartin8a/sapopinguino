package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	aiutils "sapopinguino/internal/ai"
	awsutils "sapopinguino/internal/aws"
	"sapopinguino/internal/config"
	dbutils "sapopinguino/internal/db"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/openai/openai-go"
)

func init() {
	awsutils.ConfigAWS()

	config.ReadConfig(config.ReadConfigOption{})

	awsutils.ConfigAWSGateway(&config.C.Websocket.Endpoint)

	aiutils.ConfigOpenAI()

	dbutils.ConfigDB()
}

func handler(ctx context.Context, event events.LambdaFunctionURLRequest) (*events.LambdaFunctionURLStreamingResponse, error) {

	tokens := []*aiutils.Token{}

	bodyBytes := []byte(event.Body)

	var bodyS awsutils.Body

	err := json.Unmarshal(bodyBytes, &bodyS)
	if err != nil {
		error := fmt.Errorf("Failed to unmarshal request's body: %v", err)
		log.Println(error)
		return &events.LambdaFunctionURLStreamingResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
			Body: strings.NewReader("Internal server error :/"),
		}, error
	}

	reader, writer := io.Pipe()

	go func() {
		defer writer.Close()
		tokenStreamChannel := aiutils.ChatCompletion(ctx, openai.ChatModelGPT4_1, aiutils.DeveloperPrompt, bodyS.Message)

		for res := range tokenStreamChannel {

			if res.Error != nil {
				log.Printf("Error while streaming LLM's response: %v", res.Error)

				_, err := writer.Write([]byte("event: error\ndata: {}\n\n"))
				if err != nil {
					log.Printf("Error sending error token to client: %v", err)
				}

				return
			}

			tokens = append(tokens, res.Response)

			var jsonData []byte
			jsonData, err = json.Marshal(res.Response)
			if err != nil {
				log.Printf("Error marshaling JSON: %v", err)

				_, err := writer.Write([]byte("event: error\ndata: {}\n\n"))
				if err != nil {
					log.Printf("Error sending error token to client: %v", err)
				}

				return
			}

			data := fmt.Sprintf("event: token\ndata: &s\n\n", jsonData)
			_, err = writer.Write([]byte(data))
			if err != nil {
				log.Printf("Error sending error token to client: %v", err)
				return
			}
		}

		if err == nil {
			_, err = writer.Write([]byte("event: end\ndata: {}\n\n"))
			if err != nil {
				log.Printf("Error sending error token to client: %v", err)
			}
			return
		}
	}()

	// Add data to DB

	response := events.LambdaFunctionURLStreamingResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type":  "text/event-stream",
			"Cache-Control": "no-cache, no-store, must-revalidate",
		},
		Body: reader,
	}

	return &response, nil
}

func main() {
	lambda.Start(handler)
}
