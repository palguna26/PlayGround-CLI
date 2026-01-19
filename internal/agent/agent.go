package agent

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/yourusername/playground/internal/llm"
	"github.com/yourusername/playground/internal/session"
	"github.com/yourusername/playground/internal/tools"
)

// Agent orchestrates the AI coding session
type Agent struct {
	Session  *session.Session
	Store    *session.Store
	Provider llm.Provider
	RepoRoot string
}

// AgentConfig holds configuration for the agent
type AgentConfig struct {
	MaxIterations int  // Hard limit on agent loop iterations
	Verbose       bool // Print detailed logging
}

var DefaultConfig = AgentConfig{
	MaxIterations: 10,
	Verbose:       false,
}

// getSystemPrompt returns the system prompt that defines agent behavior
func getSystemPrompt() string {
	return `You are a helpful coding assistant integrated into PlayGround, a CLI tool for AI-assisted development.

CRITICAL RULES:
1. You can NEVER write files directly
2. All code changes MUST be proposed as unified diffs via the propose_patch tool
3. Unified diffs must follow standard format with --- and +++ headers
4. Always validate your assumptions by reading files first
5. Be concise and focused on the user's goal

AVAILABLE TOOLS:
- read_file(path): Read a file's contents
- list_files(path): List files in a directory  
- git_status(): Check Git repository status
- git_diff(): See current uncommitted changes
- run_command(cmd): Execute a command (requires user approval)
- propose_patch(file_path, unified_diff): Propose a code change as a unified diff

WORKFLOW:
1. Understand the user's request
2. Explore the codebase using read_file and list_files
3. Formulate a plan
4. Propose changes as unified diffs
5. Explain what you did

Remember: You're helping the user code, not coding for them. Be helpful, safe, and transparent.`
}

// defineTools returns the tool definitions for the LLM
func defineTools() []llm.Tool {
	return []llm.Tool{
		{
			Name:        "read_file",
			Description: "Read the contents of a file from the repository",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "Relative path to the file from repository root",
					},
				},
				"required": []string{"path"},
			},
		},
		{
			Name:        "list_files",
			Description: "List files and directories at the given path",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "Relative path to directory from repository root (use '.' for root)",
					},
				},
				"required": []string{"path"},
			},
		},
		{
			Name:        "git_status",
			Description: "Get the current Git status of the repository",
			Parameters: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			Name:        "git_diff",
			Description: "Get the current Git diff (uncommitted changes)",
			Parameters: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			Name:        "run_command",
			Description: "Execute a shell command (requires user approval)",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"command": map[string]interface{}{
						"type":        "string",
						"description": "The shell command to execute",
					},
				},
				"required": []string{"command"},
			},
		},
		{
			Name:        "propose_patch",
			Description: "Propose a code change as a unified diff patch",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"file_path": map[string]interface{}{
						"type":        "string",
						"description": "Relative path to the file from repository root",
					},
					"unified_diff": map[string]interface{}{
						"type":        "string",
						"description": "Complete unified diff with --- and +++ headers",
					},
				},
				"required": []string{"file_path", "unified_diff"},
			},
		},
	}
}

// executeTool executes a tool call and returns the result
func (a *Agent) executeTool(toolCall llm.ToolCall) (string, error) {
	switch toolCall.Name {
	case "read_file":
		path, ok := toolCall.Arguments["path"].(string)
		if !ok {
			return "", fmt.Errorf("invalid path argument")
		}
		return tools.ReadFile(a.RepoRoot, path)

	case "list_files":
		path, ok := toolCall.Arguments["path"].(string)
		if !ok {
			return "", fmt.Errorf("invalid path argument")
		}
		files, err := tools.ListFiles(a.RepoRoot, path)
		if err != nil {
			return "", err
		}
		// Convert to JSON for LLM
		data, _ := json.MarshalIndent(files, "", "  ")
		return string(data), nil

	case "git_status":
		return tools.GitStatus(a.RepoRoot)

	case "git_diff":
		return tools.GitDiff(a.RepoRoot)

	case "run_command":
		cmd, ok := toolCall.Arguments["command"].(string)
		if !ok {
			return "", fmt.Errorf("invalid command argument")
		}
		return tools.RunCommand(a.RepoRoot, cmd)

	case "propose_patch":
		filePath, ok := toolCall.Arguments["file_path"].(string)
		if !ok {
			return "", fmt.Errorf("invalid file_path argument")
		}
		unifiedDiff, ok := toolCall.Arguments["unified_diff"].(string)
		if !ok {
			return "", fmt.Errorf("invalid unified_diff argument")
		}

		// Add patch to session
		newPatch := session.Patch{
			FilePath:    filePath,
			UnifiedDiff: unifiedDiff,
			CreatedAt:   time.Now(),
		}
		a.Session.PendingPatches = append(a.Session.PendingPatches, newPatch)

		return fmt.Sprintf("Patch proposed for %s. User can review with 'pg review' and apply with 'pg apply'.", filePath), nil

	default:
		return "", fmt.Errorf("unknown tool: %s", toolCall.Name)
	}
}

// logToolCall records a tool invocation in the session history
func (a *Agent) logToolCall(toolCall llm.ToolCall, result string, err error) {
	historyEntry := session.ToolCall{
		ToolName:  toolCall.Name,
		Arguments: toolCall.Arguments,
		Result:    result,
		Timestamp: time.Now(),
	}

	if err != nil {
		historyEntry.Error = err.Error()
	}

	a.Session.ToolHistory = append(a.Session.ToolHistory, historyEntry)
}
