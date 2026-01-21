package agent

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/yourusername/playground/internal/session"
)

// ChatSession represents an interactive agent chat session
type ChatSession struct {
	Agent    *Agent
	Session  *session.Session
	Store    *session.Store
	Messages []string // Chat history for display
	running  bool
}

// NewChatSession creates a new interactive chat session
func NewChatSession(agent *Agent, store *session.Store) *ChatSession {
	return &ChatSession{
		Agent:    agent,
		Session:  agent.Session,
		Store:    store,
		Messages: []string{},
		running:  true,
	}
}

// Run starts the interactive chat loop
func (cs *ChatSession) Run() error {
	cs.displayWelcome()

	reader := bufio.NewScanner(os.Stdin)

	for cs.running {
		// Display prompt
		fmt.Print("\nYou: ")

		// Read user input
		if !reader.Scan() {
			break
		}

		input := strings.TrimSpace(reader.Text())

		if input == "" {
			continue
		}

		// Check if it's a command
		if cs.isCommand(input) {
			if err := cs.handleCommand(input); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
			continue
		}

		// Send to agent with streaming
		cs.Messages = append(cs.Messages, "You: "+input)

		fmt.Print("\nAgent: ")

		// Create channel for streaming output
		outputChan := make(chan string, 10)
		var fullResponse string

		// Start streaming in goroutine
		go func() {
			cs.Agent.RunStreaming(input, AgentModeConfig, outputChan)
		}()

		// Display streaming output as it arrives
		for chunk := range outputChan {
			fmt.Print(chunk)
			fullResponse += chunk
		}

		cs.Messages = append(cs.Messages, "Agent: "+fullResponse)

		// If agent proposed patches, prompt for review
		if len(cs.Session.PendingPatches) > 0 {
			fmt.Println("\n💡 Type 'review' to see the changes, or 'apply' to accept them.")
		}

		// Save session after each interaction
		if err := cs.Store.Save(cs.Session); err != nil {
			fmt.Printf("Warning: failed to save session: %v\n", err)
		}
	}

	return nil
}

// displayWelcome shows the welcome message
func (cs *ChatSession) displayWelcome() {
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║           PlayGround Agent - Interactive Mode              ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("Session: %s\n", cs.Session.ID)
	fmt.Printf("Goal: %s\n", cs.Session.Goal)
	fmt.Println()
	fmt.Println("Available commands:")
	fmt.Println("  review   - Show pending patches")
	fmt.Println("  apply    - Apply pending patches")
	fmt.Println("  status   - Show session status")
	fmt.Println("  exit     - Exit agent mode")
	fmt.Println()
	fmt.Println("Just type naturally to chat with the agent.")
	fmt.Println("────────────────────────────────────────────────────────────")
}

// isCommand checks if input is a command
func (cs *ChatSession) isCommand(input string) bool {
	commands := []string{"review", "apply", "status", "exit", "quit", "help"}

	for _, cmd := range commands {
		if strings.EqualFold(input, cmd) {
			return true
		}
	}

	return false
}
