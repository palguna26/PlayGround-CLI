package tools

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunCommand executes a shell command with explicit user approval
// This is a security-critical function - commands must be approved
func RunCommand(repoRoot, command string) (string, error) {
	// Display command to user and request approval
	fmt.Printf("\n⚠️  The agent wants to run this command:\n")
	fmt.Printf("   %s\n\n", command)
	fmt.Printf("Allow this command? [y/N]: ")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read user input: %w", err)
	}

	response = strings.TrimSpace(strings.ToLower(response))
	if response != "y" && response != "yes" {
		return "", fmt.Errorf("command rejected by user")
	}

	// Execute command
	// Note: This is intentionally simple - runs via shell
	// Security: User has explicitly approved this specific command
	var cmd *exec.Cmd
	if strings.Contains(command, " ") {
		// Command with arguments - use sh -c
		cmd = exec.Command("sh", "-c", command)
	} else {
		cmd = exec.Command(command)
	}

	cmd.Dir = repoRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("command failed: %w", err)
	}

	return string(output), nil
}
