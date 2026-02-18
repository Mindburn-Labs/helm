# OSS-PROXY Scope

## Wedge Requirement
- POST /v1/chat/completions endpoint
- Feature-flagged, OFF by default
- Forward to upstream OpenAI-compatible LLM
- Intercept tool_calls
- Developer changes only BASE_URL
- Deterministic deny with stable reason codes
