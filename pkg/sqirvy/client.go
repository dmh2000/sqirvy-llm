// Package api provides a unified interface for interacting with various AI language models.
//
// The package supports multiple AI providers including:
// - Anthropic (Claude models)
// - Google (Gemini models)
// - OpenAI (GPT models)
//
// It provides a consistent interface for making text and JSON queries while handling
// provider-specific implementation details internally.
package api

import (
	"fmt"
)

// Provider represents supported AI providers.
// Currently supports Anthropic, Gemini, and OpenAI.
// Provider identifies which AI service provider to use
type Provider string

// Supported AI providers
const (
	Anthropic Provider = "anthropic" // Anthropic's Claude models
	Gemini    Provider = "gemini"    // Google's Gemini models
	OpenAI    Provider = "openai"    // OpenAI's GPT models
)

// AnthropicOptions contains Anthropic-specific configuration options
type AnthropicOptions struct {
	// Reserved for future Anthropic-specific settings
}

// GeminiOptions contains Google Gemini-specific configuration options
type GeminiOptions struct {
	// Reserved for future Gemini-specific settings
}

// OpenAIOptions contains OpenAI-specific configuration options
type OpenAIOptions struct {
	// Reserved for future OpenAI-specific settings
}

// Options combines all provider-specific options into a single structure.
// This allows for provider-specific configuration while maintaining a unified interface.
type Options struct {
	AnthropicOptions
	GeminiOptions
	OpenAIOptions
}

// Client provides a unified interface for AI operations.
// It abstracts away provider-specific implementations behind a common interface
// for making text and JSON queries to AI models.
type Client interface {
	QueryText(prompt string, model string, options Options) (string, error)
	QueryJSON(prompt string, model string, options Options) (string, error)
}

// NewClient creates a new AI client for the specified provider
func NewClient(provider Provider) (Client, error) {
	switch provider {
	case Anthropic:
		return &AnthropicClient{}, nil
	case Gemini:
		return &GeminiClient{}, nil
	case OpenAI:
		return &OpenAIClient{}, nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}
