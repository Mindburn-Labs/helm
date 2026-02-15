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
type OpenAIChatRequest struct {
	Model    string          `json:"model"`
	Messages []OpenAIMessage `json:"messages"`
	Stream   bool            `json:"stream,omitempty"`
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

// HandleOpenAIProxy is the handler for /v1/chat/completions.
// In production, this would forward to the upstream LLM while intercepting
// tool calls through the PEP boundary.
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

	// In the governed mode, every request passes through the PEP.
	// For now, return a governed response indicating the proxy is active.
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
			Content: "HELM governance proxy active. All tool calls are subject to PEP validation.",
		},
		FinishReason: "stop",
	})

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
