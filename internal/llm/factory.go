package llm

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// NewProvider creates an LLM provider based on environment configuration
// Priority:
// 1. LLM_PROVIDER env var (openai, gemini)
// 2. API key presence (OPENAI_API_KEY or GEMINI_API_KEY)
// 3. Defaults to OpenAI if both are set
func NewProvider() (Provider, error) {
	// Check explicit provider selection
	providerName := strings.ToLower(os.Getenv("LLM_PROVIDER"))

	switch providerName {
	case "openai":
		return NewOpenAIProvider()
	case "gemini":
		return NewGeminiProvider()
	case "":
		// Auto-detect based on API keys
		return autoDetectProvider()
	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s (supported: openai, gemini)", providerName)
	}
}

// autoDetectProvider automatically detects which provider to use based on API keys
func autoDetectProvider() (Provider, error) {
	// Check environment variables first
	hasOpenAI := os.Getenv("OPENAI_API_KEY") != ""
	hasGemini := os.Getenv("GEMINI_API_KEY") != ""

	// If not in env, check config file
	if !hasOpenAI || !hasGemini {
		// Try to load from config file (import needed)
		homeDir, err := os.UserHomeDir()
		if err == nil {
			configPath := homeDir + "/.playground/config.json"
			if data, err := os.ReadFile(configPath); err == nil {
				var config map[string]string
				if err := json.Unmarshal(data, &config); err == nil {
					if !hasGemini && config["gemini_api_key"] != "" {
						os.Setenv("GEMINI_API_KEY", config["gemini_api_key"])
						hasGemini = true
					}
					if !hasOpenAI && config["openai_api_key"] != "" {
						os.Setenv("OPENAI_API_KEY", config["openai_api_key"])
						hasOpenAI = true
					}
				}
			}
		}
	}

	if !hasOpenAI && !hasGemini {
		return nil, fmt.Errorf("no LLM API key found. Run 'pg setup' to configure or set OPENAI_API_KEY or GEMINI_API_KEY environment variable")
	}

	// Prefer Gemini if both are set (user's request context)
	if hasGemini {
		provider, err := NewGeminiProvider()
		if err == nil {
			return provider, nil
		}
		// Fall through to OpenAI if Gemini fails
		if !hasOpenAI {
			return nil, err
		}
	}

	// Default to OpenAI
	return NewOpenAIProvider()
}
