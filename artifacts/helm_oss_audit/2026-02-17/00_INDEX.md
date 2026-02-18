# HELM OSS Audit — 2026-02-17

**Verdict: CONDITIONAL SHIP** · Overall: **4.0 / 5.0** · P0 findings fixed; awaiting P1/P2 fixes.

---

| # | Document | Score | Gate |
|---|----------|-------|------|
| 01 | [Exec Summary](01_EXEC_SUMMARY.md) | — | — |
| 02 | [Scope & Threat Model](02_SCOPE_AND_THREAT_MODEL.md) | 5/5 | ✅ |
| 03 | [API Conformance](03_API_CONFORMANCE.md) | 4/5 | ✅ |
| 04 | [Execution Kernel & Replay](04_EXECUTION_KERNEL_REPLAY.md) | 4/5 | ✅ |
| 05 | [Policy & Governance](05_POLICY_AND_GOVERNANCE.md) | 4/5 | ✅ |
| 06 | [ProofGraph & Receipts](06_PROOFGRAPH_AND_RECEIPTS.md) | 4/5 | ✅ |
| 07 | [Security & Red Team](07_SECURITY_REDTEAM.md) | 4/5 | ✅ (P0s Fixed) |
| 08 | [Supply Chain & Releases](08_SUPPLY_CHAIN_RELEASES.md) | 4/5 | ✅ |
| 09 | [Perf, Reliability & SRE](09_PERF_RELIABILITY_SRE.md) | 3/5 | ✅ |
| 10 | [DX, Docs & Adoption](10_DX_DOCS_ADOPTION.md) | 4/5 | ✅ |
| 11 | [Competition Benchmark](11_COMPETITION_BENCHMARK.md) | — | — |
| 12 | [Backlog](12_BACKLOG.md) | — | — |

**Mean score: 4.0/5.0** (9 scored docs)

---

## SHIP-BLOCKING (P0) Findings

| ID | Severity | Description | File |
|----|----------|-------------|------|
| SEC-001 | LOW | ProofGraph node hash uses `json.Encoder` (Fixed) | `proofgraph/node.go` |
| SEC-002 | LOW | Evidence exporter uses `json.Marshal` (Fixed) | `evidence/exporter.go` |
| SEC-003 | LOW | Executor uses non-deterministic `time.Now()` (Fixed) | `executor/executor.go` |

## Evidence Pack: 14 files across 14 directories

| Directory | File | Status |
|-----------|------|--------|
| `build/` | `build_test_results.json` | ✅ |
| `chaos/` | `chaos_posture.json` | ✅ |
| `competition/` | `competitive_research.json` | ✅ |
| `conformance/` | `vectors_manifest.json` | ✅ |
| `fuzz/` | `fuzz_posture.json` | ✅ |
| `load/` | `load_posture.json` | ✅ |
| `logs_traces_metrics/` | `observability_evidence.json` | ✅ |
| `policy_tests/` | `coverage_report.json` | ✅ |
| `provenance/` | `ci_provenance_extract.json` | ✅ |
| `replay_vectors/` | `replay_evidence.json` | ✅ |
| `sbom/` | `sbom_extract.json` | ✅ |
| `schemas/` | `openapi_extract.json` | ✅ |
| `security_findings/` | `code_review_findings.json` | ✅ |
| `signatures/` | `crypto_evidence.json` | ✅ |

**13 documents + 14 evidence files = Complete audit pack.**
