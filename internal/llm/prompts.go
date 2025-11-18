package llm

import (
	"fmt"
	"os"
	"path/filepath"
)

// Prompts holds system and user prompts
type Prompts struct {
	System string
	User   string
}

// LoadPrompt reads a file by relative path within the project and returns its contents as a string
func LoadPrompt(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read prompt file %s: %w", path, err)
	}
	return string(data), nil
}

// LoadPrompts loads prompts for the specified mode ("clean" or "spicy")
func LoadPrompts(mode string) (Prompts, error) {
	if mode != "clean" && mode != "spicy" {
		return Prompts{}, fmt.Errorf("invalid mode: %s, must be 'clean' or 'spicy'", mode)
	}

	systemPath := filepath.Join("internal", "prompts", mode, "system.txt")
	advicePath := filepath.Join("internal", "prompts", mode, "advice.txt")

	systemPrompt, err := LoadPrompt(systemPath)
	if err != nil {
		return Prompts{}, fmt.Errorf("failed to load system prompt: %w", err)
	}

	userPrompt, err := LoadPrompt(advicePath)
	if err != nil {
		return Prompts{}, fmt.Errorf("failed to load user prompt: %w", err)
	}

	return Prompts{
		System: systemPrompt,
		User:   userPrompt,
	}, nil
}

