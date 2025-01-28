// Package api provides integration with Meta's Llama models via langchaingo.
//
// This file implements the Client interface for Meta's Llama models using
// langchaingo's OpenAI-compatible interface. It handles model initialization,
// prompt formatting, and response parsing.
package api

import (
	"context"
	"fmt"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

// MetaLlamaClient implements the Client interface for Meta's Llama models
type MetaLlamaClient struct {
	llm llms.Model // OpenAI-compatible LLM client
}

func (c *MetaLlamaClient) QueryText(prompt string, model string, options Options) (string, error) {
	if prompt == "" {
		return "", fmt.Errorf("prompt cannot be empty for text query")
	}

	// Initialize LLM if not already done
	if c.llm == nil {
		apiKey := os.Getenv("TOGETHER_API_KEY")
		if apiKey == "" {
			return "", fmt.Errorf("TOGETHER_API_KEY environment variable not set")
		}

		baseURL := os.Getenv("TOGETHER_API_BASE")
		if baseURL == "" {
			return "", fmt.Errorf("TOGETHER_API_BASE environment variable not set")
		}

		llm, err := openai.New(
			openai.WithBaseURL(os.Getenv("TOGETHER_API_BASE")),
			openai.WithToken(apiKey),
			openai.WithModel(model),
		)
		if err != nil {
			return "", fmt.Errorf("failed to create Together client: %w", err)
		}
		c.llm = llm
	}

	// Call the LLM with the prompt
	completion, err := llms.GenerateFromSinglePrompt(context.Background(), c.llm, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate completion: %w", err)
	}

	return completion, nil
}

