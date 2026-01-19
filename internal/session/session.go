package session

import (
	"time"
)

// Session represents a PlayGround coding session
type Session struct {
	ID             string     `json:"id"`
	Repo           string     `json:"repo"`            // Absolute path to repository
	Goal           string     `json:"goal"`            // User's stated goal for this session
	ContextSummary string     `json:"context_summary"` // AI-maintained summary of session progress
	PendingPatches []Patch    `json:"pending_patches"` // Diffs proposed by agent, not yet applied
	ToolHistory    []ToolCall `json:"tool_history"`    // Record of all tool invocations
	CreatedAt      time.Time  `json:"created_at"`
}

// Patch represents a proposed code change as a unified diff
type Patch struct {
	FilePath    string    `json:"file_path"`    // Relative path from repo root
	UnifiedDiff string    `json:"unified_diff"` // Complete unified diff format
	CreatedAt   time.Time `json:"created_at"`
}

// ToolCall records a tool invocation and its result
type ToolCall struct {
	ToolName  string                 `json:"tool_name"`
	Arguments map[string]interface{} `json:"arguments"`
	Result    string                 `json:"result"`
	Error     string                 `json:"error,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// Validate checks if a session is valid
func (s *Session) Validate() error {
	if s.ID == "" {
		return ErrInvalidSessionID
	}
	if s.Repo == "" {
		return ErrInvalidRepoPath
	}
	if s.Goal == "" {
		return ErrEmptyGoal
	}
	return nil
}
