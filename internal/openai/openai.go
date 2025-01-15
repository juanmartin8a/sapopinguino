package openaiutils

import (
	"context"
	"log"
	"sapopinguino/internal/config"

	openai "github.com/sashabaranov/go-openai"
)

var OpenAIClient *openai.Client

func ConfigOpenAI() {
    OpenAIClient = openai.NewClient(config.C.OpenAI.Key)
}

func ChatCompletion(context context.Context, model string, system_role string, input string) (*string, error) {
    res, err := OpenAIClient.CreateChatCompletion(
		context,
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: system_role,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: input,
				},
			},
		},
    )
    if err != nil {
        log.Printf("Error while creating chat completion with openai at ChatCompletion(): %s", err)
        return nil, err 
    }

    return &res.Choices[0].Message.Content, nil
}
