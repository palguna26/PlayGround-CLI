package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/sashabaranov/go-openai"
)

// ChatStream implements streaming response for OpenAI
func (p *OpenAIProvider) ChatStream(messages []Message, tools []Tool) (<-chan StreamChunk, error) {
	chunkChan := make(chan StreamChunk)

	// Convert our messages to OpenAI format
	chatMessages := make([]openai.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		chatMessages[i] = openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
			Name:    msg.Name,
		}
	}

	req := openai.ChatCompletionRequest{
		Model:    p.model,
		Messages: chatMessages,
		Stream:   true, // Enable streaming
	}

	// Add tools if provided
	if len(tools) > 0 {
		oaiTools := make([]openai.Tool, len(tools))
		for i, tool := range tools {
			oaiTools[i] = openai.Tool{
				Type: openai.ToolTypeFunction,
				Function: &openai.FunctionDefinition{
					Name:        tool.Name,
					Description: tool.Description,
					Parameters:  tool.Parameters,
				},
			}
		}
		req.Tools = oaiTools
	}

	// Start streaming in goroutine
	go func() {
		defer close(chunkChan)

		ctx := context.Background()
		stream, err := p.client.CreateChatCompletionStream(ctx, req)
		if err != nil {
			chunkChan <- StreamChunk{Error: fmt.Errorf("OpenAI stream error: %w", err)}
			return
		}
		defer stream.Close()

		var currentToolCall *ToolCall
		var toolCallArgs string

		for {
			response, err := stream.Recv()
			if err == io.EOF {
				// Stream finished
				if currentToolCall != nil {
					// Parse accumulated tool call
					var args map[string]interface{}
					if err := json.Unmarshal([]byte(toolCallArgs), &args); err == nil {
						currentToolCall.Arguments = args
						chunkChan <- StreamChunk{ToolCall: currentToolCall}
					}
				}
				break
			}

			if err != nil {
				chunkChan <- StreamChunk{Error: fmt.Errorf("stream receive error: %w", err)}
				return
			}

			if len(response.Choices) == 0 {
				continue
			}

			delta := response.Choices[0].Delta

			// Handle text content
			if delta.Content != "" {
				chunkChan <- StreamChunk{Content: delta.Content}
			}

			// Handle tool calls
			if len(delta.ToolCalls) > 0 {
				for _, tc := range delta.ToolCalls {
					if tc.Function.Name != "" {
						// New tool call
						if currentToolCall != nil {
							// Finish previous tool call
							var args map[string]interface{}
							if err := json.Unmarshal([]byte(toolCallArgs), &args); err == nil {
								currentToolCall.Arguments = args
								chunkChan <- StreamChunk{ToolCall: currentToolCall}
							}
						}
						currentToolCall = &ToolCall{
							ID:   tc.ID,
							Name: tc.Function.Name,
						}
						toolCallArgs = tc.Function.Arguments
					} else if currentToolCall != nil {
						// Accumulate arguments for current tool call
						toolCallArgs += tc.Function.Arguments
					}
				}
			}

			// Handle finish reason
			if response.Choices[0].FinishReason != "" {
				chunkChan <- StreamChunk{FinishReason: string(response.Choices[0].FinishReason)}
			}
		}
	}()

	return chunkChan, nil
}
