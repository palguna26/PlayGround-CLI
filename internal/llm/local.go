package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// LocalProvider implements the Provider interface using llama.cpp CLI
type LocalProvider struct {
	modelPath   string
	contextSize int
	debugMode   bool
	logDir      string
}

// NewLocalProvider creates a new local LLM provider
func NewLocalProvider(modelPath string) (*LocalProvider, error) {
	if modelPath == "" {
		return nil, fmt.Errorf("model path is required")
	}

	// Check for debug mode
	debugMode := os.Getenv("PG_DEBUG") == "1"

	// Set up log directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	logDir := filepath.Join(homeDir, ".playground", "logs")

	return &LocalProvider{
		modelPath:   modelPath,
		contextSize: 4096,
		debugMode:   debugMode,
		logDir:      logDir,
	}, nil
}

// Chat implements the Provider interface
func (p *LocalProvider) Chat(messages []Message, tools []Tool) (*Response, error) {
	// Build prompt from messages and tools
	prompt := p.buildPrompt(messages, tools)

	// Call llama.cpp CLI (llama-cli or main executable)
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "llama-cli",
		"--model", p.modelPath,
		"--prompt", prompt,
		"--ctx-size", fmt.Sprintf("%d", p.contextSize),
		"--n-predict", "2048",
		"--temp", "0.1",
		"--top-k", "40",
		"--top-p", "0.9",
		"--threads", "4",
		"--no-display-prompt",
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("llama.cpp execution failed: %w\nStderr: %s", err, stderr.String())
	}

	result := strings.TrimSpace(stdout.String())

	// Parse response for tool calls
	response := &Response{
		Content:      result,
		FinishReason: "stop",
	}

	// Try to extract tool calls from response
	toolCalls := p.extractToolCalls(result)
	if len(toolCalls) > 0 {
		response.ToolCalls = toolCalls
		response.FinishReason = "tool_calls"
	}

	return response, nil
}

// ChatStream implements streaming chat
func (p *LocalProvider) ChatStream(messages []Message, tools []Tool) (<-chan StreamChunk, error) {
	ch := make(chan StreamChunk, 10)

	go func() {
		defer close(ch)

		// Build prompt
		prompt := p.buildPrompt(messages, tools)

		// Call llama.cpp CLI with streaming
		ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
		defer cancel()

		cmd := exec.CommandContext(ctx, "llama-cli",
			"--model", p.modelPath,
			"--prompt", prompt,
			"--ctx-size", fmt.Sprintf("%d", p.contextSize),
			"--n-predict", "2048",
			"--temp", "0.1",
			"--top-k", "40",
			"--top-p", "0.9",
			"--threads", "4",
			"--no-display-prompt",
		)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			ch <- StreamChunk{Error: fmt.Errorf("failed to create stdout pipe: %w", err)}
			return
		}

		if err := cmd.Start(); err != nil {
			ch <- StreamChunk{Error: fmt.Errorf("failed to start llama.cpp: %w", err)}
			return
		}

		// Stream output
		buf := make([]byte, 1024)
		var fullResponse strings.Builder

		for {
			n, err := stdout.Read(buf)
			if n > 0 {
				chunk := string(buf[:n])
				fullResponse.WriteString(chunk)
				ch <- StreamChunk{Content: chunk}
			}
			if err != nil {
				break
			}
		}

		cmd.Wait()

		// Check for tool calls in final response
		toolCalls := p.extractToolCalls(fullResponse.String())
		if len(toolCalls) > 0 {
			for _, tc := range toolCalls {
				ch <- StreamChunk{ToolCall: &tc, FinishReason: "tool_calls"}
			}
		} else {
			ch <- StreamChunk{FinishReason: "stop"}
		}
	}()

	return ch, nil
}

// Name returns the provider name
func (p *LocalProvider) Name() string {
	return "DeepSeek-Coder-7B-Instruct-v1.5 (local)"
}

// buildPrompt constructs the prompt for DeepSeek-Coder
func (p *LocalProvider) buildPrompt(messages []Message, tools []Tool) string {
	var prompt strings.Builder

	// System message with tool definitions
	prompt.WriteString("You are a coding assistant. You can use tools by outputting JSON.\n\n")

	if len(tools) > 0 {
		prompt.WriteString("Available tools:\n")
		for _, tool := range tools {
			prompt.WriteString(fmt.Sprintf("- %s: %s\n", tool.Name, tool.Description))
		}
		prompt.WriteString("\nTo use a tool, output JSON in this format:\n")
		prompt.WriteString(`{"tool": "tool_name", "args": {"param": "value"}}` + "\n\n")
	}

	prompt.WriteString("CRITICAL RULES:\n")
	prompt.WriteString("1. NEVER write files directly\n")
	prompt.WriteString("2. ONLY propose unified diffs\n")
	prompt.WriteString("3. Explain intent BEFORE proposing changes\n")
	prompt.WriteString("4. Ask clarifying questions if ambiguous\n")
	prompt.WriteString("5. One logical change per patch\n\n")

	// Add conversation history
	for _, msg := range messages {
		switch msg.Role {
		case "system":
			prompt.WriteString(fmt.Sprintf("System: %s\n\n", msg.Content))
		case "user":
			prompt.WriteString(fmt.Sprintf("User: %s\n\n", msg.Content))
		case "assistant":
			prompt.WriteString(fmt.Sprintf("Assistant: %s\n\n", msg.Content))
		case "tool":
			prompt.WriteString(fmt.Sprintf("Tool Result (%s): %s\n\n", msg.Name, msg.Content))
		}
	}

	prompt.WriteString("Assistant: ")

	return prompt.String()
}

// extractToolCalls attempts to parse tool calls from the response
func (p *LocalProvider) extractToolCalls(response string) []ToolCall {
	var toolCalls []ToolCall

	// Look for JSON tool call patterns
	// Pattern: {"tool": "name", "args": {...}}
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "{") {
			continue
		}

		var toolCall struct {
			Tool string                 `json:"tool"`
			Args map[string]interface{} `json:"args"`
		}

		if err := json.Unmarshal([]byte(line), &toolCall); err == nil {
			if toolCall.Tool != "" {
				toolCalls = append(toolCalls, ToolCall{
					ID:        fmt.Sprintf("call_%d", time.Now().UnixNano()),
					Name:      toolCall.Tool,
					Arguments: toolCall.Args,
				})
			}
		}
	}

	return toolCalls
}

// logToFile writes content to a log file if debug mode is enabled
func (p *LocalProvider) logToFile(filename, content string) error {
	if !p.debugMode {
		return nil // Silent in production
	}

	// Create log directory if it doesn't exist
	if err := os.MkdirAll(p.logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	logPath := filepath.Join(p.logDir, filename)

	// Append timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logContent := fmt.Sprintf("=== %s ===\n%s\n\n", timestamp, content)

	// Append to file
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(logContent); err != nil {
		return fmt.Errorf("failed to write to log file: %w", err)
	}

	return nil
}

// ChatSilent is like Chat but suppresses all output in production mode
// In debug mode (PG_DEBUG=1), logs are written to ~/.playground/logs/
func (p *LocalProvider) ChatSilent(messages []Message, tools []Tool, logName string) (*Response, error) {
	// Build prompt
	prompt := p.buildPrompt(messages, tools)

	// Log prompt in debug mode
	if err := p.logToFile(logName+"_prompt.txt", prompt); err != nil {
		// Don't fail on logging errors, just continue
		fmt.Fprintf(os.Stderr, "Warning: failed to log prompt: %v\n", err)
	}

	// Call llama.cpp CLI
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "llama-cli",
		"--model", p.modelPath,
		"--prompt", prompt,
		"--ctx-size", fmt.Sprintf("%d", p.contextSize),
		"--n-predict", "2048",
		"--temp", "0.1",
		"--top-k", "40",
		"--top-p", "0.9",
		"--threads", "4",
		"--no-display-prompt",
		"--log-disable", // Disable llama.cpp's own logging
	)

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run command
	if err := cmd.Run(); err != nil {
		// Log error in debug mode
		errorMsg := fmt.Sprintf("Command failed: %v\nStderr: %s", err, stderr.String())
		p.logToFile(logName+"_error.txt", errorMsg)
		return nil, fmt.Errorf("llama-cli failed: %w", err)
	}

	output := stdout.String()

	// Log output in debug mode
	if err := p.logToFile(logName+"_output.txt", output); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to log output: %v\n", err)
	}

	// Log stderr in debug mode (even if command succeeded)
	if stderr.Len() > 0 {
		p.logToFile("llama.log", stderr.String())
	}

	// Parse tool calls
	toolCalls := p.extractToolCalls(output)

	return &Response{
		Content:   output,
		ToolCalls: toolCalls,
	}, nil
}
