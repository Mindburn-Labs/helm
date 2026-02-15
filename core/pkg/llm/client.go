package llm

import (
	"context"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Client interface {
	Chat(ctx context.Context, messages []Message, tools []ToolDefinition, options *SamplingOptions) (*Response, error)
}

type SamplingOptions struct {
	Temperature float64 `json:"temperature"`
	TopP        float64 `json:"top_p"`
	Seed        int64   `json:"seed"`
}

type ToolDefinition struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

type Response struct {
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls"`
}

type ToolCall struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}
