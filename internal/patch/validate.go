package patch

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Patch represents a proposed code change (defined in session package, but used here)
type Patch struct {
	FilePath    string
	UnifiedDiff string
}

// Validate checks if a patch is valid and can be applied
func Validate(repoRoot string, p Patch) error {
	// Parse the unified diff to extract metadata
	targetFile := filepath.Join(repoRoot, p.FilePath)

	// Check if this is a new file creation
	isNewFile := strings.Contains(p.UnifiedDiff, "--- /dev/null")

	if !isNewFile {
		// File must exist for modification
		if _, err := os.Stat(targetFile); os.IsNotExist(err) {
			return fmt.Errorf("target file does not exist: %s", p.FilePath)
		}

		// Verify context lines match current file state
		if err := validateContext(targetFile, p.UnifiedDiff); err != nil {
			return fmt.Errorf("context mismatch: %w", err)
		}
	}

	// Validate unified diff format
	if !strings.Contains(p.UnifiedDiff, "---") || !strings.Contains(p.UnifiedDiff, "+++") {
		return fmt.Errorf("invalid unified diff format: missing --- or +++ headers")
	}

	return nil
}

// validateContext checks if context lines in the diff match the current file
func validateContext(filePath string, diff string) error {
	// Read current file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	lines := strings.Split(string(content), "\n")

	// Parse diff to extract hunks
	scanner := bufio.NewScanner(strings.NewReader(diff))
	var currentLine int

	for scanner.Scan() {
		line := scanner.Text()

		// Parse hunk header: @@ -start,count +start,count @@
		if strings.HasPrefix(line, "@@") {
			// Extract old file line number
			parts := strings.Split(line, " ")
			if len(parts) >= 2 {
				oldPart := strings.TrimPrefix(parts[1], "-")
				if strings.Contains(oldPart, ",") {
					fmt.Sscanf(oldPart, "%d", &currentLine)
				} else {
					fmt.Sscanf(oldPart, "%d", &currentLine)
				}
				currentLine-- // Convert to 0-indexed
			}
			continue
		}

		// Context line (starts with space) or deletion (starts with -)
		if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "-") {
			contextContent := line[1:] // Remove prefix

			if currentLine >= 0 && currentLine < len(lines) {
				if lines[currentLine] != contextContent {
					// Context doesn't match - this means file has changed
					return fmt.Errorf("context line %d doesn't match: expected %q, got %q",
						currentLine+1, contextContent, lines[currentLine])
				}
			}

			// Only increment for context and deletions (not additions)
			if !strings.HasPrefix(line, "+") {
				currentLine++
			}
		} else if strings.HasPrefix(line, "+") {
			// Addition - doesn't affect current file line tracking
			continue
		}
	}

	return nil
}
