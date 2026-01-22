package llm

import (
	"context"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GeminiProvider implements the Provider interface using Google's Gemini API
type GeminiProvider struct {
	client *genai.Client
	model  string
}

// NewGeminiProvider creates a new Gemini provider
func NewGeminiProvider() (*GeminiProvider, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &GeminiProvider{
		client: client,
		model:  "gemini-2.5-flash", // Default model
	}, nil
}

// Name returns the provider name
func (p *GeminiProvider) Name() string {
	return "Gemini"
}

// Chat sends messages to Gemini and receives a response
func (p *GeminiProvider) Chat(messages []Message, tools []Tool) (*Response, error) {
	ctx := context.Background()
	model := p.client.GenerativeModel(p.model)

	// Configure model
	model.SetTemperature(0.7)

	// Convert tools to Gemini format
	if len(tools) > 0 {
		geminiTools := make([]*genai.Tool, 0)
		functionDecls := make([]*genai.FunctionDeclaration, len(tools))

		for i, tool := range tools {
			// Convert our parameters JSON schema to Gemini schema
			schema := &genai.Schema{
				Type: genai.TypeObject,
			}

			// Parse parameters if they exist
			if params, ok := tool.Parameters["properties"].(map[string]interface{}); ok {
				schema.Properties = make(map[string]*genai.Schema)
				for propName, propDef := range params {
					if propMap, ok := propDef.(map[string]interface{}); ok {
						propSchema := &genai.Schema{}
						if propType, ok := propMap["type"].(string); ok {
							switch propType {
							case "string":
								propSchema.Type = genai.TypeString
							case "object":
								propSchema.Type = genai.TypeObject
							case "array":
								propSchema.Type = genai.TypeArray
							}
						}
						if desc, ok := propMap["description"].(string); ok {
							propSchema.Description = desc
						}
						schema.Properties[propName] = propSchema
					}
				}
			}

			// Set required fields
			if required, ok := tool.Parameters["required"].([]string); ok {
				schema.Required = required
			} else if required, ok := tool.Parameters["required"].([]interface{}); ok {
				requiredStrs := make([]string, len(required))
				for i, r := range required {
					if s, ok := r.(string); ok {
						requiredStrs[i] = s
					}
				}
				schema.Required = requiredStrs
			}

			functionDecls[i] = &genai.FunctionDeclaration{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters:  schema,
			}
		}

		geminiTools = append(geminiTools, &genai.Tool{
			FunctionDeclarations: functionDecls,
		})
		model.Tools = geminiTools
	}

	// Start chat session
	chat := model.StartChat()

	// Convert message history to Gemini format
	// Skip system message for now as Gemini handles it differently
	var chatHistory []*genai.Content
	var systemInstruction string

	for _, msg := range messages {
		if msg.Role == "system" {
			systemInstruction = msg.Content
			continue
		}

		role := msg.Role
		if role == "assistant" {
			role = "model"
		}

		content := &genai.Content{
			Role: role,
			Parts: []genai.Part{
				genai.Text(msg.Content),
			},
		}

		chatHistory = append(chatHistory, content)
	}

	// Set system instruction if present
	if systemInstruction != "" {
		model.SystemInstruction = &genai.Content{
			Parts: []genai.Part{genai.Text(systemInstruction)},
		}
	}

	// Set chat history (exclude last user message as we'll send it separately)
	if len(chatHistory) > 1 {
		chat.History = chatHistory[:len(chatHistory)-1]
	}

	// Get last message to send
	var lastMessage string
	if len(chatHistory) > 0 {
		if textPart, ok := chatHistory[len(chatHistory)-1].Parts[0].(genai.Text); ok {
			lastMessage = string(textPart)
		}
	}

	// Send message and get response
	resp, err := chat.SendMessage(ctx, genai.Text(lastMessage))
	if err != nil {
		return nil, fmt.Errorf("Gemini API error: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no response from Gemini")
	}

	candidate := resp.Candidates[0]

	// Build response
	response := &Response{
		Content:      "",
		ToolCalls:    []ToolCall{},
		FinishReason: string(candidate.FinishReason),
	}

	// Extract content and tool calls
	for _, part := range candidate.Content.Parts {
		switch v := part.(type) {
		case genai.Text:
			response.Content += string(v)
		case genai.FunctionCall:
			// Convert function call arguments to our format
			args := make(map[string]interface{})
			for key, val := range v.Args {
				args[key] = val
			}

			response.ToolCalls = append(response.ToolCalls, ToolCall{
				ID:        v.Name, // Gemini doesn't provide separate ID
				Name:      v.Name,
				Arguments: args,
			})
		}
	}

	return response, nil
}

// Close closes the Gemini client
func (p *GeminiProvider) Close() error {
	return p.client.Close()
}
