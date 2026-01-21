package agent

import (
	"fmt"

	"github.com/yourusername/playground/internal/llm"
)

// Run executes the agent loop with the given user input
func (a *Agent) Run(userInput string, config AgentConfig) (string, error) {
	// Initialize conversation with system prompt
	messages := []llm.Message{
		{
			Role:    "system",
			Content: getSystemPrompt(config.IsAgentMode),
		},
		{
			Role:    "user",
			Content: userInput,
		},
	}

	tools := defineTools()
	iteration := 0

	for {
		iteration++

		// Hard stop condition: max iterations
		if iteration > config.MaxIterations {
			return "", fmt.Errorf("agent exceeded maximum iterations (%d)", config.MaxIterations)
		}

		if config.Verbose {
			fmt.Printf("\n[Iteration %d/%d]\n", iteration, config.MaxIterations)
		}

		// Call LLM
		response, err := a.Provider.Chat(messages, tools)
		if err != nil {
			return "", fmt.Errorf("LLM error: %w", err)
		}

		// Check finish reason
		if response.FinishReason == "stop" || response.FinishReason == "end_turn" {
			// LLM is done - return final message
			if config.Verbose {
				fmt.Println("[Agent finished]")
			}
			return response.Content, nil
		}

		// Check for tool calls
		if len(response.ToolCalls) == 0 {
			// No tool calls and not finished - this is unusual but treat as completion
			if config.Verbose {
				fmt.Println("[No tool calls, finishing]")
			}
			return response.Content, nil
		}

		// Execute tool calls
		if config.Verbose {
			fmt.Printf("[Executing %d tool(s)]\n", len(response.ToolCalls))
		}

		// Add assistant message to history
		if response.Content != "" {
			messages = append(messages, llm.Message{
				Role:    "assistant",
				Content: response.Content,
			})
		}

		// Execute each tool call and collect results
		for _, toolCall := range response.ToolCalls {
			if config.Verbose {
				fmt.Printf("  - %s\n", toolCall.Name)
			}

			result, err := a.executeTool(toolCall)

			// Log tool execution
			a.logToolCall(toolCall, result, err)

			// Add tool result to conversation
			var toolResultContent string
			if err != nil {
				toolResultContent = fmt.Sprintf("Error: %v", err)
			} else {
				toolResultContent = result
			}

			messages = append(messages, llm.Message{
				Role:    "tool",
				Content: toolResultContent,
				Name:    toolCall.Name,
			})
		}

		// Save session after each iteration (preserves state even if agent crashes)
		if err := a.Store.Save(a.Session); err != nil {
			return "", fmt.Errorf("failed to save session: %w", err)
		}

		// Continue loop
	}
}
