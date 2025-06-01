package aiutils

import (
	"log"
	"os"
)

var DeveloperPrompt string

func LoadMarkdown() {
	content, err := os.ReadFile("../../assets/developer_prompt.md")
	if err != nil {
		log.Fatalf("Error while loading the developer prompt: %s", err)
	}

	DeveloperPrompt = string(content)
}
