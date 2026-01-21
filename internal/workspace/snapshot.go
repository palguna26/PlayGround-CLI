package workspace

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// SnapshotWorkspace implements Workspace using a custom SHA-based snapshot system
// for repositories without Git.
type SnapshotWorkspace struct {
	rootDir      string
	workspaceDir string
	objectsDir   string
	snapshotsDir string
}

// SnapshotManifest represents a point-in-time snapshot
type SnapshotManifest struct {
	Label     string            `json:"label"`
	CreatedAt time.Time         `json:"created_at"`
	Files     map[string]string `json:"files"` // path -> SHA256
}

// NewSnapshotWorkspace creates a snapshot-based workspace
func NewSnapshotWorkspace(rootDir string) (*SnapshotWorkspace, error) {
	workspaceDir := filepath.Join(rootDir, ".pg", "workspace")
	objectsDir := filepath.Join(workspaceDir, "objects")
	snapshotsDir := filepath.Join(workspaceDir, "snapshots")

	// Create directories
	for _, dir := range []string{workspaceDir, objectsDir, snapshotsDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create workspace directory: %w", err)
		}
	}

	return &SnapshotWorkspace{
		rootDir:      rootDir,
		workspaceDir: workspaceDir,
		objectsDir:   objectsDir,
		snapshotsDir: snapshotsDir,
	}, nil
}

// Diff returns a unified diff comparing current file with last snapshot
func (sw *SnapshotWorkspace) Diff(filePath string) (string, error) {
	fullPath := filepath.Join(sw.rootDir, filePath)

	// Read current file
	currentContent, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Sprintf("--- /dev/null\n+++ %s\n@@ -0,0 +1 @@\n+(new file)", filePath), nil
		}
		return "", err
	}

	// Get last snapshot content
	lastContent, err := sw.getLastSnapshotContent(filePath)
	if err != nil {
		// No previous snapshot - show as new file
		lines := strings.Split(string(currentContent), "\n")
		var diff strings.Builder
		diff.WriteString(fmt.Sprintf("--- /dev/null\n+++ %s\n", filePath))
		diff.WriteString(fmt.Sprintf("@@ -0,0 +1,%d @@\n", len(lines)))
		for _, line := range lines {
			diff.WriteString("+" + line + "\n")
		}
		return diff.String(), nil
	}

	// Generate unified diff
	return generateUnifiedDiff(filePath, string(lastContent), string(currentContent)), nil
}

// Snapshot creates a named snapshot of all tracked files
func (sw *SnapshotWorkspace) Snapshot(label string) error {
	manifest := SnapshotManifest{
		Label:     label,
		CreatedAt: time.Now(),
		Files:     make(map[string]string),
	}

	// Walk the directory and snapshot all files
	err := filepath.Walk(sw.rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and hidden files/dirs
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip hidden files
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Compute SHA256
		hash := sha256.Sum256(content)
		sha := hex.EncodeToString(hash[:])

		// Store object
		objectPath := filepath.Join(sw.objectsDir, sha)
		if _, err := os.Stat(objectPath); os.IsNotExist(err) {
			if err := os.WriteFile(objectPath, content, 0644); err != nil {
				return err
			}
		}

		// Add to manifest
		relPath, _ := filepath.Rel(sw.rootDir, path)
		manifest.Files[relPath] = sha

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk directory: %w", err)
	}

	// Save manifest
	manifestPath := filepath.Join(sw.snapshotsDir, label+".json")
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(manifestPath, data, 0644)
}

// Restore restores files to a previous snapshot state
func (sw *SnapshotWorkspace) Restore(label string) error {
	// Load manifest
	manifestPath := filepath.Join(sw.snapshotsDir, label+".json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("snapshot not found: %s", label)
	}

	var manifest SnapshotManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return err
	}

	// Restore each file
	for filePath, sha := range manifest.Files {
		objectPath := filepath.Join(sw.objectsDir, sha)
		content, err := os.ReadFile(objectPath)
		if err != nil {
			return fmt.Errorf("object not found: %s", sha)
		}

		fullPath := filepath.Join(sw.rootDir, filePath)

		// Ensure directory exists
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return err
		}

		if err := os.WriteFile(fullPath, content, 0644); err != nil {
			return err
		}
	}

	return nil
}

// ListSnapshots returns available snapshots
func (sw *SnapshotWorkspace) ListSnapshots() ([]SnapshotInfo, error) {
	entries, err := os.ReadDir(sw.snapshotsDir)
	if err != nil {
		return nil, err
	}

	var snapshots []SnapshotInfo
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		// Load manifest for metadata
		manifestPath := filepath.Join(sw.snapshotsDir, entry.Name())
		data, err := os.ReadFile(manifestPath)
		if err != nil {
			continue
		}

		var manifest SnapshotManifest
		if err := json.Unmarshal(data, &manifest); err != nil {
			continue
		}

		snapshots = append(snapshots, SnapshotInfo{
			Label:     manifest.Label,
			CreatedAt: manifest.CreatedAt.Format("2006-01-02 15:04:05"),
			Files:     len(manifest.Files),
		})
	}

	return snapshots, nil
}

// IsGitBacked returns false for snapshot workspaces
func (sw *SnapshotWorkspace) IsGitBacked() bool {
	return false
}

// GetRoot returns the workspace root directory
func (sw *SnapshotWorkspace) GetRoot() string {
	return sw.rootDir
}

// getLastSnapshotContent retrieves a file's content from the most recent snapshot
func (sw *SnapshotWorkspace) getLastSnapshotContent(filePath string) ([]byte, error) {
	entries, err := os.ReadDir(sw.snapshotsDir)
	if err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("no snapshots available")
	}

	// Find most recent snapshot (by modification time)
	var latestManifest *SnapshotManifest
	var latestTime time.Time

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		manifestPath := filepath.Join(sw.snapshotsDir, entry.Name())
		data, err := os.ReadFile(manifestPath)
		if err != nil {
			continue
		}

		var manifest SnapshotManifest
		if err := json.Unmarshal(data, &manifest); err != nil {
			continue
		}

		if manifest.CreatedAt.After(latestTime) {
			latestTime = manifest.CreatedAt
			latestManifest = &manifest
		}
	}

	if latestManifest == nil {
		return nil, fmt.Errorf("no valid snapshots found")
	}

	sha, ok := latestManifest.Files[filePath]
	if !ok {
		return nil, fmt.Errorf("file not in snapshot: %s", filePath)
	}

	objectPath := filepath.Join(sw.objectsDir, sha)
	return os.ReadFile(objectPath)
}

// generateUnifiedDiff creates a simple unified diff between two strings
func generateUnifiedDiff(filePath, old, new string) string {
	oldLines := strings.Split(old, "\n")
	newLines := strings.Split(new, "\n")

	var diff strings.Builder
	diff.WriteString(fmt.Sprintf("--- %s\n+++ %s\n", filePath, filePath))

	// Simple line-by-line diff (not optimal, but functional)
	maxLines := len(oldLines)
	if len(newLines) > maxLines {
		maxLines = len(newLines)
	}

	diff.WriteString(fmt.Sprintf("@@ -1,%d +1,%d @@\n", len(oldLines), len(newLines)))

	for i := 0; i < maxLines; i++ {
		if i < len(oldLines) && i < len(newLines) {
			if oldLines[i] != newLines[i] {
				diff.WriteString("-" + oldLines[i] + "\n")
				diff.WriteString("+" + newLines[i] + "\n")
			} else {
				diff.WriteString(" " + oldLines[i] + "\n")
			}
		} else if i < len(oldLines) {
			diff.WriteString("-" + oldLines[i] + "\n")
		} else {
			diff.WriteString("+" + newLines[i] + "\n")
		}
	}

	return diff.String()
}
