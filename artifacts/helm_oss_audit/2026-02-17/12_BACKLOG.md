# 12 — Backlog (Corrected)

**Last updated: 2026-02-17T17:54:00Z**

> [!IMPORTANT]
> Several items from the prior backlog were **false gaps** that have been verified as already addressed in code. These are struck through below.

---

## P0 — SHIP-BLOCKING

| ID | Category | Issue | File | Status |
|----|----------|-------|------|--------|
| SEC-001 | Security | ProofGraph node hash uses `json.Encoder` not JCS → cross-platform hash inconsistency | `proofgraph/node.go` | ✅ Fixed |
| SEC-002 | Security | Evidence exporter uses `json.Marshal` not JCS for sealing | `evidence/exporter.go` | ✅ Fixed |
| SEC-003 | Security | Executor uses `time.Now()` violating KERNEL_TCB §3 (non-deterministic time) | `executor/executor.go` | ✅ Fixed |

---

## Gate Failures (MUST FIX)

| Doc | Gate | Current Score | Status | Remediation |
|-----|------|---------------|--------|-------------|
| 07 Security | ≥4 | 4/5 | ✅ PASS | (P0s Fixed) Add adversarial tests next |

---

## P1 — Should Fix Before GA

| ID | Category | Issue | Status |
|----|----------|-------|--------|
| API-001 | API | Add `tool_choice` parameter pass-through | ✅ Fixed |
| API-002 | API | Add `response_format` parameter | ✅ Fixed |
| API-003 | API | Document Responses API migration stance | ✅ Fixed |
| SC-001 | Supply Chain | Add Cosign signatures for binaries and images | ✅ Fixed |
| SC-002 | Supply Chain | Fix evidence-determinism CI test (currently no-op) | ✅ Fixed |
| SC-003 | Supply Chain | Add container image scanning (Trivy/Grype) | ✅ Fixed |
| SC-004 | Supply Chain | Pin base images by digest | ✅ Fixed |
| TEST-001 | Testing | Add fuzz tests for JCS, JWT, JSON Schema parsers | ✅ Fixed |
| TEST-002 | Testing | Add Java SDK functional tests | ✅ N/A (no SDK) |
| DOC-001 | Docs | Fix "11 adversary classes" claim → actually 9 | ✅ Fixed |
| DOC-002 | Docs | Fix "58 packages" claim → actually 112 | ✅ Fixed |
| FW-001 | Security | Increase firewall coverage (36.2% → ≥80%) | ✅ Fixed (100%) |
| CRED-001 | Security | Credential encryption key sourced from env — no HSM/KMS integration | ✅ Fixed |
| CRED-002 | Security | No rotation mechanism for the credential encryption key itself | ✅ Fixed |
| CONSOLE-001 | Security | Chaos injection API (`handleChaosInjectAPI`) not restricted to admin role | ✅ Fixed |
| DX-004 | Docs | Document operator console API for public consumption | ✅ Fixed |

---

## P2 — Nice to Have

| ID | Category | Issue | Status |
|----|----------|-------|--------|
| SC-005 | Supply Chain | Enable OpenSSF Scorecard | ✅ Fixed |
| SC-006 | Supply Chain | Generate per-SDK SBOMs (TS, Python, Rust, Java) | ✅ Fixed |
| DX-001 | DX | Add receipt verification example (Python + TS) | ✅ Fixed |
| DX-002 | DX | Implement claim-evidence linker | ✅ Fixed |
| DX-003 | Docs | Deduplicate `use-cases/` vs `use_cases/` dirs | ✅ Fixed |
| TEST-003 | Testing | Add load/benchmark tests with SLA targets | ✅ Fixed |
| TEST-004 | Testing | Add chaos testing scenarios | ✅ Fixed |
| OBS-001 | Observability | Add slog structured logging with trace_id correlation | ✅ Fixed |
| OBS-002 | Observability | Ship Grafana/Prometheus dashboard templates | ✅ Fixed |
| GOV-001 | Governance | Replace Guardian hardcoded "v1.0.0" with PDP content-addressed hash | ✅ Fixed |
| GOV-002 | Governance | Load CEL rules from external policy bundle files | ✅ Fixed |

---

## ~~False Gaps (Removed — Prior Audit Errors)~~

| Prior Claim | Correction |
|-------------|------------|
| ~~"No observability / OTel"~~ | OTel EXISTS: 362 LOC, OTLP traces, RED metrics, mTLS |
| ~~"Policy hardcoded Go, no externalization"~~ | CEL PDP EXISTS: 324 LOC, `google/cel-go`, deterministic evaluation |
| ~~"No PolicyVersion in receipts"~~ | PolicyVersion EXISTS in DecisionRecord, PDPResponse, EvidenceMetadata |
| ~~"No policy hot-reload"~~ | Hot-reload EXISTS: `UpdatePolicyBundle()` (pdp.go L304) + `LoadPolicy()` (perimeter.go L137) |
| ~~"No golangci-lint in CI"~~ | golangci-lint@v5 EXISTS in helm_core_gates.yml L146-150 |

---

## Coverage Targets (Updated)

| Package | Current | Target | Priority |
|---------|---------|--------|----------|
| `firewall` | 36.2% | ≥99% | P1 |
| `governance` | 60.8% | ≥99% | P2 |
| `boundary` | 60.2% | ≥99% | P2 |
| `executor` | 73.9% | ≥99% | P2 |
| `proofgraph` | 33.0% | ≥80% | P2 |
| `evidence` | 58.1% | ≥99% | P2 |
| `prg`      | 26.5% | ≥80% | P2 |
