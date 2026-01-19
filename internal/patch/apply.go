package patch

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Apply applies a validated patch to the repository
// This is the ONLY place in the codebase that modifies files
func Apply(repoRoot string, p Patch) error {
	// Validate first - safety check
	if err := Validate(repoRoot, p); err != nil {
		return fmt.Errorf("patch validation failed: %w", err)
	}

	targetFile := filepath.Join(repoRoot, p.FilePath)

	// Create backup for rollback
	var backupPath string
	var needsBackup bool

	if _, err := os.Stat(targetFile); err == nil {
		needsBackup = true
		backupPath = targetFile + ".backup"

		if err := copyFile(targetFile, backupPath); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
	}

	// Apply patch using the `patch` command
	// Write diff to temporary file
	tempDiff, err := os.CreateTemp("", "pg-patch-*.diff")
	if err != nil {
		cleanup(backupPath)
		return fmt.Errorf("failed to create temp diff file: %w", err)
	}
	defer os.Remove(tempDiff.Name())

	if _, err := tempDiff.WriteString(p.UnifiedDiff); err != nil {
		cleanup(backupPath)
		return fmt.Errorf("failed to write diff: %w", err)
	}
	tempDiff.Close()

	// Apply patch
	cmd := exec.Command("patch", "-p0", "-i", tempDiff.Name())
	cmd.Dir = repoRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Patch failed - restore backup
		if needsBackup {
			if restoreErr := copyFile(backupPath, targetFile); restoreErr != nil {
				return fmt.Errorf("patch failed and backup restoration failed: %w (original error: %v)",
					restoreErr, err)
			}
		}
		cleanup(backupPath)
		return fmt.Errorf("patch application failed: %w\nOutput: %s", err, string(output))
	}

	// Success - remove backup
	cleanup(backupPath)

	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, data, 0644)
}

// cleanup removes a file, ignoring errors
func cleanup(path string) {
	if path != "" {
		os.Remove(path)
	}
}
