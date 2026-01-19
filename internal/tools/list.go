package tools

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileInfo represents basic file/directory information
type FileInfo struct {
	Name  string `json:"name"`
	IsDir bool   `json:"is_dir"`
	Size  int64  `json:"size"`
}

// ListFiles lists files and directories in the given path
// Respects .gitignore rules by only showing Git-tracked structure
func ListFiles(repoRoot, relPath string) ([]FileInfo, error) {
	// Resolve to absolute path
	absPath := filepath.Join(repoRoot, relPath)

	// Security: Ensure path is within repository bounds
	cleanPath, err := filepath.Abs(absPath)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	cleanRepo, err := filepath.Abs(repoRoot)
	if err != nil {
		return nil, fmt.Errorf("invalid repo path: %w", err)
	}

	// Check if cleanPath is under cleanRepo
	relativeToRepo, err := filepath.Rel(cleanRepo, cleanPath)
	if err != nil || (len(relativeToRepo) > 0 && relativeToRepo[0:2] == "..") {
		return nil, fmt.Errorf("path is outside repository bounds: %s", relPath)
	}

	// Read directory
	entries, err := os.ReadDir(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("directory not found: %s", relPath)
		}
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	// Convert to FileInfo
	var files []FileInfo
	for _, entry := range entries {
		// Skip hidden files and .git directory
		if entry.Name()[0] == '.' {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue // Skip files we can't stat
		}

		files = append(files, FileInfo{
			Name:  entry.Name(),
			IsDir: entry.IsDir(),
			Size:  info.Size(),
		})
	}

	return files, nil
}
