package agent

import (
	"fmt"

	"github.com/yourusername/playground/internal/llm"
)

// Run executes the agent loop with the given user input
func (a *Agent) Run(userInput string, config AgentConfig) (string, error) {
	// Initialize conversation with system prompt
	messages := []llm.Message{
		{Role: "system", Content: getSystemPrompt(config.IsAgentMode)},
		{Role: "user", Content: userInput},
	}

	tools := defineTools()
	
	for iteration := 1; iteration <= config.MaxIterations; iteration++ {
		if config.Verbose {
			fmt.Printf("\n[Iteration %d/%d]\n", iteration, config.MaxIterations)
		}

		// Call LLM
		response, err := a.Provider.Chat(messages, tools)
		if err != nil {
			return "", fmt.Errorf("LLM error: %w", err)
		}

		// Check if done (stop condition or no tool calls)
		if response.FinishReason == "stop" || response.FinishReason == "end_turn" || len(response.ToolCalls) == 0 {
			if config.Verbose {
				fmt.Println("[Agent finished]")
			}
			return response.Content, nil
		}

		// Add assistant message if present
		if response.Content != "" {
			messages = append(messages, llm.Message{Role: "assistant", Content: response.Content})
		}

		// Execute tool calls
		if config.Verbose {
			fmt.Printf("[Executing %d tool(s)]\n", len(response.ToolCalls))
		}

		toolResults := make([]llm.Message, 0, len(response.ToolCalls))
		for _, toolCall := range response.ToolCalls {
			if config.Verbose {
				fmt.Printf("  - %s\n", toolCall.Name)
			}

			result, err := a.executeTool(toolCall)
			a.logToolCall(toolCall, result, err)

			// Prepare tool result
			content := result
			if err != nil {
				content = fmt.Sprintf("Error: %v", err)
			}

			toolResults = append(toolResults, llm.Message{
				Role:    "tool",
				Content: content,
				Name:    toolCall.Name,
			})
		}

		// Batch append tool results
		messages = append(messages, toolResults...)

		// Save session after iteration
		if err := a.Store.Save(a.Session); err != nil {
			return "", fmt.Errorf("failed to save session: %w", err)
		}
	}

	return "", fmt.Errorf("agent exceeded maximum iterations (%d)", config.MaxIterations)
}