package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/playground/internal/agent"
	"github.com/yourusername/playground/internal/llm"
	"github.com/yourusername/playground/internal/session"
)

var askCmd = &cobra.Command{
	Use:   "ask [question]",
	Short: "Ask the AI agent a question or give it a task",
	Long: `Invoke the AI agent with a question or task within the current session.
The agent can use tools to read files, check git status, and propose code changes.

Example:
  pg ask "what auth do we currently use?"
  pg ask "add jwt middleware"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		question := args[0]

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
			return fmt.Errorf("no active session. Start one with: pg start \"<goal>\"")
		}

		// Load session
		sess, err := store.Load(sessionID)
		if err != nil {
			return fmt.Errorf("failed to load session: %w", err)
		}

		// Load config to get model path
		config, err := LoadConfig()
		if err != nil || config == nil || config.ModelPath == "" {
			return fmt.Errorf("no model configured. Run: pg setup")
		}

		// Create local LLM provider
		provider, err := llm.NewLocalProvider(config.ModelPath)
		if err != nil {
			return fmt.Errorf("failed to load local model: %w", err)
		}

		fmt.Printf("Using: %s\n", provider.Name())

		// Create agent
		agentInstance := &agent.Agent{
			Session:  sess,
			Store:    store,
			Provider: provider,
			RepoRoot: repoRoot,
		}

		// Run agent
		fmt.Printf("🤖 Agent working...\n\n")

		response, err := agentInstance.Run(question, agent.DefaultConfig)
		if err != nil {
			return fmt.Errorf("agent error: %w", err)
		}

		// Display response
		fmt.Printf("Agent: %s\n", response)

		// Save final session state
		if err := store.Save(sess); err != nil {
			return fmt.Errorf("failed to save session: %w", err)
		}

		return nil
	},
}
