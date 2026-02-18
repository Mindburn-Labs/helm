# HELM OSS Wedge Gap Audit - 2026-02-16

**Auditor**: AGC Gap Audit Agent
**Scope**: HELM OSS wedge requirements (Steps 1-5) vs current repo state
**Repo**: Mindburn-Labs/helm (public), commit at HEAD
**Date**: 2026-02-16

## Audit Structure

| File | Purpose |
|------|---------|
| [01_GATEBOARD.md](01_GATEBOARD.md) | Go/no-go status for each wedge step |
| [02_RISK_REGISTER.md](02_RISK_REGISTER.md) | Risks ranked by severity |
| [03_DECISIONS.md](03_DECISIONS.md) | Architecture decisions and rationale |
| [04_INTEGRATION_MATRIX.md](04_INTEGRATION_MATRIX.md) | Component wiring correctness |
| [05_SOTA_BASELINE.md](05_SOTA_BASELINE.md) | Competitive moat assessment |
| [06_OSS_WEDGE_REQUIREMENTS.md](06_OSS_WEDGE_REQUIREMENTS.md) | Full requirements checklist |
| [07_GAP_MATRIX.md](07_GAP_MATRIX.md) | Gap-by-gap matrix with severity |
| [08_BACKLOG_PLAN.md](08_BACKLOG_PLAN.md) | Phased implementation backlog |
| [09_VERIFICATION_PLAN.md](09_VERIFICATION_PLAN.md) | Verification commands and expected outputs |
| [10_EVIDENCE/](10_EVIDENCE/) | Run logs, diffs, hash outputs |

## Module Deep-Dives

| Module | Path |
|--------|------|
| OSS-PROXY | [modules/OSS-PROXY/](modules/OSS-PROXY/) |
| OSS-EXEC-LOOP | [modules/OSS-EXEC-LOOP/](modules/OSS-EXEC-LOOP/) |
| OSS-PEP-GUARDIAN | [modules/OSS-PEP-GUARDIAN/](modules/OSS-PEP-GUARDIAN/) |
| OSS-EXECUTOR | [modules/OSS-EXECUTOR/](modules/OSS-EXECUTOR/) |
| OSS-PROOFGRAPH | [modules/OSS-PROOFGRAPH/](modules/OSS-PROOFGRAPH/) |
| OSS-TRUST-REGISTRY | [modules/OSS-TRUST-REGISTRY/](modules/OSS-TRUST-REGISTRY/) |
| OSS-SCHEMAS-MANIFEST | [modules/OSS-SCHEMAS-MANIFEST/](modules/OSS-SCHEMAS-MANIFEST/) |
| OSS-BOUNDED-COMPUTE | [modules/OSS-BOUNDED-COMPUTE/](modules/OSS-BOUNDED-COMPUTE/) |
| OSS-APPROVALS | [modules/OSS-APPROVALS/](modules/OSS-APPROVALS/) |
| OSS-EVIDENCEPACK | [modules/OSS-EVIDENCEPACK/](modules/OSS-EVIDENCEPACK/) |
| OSS-REPLAY | [modules/OSS-REPLAY/](modules/OSS-REPLAY/) |
| OSS-CONFORMANCE | [modules/OSS-CONFORMANCE/](modules/OSS-CONFORMANCE/) |
| OSS-PROFILES | [modules/OSS-PROFILES/](modules/OSS-PROFILES/) |
| OSS-MCP | [modules/OSS-MCP/](modules/OSS-MCP/) |
| OSS-UI | [modules/OSS-UI/](modules/OSS-UI/) |
| OSS-CI-GATES | [modules/OSS-CI-GATES/](modules/OSS-CI-GATES/) |
| OSS-DOCS-ADOPTION | [modules/OSS-DOCS-ADOPTION/](modules/OSS-DOCS-ADOPTION/) |

## Top 10 P0 Gaps (Summary)

1. **Proxy execution loop is not wired**: `proxy_cmd.go` intercepts tool_calls but does not route through KernelBridge -> Guardian -> Executor boundary
2. **No real upstream forwarding in proxy handler**: `openai_proxy.go` returns static response, no LLM forwarding
3. **No tool_call iteration loop**: proxy does not re-invoke LLM with tool outputs until completion
4. **ProofGraph not integrated into proxy**: receipts in proxy use flat JSONL, not ProofGraph DAG
5. **MCP gateway is a stub**: returns "requires governance approval" for all calls
6. **Regional profiles are a single file**, not separate `profile_{us,eu,ru,cn}.yaml` with full config loading enforcement
7. **No control-room UI exists**: no dashboard, approvals inbox, or ProofGraph timeline
8. **Conformance output not deterministic bytes**: no canonicalization of JSON output, no environment fingerprinting
9. **EvidencePack export not wired to CLI**: `helm export pack` is internal function, not integrated with session state
10. **Replay lacks actual effect re-execution**: tape verification only, no offline re-execution
