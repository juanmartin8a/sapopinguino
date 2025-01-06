package main

import (
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

func handleRequest() {

}

func main() {
    lambda.Start(handleRequest)
}
