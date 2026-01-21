package workspace

import (
	"os/exec"
	"path/filepath"
	"strings"
)

// Workspace provides a unified interface for version control operations,
// whether backed by Git or a custom snapshot system.
type Workspace interface {
	// Diff returns the unified diff for a file against last snapshot/commit
	Diff(filePath string) (string, error)

	// Snapshot creates a named snapshot of current state
	Snapshot(label string) error

	// Restore restores files to a previous snapshot state
	Restore(label string) error

	// ListSnapshots returns list of available snapshots
	ListSnapshots() ([]SnapshotInfo, error)

	// IsGitBacked returns true if using Git, false if using snapshots
	IsGitBacked() bool

	// GetRoot returns the workspace root directory
	GetRoot() string
}

// SnapshotInfo contains metadata about a snapshot
type SnapshotInfo struct {
	Label     string
	CreatedAt string
	Files     int
}

// NewWorkspace creates an appropriate workspace based on whether Git exists
func NewWorkspace(path string) (Workspace, error) {
	// Check if we're in a Git repository
	if isGitRepo(path) {
		root, err := getGitRoot(path)
		if err != nil {
			return nil, err
		}
		return NewGitWorkspace(root), nil
	}

	// Fall back to snapshot-based workspace
	return NewSnapshotWorkspace(path)
}

// isGitRepo checks if the directory is part of a Git repository
func isGitRepo(path string) bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = path
	err := cmd.Run()
	return err == nil
}

// getGitRoot finds the root directory of the Git repository
func getGitRoot(path string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = path
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// GetWorkspaceDir returns the .pg/workspace directory path
func GetWorkspaceDir(root string) string {
	return filepath.Join(root, ".pg", "workspace")
}
