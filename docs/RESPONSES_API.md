# Responses API Migration Stance

## Position

HELM currently implements the **OpenAI Chat Completions API** (`/v1/chat/completions`) as its proxy compatibility layer. This is the industry-standard API used by the majority of AI agent frameworks.

## Responses API (2025+)

OpenAI introduced the Responses API as a successor to Chat Completions, adding native tool orchestration, multi-turn state management, and built-in retrieval. 

### HELM's Stance

| Aspect | Position |
|--------|----------|
| Chat Completions API | **Fully supported** — indefinite support, no deprecation planned |
| Responses API | **Tracked** — monitoring adoption across the ecosystem |
| Migration timeline | **No near-term migration** — will implement when ecosystem adoption reaches critical mass |
| Breaking changes | **None** — Chat Completions compatibility is a core contract |

### Rationale

1. **HELM's value is orthogonal** — HELM's security layer (Guardian, receipts, ProofGraph) operates at a layer below the API surface. The same governance applies regardless of the API format.
2. **Ecosystem inertia** — As of early 2026, Chat Completions remains the dominant API across providers (Anthropic, Google, Mistral, Cohere).
3. **Clean separation** — When/if HELM adds Responses API support, it will be an additional endpoint, not a replacement for Chat Completions.

### What This Means for Integrators

- **No action required** — Continue using `/v1/chat/completions`
- **Tool use** — Pass `tools` array in the request body; HELM's Guardian validates and signs tool calls before execution
- **Future-proofing** — HELM will announce Responses API support with a minimum 6-month notice period

## Missing Parameters (Tracked)

The following Chat Completions parameters are not yet proxied by HELM:

| Parameter | Status | Priority |
|-----------|--------|----------|
| `tool_choice` | Planned | P1 |
| `parallel_tool_calls` | Planned | P1 |
| `response_format` | Planned | P1 |
| `logprobs` | Not planned | P2 |
| `top_logprobs` | Not planned | P2 |

These parameters are pass-through to upstream providers and do not affect HELM's security guarantees. Implementation is tracked in the API conformance backlog.
