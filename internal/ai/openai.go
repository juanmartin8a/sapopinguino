package aiutils

import (
	"context"
	"encoding/json"
	"fmt"

	"sapopinguino/internal/config"

	"strings"

	openai "github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"github.com/openai/openai-go/v2/responses"
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

func StreamResponse(context context.Context, model string, developer_prompt string, input string) <-chan TokenStreamRes {

	tokenStreamChannel := make(chan TokenStreamRes)

	go func() {
		stream := OpenAIClient.Responses.NewStreaming(
			context,
			responses.ResponseNewParams{
				Model:        model,
				Instructions: openai.String(developer_prompt),
				Input: responses.ResponseNewParamsInputUnion{
					OfString: openai.String(input),
				},
				Reasoning: openai.ReasoningParam{
					Effort: openai.ReasoningEffortMinimal,
				},
			},
		)

		isInTokensArray := false
		inQuotation := false
		buildingToken := false
		token := ""
		bbq := false // bbq stands for "backslash before quotation"

		for stream.Next() {
			data := stream.Current()
			// if data.Delta. {
			aiToken := data.Delta
			fmt.Println(aiToken)

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
				// }
			}
		}

		if err := stream.Err(); err != nil {
			tokenStreamChannel <- TokenStreamRes{
				Response: nil,
				Error:    fmt.Errorf("Error while or during LLM's response stream : %v", err),
			}
		}

		close(tokenStreamChannel)
	}()

	return tokenStreamChannel
}
