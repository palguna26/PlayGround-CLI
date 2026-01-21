package agent

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/yourusername/playground/internal/patch"
)

// handleCommand processes in-chat commands
func (cs *ChatSession) handleCommand(input string) error {
	cmd := strings.ToLower(strings.TrimSpace(input))

	switch cmd {
	case "review":
		return cs.handleReview()
	case "apply":
		return cs.handleApply()
	case "status":
		return cs.handleStatus()
	case "help":
		return cs.handleHelp()
	case "exit", "quit":
		return cs.handleExit()
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
}

// handleReview shows all pending patches
func (cs *ChatSession) handleReview() error {
	if len(cs.Session.PendingPatches) == 0 {
		fmt.Println("No pending patches to review.")
		return nil
	}

	fmt.Printf("\n═══════════════════════════════════════\n")
	fmt.Printf("  Pending Patches: %d\n", len(cs.Session.PendingPatches))
	fmt.Printf("═══════════════════════════════════════\n\n")

	for i, p := range cs.Session.PendingPatches {
		fmt.Printf("═══ Patch %d/%d ═══\n", i+1, len(cs.Session.PendingPatches))
		fmt.Printf("File: %s\n", p.FilePath)
		fmt.Printf("Created: %s\n\n", p.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println(p.UnifiedDiff)
		fmt.Println()
	}

	return nil
}

// handleApply applies pending patches with user confirmation
func (cs *ChatSession) handleApply() error {
	if len(cs.Session.PendingPatches) == 0 {
		fmt.Println("No pending patches to apply.")
		return nil
	}

	fmt.Printf("\nAbout to apply %d patch(es) to the repository.\n", len(cs.Session.PendingPatches))
	fmt.Print("Apply all patches? [y/N]: ")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	response = strings.TrimSpace(strings.ToLower(response))
	if response != "y" && response != "yes" {
		fmt.Println("Patch application cancelled.")
		return nil
	}

	// Apply patches
	applied := 0
	for i, p := range cs.Session.PendingPatches {
		fmt.Printf("Applying patch %d/%d: %s... ", i+1, len(cs.Session.PendingPatches), p.FilePath)

		patchToApply := patch.Patch{
			FilePath:    p.FilePath,
			UnifiedDiff: p.UnifiedDiff,
		}

		if err := patch.Apply(cs.Agent.RepoRoot, patchToApply); err != nil {
			fmt.Printf("❌ FAILED\n")
			fmt.Printf("Error: %v\n", err)
			fmt.Printf("\nApplied %d/%d patches before failure.\n", applied, len(cs.Session.PendingPatches))
			return fmt.Errorf("patch application failed")
		}

		fmt.Printf("✓\n")
		applied++
	}

	// Clear pending patches
	cs.Session.PendingPatches = nil
	if err := cs.Store.Save(cs.Session); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	fmt.Printf("\n✅ Successfully applied %d patch(es)\n", applied)
	return nil
}

// handleStatus displays current session status
func (cs *ChatSession) handleStatus() error {
	fmt.Printf("\n═══════════════════════════════════════\n")
	fmt.Printf("  Session Status\n")
	fmt.Printf("═══════════════════════════════════════\n\n")

	fmt.Printf("Session ID: %s\n", cs.Session.ID)
	fmt.Printf("Goal: %s\n", cs.Session.Goal)
	fmt.Printf("Repository: %s\n", cs.Session.Repo)
	fmt.Printf("Created: %s\n\n", cs.Session.CreatedAt.Format("2006-01-02 15:04:05"))

	fmt.Printf("Pending Patches: %d\n", len(cs.Session.PendingPatches))
	fmt.Printf("Tool Calls: %d\n", len(cs.Session.ToolHistory))

	// Show recent tool calls
	if len(cs.Session.ToolHistory) > 0 {
		fmt.Println("\nRecent Tool Calls:")
		start := len(cs.Session.ToolHistory) - 5
		if start < 0 {
			start = 0
		}

		for i := start; i < len(cs.Session.ToolHistory); i++ {
			call := cs.Session.ToolHistory[i]
			status := "✓"
			if call.Error != "" {
				status = "✗"
			}
			fmt.Printf("  %s %s - %s\n", status, call.ToolName, call.Timestamp.Format("15:04:05"))
		}
	}

	if cs.Session.ContextSummary != "" {
		fmt.Printf("\nContext: %s\n", cs.Session.ContextSummary)
	}

	return nil
}

// handleHelp displays available commands
func (cs *ChatSession) handleHelp() error {
	fmt.Println("\n═══════════════════════════════════════")
	fmt.Println("  Available Commands")
	fmt.Println("═══════════════════════════════════════\n")
	fmt.Println("  review   - Show all pending patches as diffs")
	fmt.Println("  apply    - Apply pending patches to files")
	fmt.Println("  status   - Display current session status")
	fmt.Println("  help     - Show this help message")
	fmt.Println("  exit     - Exit agent mode and save session")
	fmt.Println()
	fmt.Println("Or just chat naturally! The agent will help you code.")

	return nil
}

// handleExit exits the chat session
func (cs *ChatSession) handleExit() error {
	fmt.Println("\n👋 Exiting agent mode...")

	// Save session one final time
	if err := cs.Store.Save(cs.Session); err != nil {
		fmt.Printf("Warning: failed to save session: %v\n", err)
	}

	fmt.Printf("Session %s saved. Resume with: pg agent --resume %s\n", cs.Session.ID, cs.Session.ID)
	cs.running = false

	return nil
}
