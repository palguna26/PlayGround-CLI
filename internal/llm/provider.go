package llm

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`           // "system", "user", "assistant", "tool"
	Content string `json:"content"`        // Message content
	Name    string `json:"name,omitempty"` // For tool responses
}

// Tool represents a tool that the LLM can call
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"` // JSON Schema
}

// ToolCall represents a request from the LLM to invoke a tool
type ToolCall struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// Response represents the LLM's response
type Response struct {
	Content      string     `json:"content"`       // Text response
	ToolCalls    []ToolCall `json:"tool_calls"`    // Requested tool invocations
	FinishReason string     `json:"finish_reason"` // "stop", "tool_calls", "length", etc.
}

// StreamChunk represents a chunk of streaming response
type StreamChunk struct {
	Content      string    // Text content delta
	ToolCall     *ToolCall // Tool call if present
	FinishReason string    // Finish reason if stream is ending
	Error        error     // Error if something went wrong
}

// Provider is the interface that all LLM providers must implement
type Provider interface {
	// Chat sends messages to the LLM and receives a response
	// Tools are optional - pass nil if not using tools
	Chat(messages []Message, tools []Tool) (*Response, error)

	// ChatStream sends messages and returns a streaming response channel
	ChatStream(messages []Message, tools []Tool) (<-chan StreamChunk, error)

	// Name returns the provider name (for logging)
	Name() string
}
