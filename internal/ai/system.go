package aiutils

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sapopinguino/internal/config"
)

var DeveloperPrompt string

var JsonSchema map[string]any

func LoadDeveloperPrompt() {
	bytes, err := os.ReadFile(
		filepath.Join(config.RootDir(), "../assets/developer_prompt.md"),
	)
	if err != nil {
		log.Fatalf("Error while loading the developer prompt: %s", err)
	}

	DeveloperPrompt = string(bytes)
}

func LoadJsonSchema() {
	bytes, err := os.ReadFile(
		filepath.Join(config.RootDir(), "../assets/json_schema.json"),
	)
	if err != nil {
		log.Fatalf("Error while loading the json schema: %s", err)
	}

	if err := json.Unmarshal(bytes, &JsonSchema); err != nil {
		log.Fatalf("Error while unmarshaling the json schema: %s", err)
	}
}
