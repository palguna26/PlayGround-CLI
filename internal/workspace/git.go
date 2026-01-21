package workspace

import (
	"fmt"
	"os/exec"
	"strings"
)

// GitWorkspace implements Workspace using Git for version control
type GitWorkspace struct {
	repoRoot string
}

// NewGitWorkspace creates a Git-backed workspace
func NewGitWorkspace(repoRoot string) *GitWorkspace {
	return &GitWorkspace{repoRoot: repoRoot}
}

// Diff returns the unified diff for a file
func (gw *GitWorkspace) Diff(filePath string) (string, error) {
	cmd := exec.Command("git", "diff", filePath)
	cmd.Dir = gw.repoRoot
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git diff failed: %w", err)
	}
	return string(output), nil
}

// Snapshot creates a Git stash with a label
func (gw *GitWorkspace) Snapshot(label string) error {
	// First check if there are changes to stash
	statusCmd := exec.Command("git", "status", "--porcelain")
	statusCmd.Dir = gw.repoRoot
	statusOutput, _ := statusCmd.Output()

	if len(strings.TrimSpace(string(statusOutput))) == 0 {
		// No changes to snapshot
		return nil
	}

	// Create stash with message
	cmd := exec.Command("git", "stash", "push", "-m", "pg-snapshot: "+label)
	cmd.Dir = gw.repoRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git stash failed: %w\n%s", err, output)
	}
	return nil
}

// Restore restores files from a Git stash
func (gw *GitWorkspace) Restore(label string) error {
	// List stashes to find the one with matching label
	listCmd := exec.Command("git", "stash", "list")
	listCmd.Dir = gw.repoRoot
	output, err := listCmd.Output()
	if err != nil {
		return fmt.Errorf("git stash list failed: %w", err)
	}

	// Find stash with matching label
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "pg-snapshot: "+label) {
			// Extract stash reference (e.g., "stash@{0}")
			parts := strings.SplitN(line, ":", 2)
			if len(parts) < 1 {
				continue
			}
			stashRef := strings.TrimSpace(parts[0])

			// Pop the stash
			popCmd := exec.Command("git", "stash", "pop", stashRef)
			popCmd.Dir = gw.repoRoot
			popOutput, err := popCmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("git stash pop failed: %w\n%s", err, popOutput)
			}
			return nil
		}
	}

	return fmt.Errorf("snapshot not found: %s", label)
}

// ListSnapshots returns available Git stashes
func (gw *GitWorkspace) ListSnapshots() ([]SnapshotInfo, error) {
	cmd := exec.Command("git", "stash", "list")
	cmd.Dir = gw.repoRoot
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git stash list failed: %w", err)
	}

	var snapshots []SnapshotInfo
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "pg-snapshot:") {
			// Parse label from stash message
			if idx := strings.Index(line, "pg-snapshot:"); idx != -1 {
				label := strings.TrimSpace(line[idx+len("pg-snapshot:"):])
				snapshots = append(snapshots, SnapshotInfo{
					Label: label,
				})
			}
		}
	}

	return snapshots, nil
}

// IsGitBacked returns true for Git workspaces
func (gw *GitWorkspace) IsGitBacked() bool {
	return true
}

// GetRoot returns the repository root
func (gw *GitWorkspace) GetRoot() string {
	return gw.repoRoot
}
