// Package api provides the OpenAI-compatible proxy endpoint for HELM.
// Enabled via HELM_ENABLE_OPENAI_PROXY=1, this intercepts tool calls
// through the PEP boundary, enforcing governance on every operation.
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// OpenAIProxyConfig configures the OpenAI-compatible proxy.
type OpenAIProxyConfig struct {
	UpstreamURL  string `json:"upstream_url"`
	DefaultModel string `json:"default_model"`
}

// OpenAIMessage represents a message in the OpenAI chat format.
type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIChatRequest is the OpenAI-compatible request format.
// API-001/002: Includes tool_choice, parallel_tool_calls, and response_format
// for upstream provider pass-through.
type OpenAIChatRequest struct {
	Model             string          `json:"model"`
	Messages          []OpenAIMessage `json:"messages"`
	Stream            bool            `json:"stream,omitempty"`
	ToolChoice        any             `json:"tool_choice,omitempty"`         // API-001: "auto", "none", "required", or {"type":"function","function":{"name":"..."}}
	ParallelToolCalls *bool           `json:"parallel_tool_calls,omitempty"` // API-001: Enable/disable parallel tool execution
	ResponseFormat    any             `json:"response_format,omitempty"`     // API-002: {"type":"json_object"} or {"type":"json_schema","json_schema":{...}}
	MaxTokens         *int            `json:"max_tokens,omitempty"`
	Temperature       *float64        `json:"temperature,omitempty"`
	TopP              *float64        `json:"top_p,omitempty"`
}

// OpenAIChatResponse is the OpenAI-compatible response format.
type OpenAIChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int           `json:"index"`
		Message      OpenAIMessage `json:"message"`
		FinishReason string        `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// HandleOpenAIProxy is the handler for /v1/chat/completions in server mode.
// This is an in-process stub. For governed proxy mode with upstream forwarding,
// use `helm proxy --upstream <url>` which provides full governance: Guardian,
// ProofGraph, budget enforcement, and receipt generation.
func HandleOpenAIProxy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteMethodNotAllowed(w)
		return
	}

	var req OpenAIChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteBadRequest(w, "Invalid request body")
		return
	}

	if req.Model == "" {
		req.Model = "gpt-4"
	}

	// In-process mode: return stub response directing to CLI proxy.
	// For full governance, use: helm proxy --upstream <llm-api-url>
	resp := OpenAIChatResponse{
		ID:      fmt.Sprintf("chatcmpl-helm-%d", time.Now().UnixNano()),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   req.Model,
	}
	resp.Choices = append(resp.Choices, struct {
		Index        int           `json:"index"`
		Message      OpenAIMessage `json:"message"`
		FinishReason string        `json:"finish_reason"`
	}{
		Index: 0,
		Message: OpenAIMessage{
			Role:    "assistant",
			Content: "HELM in-process proxy stub. For governed proxy mode with tool call interception, budget enforcement, and ProofGraph receipts, use: helm proxy --upstream <your-llm-api-url>",
		},
		FinishReason: "stop",
	})

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
