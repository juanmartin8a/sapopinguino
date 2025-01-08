package main

import (
	"context"
	"encoding/json"
	"log"
	awsutils "sapopinguino/internal/aws"
	"sapopinguino/internal/config"
	dbutils "sapopinguino/internal/db"
	openaiutils "sapopinguino/internal/openai"

	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	config.ReadConfig(config.ReadConfigOption{})

    awsutils.ConfigAWS()

	openaiutils.ConfigOpenAI()

    dbutils.ConfigDB()
}

func handleRequest(ctx context.Context, event json.RawMessage) error {
    log.Println("handle called :D")
    return nil
}

func main() {
    lambda.Start(handleRequest)
}
