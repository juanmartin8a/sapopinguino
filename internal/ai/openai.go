package aiutils

import (
	"context"
	"encoding/json"
	// "encoding/json"
	"log"
	"sapopinguino/internal/config"

	// openai "github.com/sashabaranov/go-openai"
	"strings"

	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

var OpenAIClient *openai.Client

func ConfigOpenAI() {
    OpenAIClient = openai.NewClient(
        option.WithAPIKey(config.C.OpenAI.Key),
    )
}

func ChatCompletion(context context.Context, model string, system_role string, input string) (*string, error) {
    stream := OpenAIClient.Chat.Completions.NewStreaming(
		context,
		openai.ChatCompletionNewParams{
			Model: openai.F(openai.ChatModelGPT4oMini),
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
    bbq := false // bbq stands for "backslash befoer quotation"

    acc := openai.ChatCompletionAccumulator{}

    for stream.Next() {
        chunk := stream.Current()
        acc.AddChunk(chunk)

        if content, ok := acc.JustFinishedContent(); ok {
            println("Content stream finished:", content)
        }

        // if using tool calls
        if tool, ok := acc.JustFinishedToolCall(); ok {
            println("Tool call stream finished:", tool.Index, tool.Name, tool.Arguments)
        }

        if refusal, ok := acc.JustFinishedRefusal(); ok {
            println("Refusal stream finished:", refusal)
        }

        if len(chunk.Choices) > 0 {
            aiToken := chunk.Choices[0].Delta.Content
            // println(aiToken)

            if !isInTokensArray {
                if strings.Contains(aiToken, "[") {
                    isInTokensArray = true; 
                    // token = "";
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
                            } else {
                                log.Println("toro")
                                log.Println(tokenS.Type)
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

    sapotoro := acc.Choices[0].Message.Content
    log.Println(sapotoro)

    // return &res.Choices[0].Message.Content, nil
    return nil, nil
}
