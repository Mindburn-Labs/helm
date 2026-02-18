# Gap Matrix

## Legend

- **Exists**: Y = yes, N = no, P = partial
- **Quality**: SOTA = state-of-the-art, OK = functional, Partial = incomplete, Placeholder = stub/non-functional
- **Severity**: P0 = ship blocker, P1 = demo/credibility blocker, P2 = polish

| # | Requirement | Exists | Quality | Evidence | Owner Module | Severity | Proposed Fix | DoD |
|---|------------|--------|---------|----------|-------------|----------|--------------|-----|
| 1 | POST /v1/chat/completions with upstream forwarding | P | Partial | [CODE] `proxy_cmd.go:140-250` forwards; `openai_proxy.go:50-92` is stub | OSS-PROXY | P0 | [MODIFY] Converge proxy handler to use proxy_cmd logic | Handler forwards to upstream and returns real LLM response |
| 2 | Tool call interception through KernelBridge -> Guardian -> Executor | N | - | [CODE] `proxy_cmd.go` intercepts but does not call Guardian/Executor | OSS-EXEC-LOOP | P0 | [MODIFY] `proxy_cmd.go` wire `guardian.EvaluateDecision` + `executor.Execute` | Tool calls go through full governance pipeline |
| 3 | LLM re-invocation loop | N | - | No iteration logic in any proxy path | OSS-EXEC-LOOP | P0 | [NEW] Implement tool_call -> execute -> reinvoke loop in proxy | Agent completes multi-turn tool use |
| 4 | ProofGraph integrated into proxy | N | - | [CODE] `proxy_cmd.go` uses JSONL, not `proofgraph.Graph` | OSS-PROOFGRAPH | P0 | [MODIFY] Replace JSONL receiptStore with ProofGraph | Proxy receipts are DAG nodes with chain validation |
| 5 | Pinned schema validation for tool args | P | Partial | [CODE] `proxy_cmd.go:86-104` validates JSON only | OSS-SCHEMAS-MANIFEST | P0 | [MODIFY] Wire `manifest.ValidateToolArgs` | Args validated against pinned schemas before execution |
| 6 | Pinned schema validation for tool outputs | N | - | No output validation in proxy | OSS-SCHEMAS-MANIFEST | P1 | [NEW] Wire `executor.OutputSchemaRegistry.LookupOutput` | Outputs validated for drift after execution |
| 7 | Proxy reason codes from stable taxonomy | N | - | [CODE] Proxy uses ad-hoc status strings | OSS-PROXY | P0 | [MODIFY] Map outcomes to `conform.Reason*` codes | All proxy denials use canonical reason codes |
| 8 | Budget enforcement in proxy | N | - | [CODE] `budget/enforcer.go` exists but not wired | OSS-BOUNDED-COMPUTE | P0 | [MODIFY] Wire `budget.SimpleEnforcer` into proxy | Budget violations produce deterministic deny |
| 9 | Max iterations / wallclock limits | P | Partial | [CODE] `proxy_cmd.go` has `--max-tokens` | OSS-BOUNDED-COMPUTE | P1 | [MODIFY] Add `--max-iterations`, `--max-wallclock` flags | Strict loop termination enforced |
| 10 | MCP gateway wired to kernel | N | Placeholder | [CODE] `mcp/gateway.go:75-98` returns static error | OSS-MCP | P1 | [MODIFY] Wire `handleExecute` to Guardian -> Executor | MCP tool calls governed |
| 11 | Feature flag OFF by default | Y | OK | [CODE] `HELM_ENABLE_OPENAI_PROXY` env var | OSS-PROXY | - | None | N/A |
| 12 | Regional profile separate files | P | Partial | [CODE] `config/profiles/regional.yaml` (single file) | OSS-PROFILES | P1 | [NEW] Split to `profile_{us,eu,ru,cn}.yaml` | 4 separate files with full config |
| 13 | Profile config loader | N | - | [CODE] `config/config.go` is 945B minimal | OSS-PROFILES | P1 | [NEW] Implement YAML profile loader | `helm --profile eu` loads correct profile |
| 14 | Outbound networking / allowlists / island mode | N | - | Not in profile YAML | OSS-PROFILES | P1 | [MODIFY] Add networking, allowlist, island fields | Profile controls applied at runtime |
| 15 | Crypto policy / retention in profiles | P | Partial | [CODE] `encryption` field exists | OSS-PROFILES | P2 | [MODIFY] Add crypto allowlist + retention fields | Full profile spec |
| 16 | Control-room UI | N | - | No UI code in repo | OSS-UI | P1 | [NEW] Minimal React dashboard | P0 dashboard, approvals, ProofGraph timeline |
| 17 | Rogue agent demo script | N | - | No dedicated script | OSS-DOCS-ADOPTION | P1 | [NEW] `scripts/demo/rogue_agent.sh` | Turnkey budget deny reproduction |
| 18 | `helm conform --level L1\|L2` | N | - | [CODE] Uses `--profile` not `--level` | OSS-CONFORMANCE | P1 | [MODIFY] Add `--level` alias | L1/L2 maps to gate subsets |
| 19 | Deterministic conformance output bytes | P | Partial | [CODE] `json.MarshalIndent` not canonicalized | OSS-CONFORMANCE | P1 | [MODIFY] Use JCS + add git commit + env fingerprint | Same input -> same bytes (SHA-256 verified) |
| 20 | Conformance CI gate | P | Partial | [CODE] CI runs tests but not `helm conform` | OSS-CI-GATES | P1 | [MODIFY] Add `helm conform --profile CORE --json` job | CI fails on conformance regression |
| 21 | `helm export pack <session_id>` CLI | P | Partial | [CODE] Function exists, CLI wiring incomplete | OSS-EVIDENCEPACK | P0 | [MODIFY] Wire into CLI with session lookup | `helm export pack abc123` produces tar.gz |
| 22 | EvidencePack bundles ProofGraph | N | - | Pack has file hashes only | OSS-EVIDENCEPACK | P1 | [MODIFY] Include ProofGraph JSON in pack | Pack is self-contained for offline verify |
| 23 | Offline replay re-execution | N | - | [CODE] Replay checks tapes only | OSS-REPLAY | P1 | [MODIFY] Add effect re-execution with deterministic driver | Replay reproduces decisions offline |
| 24 | Use case docs with assertions | P | Placeholder | [CODE] UC docs are ~120B pointer files | OSS-DOCS-ADOPTION | P2 | [MODIFY] Expand with input/output/assertions | Each UC is reproducible |
| 25 | Streaming (SSE) proxy support | N | - | `Stream` field exists but not implemented | OSS-PROXY | P2 | [NEW] SSE passthrough with chunk governance | Streaming clients work |
| 26 | Proxy threat model | P | Partial | [CODE] `docs/THREAT_MODEL.md` exists | OSS-DOCS-ADOPTION | P2 | [MODIFY] Add proxy-specific section | Threat model covers proxy surface |
| 27 | ProofGraph persistence | N | - | In-memory only + store interface | OSS-PROOFGRAPH | P1 | [NEW] Filesystem-backed store | ProofGraph survives restarts |
| 28 | EvidencePack includes trust roots | N | - | Trust roots not in pack | OSS-EVIDENCEPACK | P2 | [MODIFY] Bundle signing key + trust anchors | Verifiable without online registry |
