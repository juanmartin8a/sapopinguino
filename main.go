package main

import (
	"context"
	"log"
	aiutils "sapopinguino/internal/ai"
	awsutils "sapopinguino/internal/aws"
	"sapopinguino/internal/config"
	dbutils "sapopinguino/internal/db"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sashabaranov/go-openai"
)

func init() {
    awsutils.ConfigAWS()

	config.ReadConfig(config.ReadConfigOption{})

	aiutils.ConfigOpenAI()

    dbutils.ConfigDB()
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    // gpt 4o-mini call
    // return res to client
    // store data in DB

    _, err := aiutils.ChatCompletion(ctx, openai.GPT4o, aiutils.SystemRoleContent, `{
        "input_language": "English",
        "target_language": "Spanish",
        "input": "Hi! Hello my friends. abc, easy as do re mi, or as simple as 123, abc 123 baby you and me girl"
    }`)
    if err != nil {
        log.Println(err)
    }
    // log.Println(*res)

    // db
    // separate input by words,
    // separate output by words
    // merge input and output into a map
    // 

    // for i, token := range res.Tokens {
    //     if token.Type == "word" {
    //          
    //     } else {
    //
    //     }
    // }






	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "\"Hello from Lambda!\"",
	}
	return response, nil
}


func main() {
    lambda.Start(handler)
}
