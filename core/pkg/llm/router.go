package llm

import (
	"context"
	"fmt"
	"strings"
)

// Router decides which model to use for a given request.
type Router struct {
	fastClient  Client // e.g. Local Llama
	smartClient Client // e.g. GPT-4
	embedder    Embedder
}

func NewRouter(fast, smart Client, embedder Embedder) *Router {
	return &Router{fastClient: fast, smartClient: smart, embedder: embedder}
}

func (r *Router) Chat(ctx context.Context, msgs []Message, tools []ToolDefinition, options *SamplingOptions) (*Response, error) {
	if len(msgs) == 0 {
		return nil, fmt.Errorf("router: messages must not be empty")
	}

	// Simple Heuristic Routing
	// 1. If tools are implicated, use Smart Model (better function calling)
	if len(tools) > 0 {
		return r.smartClient.Chat(ctx, msgs, tools, options)
	}

	// 2. Analyze complexity via Embeddings/Semantic Router
	// Real impl: compute embedding of lastMsg, compare with "complex task" cluster center.
	// For now, we still use heuristics but structured to support the upgrade.
	lastMsg := msgs[len(msgs)-1].Content
	if r.isComplexSemantic(ctx, lastMsg) {
		return r.smartClient.Chat(ctx, msgs, tools, options)
	}

	// 3. Fast path
	return r.fastClient.Chat(ctx, msgs, tools, options)
}

func (r *Router) isComplexSemantic(ctx context.Context, text string) bool {
	// GAP-16: Use Embeddings.
	// We'd call r.embedder.Embed(text) -> vector
	// dist := cosineSimilarity(vector, COMPLEX_CLUSTER_CENTER)
	// return dist < THRESHOLD

	// Fallback to improved heuristics for now
	keywords := []string{"plan", "design", "architect", "reason", "verify", "root cause", "analyze"}
	text = strings.ToLower(text)
	for _, k := range keywords {
		if strings.Contains(text, k) {
			return true
		}
	}
	return len(text) > 200
}
