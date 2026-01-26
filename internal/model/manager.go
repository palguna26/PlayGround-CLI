package model

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ModelManager handles model discovery and validation
type ModelManager struct {
	modelsDir string
}

// ModelInfo contains metadata about a model
type ModelInfo struct {
	Name      string
	Path      string
	SizeBytes int64
	SizeGB    float64
	IsValid   bool
}

// NewModelManager creates a new model manager
func NewModelManager(modelsDir string) *ModelManager {
	if modelsDir == "" {
		homeDir, _ := os.UserHomeDir()
		modelsDir = filepath.Join(homeDir, ".playground", "models")
	}
	return &ModelManager{modelsDir: modelsDir}
}

// GetDefaultModelsDir returns the default models directory
func GetDefaultModelsDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".playground", "models")
}

// FindModel searches for a model by name in the models directory
func (m *ModelManager) FindModel(name string) (string, error) {
	// If name is already an absolute path, validate and return it
	if filepath.IsAbs(name) {
		if err := m.ValidateModel(name); err != nil {
			return "", err
		}
		return name, nil
	}

	// Search in models directory
	modelPath := filepath.Join(m.modelsDir, name)
	if err := m.ValidateModel(modelPath); err == nil {
		return modelPath, nil
	}

	// Try with .gguf extension if not present
	if !strings.HasSuffix(name, ".gguf") {
		modelPath = filepath.Join(m.modelsDir, name+".gguf")
		if err := m.ValidateModel(modelPath); err == nil {
			return modelPath, nil
		}
	}

	return "", fmt.Errorf("model not found: %s", name)
}

// ValidateModel checks if a model file exists and is valid
func (m *ModelManager) ValidateModel(path string) error {
	// Check if file exists
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("model file does not exist: %s", path)
		}
		return fmt.Errorf("failed to stat model file: %w", err)
	}

	// Check if it's a file (not a directory)
	if info.IsDir() {
		return fmt.Errorf("path is a directory, not a file: %s", path)
	}

	// Check file extension
	if !strings.HasSuffix(strings.ToLower(path), ".gguf") {
		return fmt.Errorf("model file must have .gguf extension: %s", path)
	}

	// Check file size (should be at least 100MB for a valid model)
	minSize := int64(100 * 1024 * 1024) // 100MB
	if info.Size() < minSize {
		return fmt.Errorf("model file too small (< 100MB), may be corrupted: %s", path)
	}

	// Check file size (warn if > 20GB, likely wrong file)
	maxSize := int64(20 * 1024 * 1024 * 1024) // 20GB
	if info.Size() > maxSize {
		return fmt.Errorf("model file too large (> 20GB), may be wrong file: %s", path)
	}

	return nil
}

// GetModelInfo returns metadata about a model
func (m *ModelManager) GetModelInfo(path string) (*ModelInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to stat model: %w", err)
	}

	sizeBytes := info.Size()
	sizeGB := float64(sizeBytes) / (1024 * 1024 * 1024)

	isValid := m.ValidateModel(path) == nil

	return &ModelInfo{
		Name:      filepath.Base(path),
		Path:      path,
		SizeBytes: sizeBytes,
		SizeGB:    sizeGB,
		IsValid:   isValid,
	}, nil
}

// ListModels returns all models in the models directory
func (m *ModelManager) ListModels() ([]*ModelInfo, error) {
	// Ensure models directory exists
	if err := os.MkdirAll(m.modelsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create models directory: %w", err)
	}

	entries, err := os.ReadDir(m.modelsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read models directory: %w", err)
	}

	var models []*ModelInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if !strings.HasSuffix(strings.ToLower(entry.Name()), ".gguf") {
			continue
		}

		modelPath := filepath.Join(m.modelsDir, entry.Name())
		info, err := m.GetModelInfo(modelPath)
		if err != nil {
			continue // Skip invalid models
		}

		models = append(models, info)
	}

	return models, nil
}

// EnsureModelsDir creates the models directory if it doesn't exist
func (m *ModelManager) EnsureModelsDir() error {
	return os.MkdirAll(m.modelsDir, 0755)
}
