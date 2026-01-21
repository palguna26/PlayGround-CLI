package llm

import (
	"context"
	"fmt"
	"io"

	"github.com/google/generative-ai-go/genai"
)

// ChatStream implements streaming response for Gemini
func (p *GeminiProvider) ChatStream(messages []Message, tools []Tool) (<-chan StreamChunk, error) {
	chunkChan := make(chan StreamChunk)

	ctx := context.Background()
	model := p.client.GenerativeModel(p.model)

	// Configure model
	model.SetTemperature(0.7)

	// Convert tools to Gemini format
	if len(tools) > 0 {
		geminiTools := make([]*genai.Tool, 0)
		functionDecls := make([]*genai.FunctionDeclaration, len(tools))

		for i, tool := range tools {
			schema := &genai.Schema{
				Type: genai.TypeObject,
			}

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

	// Convert message history
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

	if systemInstruction != "" {
		model.SystemInstruction = &genai.Content{
			Parts: []genai.Part{genai.Text(systemInstruction)},
		}
	}

	if len(chatHistory) > 1 {
		chat.History = chatHistory[:len(chatHistory)-1]
	}

	var lastMessage string
	if len(chatHistory) > 0 {
		if textPart, ok := chatHistory[len(chatHistory)-1].Parts[0].(genai.Text); ok {
			lastMessage = string(textPart)
		}
	}

	// Start streaming in goroutine
	go func() {
		defer close(chunkChan)

		iter := chat.SendMessageStream(ctx, genai.Text(lastMessage))

		for {
			resp, err := iter.Next()
			if err == io.EOF {
				break
			}

			if err != nil {
				chunkChan <- StreamChunk{Error: fmt.Errorf("Gemini stream error: %w", err)}
				return
			}

			if len(resp.Candidates) == 0 {
				continue
			}

			candidate := resp.Candidates[0]

			// Extract content and tool calls
			for _, part := range candidate.Content.Parts {
				switch v := part.(type) {
				case genai.Text:
					chunkChan <- StreamChunk{Content: string(v)}
				case genai.FunctionCall:
					args := make(map[string]interface{})
					for key, val := range v.Args {
						args[key] = val
					}

					chunkChan <- StreamChunk{
						ToolCall: &ToolCall{
							ID:        v.Name,
							Name:      v.Name,
							Arguments: args,
						},
					}
				}
			}

			// Handle finish reason
			if candidate.FinishReason != 0 {
				chunkChan <- StreamChunk{FinishReason: string(candidate.FinishReason)}
			}
		}
	}()

	return chunkChan, nil
}
