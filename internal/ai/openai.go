package aiutils

import (
	"context"
	"encoding/json"
	"fmt"

	// "encoding/json"
	"log"
	"sapopinguino/internal/config"

	// openai "github.com/sashabaranov/go-openai"
	"strings"

	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

var OpenAIClient *openai.Client

type TokenStreamRes struct {
	Response *Token 
	Error  error
}

func ConfigOpenAI() {
    OpenAIClient = openai.NewClient(
        option.WithAPIKey(config.C.OpenAI.Key),
    )
}

func ChatCompletion(context context.Context, model string, system_role string, input string) <-chan TokenStreamRes {

    tokenStreamChannel := make(chan TokenStreamRes)

    go func() {
        stream := OpenAIClient.Chat.Completions.NewStreaming(
            context,
            openai.ChatCompletionNewParams{
                Model: openai.F(model),
                Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
                    openai.SystemMessage(system_role),
                    openai.UserMessage(input),
                }),
            },
        )

        isInTokensArray := false
        inQuotation := false 
        buildingToken := false
        token := ""
        bbq := false // bbq stands for "backslash before quotation"

        // acc := openai.ChatCompletionAccumulator{}

        for stream.Next() {
            chunk := stream.Current()
            // acc.AddChunk(chunk)

            // if content, ok := acc.JustFinishedContent(); ok {
            //     println("Content stream finished:", content)
            // }
            //
            // // if using tool calls
            // if tool, ok := acc.JustFinishedToolCall(); ok {
            //     println("Tool call stream finished:", tool.Index, tool.Name, tool.Arguments)
            // }
            //
            // if refusal, ok := acc.JustFinishedRefusal(); ok {
            //     println("Refusal stream finished:", refusal)
            // }

            if len(chunk.Choices) > 0 {
                aiToken := chunk.Choices[0].Delta.Content

                if !isInTokensArray {
                    if strings.Contains(aiToken, "[") {
                        isInTokensArray = true; 
                    }
                } else {
                    for _, r := range aiToken {
                        if r == '"' && !bbq {
                            inQuotation = !inQuotation;
                        }
                        if r == '\\' {
                            bbq = true; 
                        } else {
                            if (bbq == true) {
                                bbq = false;
                            }
                        }
                        if !inQuotation {
                            if r == '{' {
                                buildingToken = true   
                            } else if r == '}' {
                                buildingToken = false
                                token += string(r);

                                log.Println("Token: ")
                                log.Println(token)

                                tokenBytes := []byte(token)

                                var tokenS Token

                                err := json.Unmarshal(tokenBytes, &tokenS)
                                if err != nil {
                                    log.Println("Error:", err)
                                    tokenStreamChannel <- TokenStreamRes{
                                        Response: nil,
                                        Error: fmt.Errorf("Error while unmarshalling token"),
                                    }
                                }

                                tokenStreamChannel <- TokenStreamRes{
                                    Response: &tokenS,
                                    Error: nil,
                                }

                                token = "";
                            }
                        }
                        if buildingToken {
                            token += string(r);
                        }
                    }
                }
            }
        }

        if err := stream.Err(); err != nil {
            panic(err)
        }
    }()

    // sapotoro := acc.Choices[0].Message.Content
    // log.Println(sapotoro)

    // return &res.Choices[0].Message.Content, nil
    return tokenStreamChannel
}
