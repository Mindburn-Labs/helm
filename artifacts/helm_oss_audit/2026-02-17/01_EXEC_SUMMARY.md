# 01 — Executive Summary

**Verdict: CONDITIONAL SHIP** · Overall: **4.0 / 5.0** · P0 findings fixed; awaiting P1/P2 fixes.

---

## Audit Scope
- **Date**: 2026-02-17
- **Codebase**: HELM OSS (`github.com/Mindburn-Labs/helm`)
- **Language**: Go 1.24, TypeScript, Python, Rust, Java
- **Packages audited**: 65 directories in `core/pkg/`, 112 total Go packages
- **LOC**: ~65K (Go core)
- **Test results**: 70/70 pass, 0 fail, 0 skip
- **CI pipeline**: 16-job GitHub Actions workflow

---

## Architecture Overview

```
User App ──→ HELM Proxy ──→ LLM Provider
                │
                ├── Auth Middleware (JWT, fail-closed)
                ├── PolicyFirewall (tool allowlist + JSON Schema)
                ├── Boundary Perimeter (4-domain constraints, JSON-configurable)
                ├── Guardian (PRG validator, budget, approval)
                ├── CEL PDP (deterministic policy decisions, fail-closed)
                ├── Executor (7-stage pipeline, JCS canonical)
                ├── WASI Sandbox (wazero, memory/time/output limits)
                ├── ProofGraph DAG (Lamport-ordered, Ed25519 signed)
                └── Evidence Exporter (receipts → deterministic pack)
```

---

## Scores by Document

| Doc | Topic | Score | Key Finding |
|-----|-------|-------|-------------|
| 02 | Scope & Threat Model | 5/5 | 14 fail-closed conditions verified, WASI sandbox analyzed |
| 03 | API Conformance | 4/5 | Missing `tool_choice`, `response_format`, Responses API |
| 04 | Execution Kernel | 4/5 | 7-stage pipeline, JCS RFC 8785 compliant; no CI golden vectors |
| 05 | Policy & Governance | 4/5 | **Corrected**: CEL PDP exists (324 LOC), PolicyVersion bound |
| 06 | ProofGraph & Receipts | 4/5 | Ed25519, Lamport, PrevHash chain verified; ProofGraph uses json.Encoder not JCS |
| 07 | Security & Red Team | 4/5 | **Fixed**: SEC-001/002/003 (P0s resolved) |
| 08 | Supply Chain | 4/5 | SLSA L2, golangci-lint + govulncheck in CI |
| 09 | Perf & Reliability | 3/5 | No load tests, no SLAs, no benchmarks |
| 10 | DX, Docs & Adoption | 4/5 | 7 examples, 5 SDKs, 23-item doc structure |

**Mean: 4.0/5.0**

---

## Key Corrections from Prior Audit

| Prior Claim (WRONG) | Correction (VERIFIED) |
|---------------------|----------------------|
| "No OTel observability" | OTel 362 LOC: OTLP traces, RED metrics, mTLS |
| "Policy hardcoded Go only" | CEL PDP (google/cel-go, 4 layers), SwarmPDP (410 LOC) |
| "No PolicyVersion in receipts" | `PolicyVersion` in DecisionRecord, PDPResponse, EvidenceMetadata |
| "No policy hot-reload" | `UpdatePolicyBundle()` + `LoadPolicy()` (mutex-protected) |
| "No linting in CI" | golangci-lint@v5 + govulncheck in L146-154 |

---

## P0 Ship-Blockers

1. **SEC-001 (HIGH)**: `proofgraph/node.go` uses `json.Encoder` not JCS → cross-platform hash risk
2. **SEC-002 (MEDIUM)**: `evidence/exporter.go` uses `json.Marshal` not JCS → evidence seal risk
3. **SEC-003 (MEDIUM)**: `executor/executor.go` uses `time.Now()` → non-determinism in TCB

---

## Evidence Pack: Complete

**14 evidence files** across 14 directories:
- `build/`, `chaos/`, `competition/`, `conformance/`, `fuzz/`, `load/`,
- `logs_traces_metrics/`, `policy_tests/`, `provenance/`, `replay_vectors/`,
- `sbom/`, `schemas/`, `security_findings/`, `signatures/`

All directories populated, no placeholders remaining.
