package tools

import (
	"fmt"
	"os/exec"
	"strings"
)

// GitStatus returns the output of `git status`
func GitStatus(repoRoot string) (string, error) {
	cmd := exec.Command("git", "status")
	cmd.Dir = repoRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git status failed: %w\nOutput: %s", err, string(output))
	}

	return strings.TrimSpace(string(output)), nil
}

// GitDiff returns the output of `git diff`
func GitDiff(repoRoot string) (string, error) {
	cmd := exec.Command("git", "diff")
	cmd.Dir = repoRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git diff failed: %w\nOutput: %s", err, string(output))
	}

	return string(output), nil
}

// GitDiffStaged returns the output of `git diff --staged`
func GitDiffStaged(repoRoot string) (string, error) {
	cmd := exec.Command("git", "diff", "--staged")
	cmd.Dir = repoRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git diff --staged failed: %w\nOutput: %s", err, string(output))
	}

	return string(output), nil
}
