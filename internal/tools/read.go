package tools

import (
	"fmt"
	"os"
	"path/filepath"
)

// ReadFile reads a file's contents from the repository
// Security: validates path is within repo bounds
func ReadFile(repoRoot, relPath string) (string, error) {
	// Resolve to absolute path
	absPath := filepath.Join(repoRoot, relPath)

	// Security: Ensure path is within repository bounds
	cleanPath, err := filepath.Abs(absPath)
	if err != nil {
		return "", fmt.Errorf("invalid path: %w", err)
	}

	cleanRepo, err := filepath.Abs(repoRoot)
	if err != nil {
		return "", fmt.Errorf("invalid repo path: %w", err)
	}

	// Check if cleanPath is under cleanRepo
	relativeToRepo, err := filepath.Rel(cleanRepo, cleanPath)
	if err != nil || len(relativeToRepo) > 0 && relativeToRepo[0:2] == ".." {
		return "", fmt.Errorf("path is outside repository bounds: %s", relPath)
	}

	// Read file
	content, err := os.ReadFile(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("file not found: %s", relPath)
		}
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return string(content), nil
}
