package api

import (
	"os"
	"strings"
	"testing"
)

func TestAnthropicClient_QueryText(t *testing.T) {
	if os.Getenv("ANTHROPIC_API_KEY") == "" {
		t.Skip("ANTHROPIC_API_KEY not set")
	}

	client := &AnthropicClient{}

	tests := []struct {
		name   string
		prompt string
	}{
		{
			name:   "Basic prompt",
			prompt: "Say 'Hello, World!'",
		},
		{
			name:   "Empty prompt",
			prompt: "",
		},
	}

	tt := tests[0]
	t.Run(tt.name, func(t *testing.T) {
		got, err := client.QueryText(tt.prompt, "claude-3-sonnet-20240229", Options{})
		if err != nil {
			t.Errorf("AnthropicClient.QueryText() error = %v", err)
			return
		}
		if !strings.Contains(got, "Hello") {
			t.Errorf("AnthropicClient.QueryText() = %v, expected response containing 'Hello'", got)
		}
	})

	tt = tests[1]
	t.Run(tt.name, func(t *testing.T) {
		_, err := client.QueryText(tt.prompt, "claude-3-sonnet-20240229", Options{})
		if err == nil {
			t.Errorf("AnthropicClient.QueryText() empty prompt should have failed")
			return
		}
	})
}

func TestAnthropicClient_QueryJSON(t *testing.T) {
	if os.Getenv("ANTHROPIC_API_KEY") == "" {
		t.Skip("ANTHROPIC_API_KEY not set")
	}

	client := &AnthropicClient{}

	tests := []struct {
		name    string
		prompt  string
		wantErr bool
	}{
		{
			name:   "JSON request",
			prompt: "Return a JSON object with a greeting field containing 'Hello, World!'",
		},
		{
			name:   "Empty prompt",
			prompt: "",
		},
	}

	tt := tests[0]
	t.Run(tt.name, func(t *testing.T) {
		got, err := client.QueryJSON(tt.prompt, "claude-3-sonnet-20240229", Options{})
		if err != nil {
			t.Errorf("AnthropicClient.QueryJSON() error = %v", err)
			return
		}
		if !strings.Contains(got, "{") {
			t.Errorf("AnthropicClient.QueryJSON() = %v, expected JSON response", got)
		}
	})

	tt = tests[1]
	t.Run(tt.name, func(t *testing.T) {
		_, err := client.QueryJSON(tt.prompt, "claude-3-sonnet-20240229", Options{})
		if err == nil {
			t.Errorf("AnthropicClient.QueryJSON() empty prompt should have failed")
			return
		}
	})
}
