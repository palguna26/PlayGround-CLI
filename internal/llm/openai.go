package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

// OpenAIProvider implements the Provider interface using OpenAI's API
type OpenAIProvider struct {
	client *openai.Client
	model  string
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider() (*OpenAIProvider, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	client := openai.NewClient(apiKey)

	return &OpenAIProvider{
		client: client,
		model:  "gpt-4", // Default model, can be configured
	}, nil
}

// Name returns the provider name
func (p *OpenAIProvider) Name() string {
	return "OpenAI"
}

// Chat sends messages to OpenAI and receives a response
func (p *OpenAIProvider) Chat(messages []Message, tools []Tool) (*Response, error) {
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

	// Make API call
	resp, err := p.client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	choice := resp.Choices[0]

	// Convert response
	response := &Response{
		Content:      choice.Message.Content,
		FinishReason: string(choice.FinishReason),
		ToolCalls:    []ToolCall{},
	}

	// Parse tool calls if present
	for _, tc := range choice.Message.ToolCalls {
		var args map[string]interface{}
		if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
			return nil, fmt.Errorf("failed to parse tool arguments: %w", err)
		}

		response.ToolCalls = append(response.ToolCalls, ToolCall{
			ID:        tc.ID,
			Name:      tc.Function.Name,
			Arguments: args,
		})
	}

	return response, nil
}
