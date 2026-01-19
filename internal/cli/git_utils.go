package cli

import (
	"os"
	"os/exec"
	"path/filepath"
)

// isGitRepo checks if the given directory is within a Git repository
func isGitRepo(dir string) bool {
	gitDir := filepath.Join(dir, ".git")
	if _, err := os.Stat(gitDir); err == nil {
		return true
	}

	// Check parent directory recursively
	parent := filepath.Dir(dir)
	if parent == dir {
		return false // Reached root
	}

	return isGitRepo(parent)
}

// getGitRoot finds the root directory of the Git repository
func getGitRoot(dir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = dir

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	root := string(output)
	// Remove trailing newline (works for both \n and \r\n)
	if len(root) > 0 && root[len(root)-1] == '\n' {
		root = root[:len(root)-1]
	}
	if len(root) > 0 && root[len(root)-1] == '\r' {
		root = root[:len(root)-1]
	}

	return root, nil
}
