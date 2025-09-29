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

func ConfigOpenAI(c *config.Config) error {

	err := LoadDeveloperPrompt()
	if err != nil {
		return err
	}

	err = LoadJsonSchema()
	if err != nil {
		return err
	}

	openaiClient := openai.NewClient(
		option.WithAPIKey(c.OpenAIKey()),
	)
	OpenAIClient = &openaiClient

	return nil
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
				Text: responses.ResponseTextConfigParam{
					Format: responses.ResponseFormatTextConfigUnionParam{
						OfJSONSchema: &responses.ResponseFormatTextJSONSchemaConfigParam{
							Name:   "sapopinguino_transliteration",
							Strict: openai.Bool(true),
							Schema: JsonSchema,
						},
					},
				},
			},
		)

		isInTokensArray := false
		inQuotation := false
		buildingToken := false
		token := ""

		for stream.Next() {
			data := stream.Current()
			// if data.Delta. {
			aiToken := data.Delta

			if !isInTokensArray {
				if strings.Contains(aiToken, "[") {
					isInTokensArray = true
				}
			} else {
				for _, r := range aiToken {
					if !inQuotation {
						if r == '[' {
							buildingToken = true
						} else if r == ']' && buildingToken == true {
							buildingToken = false
							token += string(r)

							tokenBytes := []byte(token)

							var tokenSlice []string

							err := json.Unmarshal(tokenBytes, &tokenSlice)
							if err != nil {
								tokenStreamChannel <- TokenStreamRes{
									Response: nil,
									Error:    fmt.Errorf("Error while unmarshalling token: %v", err),
								}
							}

							var tokenStruct Token

							if tokenSlice[0] == "word" {
								tokenStruct = Token{
									Type:          tokenSlice[0],
									Input:         tokenSlice[1],
									Transcription: tokenSlice[2],
									Output:        tokenSlice[3],
								}
							} else {
								tokenStruct = Token{
									Type:  tokenSlice[0],
									Value: tokenSlice[1],
								}
							}

							tokenStreamChannel <- TokenStreamRes{
								Response: &tokenStruct,
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
