package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/playground/internal/session"
)

var resumeCmd = &cobra.Command{
	Use:   "resume [session-id]",
	Short: "Resume a previous session",
	Long: `Resume a previous PlayGround session by its ID.
Loads the session state and makes it the active session.

Example:
  pg resume pg-12`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sessionID := args[0]

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

		// Load the specified session
		sess, err := store.Load(sessionID)
		if err != nil {
			return fmt.Errorf("failed to load session: %w", err)
		}

		// Verify the session's repo matches current repo
		if sess.Repo != repoRoot {
			return fmt.Errorf("session %s is for repository %s, but you are in %s",
				sessionID, sess.Repo, repoRoot)
		}

		// Set as active session
		if err := store.SetActiveSessionID(sessionID); err != nil {
			return fmt.Errorf("failed to set active session: %w", err)
		}

		fmt.Printf("✓ Resumed session: %s\n", sessionID)
		fmt.Printf("  Goal: %s\n", sess.Goal)
		fmt.Printf("  Pending patches: %d\n", len(sess.PendingPatches))

		return nil
	},
}
