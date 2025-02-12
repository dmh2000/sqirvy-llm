// Package api provides integration with deepseek's GPT models.
//
// This file implements the Client interface for deepseek's API, supporting
// both text and JSON queries to GPT models. It handles authentication,
// request formatting, and response parsing specific to deepseek's requirements.
package sqirvy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// DeepSeekClient implements the Client interface for deepseek's API
type DeepSeekClient struct {
	apiKey string       // deepseek API authentication key
	client *http.Client // HTTP client for making API requests
}

// deepseekRequest represents the structure of a request to deepseek's chat completion API
type deepseekRequest struct {
	Model          string            `json:"model"`                           // Model identifier
	Messages       []deepseekMessage `json:"messages"`                        // Conversation messages
	MaxTokens      int               `json:"max_completion_tokens,omitempty"` // Max response length
	ResponseFormat string            `json:"response_format,omitempty"`       // Desired response format
	Temperature    float32           `json:"temperature,omitempty"`           // Controls the randomness of the output
}

type deepseekMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type deepseekResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (c *DeepSeekClient) QueryText(prompt string, model string, options Options) (string, error) {
	if prompt == "" {
		return "", fmt.Errorf("prompt cannot be empty for text query")
	}

	// Initialize HTTP client and API key if not already done
	if c.client == nil {
		c.client = &http.Client{}
		// Get API key from environment variable
		c.apiKey = os.Getenv("DEEPSEEK_API_KEY")
		if c.apiKey == "" {
			return "", fmt.Errorf("DEEPSEEK_API_KEY environment variable not set")
		}
	}

	// validate temperature
	if options.Temperature < 0.0 {
		options.Temperature = 0.0
	}
	if options.Temperature > 100.0 {
		return "", fmt.Errorf("temperature must be between 1 and 100")
	}
	// scale Temperature for deepseek 0..2.0
	options.Temperature = (options.Temperature * 2) / 100.0

	// Set default max tokens if not specified
	maxTokens := options.MaxTokens
	if maxTokens == 0 {
		maxTokens = MaxTokensDefault
	}

	// Construct the request body with the prompt as a user message
	reqBody := deepseekRequest{
		Model: model,
		Messages: []deepseekMessage{
			{Role: "user", Content: prompt},
		},
		MaxTokens:   int(maxTokens),      // Limit response length
		Temperature: options.Temperature, // Set temperature
	}

	// Send request and return response
	return c.makeRequest(reqBody)
}

func (c *DeepSeekClient) makeRequest(reqBody deepseekRequest) (string, error) {
	// Convert request body to JSON
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create new HTTP request with JSON body
	endpoint := ""
	if base := os.Getenv("DEEPSEEK_BASE_URL"); base != "" {
		endpoint = base
	} else {
		return "", fmt.Errorf("DEEPSEEK_BASE_URL environment variable not set")
	}
	endpoint += "/chat/completions"

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Send the request
	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response JSON
	var deepseekResp deepseekResponse
	if err := json.Unmarshal(body, &deepseekResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Ensure we got at least one choice back
	if len(deepseekResp.Choices) == 0 {
		return "", fmt.Errorf("no content in response")
	}

	// Return the content of the first choice
	return deepseekResp.Choices[0].Message.Content, nil
}
