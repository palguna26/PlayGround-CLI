package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/playground/internal/session"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the current session status",
	Long: `Display information about the current active session including:
- Session ID
- Goal
- Number of pending patches
- Recent tool history

Example:
  pg status`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}

		// Find Git repository root
		repoRoot, err := getGitRoot(cwd)
		if err != nil {
			return fmt.Errorf("not in a git repository: %w", err)
		}

		// Create session store
		store, err := session.NewStore(repoRoot)
		if err != nil {
			return fmt.Errorf("failed to create session store: %w", err)
		}

		// Get active session ID
		sessionID, err := store.GetActiveSessionID()
		if err != nil {
			return fmt.Errorf("failed to get active session: %w", err)
		}

		if sessionID == "" {
			fmt.Println("No active session")
			fmt.Println("\nStart a new session with: pg start \"<goal>\"")
			return nil
		}

		// Load session
		sess, err := store.Load(sessionID)
		if err != nil {
			return fmt.Errorf("failed to load session: %w", err)
		}

		// Display session information
		fmt.Printf("Session: %s\n", sess.ID)
		fmt.Printf("Goal: %s\n", sess.Goal)
		fmt.Printf("Repository: %s\n", sess.Repo)
		fmt.Printf("Created: %s\n", sess.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("\nPending Patches: %d\n", len(sess.PendingPatches))
		fmt.Printf("Tool History: %d calls\n", len(sess.ToolHistory))

		// Show recent tool history (last 5)
		if len(sess.ToolHistory) > 0 {
			fmt.Println("\nRecent Tool Calls:")
			start := len(sess.ToolHistory) - 5
			if start < 0 {
				start = 0
			}

			for i := start; i < len(sess.ToolHistory); i++ {
				call := sess.ToolHistory[i]
				status := "✓"
				if call.Error != "" {
					status = "✗"
				}
				fmt.Printf("  %s %s - %s\n", status, call.ToolName, call.Timestamp.Format("15:04:05"))
			}
		}

		if sess.ContextSummary != "" {
			fmt.Printf("\nContext: %s\n", sess.ContextSummary)
		}

		return nil
	},
}
