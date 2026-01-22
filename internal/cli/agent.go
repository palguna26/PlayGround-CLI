package cli

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourusername/playground/internal/agent"
	"github.com/yourusername/playground/internal/llm"
	"github.com/yourusername/playground/internal/session"
	"github.com/yourusername/playground/internal/workspace"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Start interactive agent chat mode",
	Long: `Launch an interactive chat session with the PlayGround agent.

This provides a conversational interface similar to Claude Code, but with
PlayGround's safety guarantees: diff-only changes, explicit review/apply,
and full control over every modification.

Git is optional - if not in a Git repository, PlayGround will use its
own snapshot system for versioning.

Example:
  pg agent                    # Start new agent session
  pg agent --resume pg-5      # Resume previous session

In agent mode, you can:
  • Chat naturally with the AI
  • Type 'review' to see proposed changes
  • Type 'apply' to accept changes
  • Type 'undo' to rollback changes
  • Type 'status' for session info
  • Type 'exit' to quit`,
	RunE: func(cmd *cobra.Command, args []string) error {
		resumeSession, _ := cmd.Flags().GetString("resume")

		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}

		// Create workspace (auto-detects Git or uses snapshots)
		ws, err := workspace.NewWorkspace(cwd)
		if err != nil {
			return fmt.Errorf("failed to initialize workspace: %w", err)
		}

		workspaceRoot := ws.GetRoot()

		if ws.IsGitBacked() {
			fmt.Println("📂 Using Git for version control")
		} else {
			fmt.Println("📂 Using snapshot-based version control (no Git)")
		}

		// Create session store
		store, err := session.NewStore(workspaceRoot)
		if err != nil {
			return fmt.Errorf("failed to create session store: %w", err)
		}

		var sess *session.Session

		if resumeSession != "" {
			// Resume existing session
			sess, err = store.Load(resumeSession)
			if err != nil {
				return fmt.Errorf("failed to load session %s: %w", resumeSession, err)
			}

			if sess.Repo != workspaceRoot {
				return fmt.Errorf("session %s is for repository %s, but you are in %s",
					resumeSession, sess.Repo, workspaceRoot)
			}

			store.SetActiveSessionID(resumeSession)
		} else {
			// Create new session
			sessionID, err := store.GenerateSessionID()
			if err != nil {
				return fmt.Errorf("failed to generate session ID: %w", err)
			}

			// Prompt for goal
			fmt.Print("What's your goal for this session? ")
			scanner := bufio.NewScanner(os.Stdin)
			var goal string
			if scanner.Scan() {
				goal = scanner.Text()
			}

			if goal == "" {
				goal = "Interactive coding session"
			}

			sess = &session.Session{
				ID:             sessionID,
				Repo:           workspaceRoot,
				Goal:           goal,
				ContextSummary: "",
				PendingPatches: []session.Patch{},
				ToolHistory:    []session.ToolCall{},
				CreatedAt:      time.Now(),
			}

			if err := store.Save(sess); err != nil {
				return fmt.Errorf("failed to save session: %w", err)
			}

			store.SetActiveSessionID(sessionID)
		}

		// Create LLM provider
		provider, err := llm.NewProvider()
		if err != nil {
			return fmt.Errorf("failed to create LLM provider: %w", err)
		}

		fmt.Printf("Using LLM provider: %s\n", provider.Name())

		// Create agent with agent mode prompt
		agentInstance := &agent.Agent{
			Session:  sess,
			Store:    store,
			Provider: provider,
			RepoRoot: workspaceRoot,
		}

		// Override system prompt for agent mode
		// (This will require modifying the agent loop to accept mode)

		// Create and run chat session
		chatSession := agent.NewChatSession(agentInstance, store)
		return chatSession.Run()
	},
}

func init() {
	agentCmd.Flags().String("resume", "", "Resume a previous agent session by ID")
}
