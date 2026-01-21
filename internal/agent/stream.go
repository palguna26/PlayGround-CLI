package agent

import (
	"fmt"
	"strings"
	"time"

	"github.com/yourusername/playground/internal/llm"
)

// RunStreaming executes the agent loop with streaming responses
func (a *Agent) RunStreaming(userInput string, config AgentConfig, outputChan chan<- string) error {
	defer close(outputChan)

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

		// Hard stop condition
		if iteration > config.MaxIterations {
			outputChan <- fmt.Sprintf("\n[Error: exceeded max iterations (%d)]", config.MaxIterations)
			return fmt.Errorf("agent exceeded maximum iterations (%d)", config.MaxIterations)
		}

		// Get streaming response
		streamChan, err := a.Provider.ChatStream(messages, tools)
		if err != nil {
			outputChan <- fmt.Sprintf("\n[Error: %v]", err)
			return fmt.Errorf("LLM error: %w", err)
		}

		var fullContent strings.Builder
		var toolCalls []llm.ToolCall
		var finishReason string

		// Process streaming chunks
		for chunk := range streamChan {
			if chunk.Error != nil {
				outputChan <- fmt.Sprintf("\n[Error: %v]", chunk.Error)
				return chunk.Error
			}

			if chunk.Content != "" {
				outputChan <- chunk.Content
				fullContent.WriteString(chunk.Content)
			}

			if chunk.ToolCall != nil {
				toolCalls = append(toolCalls, *chunk.ToolCall)
			}

			if chunk.FinishReason != "" {
				finishReason = chunk.FinishReason
			}
		}

		// newline after streaming complete
		if fullContent.Len() > 0 {
			outputChan <- "\n"
		}

		// Check finish reason
		if finishReason == "stop" || finishReason == "end_turn" {
			return nil
		}

		// Check for tool calls
		if len(toolCalls) == 0 {
			return nil
		}

		// Add assistant message to history
		if fullContent.Len() > 0 {
			messages = append(messages, llm.Message{
				Role:    "assistant",
				Content: fullContent.String(),
			})
		}

		// Execute tool calls
		for _, toolCall := range toolCalls {
			result, err := a.executeTool(toolCall)
			a.logToolCall(toolCall, result, err)

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

		// Save session after each iteration
		if err := a.Store.Save(a.Session); err != nil {
			return fmt.Errorf("failed to save session: %w", err)
		}

		// Small delay before next iteration
		time.Sleep(100 * time.Millisecond)
	}
}
