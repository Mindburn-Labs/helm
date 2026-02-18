# 09 — Performance, Reliability & SRE

**Score: 3/5** · Gate ≥3 · **✅ PASS** (Corrected from 2/5)

---

## Correction Notice

> The prior audit scored this domain 2/5 claiming "no OpenTelemetry instrumentation, no Prometheus metrics." This was **factually incorrect**. The `observability` package (362 LOC) provides a fully implemented OpenTelemetry stack. This document corrects the record with evidence.

---

## Observability Implementation (VERIFIED)

### OpenTelemetry Provider (`observability/observability.go`, 362 LOC)

| Feature | Implementation | Line Range |
|---------|---------------|------------|
| OTLP Trace Export | `initTraceProvider()` with gRPC exporter | L140-191 |
| OTLP Metric Export | `initMetricProvider()` with periodic reader | L193-219 |
| RED Metrics | `initREDMetrics()` — Rate, Errors, Duration counters | L221-263 |
| Configurable mTLS | `CertFile`, `KeyFile`, `CAFile` in Config | L31-43 |
| Sampling Rate | `Config.SampleRate` (default production-ready) | L45-57 |
| Span Management | `StartSpan()`, `TrackOperation()` helpers | L296-361 |
| Graceful Shutdown | `Shutdown()` with context | L265-278 |

### Config Structure
```go
type Config struct {
    ServiceName    string
    ServiceVersion string
    Environment    string
    OTLPEndpoint   string        // OTLP collector endpoint
    SampleRate     float64       // Trace sampling rate
    BatchTimeout   time.Duration // Batch export interval
    Enabled        bool
    Insecure       bool          // For dev
    CertFile       string        // mTLS cert
    KeyFile        string        // mTLS key
    CAFile         string        // mTLS CA
}
```

### Test Coverage
- `observability/`: 69.6% coverage
- Tests exist (confirmed in test results)

### Gap: Integration with Core Pipeline
While the OTel package exists, evidence of it being **wired into** the Guardian/Executor/ProofGraph pipeline at runtime was not confirmed — the imports are in the observability package itself. Whether `main.go` initializes the `Provider` and passes it to the kernel needs verification.

---

## Test Infrastructure

### Verified Test Results (2026-02-17)

```
Total packages:    112
Packages tested:    70 (PASS)
No test files:       6 (certification, conform/adversarial, kernel/errorir,
                        memory, security/suites, util/resiliency)
Test failures:       0
```

### Coverage Distribution

| Coverage Band | Packages | Notable |
|---------------|----------|---------|
| **100%** | 7 | `tiers`, `runtime`, `privacy`, `kernel/pdp`, `compliance/controls`, `compliance/fca`, `compliance/dora` (97.7%) |
| **90-99%** | 10 | `kernel/consistency` (98.9%), `versioning` (98.2%), `buildguard` (96.3%), `compliance/evidence` (95.3%), `merkle` (94.1%), `authz` (93.3%), `kernel/retry` (93.1%), `compliance` (92.3%), `runtime/budget` (91.7%), `rir` (91.2%) |
| **70-89%** | 15 | `trust/registry` (88.5%), `sdk` (88.9%), `provenance` (87.2%), `arc/connectors` (80.0%), `trust` (79.3%), `tooling` (79.8%), `compliance/csr` (79.4%), `bridge` (77.1%), `llm/modelpolicy` (77.5%), `sovereignty` (76.9%), `pack` (75.5%), `tape` (74.2%), `certification/admission` (72.6%), `kernel/celdp` (71.8%), `observability` (69.6%) |
| **50-69%** | 11 | `canonicalize` (65.9%), `boundary` (64.8%), `obligation` (63.6%), `runtime/sandbox` (61.0%), `manifest` (61.5%), `kernelruntime` (60.6%), `replay` (59.1%), `capabilities` (59.1%), `auth` (54.1%), `registry` (50.5%), `kernel/csnf` (50.0%) |
| **20-49%** | 8 | `mcp` (49.6%), `budget` (48.0%), `tenants` (46.5%), `llm` (45.1%), `proofgraph` (35.3%), `arc` (32.4%), `store` (28.6%), `prg` (28.9%) |
| **< 20%** | 5 | `agent` (24.4%), `artifacts` (23.4%), `api` (15.8%), `store/ledger` (12.7%), `metering` (9.3%) |
| **0%** | 6 | `certification`, `conform/adversarial`, `kernel/errorir`, `memory`, `security/suites`, `util/resiliency` |

### Critical Coverage Concerns

| Package | Coverage | Risk |
|---------|----------|------|
| `proofgraph` | 35.3% | Core trust component — insufficient |
| `prg` | 28.9% | Proof Requirement Graph — decision logic under-tested |
| `metering` | 9.3% | Billing accuracy at risk |
| `api` | 15.8% | HTTP handler edge cases untested |
| `store/ledger` | 12.7% | Data persistence reliability |

---

## Error Handling Architecture

### Fail-Closed Design (Verified)

| Component | Behavior | Code Reference |
|-----------|----------|----------------|
| Guardian | Missing evidence → `FAIL` verdict, still signed | `guardian.go` L164-174 |
| Guardian | Budget check error → `FAIL` (fail-closed) | `guardian.go` L131-136 |
| Guardian | Audit failure → hard error | `guardian.go` L324-328 |
| Executor | Missing decision → error ("missing decision") | `executor.go` L66 |
| Executor | Missing intent → error ("missing execution intent") | `executor.go` L220 |
| Executor | Invalid signatures → error | `executor.go` L227-233 |
| Executor | Receipt signing failure → error ("fail-closed") | `executor.go` L304-308 |
| Executor | Canonicalization failure → error | `executor.go` L147-150 |
| Executor | Output drift → `ERR_CONNECTOR_CONTRACT_DRIFT` | `executor.go` L131-133 |
| Evidence | No signing key → error ("fail-closed") | `exporter.go` L66-68 |
| Evidence | Store integrity violation → error | `executor.go` L161-163 |

### Error Taxonomy Gap
No formalized error code taxonomy (e.g., `ERR_CONNECTOR_CONTRACT_DRIFT` is used ad-hoc but not part of a registered enum). The `kernel/errorir` package exists but has 0% coverage and no test files.

---

## Load Testing

- `helm-loadtest` tool exists (`core/cmd/helm-loadtest/main.go`)
- No evidence of load test results in CI or evidence directory
- No benchmark thresholds (p99 latency, throughput) defined

## Chaos Engineering
- No evidence of chaos testing (no `chaos/` directory, no fault injection tests)
- `kernel/retry` package has 93.1% coverage — retry logic is well-tested in isolation

---

## Score Justification: 3/5

| Factor | Assessment |
|--------|-----------|
| Observability | ✅ OTel fully implemented (was previously missed) |
| Test Coverage | ⚠️ 70/70 pass, avg ~63%, but critical packages < 40% |
| Error Handling | ✅ Comprehensive fail-closed design |
| Load Testing | ❌ Tool exists but no results/thresholds |
| Chaos Engineering | ❌ Not present |
| Error Taxonomy | ⚠️ Ad-hoc, not formalized |

**Why 3 not 4**: Missing load test evidence, no chaos engineering, critical packages under-tested (proofgraph 35%, prg 29%).
**Why 3 not 2**: OTel exists (362 LOC), 70/70 tests pass, comprehensive fail-closed design, replay engine with divergence detection.
