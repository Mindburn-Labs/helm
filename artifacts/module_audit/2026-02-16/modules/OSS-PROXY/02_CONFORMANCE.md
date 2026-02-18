# OSS-PROXY Conformance

## Passes
- [x] POST /v1/chat/completions endpoint exists
- [x] Feature flag OFF by default
- [x] Upstream forwarding works (CLI path)
- [x] Tool call interception works (CLI path)
- [x] Receipt creation with LamportClock + PrevHash
- [x] Ed25519 signing available

## Fails
- [ ] Handler path (openai_proxy.go) does not forward to upstream
- [ ] Tool calls not routed through Guardian -> Executor
- [ ] No pinned schema validation (only JSON parse + JCS hash)
- [ ] No LLM re-invocation loop
- [ ] Reason codes are ad-hoc strings, not from stable taxonomy
- [ ] Receipts in JSONL, not ProofGraph DAG
- [ ] Request/response types missing tool_calls, tools fields
- [ ] No budget enforcement
- [ ] No SSE streaming support
