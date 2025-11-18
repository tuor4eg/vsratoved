package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/tuor4eg/vsratoved/internal/config"
)

// OpenRouterRequest represents the request structure for OpenRouter API
type OpenRouterRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens,omitempty"`
}

// Message represents a single message in the conversation
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenRouterResponse represents the response structure from OpenRouter API
type OpenRouterResponse struct {
	Choices []Choice `json:"choices"`
}

// Choice represents a choice in the response
type Choice struct {
	Message Message `json:"message"`
}

// AdviceResponse represents parsed advice with author
type AdviceResponse struct {
	Author string
	Advice string
}

func ErrorFallbackMessage() string {
	return "Почему-то совет не доехал, может денег нет, но ты всё равно держись"
}

// GetWeirdAdvice generates weird advice using OpenRouter API
func GetWeirdAdvice(ctx context.Context, mode string) (*AdviceResponse, error) {
	// Check if config is loaded
	if config.C == nil {
		return nil, fmt.Errorf("config is not loaded, call config.Load() first")
	}

	apiKey := config.C.OpenRouterAPIKey
	apiURL := config.C.OpenRouterAPIURL
	model := config.C.OpenRouterModel

	// Load prompts
	prompts, err := LoadPrompts(mode)
	if err != nil {
		return nil, fmt.Errorf("failed to load prompts: %w", err)
	}

	// Build request
	reqBody := OpenRouterRequest{
		Model: model,
		Messages: []Message{
			{Role: "system", Content: prompts.System},
			{Role: "user", Content: prompts.User},
		},
		MaxTokens: 1000,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON response
	var openRouterResp OpenRouterResponse
	if err := json.Unmarshal(body, &openRouterResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Extract content from choices[0].message.content
	if len(openRouterResp.Choices) == 0 {
		return nil, fmt.Errorf("response contains no choices")
	}

	content := openRouterResp.Choices[0].Message.Content
	if content == "" {
		return nil, fmt.Errorf("response content is empty")
	}

	// Parse content to extract author and advice
	return parseAdviceResponse(content), nil
}

// parseAdviceResponse parses LLM response and extracts author and advice
func parseAdviceResponse(content string) *AdviceResponse {
	lines := strings.Split(strings.TrimSpace(content), "\n")
	var author, advice string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Автор:") {
			author = strings.TrimSpace(strings.TrimPrefix(line, "Автор:"))
		} else if strings.HasPrefix(line, "Совет:") {
			advice = strings.TrimSpace(strings.TrimPrefix(line, "Совет:"))
		}
	}

	return &AdviceResponse{
		Author: author,
		Advice: advice,
	}
}
