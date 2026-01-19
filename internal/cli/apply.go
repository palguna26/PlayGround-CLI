package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/playground/internal/patch"
	"github.com/yourusername/playground/internal/session"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply pending patches to the repository",
	Long: `Apply all validated pending patches to the repository files.
Requires user confirmation before applying changes.
All patches are validated before application to ensure safety.

Example:
  pg apply`,
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

		// Check if there are patches to apply
		if len(sess.PendingPatches) == 0 {
			fmt.Println("No pending patches to apply")
			return nil
		}

		// Request user confirmation
		fmt.Printf("About to apply %d patch(es) to the repository.\n", len(sess.PendingPatches))
		fmt.Printf("Apply all patches? [y/N]: ")

		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read user input: %w", err)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Patch application cancelled")
			return nil
		}

		// Apply patches
		applied := 0
		for i, p := range sess.PendingPatches {
			fmt.Printf("Applying patch %d/%d: %s... ", i+1, len(sess.PendingPatches), p.FilePath)

			// Convert session.Patch to patch.Patch
			patchToApply := patch.Patch{
				FilePath:    p.FilePath,
				UnifiedDiff: p.UnifiedDiff,
			}

			if err := patch.Apply(repoRoot, patchToApply); err != nil {
				fmt.Printf("❌ FAILED\n")
				fmt.Printf("Error: %v\n", err)
				fmt.Printf("\nApplied %d/%d patches before failure.\n", applied, len(sess.PendingPatches))
				return fmt.Errorf("patch application failed")
			}

			fmt.Printf("✓\n")
			applied++
		}

		// Clear pending patches from session
		sess.PendingPatches = []session.Patch{}
		if err := store.Save(sess); err != nil {
			return fmt.Errorf("failed to save session: %w", err)
		}

		fmt.Printf("\n✓ Successfully applied %d patch(es)\n", applied)
		return nil
	},
}
