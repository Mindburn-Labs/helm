package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type OpenAIClient struct {
	apiKey string
	model  string
}

func NewOpenAIClient(apiKey, model string) *OpenAIClient {
	return &OpenAIClient{
		apiKey: apiKey,
		model:  model,
	}
}

// Internal structure for OpenAI API "tools" array
type openAITool struct {
	Type     string         `json:"type"`
	Function ToolDefinition `json:"function"`
}

type openAIRequest struct {
	Model       string       `json:"model"`
	Messages    []Message    `json:"messages"`
	Tools       []openAITool `json:"tools,omitempty"`
	Temperature float64      `json:"temperature,omitempty"`
	TopP        float64      `json:"top_p,omitempty"`
	Seed        int64        `json:"seed,omitempty"`
}

type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content   string `json:"content"`
			ToolCalls []struct {
				ID       string `json:"id"`
				Function struct {
					Name      string `json:"name"`
					Arguments string `json:"arguments"`
				} `json:"function"`
			} `json:"tool_calls"`
		} `json:"message"`
	} `json:"choices"`
}

func (c *OpenAIClient) Chat(ctx context.Context, msgs []Message, tools []ToolDefinition, options *SamplingOptions) (*Response, error) {
	var oaiTools []openAITool
	for _, t := range tools {
		oaiTools = append(oaiTools, openAITool{
			Type:     "function",
			Function: t,
		})
	}

	reqBody := openAIRequest{
		Model:    c.model,
		Messages: msgs,
		Tools:    oaiTools,
	}

	if options != nil {
		reqBody.Temperature = options.Temperature
		reqBody.TopP = options.TopP
		reqBody.Seed = options.Seed
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("openai: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("openai: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("openai error: %d", resp.StatusCode)
	}

	var oaiResp openAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&oaiResp); err != nil {
		return nil, err
	}

	if len(oaiResp.Choices) == 0 {
		return nil, fmt.Errorf("openai: empty choices in response")
	}
	choice := oaiResp.Choices[0].Message

	var toolCalls []ToolCall
	for _, tc := range choice.ToolCalls {
		var args map[string]any
		_ = json.Unmarshal([]byte(tc.Function.Arguments), &args) // Best effort parsing
		toolCalls = append(toolCalls, ToolCall{
			ID:        tc.ID,
			Name:      tc.Function.Name,
			Arguments: args,
		})
	}

	return &Response{
		Content:   choice.Content,
		ToolCalls: toolCalls,
	}, nil
}
