package aiutils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sapopinguino/internal/config"
)

var DeveloperPrompt string

var JsonSchema map[string]any

func LoadDeveloperPrompt() error {

	bytes, err := os.ReadFile(
		filepath.Join(config.RootDir(), "assets/developer_prompt.md"),
	)
	if err != nil {
		return fmt.Errorf("Error while loading the developer prompt: %s", err)
	}

	DeveloperPrompt = string(bytes)

	return nil
}

func LoadJsonSchema() error {

	bytes, err := os.ReadFile(
		filepath.Join(config.RootDir(), "assets/json_schema.json"),
	)
	if err != nil {
		return fmt.Errorf("Error while loading the json schema: %s", err)
	}

	if err := json.Unmarshal(bytes, &JsonSchema); err != nil {
		return fmt.Errorf("Error while unmarshaling the json schema: %s", err)
	}

	return nil
}
