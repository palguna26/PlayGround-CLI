package model

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// DownloadHelper provides utilities for downloading models
type DownloadHelper struct {
	manager *ModelManager
}

// NewDownloadHelper creates a new download helper
func NewDownloadHelper(manager *ModelManager) *DownloadHelper {
	return &DownloadHelper{manager: manager}
}

// DownloadDeepSeekCoder downloads DeepSeek-Coder-7B-Instruct v1.5 Q4_K_M
func (d *DownloadHelper) DownloadDeepSeekCoder(progressCallback func(downloaded, total int64)) error {
	// Model URL from HuggingFace
	modelURL := "https://huggingface.co/TheBloke/deepseek-coder-7B-instruct-v1.5-GGUF/resolve/main/deepseek-coder-7b-instruct-v1.5.Q4_K_M.gguf"
	modelName := "deepseek-coder-7b-instruct-v1.5.Q4_K_M.gguf"

	// Ensure models directory exists
	if err := d.manager.EnsureModelsDir(); err != nil {
		return fmt.Errorf("failed to create models directory: %w", err)
	}

	modelPath := filepath.Join(d.manager.modelsDir, modelName)

	// Check if already downloaded
	if err := d.manager.ValidateModel(modelPath); err == nil {
		return fmt.Errorf("model already exists: %s", modelPath)
	}

	// Download the model
	fmt.Printf("Downloading model from: %s\n", modelURL)
	fmt.Printf("Destination: %s\n", modelPath)
	fmt.Println("This may take a while (~4GB download)...")

	resp, err := http.Get(modelURL)
	if err != nil {
		return fmt.Errorf("failed to download model: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %s", resp.Status)
	}

	// Create temporary file
	tmpPath := modelPath + ".tmp"
	out, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Download with progress tracking
	totalSize := resp.ContentLength
	downloaded := int64(0)

	buf := make([]byte, 32*1024) // 32KB buffer
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			if _, writeErr := out.Write(buf[:n]); writeErr != nil {
				os.Remove(tmpPath)
				return fmt.Errorf("failed to write to file: %w", writeErr)
			}
			downloaded += int64(n)

			if progressCallback != nil {
				progressCallback(downloaded, totalSize)
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			os.Remove(tmpPath)
			return fmt.Errorf("download error: %w", err)
		}
	}

	// Close file before rename
	out.Close()

	// Rename temp file to final name
	if err := os.Rename(tmpPath, modelPath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to finalize download: %w", err)
	}

	// Validate downloaded model
	if err := d.manager.ValidateModel(modelPath); err != nil {
		os.Remove(modelPath)
		return fmt.Errorf("downloaded model is invalid: %w", err)
	}

	fmt.Printf("\n✅ Model downloaded successfully: %s\n", modelPath)
	return nil
}

// GetRecommendedModel returns the path to the recommended model
func (d *DownloadHelper) GetRecommendedModel() string {
	return filepath.Join(d.manager.modelsDir, "deepseek-coder-7b-instruct-v1.5.Q4_K_M.gguf")
}
