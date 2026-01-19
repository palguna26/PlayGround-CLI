package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourusername/playground/internal/session"
)

var startCmd = &cobra.Command{
	Use:   "start [goal]",
	Short: "Start a new coding session with a goal",
	Long: `Start a new PlayGround session in the current repository.
This creates a new session with the specified goal and generates a unique session ID.

Example:
  pg start "add jwt auth to api"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		goal := args[0]

		// Get current working directory
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}

		// Verify we're in Git repository
		if !isGitRepo(cwd) {
			return fmt.Errorf("not a git repository (or any of the parent directories)")
		}

		// Find Git repository root
		repoRoot, err := getGitRoot(cwd)
		if err != nil {
			return fmt.Errorf("failed to find git repository root: %w", err)
		}

		// Create session store
		store, err := session.NewStore(repoRoot)
		if err != nil {
			return fmt.Errorf("failed to create session store: %w", err)
		}

		// Generate new session ID
		sessionID, err := store.GenerateSessionID()
		if err != nil {
			return fmt.Errorf("failed to generate session ID: %w", err)
		}

		// Create new session
		newSession := &session.Session{
			ID:             sessionID,
			Repo:           repoRoot,
			Goal:           goal,
			ContextSummary: "",
			PendingPatches: []session.Patch{},
			ToolHistory:    []session.ToolCall{},
			CreatedAt:      time.Now(),
		}

		// Save session
		if err := store.Save(newSession); err != nil {
			return fmt.Errorf("failed to save session: %w", err)
		}

		// Set as active session
		if err := store.SetActiveSessionID(sessionID); err != nil {
			return fmt.Errorf("failed to set active session: %w", err)
		}

		fmt.Printf("✓ Started new session: %s\n", sessionID)
		fmt.Printf("  Goal: %s\n", goal)
		fmt.Printf("  Repo: %s\n", repoRoot)

		return nil
	},
}
