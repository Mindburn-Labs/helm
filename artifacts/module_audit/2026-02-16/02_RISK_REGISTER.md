# Risk Register

## P0 - Ship Blockers

| ID | Risk | Impact | Location | Mitigation |
|----|------|--------|----------|------------|
| R-001 | Proxy does not route tool calls through Guardian/Executor | Wedge Step 1 fails. No proof that governance actually runs. | `core/cmd/helm/proxy_cmd.go:200-280` intercepts but does not call `guardian.EvaluateDecision` | Wire `runProxyCmd` to invoke Guardian -> Executor pipeline for each tool_call |
| R-002 | No LLM re-invocation loop | Agent workflows break. Only single-turn passthrough works. | `proxy_cmd.go` does one forward + intercept, no iteration | Implement agentic loop: forward -> intercept tool_calls -> execute -> re-invoke LLM with results |
| R-003 | ProofGraph not integrated into proxy receipts | Receipts are flat JSONL, not cryptographic DAG. Audit chain lacks structural integrity. | `proxy_cmd.go:45-80` uses `receiptStore` (JSONL), not `proofgraph.Graph` | Replace receiptStore with ProofGraph.Append for each proxy receipt |
| R-004 | MCP gateway is a stub | MCP interop promise is not fulfilled. All tool calls return governance error. | `core/pkg/mcp/gateway.go:75-98` | Wire handleExecute to Guardian -> Executor boundary |
| R-005 | Region profiles not enforced | Step 4 fails. Config is documentation, not runtime enforcement. | `core/pkg/config/profiles/regional.yaml` + `config.go` (945B) | Implement profile loader, split to separate files, wire enforcement |
| R-006 | No control-room UI | Step 2 demo cannot visually prove value to skeptics | No UI directory exists | Build minimal React/Next.js dashboard with P0, approvals inbox, ProofGraph timeline |
| R-007 | Conformance output not byte-deterministic | Step 3 fails. Two runs with same input may produce different bytes. | `conform/engine.go:143` uses `json.MarshalIndent` | Use JCS canonicalization + include git commit + env fingerprint |
| R-008 | EvidencePack export not wired to session state | `helm export pack` cannot be run as CLI command with session_id | `export_pack.go` is library function, not CLI subcommand | Wire ExportPack into CLI with session_id -> receipt lookup logic |
| R-009 | Replay does not re-execute effects | Offline replay only verifies tape integrity, not effect reproduction | `replay_cmd.go` loads and checks tapes only | Implement effect replay with deterministic driver + hash comparison |
| R-010 | Proxy reason codes not from stable taxonomy | Proxy uses ad-hoc status strings, not `conform.Reason*` constants | `proxy_cmd.go` receipt status is freeform | Map all proxy deny/allow outcomes to `conform.Reason*` codes |

## P1 - Demo/Credibility Blockers

| ID | Risk | Impact | Location | Mitigation |
|----|------|--------|----------|------------|
| R-011 | Use case docs are pointer files (~120B) | UCs lack substance. Cannot demonstrate reproducible scenarios. | `docs/use-cases/UC-001.md` through `UC-012.md` | Expand each UC to include expected input, expected output, deterministic assertions |
| R-012 | No rogue agent turnkey demo | Step 2 rogue agent reproduction requires manual setup | No dedicated script | Create `scripts/demo/rogue_agent.sh` that exercises budget deny path |
| R-013 | Conformance not run in CI | CI runs tests but not `helm conform --profile CORE` | `.github/workflows/helm_core_gates.yml` has use-cases job but no conform job | Add `helm conform --profile CORE --json` to CI |
| R-014 | Schema validation not integrated in proxy | Tool args validated as JSON but not against pinned tool schemas | `proxy_cmd.go:86-104` does JCS validation only | Wire `manifest.ValidateToolArgs` into proxy loop |
| R-015 | Budget enforcer not wired into proxy | Budget checks exist but proxy does not invoke them | `budget/enforcer.go` + `proxy_cmd.go` | Wire budget.SimpleEnforcer.Check into proxy tool execution path |

## P2 - Polish/Completeness

| ID | Risk | Impact | Location | Mitigation |
|----|------|--------|----------|------------|
| R-016 | Config loader is minimal (945B) | Cannot dynamically load profiles or region config | `core/pkg/config/config.go` | Implement full config loader with YAML parsing and env overrides |
| R-017 | MCP catalog has no real tools | Catalog returns empty or mock tools | `core/pkg/mcp/catalog.go` | Populate with HELM-governed tool references |
| R-018 | No streaming (SSE) support in proxy | Modern SDK clients expect SSE responses | `openai_proxy.go` has Stream field but no SSE implementation | Add SSE passthrough with per-chunk governance |
| R-019 | EvidencePack does not bundle ProofGraph | Pack contains file hashes but not the full DAG | `export_pack.go` | Include ProofGraph JSON serialization in pack |
| R-020 | No threat model for proxy endpoint | Security posture incomplete for OSS surface | `docs/THREAT_MODEL.md` exists but does not cover proxy | Add proxy-specific threat model section |
