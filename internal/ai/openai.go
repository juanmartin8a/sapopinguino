package aiutils

import (
	"context"
	"encoding/json"
	"fmt"

	"sapopinguino/internal/config"

	"strings"

	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

var OpenAIClient *openai.Client

type TokenStreamRes struct {
	Response *Token
	Error    error
}

func ConfigOpenAI() {

	LoadMarkdown()

	openaiClient := openai.NewClient(
		option.WithAPIKey(config.C.OpenAI.Key),
	)
	OpenAIClient = &openaiClient
}

func ChatCompletion(context context.Context, model string, developer_prompt string, input string) <-chan TokenStreamRes {

	tokenStreamChannel := make(chan TokenStreamRes)

	go func() {
		stream := OpenAIClient.Chat.Completions.NewStreaming(
			context,
			openai.ChatCompletionNewParams{
				Model: model,
				Messages: []openai.ChatCompletionMessageParamUnion{
					openai.DeveloperMessage(developer_prompt),
					openai.UserMessage(input),
				},
			},
		)

		isInTokensArray := false
		inQuotation := false
		buildingToken := false
		token := ""
		bbq := false // bbq stands for "backslash before quotation"

		for stream.Next() {
			chunk := stream.Current()

			if len(chunk.Choices) > 0 {
				aiToken := chunk.Choices[0].Delta.Content

				if !isInTokensArray {
					if strings.Contains(aiToken, "[") {
						isInTokensArray = true
					}
				} else {
					for _, r := range aiToken {
						if r == '"' && !bbq {
							inQuotation = !inQuotation
						}
						if r == '\\' {
							bbq = true
						} else {
							if bbq == true {
								bbq = false
							}
						}
						if !inQuotation {
							if r == '{' {
								buildingToken = true
							} else if r == '}' && buildingToken == true {
								buildingToken = false
								token += string(r)

								tokenBytes := []byte(token)

								var tokenS Token

								err := json.Unmarshal(tokenBytes, &tokenS)
								if err != nil {
									tokenStreamChannel <- TokenStreamRes{
										Response: nil,
										Error:    fmt.Errorf("Error while unmarshalling token: %v", err),
									}
								}

								tokenStreamChannel <- TokenStreamRes{
									Response: &tokenS,
									Error:    nil,
								}

								token = ""
							}
						}
						if buildingToken {
							token += string(r)
						}
					}
				}
			}
		}

		if err := stream.Err(); err != nil {
			tokenStreamChannel <- TokenStreamRes{
				Response: nil,
				Error:    fmt.Errorf("Error while or during LLM's response strean : %v", err),
			}
		}

		close(tokenStreamChannel)
	}()

	return tokenStreamChannel
}
