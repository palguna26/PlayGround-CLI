package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/playground/internal/session"
)

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Review pending patches in the current session",
	Long: `Display all pending patches that have been proposed by the agent.
Shows the unified diff for each patch so you can review changes before applying.

Example:
  pg review`,
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

		// Get active session
		sessionID, err := store.GetActiveSessionID()
		if err != nil {
			return fmt.Errorf("failed to get active session: %w", err)
		}

		if sessionID == "" {
			return fmt.Errorf("no active session")
		}

		// Load session
		sess, err := store.Load(sessionID)
		if err != nil {
			return fmt.Errorf("failed to load session: %w", err)
		}

		// Display patches
		if len(sess.PendingPatches) == 0 {
			fmt.Println("No pending patches")
			return nil
		}

		fmt.Printf("Session: %s\n", sess.ID)
		fmt.Printf("Pending patches: %d\n\n", len(sess.PendingPatches))

		for i, patch := range sess.PendingPatches {
			fmt.Printf("═══ Patch %d/%d ═══\n", i+1, len(sess.PendingPatches))
			fmt.Printf("File: %s\n", patch.FilePath)
			fmt.Printf("Created: %s\n\n", patch.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Println(patch.UnifiedDiff)
			fmt.Println()
		}

		return nil
	},
}
