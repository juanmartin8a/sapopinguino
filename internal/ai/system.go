package aiutils

import (
	"encoding/json"
	"log"
	"os"
)

var DeveloperPrompt string

var JSON_Schema map[string]any

func LoadDeveloperPrompt() {
	bytes, err := os.ReadFile("../../assets/developer_prompt.md")
	if err != nil {
		log.Fatalf("Error while loading the developer prompt: %s", err)
	}

	DeveloperPrompt = string(bytes)
}

func LoadJsonSchema() {
	bytes, err := os.ReadFile("../../assets/json_schema.json")
	if err != nil {
		log.Fatalf("Error while loading the json schema: %s", err)
	}

	if err := json.Unmarshal(bytes, &JSON_Schema); err != nil {
		log.Fatalf("Error while unmarshaling the json schema: %s", err)
	}
}
