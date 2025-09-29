// dev.go
//go:build !prod && !dev

package main

import (
	"context"
	"log"
	aiutils "sapopinguino/internal/ai"
	"sapopinguino/internal/config"

	"github.com/openai/openai-go/v2"
)

func main() {
	config.LoadDotEnv()

	c, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	aiutils.ConfigOpenAI(c)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	input := `{"input_language":"English","target_language":"Spanish","input":"Hi! How are you doing!? I heeard that you won the comp. That's amazing!!"}`

	tokenStreamChannel := aiutils.StreamResponse(ctx, openai.ChatModelGPT5, aiutils.DeveloperPrompt, input)

	for res := range tokenStreamChannel {
		if res.Error != nil {
			log.Printf("error: %v", res.Error)
		}

		log.Printf("res: %v", res.Response)
	}
}
