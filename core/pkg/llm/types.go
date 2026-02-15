// Package llm provides LLM integration types.
// This file contains shared type definitions used by the llm package.
// NOTE: Many implementations have been removed as dead code.
// Only essential type definitions are retained here.
package llm

import "context"

// Embedder is the interface for creating text embeddings.
type Embedder interface {
	Embed(ctx context.Context, text string) ([]float32, error)
}

// ModelFingerprint contains versioning and identity information for a model invocation.
type ModelFingerprint struct {
	ProviderID       string   `json:"provider_id"`
	ModelID          string   `json:"model_id"`
	ModelVersion     string   `json:"model_version,omitempty"`
	SystemPromptHash string   `json:"system_prompt_hash,omitempty"`
	ToolManifestHash string   `json:"tool_manifest_hash,omitempty"`
	ConfigHash       string   `json:"config_hash,omitempty"`
	ParameterSet     []string `json:"parameter_set,omitempty"`
}
