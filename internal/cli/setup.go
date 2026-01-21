package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// Config represents the PlayGround configuration
type Config struct {
	GeminiAPIKey string `json:"gemini_api_key,omitempty"`
	OpenAIAPIKey string `json:"openai_api_key,omitempty"`
	LLMProvider  string `json:"llm_provider,omitempty"` // "openai" or "gemini"
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Interactive setup wizard for PlayGround",
	Long: `Configure PlayGround with your API keys and preferences.
This creates a config file in your home directory that persists across sessions.

Example:
  pg setup`,
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)

		fmt.Println("╔════════════════════════════════════════╗")
		fmt.Println("║   PlayGround CLI - Setup Wizard       ║")
		fmt.Println("╚════════════════════════════════════════╝")
		fmt.Println()

		// Load existing config if present
		config, _ := LoadConfig()
		if config == nil {
			config = &Config{}
		}

		// Display current configuration
		if config.GeminiAPIKey != "" || config.OpenAIAPIKey != "" {
			fmt.Println("📋 Current Configuration:")
			if config.GeminiAPIKey != "" {
				fmt.Printf("  • Gemini API Key: %s\n", maskAPIKey(config.GeminiAPIKey))
			}
			if config.OpenAIAPIKey != "" {
				fmt.Printf("  • OpenAI API Key: %s\n", maskAPIKey(config.OpenAIAPIKey))
			}
			if config.LLMProvider != "" {
				fmt.Printf("  • Preferred Provider: %s\n", config.LLMProvider)
			}
			fmt.Println()
		}

		// Provider selection
		fmt.Println("Which LLM provider would you like to use?")
		fmt.Println("  1) Google Gemini (Recommended - generous free tier)")
		fmt.Println("  2) OpenAI (GPT-4, GPT-3.5)")
		fmt.Println("  3) Both (configure both, auto-select)")
		fmt.Print("\nChoice [1-3]: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		setupGemini := choice == "1" || choice == "3"
		setupOpenAI := choice == "2" || choice == "3"

		// Gemini setup
		if setupGemini {
			fmt.Println("\n🔹 Gemini Configuration")
			fmt.Println("Get your API key from: https://makersuite.google.com/app/apikey")
			fmt.Print("Enter Gemini API Key (or press Enter to skip): ")

			apiKey, _ := reader.ReadString('\n')
			apiKey = strings.TrimSpace(apiKey)

			if apiKey != "" {
				config.GeminiAPIKey = apiKey
				fmt.Println("✓ Gemini API key saved")
			}
		}

		// OpenAI setup
		if setupOpenAI {
			fmt.Println("\n🔹 OpenAI Configuration")
			fmt.Println("Get your API key from: https://platform.openai.com/api-keys")
			fmt.Print("Enter OpenAI API Key (or press Enter to skip): ")

			apiKey, _ := reader.ReadString('\n')
			apiKey = strings.TrimSpace(apiKey)

			if apiKey != "" {
				config.OpenAIAPIKey = apiKey
				fmt.Println("✓ OpenAI API key saved")
			}
		}

		// Provider preference (if both configured)
		if config.GeminiAPIKey != "" && config.OpenAIAPIKey != "" {
			fmt.Println("\n🔹 Provider Preference")
			fmt.Println("Both providers are configured. Which should be preferred?")
			fmt.Println("  1) Gemini (faster, more generous free tier)")
			fmt.Println("  2) OpenAI (GPT-4 capabilities)")
			fmt.Print("\nChoice [1-2]: ")

			pref, _ := reader.ReadString('\n')
			pref = strings.TrimSpace(pref)

			if pref == "2" {
				config.LLMProvider = "openai"
			} else {
				config.LLMProvider = "gemini"
			}
		} else if config.GeminiAPIKey != "" {
			config.LLMProvider = "gemini"
		} else if config.OpenAIAPIKey != "" {
			config.LLMProvider = "openai"
		}

		// Save configuration
		if err := SaveConfig(config); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		// Success message
		fmt.Println("\n✅ Configuration saved successfully!")
		fmt.Printf("📁 Config location: %s\n", GetConfigPath())
		fmt.Println("\n🚀 You're ready to go! Try:")
		fmt.Println("   pg start \"your goal here\"")
		fmt.Println("   pg ask \"your question\"")

		return nil
	},
}

// LoadConfig loads the configuration from disk
func LoadConfig() (*Config, error) {
	configPath := GetConfigPath()

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No config file yet
		}
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig saves the configuration to disk
func SaveConfig(config *Config) error {
	configPath := GetConfigPath()

	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600) // 0600 for security (user-only read/write)
}

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	return filepath.Join(homeDir, ".playground", "config.json")
}

// maskAPIKey masks an API key for display
func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "****" + key[len(key)-4:]
}

// GetAPIKeys returns API keys from config file, with env var override
func GetAPIKeys() (geminiKey, openaiKey, provider string) {
	// First check environment variables (highest priority)
	geminiKey = os.Getenv("GEMINI_API_KEY")
	openaiKey = os.Getenv("OPENAI_API_KEY")
	provider = os.Getenv("LLM_PROVIDER")

	// If not in env, load from config file
	if geminiKey == "" || openaiKey == "" || provider == "" {
		config, err := LoadConfig()
		if err == nil && config != nil {
			if geminiKey == "" {
				geminiKey = config.GeminiAPIKey
			}
			if openaiKey == "" {
				openaiKey = config.OpenAIAPIKey
			}
			if provider == "" {
				provider = config.LLMProvider
			}
		}
	}

	return
}
