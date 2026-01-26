package agent

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// ToolCallRequest represents a parsed tool call from the model
type ToolCallRequest struct {
	Tool string                 `json:"tool"`
	Args map[string]interface{} `json:"args"`
}

// ParseToolCalls extracts tool calls from model output
// Handles both explicit JSON and markdown-wrapped JSON
func ParseToolCalls(output string) ([]ToolCallRequest, error) {
	var toolCalls []ToolCallRequest

	// Strategy 1: Look for explicit JSON objects
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and non-JSON lines
		if line == "" || !strings.HasPrefix(line, "{") {
			continue
		}

		// Try to parse as JSON
		var call ToolCallRequest
		if err := json.Unmarshal([]byte(line), &call); err == nil {
			if call.Tool != "" {
				toolCalls = append(toolCalls, call)
			}
		}
	}

	// Strategy 2: Extract JSON from markdown code blocks
	if len(toolCalls) == 0 {
		jsonBlocks := extractJSONFromMarkdown(output)
		for _, block := range jsonBlocks {
			var call ToolCallRequest
			if err := json.Unmarshal([]byte(block), &call); err == nil {
				if call.Tool != "" {
					toolCalls = append(toolCalls, call)
				}
			}
		}
	}

	// Strategy 3: Natural language fallback (basic pattern matching)
	if len(toolCalls) == 0 {
		naturalCalls := parseNaturalLanguageTools(output)
		toolCalls = append(toolCalls, naturalCalls...)
	}

	return toolCalls, nil
}

// extractJSONFromMarkdown finds JSON objects in markdown code blocks
func extractJSONFromMarkdown(text string) []string {
	var jsonBlocks []string

	// Pattern: ```json ... ``` or ``` ... ```
	re := regexp.MustCompile("```(?:json)?\\s*\\n?([^`]+)```")
	matches := re.FindAllStringSubmatch(text, -1)

	for _, match := range matches {
		if len(match) > 1 {
			content := strings.TrimSpace(match[1])
			if strings.HasPrefix(content, "{") {
				jsonBlocks = append(jsonBlocks, content)
			}
		}
	}

	return jsonBlocks
}

// parseNaturalLanguageTools attempts to extract tool calls from natural language
// This is a fallback for when the model doesn't use JSON format
func parseNaturalLanguageTools(text string) []ToolCallRequest {
	var calls []ToolCallRequest

	text = strings.ToLower(text)

	// Pattern: "read file X" or "read X"
	if strings.Contains(text, "read") && (strings.Contains(text, "file") || strings.Contains(text, ".go") || strings.Contains(text, ".js")) {
		// Try to extract filename
		re := regexp.MustCompile(`read\s+(?:file\s+)?([a-zA-Z0-9_./\\-]+\.[a-z]+)`)
		matches := re.FindStringSubmatch(text)
		if len(matches) > 1 {
			calls = append(calls, ToolCallRequest{
				Tool: "read_file",
				Args: map[string]interface{}{"path": matches[1]},
			})
		}
	}

	// Pattern: "list files" or "show directory"
	if strings.Contains(text, "list") && strings.Contains(text, "file") {
		calls = append(calls, ToolCallRequest{
			Tool: "list_files",
			Args: map[string]interface{}{"path": "."},
		})
	}

	// Pattern: "git status" or "check status"
	if strings.Contains(text, "git") && strings.Contains(text, "status") {
		calls = append(calls, ToolCallRequest{
			Tool: "git_status",
			Args: map[string]interface{}{},
		})
	}

	// Pattern: "git diff" or "show changes"
	if strings.Contains(text, "git") && strings.Contains(text, "diff") {
		calls = append(calls, ToolCallRequest{
			Tool: "git_diff",
			Args: map[string]interface{}{},
		})
	}

	return calls
}

// ValidateToolCall checks if a tool call is valid
func ValidateToolCall(call ToolCallRequest) error {
	validTools := map[string][]string{
		"read_file":     {"path"},
		"list_files":    {"path"},
		"git_status":    {},
		"git_diff":      {},
		"run_command":   {"cmd"},
		"propose_patch": {"file_path", "unified_diff"},
	}

	requiredArgs, exists := validTools[call.Tool]
	if !exists {
		return fmt.Errorf("unknown tool: %s", call.Tool)
	}

	for _, arg := range requiredArgs {
		if _, ok := call.Args[arg]; !ok {
			return fmt.Errorf("missing required argument '%s' for tool '%s'", arg, call.Tool)
		}
	}

	return nil
}
