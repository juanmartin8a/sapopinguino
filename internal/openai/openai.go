package openaiutils

import (
	"sapopinguino/internal/config"

	openai "github.com/sashabaranov/go-openai"
)

var OpenAIClient *openai.Client

func ConfigOpenAI() {
    OpenAIClient = openai.NewClient(config.C.OpenAI.Key)
}
