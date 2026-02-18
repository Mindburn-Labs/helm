# 03 — API Conformance

**Score: 4/5** · Gate ≥3 · **✅ PASS**

---

## OpenAI-Compatible Proxy (Verified)

### Endpoint Compatibility

| OpenAI Endpoint | HELM Equivalent | Status | Evidence |
|----------------|-----------------|--------|----------|
| `POST /v1/chat/completions` | `POST /v1/chat/completions` | ✅ | OpenAPI spec |
| Streaming (SSE) | `stream=true` | ⚠️ SSE buffered-then-validated | T9 notes eventual consistency |
| `POST /v1/responses` | — | ❌ Not implemented | OpenAI Responses API (March 2025) |
| Tool call shape | `tool_calls[].function.{name,arguments}` | ✅ | Matches OpenAI |
| Tool message role | `role: "tool"` + `tool_call_id` | ✅ | In OpenAPI spec |
| Multi-step tool loop | Auto-loop until `finish_reason: "stop"` | ✅ | UC-012 |

### Request Schema

| Parameter | Supported | Notes |
|-----------|-----------|-------|
| `model` | ✅ | Required string |
| `messages` | ✅ | Full role/content/tool_call_id |
| `tools` | ✅ | `type: function` + `function.{name, description, parameters}` |
| `temperature` | ✅ | |
| `max_tokens` | ✅ | |
| `stream` | ✅ | |
| `tool_choice` | ❌ | **Not in OpenAPI spec** — confirmed via grep |
| `parallel_tool_calls` | ❌ | Not supported |
| `response_format` | ❌ | Not supported |
| `top_p`, `frequency_penalty` | ❌ | Not supported |
| `logprobs`, `top_logprobs` | ❌ | Not supported |
| `seed` | ❌ | Not supported |

### Response Schema

| Field | Supported | Notes |
|-------|-----------|-------|
| `id`, `object`, `created`, `model` | ✅ | Standard fields |
| `choices[].message.{role, content, tool_calls}` | ✅ | |
| `choices[].finish_reason` | ✅ | |
| `usage.{prompt_tokens, completion_tokens, total_tokens}` | ✅ | |
| `X-Helm-Receipt-ID` | ✅ | HELM-specific header |
| `X-Helm-Output-Hash` | ✅ | HELM-specific header |
| `X-Helm-Lamport-Clock` | ✅ | HELM-specific header |

### Error Shape

```json
{
  "error": {
    "message": "...",
    "type": "invalid_request|authentication_error|permission_denied|not_found|internal_error",
    "code": "...",
    "reason_code": "DENY_*"
  }
}
```

OpenAI-compatible (`error.{message, type, code}`) with HELM-specific `reason_code` extension. Good practice — no breaking changes for standard clients.

---

## Authentication (Verified: `auth/middleware.go`, 143 LOC)

### JWT Middleware
- Bearer token extraction from Authorization header
- `HelmClaims` structure: `TenantID`, `Roles`, `RegisteredClaims` (L20-24)
- Fail-closed: nil validator → all non-public requests rejected (L103)
- Public paths excluded: `/health`, `/readiness`, `/startup`, `/api/auth/*` (L50-61)
- Principal injection into context for downstream authorization (L119-126)

### Rate Limiting
- `auth/ratelimit.go` exists with tests (confirmed in codebase)

---

## Responses API Migration Posture

> [!CAUTION]
> OpenAI GPT-4o API support **ended Feb 16, 2026**. Assistants API sunset: **Aug 26, 2026**. Chat Completions API continues "indefinitely."

### Impact Assessment
| Item | Status |
|------|--------|
| Chat Completions API deprecation risk | LOW — OpenAI says "indefinitely" |
| Responses API compatibility | ❌ Not implemented |
| Model name updates needed | YES — `gpt-4o` → successor models |
| MCP tool routing via proxy | ✅ MCP support exists (`mcp/` package) |

### Recommendation
Document Responses API posture. The Chat Completions path is safe for now, but adding Responses API compatibility would expand HELM's market.

---

## Strict Mode Assessment (Source-Verified)

| Feature | Status | Code Evidence |
|---------|--------|---------------|
| Reject undeclared tools | ✅ Default deny | `firewall.go` L58 |
| Reject non-JSON tool args | ✅ JCS enforces valid JSON | `canonicalize/jcs.go` |
| Enforce schema IDs | ✅ Pinned schemas with SHA-256 | `manifest/` package |
| Enforce max iterations | ⚠️ Not documented | No config option found |
| Enforce wallclock budget | ✅ WASI time limits | `sandbox.go` L99-102 |
| Enforce token budget | ✅ Budget package | `budget/budget.go` |
| Enforce cost budget | ✅ | `budget/budget.go` |

---

## Golden Vectors

| Vector | Description | Status |
|--------|-------------|--------|
| CV-001 | Simple chat, no tools | Needs creation |
| CV-002 | Single tool call, valid schema | Needs creation |
| CV-003 | Tool call with schema mismatch → DENY | UC-002 exercises |
| CV-004 | Multi-step tool loop | UC-012 exercises |
| CV-005 | Streaming with tool call | Needs creation |
| CV-006 | Budget exceeded mid-execution | UC-005 exercises |
| CV-007 | Unknown tool → DENY | UC-001 exercises |

---

## Score: 4/5

**Justification:**
- ✅ Core OpenAI tool calling semantics correct
- ✅ Error shapes OpenAI-compatible with extension
- ✅ 12 use cases provide conformance evidence
- ✅ JWT auth with fail-closed behavior
- ❌ Missing `tool_choice`, `parallel_tool_calls`, `response_format`
- ❌ No Responses API migration stance
- ❌ Max tool iteration not configurable

### To reach 5/5:
1. Add `tool_choice` parameter support
2. Document Responses API posture
3. Add max tool iteration configuration
4. Create and store golden vectors under `evidence/conformance/`
