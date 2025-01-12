package main

import (
	"context"
	"log"
	awsutils "sapopinguino/internal/aws"
	"sapopinguino/internal/config"
	dbutils "sapopinguino/internal/db"
	openaiutils "sapopinguino/internal/openai"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
    awsutils.ConfigAWS()

	config.ReadConfig(config.ReadConfigOption{})

	openaiutils.ConfigOpenAI()

    dbutils.ConfigDB()
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    log.Println("handle called :D")
	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "\"Hello from Lambda!\"",
	}
	return response, nil
}


func main() {
    lambda.Start(handler)
}
