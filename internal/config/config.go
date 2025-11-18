package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	OpenRouterAPIKey string
	OpenRouterAPIURL string
	OpenRouterModel  string
	TelegramBotToken string
}

var C *Config

// Load loads configuration from .env file and environment variables
func Load() error {
	// Load .env file if it exists
	_ = godotenv.Load()

	C = &Config{
		OpenRouterAPIKey: os.Getenv("OPENROUTER_API_KEY"),
		OpenRouterAPIURL: os.Getenv("OPENROUTER_API_URL"),
		OpenRouterModel:  os.Getenv("OPENROUTER_MODEL"),
		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
	}

	// Validate required fields
	if C.OpenRouterAPIKey == "" {
		return fmt.Errorf("OPENROUTER_API_KEY is required")
	}
	if C.OpenRouterAPIURL == "" {
		return fmt.Errorf("OPENROUTER_API_URL is required")
	}
	if C.OpenRouterModel == "" {
		return fmt.Errorf("OPENROUTER_MODEL is required")
	}

	return nil
}
