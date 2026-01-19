package cli

import (
	"github.com/spf13/cobra"
)

const version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:   "pg",
	Short: "PlayGround - A session-based CLI for AI-assisted coding",
	Long: `PlayGround (pg) is a local-first, session-based CLI that allows AI agents 
to safely reason about and modify code using diff-only, reviewable patches.

This is an open runtime for coding agents, not a chat app or IDE replacement.`,
	Version: version,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Register subcommands
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(askCmd)
	rootCmd.AddCommand(reviewCmd)
	rootCmd.AddCommand(applyCmd)
	rootCmd.AddCommand(resumeCmd)
}
