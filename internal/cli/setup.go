package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/playground/internal/model"
)

// Config represents the PlayGround configuration
type Config struct {
	ModelPath string `json:"model_path,omitempty"` // Path to local GGUF model
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Configure local model for PlayGround",
	Long: `Set up PlayGround with a local DeepSeek-Coder model.

This wizard will help you:
  1. Specify the path to your local GGUF model
  2. Validate the model file exists
  3. Check system requirements (RAM)

Example:
  pg setup`,
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)

		fmt.Println("╔════════════════════════════════════════╗")
		fmt.Println("║   PlayGround CLI - Local Model Setup  ║")
		fmt.Println("╚════════════════════════════════════════╝")
		fmt.Println()

		// Load existing config if present
		config, _ := LoadConfig()
		if config == nil {
			config = &Config{}
		}

		// Display current configuration
		if config.ModelPath != "" {
			fmt.Println("📋 Current Configuration:")
			fmt.Printf("  • Model Path: %s\n", config.ModelPath)
			fmt.Println()
		}

		// Model path configuration
		fmt.Println("🤖 Local Model Configuration")
		fmt.Println()
		fmt.Println("PlayGround uses DeepSeek-Coder-7B-Instruct v1.5 (GGUF format)")
		fmt.Println("Recommended quantization: Q4_K_M")
		fmt.Println()

		defaultPath := filepath.Join(getHomeDir(), ".playground", "models", "deepseek-coder-7b-instruct-v1.5.Q4_K_M.gguf")
		fmt.Printf("Default path: %s\n", defaultPath)
		fmt.Println()

		fmt.Print("Enter model path (or press Enter for default): ")
		modelPath, _ := reader.ReadString('\n')
		modelPath = strings.TrimSpace(modelPath)

		if modelPath == "" {
			modelPath = defaultPath
		}

		// Expand ~ to home directory
		if strings.HasPrefix(modelPath, "~") {
			modelPath = filepath.Join(getHomeDir(), modelPath[1:])
		}

		// Validate model file exists
		if _, err := os.Stat(modelPath); os.IsNotExist(err) {
			fmt.Println()
			fmt.Println("⚠️  Model file not found!")
			fmt.Println()

			// Offer automatic download
			fmt.Println("Would you like to download the model automatically?")
			fmt.Println("  Model: DeepSeek-Coder-7B-Instruct v1.5 (Q4_K_M)")
			fmt.Println("  Size: ~4GB")
			fmt.Println("  Source: HuggingFace")
			fmt.Println()
			fmt.Print("Download now? [Y/n]: ")

			downloadChoice, _ := reader.ReadString('\n')
			downloadChoice = strings.TrimSpace(strings.ToLower(downloadChoice))

			if downloadChoice == "" || downloadChoice == "y" || downloadChoice == "yes" {
				// Create models directory
				modelsDir := filepath.Dir(modelPath)
				if err := os.MkdirAll(modelsDir, 0755); err != nil {
					return fmt.Errorf("failed to create models directory: %w", err)
				}

				// Import model package
				// Note: This requires adding the import at the top of the file
				fmt.Println()
				fmt.Println("📥 Downloading model...")
				fmt.Println("This may take 10-30 minutes depending on your connection.")
				fmt.Println()

				// Use the download helper
				manager := model.NewModelManager(modelsDir)
				downloader := model.NewDownloadHelper(manager)

				// Progress callback
				lastPercent := -1
				progressCallback := func(downloaded, total int64) {
					if total > 0 {
						percent := int(float64(downloaded) / float64(total) * 100)
						if percent != lastPercent && percent%5 == 0 {
							fmt.Printf("Progress: %d%% (%.2f GB / %.2f GB)\n",
								percent,
								float64(downloaded)/(1024*1024*1024),
								float64(total)/(1024*1024*1024))
							lastPercent = percent
						}
					}
				}

				if err := downloader.DownloadDeepSeekCoder(progressCallback); err != nil {
					fmt.Println()
					fmt.Println("❌ Download failed:", err)
					fmt.Println()
					fmt.Println("Manual download instructions:")
					fmt.Println("  1. Visit: https://huggingface.co/TheBloke/deepseek-coder-7B-instruct-v1.5-GGUF")
					fmt.Println("  2. Download: deepseek-coder-7b-instruct-v1.5.Q4_K_M.gguf (~4GB)")
					fmt.Printf("  3. Save to: %s\n", modelPath)
					fmt.Println()
					return fmt.Errorf("model download failed")
				}

				// Update modelPath to the downloaded file
				modelPath = downloader.GetRecommendedModel()
				fmt.Println()
				fmt.Println("✅ Model downloaded successfully!")
			} else {
				fmt.Println()
				fmt.Println("Manual download instructions:")
				fmt.Println("  1. Visit: https://huggingface.co/TheBloke/deepseek-coder-7B-instruct-v1.5-GGUF")
				fmt.Println("  2. Download: deepseek-coder-7b-instruct-v1.5.Q4_K_M.gguf (~4GB)")
				fmt.Printf("  3. Save to: %s\n", modelPath)
				fmt.Println()
				fmt.Println("Then run 'pg setup' again.")
				return fmt.Errorf("model file not found: %s", modelPath)
			}
		}

		// Check file size (should be ~3-5GB for Q4_K_M)
		fileInfo, _ := os.Stat(modelPath)
		fileSizeGB := float64(fileInfo.Size()) / (1024 * 1024 * 1024)

		fmt.Println()
		fmt.Printf("✓ Model found: %.2f GB\n", fileSizeGB)

		// RAM warning
		fmt.Println()
		fmt.Println("💾 System Requirements:")
		fmt.Println("  • Minimum RAM: 8 GB")
		fmt.Println("  • Recommended: 16 GB")
		fmt.Println("  • Model will run on CPU")
		fmt.Println()

		// Save configuration
		config.ModelPath = modelPath
		if err := SaveConfig(config); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		// Success message
		fmt.Println("✅ Configuration saved successfully!")
		fmt.Printf("📁 Config location: %s\n", GetConfigPath())
		fmt.Println()
		fmt.Println("🚀 You're ready to go! Try:")
		fmt.Println("   pg agent")
		fmt.Println("   pg start \"your goal here\"")

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
	return filepath.Join(getHomeDir(), ".playground", "config.json")
}

// getHomeDir returns the user's home directory
func getHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	return homeDir
}
