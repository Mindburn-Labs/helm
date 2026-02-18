# OSS-PROXY Reality

## What Exists

### 1. HTTP Handler (Stub)
- **File**: `core/pkg/api/openai_proxy.go` (93 lines)
- **Status**: Stub. Returns static "HELM governance proxy active" response.
- **No upstream forwarding**, no tool_call interception, no receipt creation.
- Registered at `/v1/chat/completions` when `HELM_ENABLE_OPENAI_PROXY=1` [CODE]

### 2. CLI Proxy Command (Functional)
- **File**: `core/cmd/helm/proxy_cmd.go` (388 lines)
- **Status**: Functional proxy with significant capability.
- Features:
  - `httputil.ReverseProxy` upstream forwarding (L140-250) [CODE]
  - Tool call interception from LLM response (L200-280) [CODE]
  - PEP validation: JCS canonicalization + SHA-256 hash of tool args (L86-104) [CODE]
  - Receipt chain: `proxyReceipt` with LamportClock, PrevHash, Ed25519 signature (L27-42) [CODE]
  - JSONL receipt persistence via `receiptStore` (L45-80) [CODE]
  - `--upstream`, `--port`, `--sign`, `--max-tokens` flags [CODE]

### 3. Feature Flag
- **Env var**: `HELM_ENABLE_OPENAI_PROXY` checked in `main.go`
- **Default**: OFF (flag not set) [CODE]

### 4. Request/Response Types
- `OpenAIChatRequest`: model, messages, stream fields [CODE]
- `OpenAIChatResponse`: id, object, created, model, choices, usage [CODE]
- Missing: `tools`, `tool_calls`, `tool_choice`, `function_call` fields

## Code References
| Component | File | Lines |
|-----------|------|-------|
| Stub handler | `core/pkg/api/openai_proxy.go` | 50-92 |
| CLI proxy | `core/cmd/helm/proxy_cmd.go` | 106-387 |
| Receipt type | `core/cmd/helm/proxy_cmd.go` | 27-42 |
| Receipt store | `core/cmd/helm/proxy_cmd.go` | 45-80 |
| JCS validation | `core/cmd/helm/proxy_cmd.go` | 86-104 |
