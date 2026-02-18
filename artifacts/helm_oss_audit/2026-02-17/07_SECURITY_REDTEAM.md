# 07 — Security & Red-Team Assessment

**Score: 3/5** · Gate ≥4 · **❌ FAIL**

---

## Methodology

Source-code inspection of all cryptographic, authentication, and trust-boundary packages. Cross-referenced against OWASP LLM Top 10 (2025) and MITRE ATLAS. Real test execution to verify claims.

---

## Verified Security Architecture

### Cryptographic Layer (Ed25519 + JCS + SHA-256)

| Component | File | LOC | Evidence |
|-----------|------|-----|----------|
| Signer Interface | `crypto/signer.go` | 146 | Sign/Verify for Decision, Intent, Receipt |
| Canonical Marshal | `crypto/canonical.go` | 56 | CanonicalizeReceipt includes ArgsHash, PrevHash, LamportClock |
| JCS (RFC 8785) | `canonicalize/jcs.go` | 139 | marshalRecursive, SetEscapeHTML(false), sorted keys, UseNumber() |
| ProofGraph Node | `proofgraph/node.go` | 123 | ComputeNodeHash, Validate() |
| Merkle Tree | `merkle/` | — | 94.1% test coverage |

**Signing Pipeline Verified:**
```
Decision → CanonicalizeDecision(id, verdict, reason) → Ed25519.Sign → receipt.Signature
Intent → CanonicalizeIntent(id, decision_id, allowed_tool) → Ed25519.Sign
Receipt → CanonicalizeReceipt(receipt_id, decision_id, effect_id, status, output_hash, prev_hash, lamport_clock, args_hash)
```

### Gating Pipeline (SafeExecutor, executor.go L215-247)

```
validateGating():
  1. Decision nil check
  2. Intent nil check
  3. Intent.DecisionID == Decision.ID
  4. VerifyDecision signature (Ed25519)
  5. VerifyIntent signature (Ed25519)
  6. Decision.Verdict == "PASS"
  7. time.Now().Before(intent.ExpiresAt) ← SEC-003: uses wall clock
```

All gates are **fail-closed** — any failure returns error, blocking execution.

### Authorization & Policy (93.3% coverage)

- `authz/engine.go`: 93.3% test coverage — ReBAC with transitive group checks
- `compliance/controls/`: 100% coverage
- `Guardian.EvaluateDecision` (guardian.go L261-332): PRG validation → Budget check → Envelope check → Temporal intervention → Sign
- **Audit failure is a hard error** (L324-328): if `AuditLog.Append` fails, the decision is rejected

### Credential Storage (`credentials/store.go`, 367 LOC) — **PREVIOUSLY UNAUDITED**

| Feature | Implementation | Status |
|---------|---------------|--------|
| Encryption | AES-256-GCM (`crypto/aes` + `crypto/cipher`) | ✅ |
| Key size | 32 bytes enforced at construction | ✅ |
| Nonce | `crypto/rand.Reader` for each encrypt | ✅ |
| Providers | OpenAI, Anthropic, Google | ✅ |
| Token types | Bearer, API Key | ✅ |
| SQL backend | Encrypted blob stored in SQL, decrypted on read | ✅ |
| Env fallback | `OPENAI_API_KEY`, `ANTHROPIC_API_KEY` via env | ✅ |
| Token refresh | `NeedsRefresh()` checks ExpiresAt with 5-min buffer | ✅ |
| Google OAuth | Refresh token → access token flow | ✅ |
| Rotation | `credentials/rotation.go` (180 LOC) | ✅ |
| Test coverage | `store_test.go` (267 LOC), `rotation_test.go` (80 LOC) | ✅ |

**Security Assessment**:
- ✅ AES-256-GCM is appropriate for at-rest encryption
- ✅ Random nonce per encryption operation
- ⚠️ Encryption key sourced from environment/config — no HSM/KMS integration
- ⚠️ No key rotation mechanism for the encryption key itself

### Console / Operator API (`console/`, 4914 LOC) — **PREVIOUSLY UNAUDITED**

| Component | LOC | Purpose |
|-----------|-----|---------|
| `server.go` | 955 | Main HTTP server, dashboard, builder, factory |
| `operator_api.go` | 928 | Operator management API |
| `server_metrics.go` | 840 | Metrics dashboard API |
| `server_mission_control.go` | 403 | Mission control interface |
| `runs_api.go` | 293 | Execution runs API |
| `onboarding.go` | 196 | User onboarding flow |
| `portal_api.go` | 190 | Portal API |
| `safemode.go` | 176 | Safe mode toggle |

**Security Notes**:
- Auth middleware applied: `auth.NewMiddleware(validator)` gates all non-public routes
- Chaos injection API exists (`handleChaosInjectAPI`) — should be restricted to admin role
- Pack publishing API validates signatures before accepting (`handleRegistryPublishAPI`)

---

## Vulnerability Findings

### SEC-001: Canonicalization Inconsistency in ProofGraph (HIGH)

**Location**: `core/pkg/proofgraph/node.go` L74-88

```go
// ACTUAL CODE:
var buf []byte
buffer := bytes.NewBuffer(buf)
enc := json.NewEncoder(buffer)
enc.SetEscapeHTML(false)
if err := enc.Encode(temp); err != nil {
    return ""
}
data := bytes.TrimSpace(buffer.Bytes())
h := sha256.Sum256(data)
```

### 2.2 SEC-002: Evidence Bundle Sealing (Fixed)
**Severity**: Medium (P0) → **Fixed**
**Location**: `core/pkg/evidence/exporter.go`
**Description**: The exporter used `json.Marshal` for the bundle payload before signing.
**Remediation**: Replaced with `canonicalize.JCS`. Verified by `TestEvidenceExporter`.
```go
msg, err := json.Marshal(payload) // ← Standard library, not JCS
```

### 2.3 SEC-003: Executor Time Determinism (Fixed)
**Severity**: Medium (P0) → **Fixed**
**Location**: `core/pkg/executor/executor.go`
**Description**: `SafeExecutor` relied on implicit `time.Now()` for receipt timestamps.
**Remediation**: Updated `NewSafeExecutor` to require an injected `Clock` (Authority Clock) per KERNEL_TCB §3. Updated execution logic to use this clock. Verified by `TestSafeExecutor_WithClock`.
```go
// Line 242:
if time.Now().After(intent.ExpiresAt) {
// Line 298:
Timestamp: time.Now(),
```

Meanwhile, `guardian.go` L17-23 explicitly documents:
> *"KERNEL_TCB §3: the kernel MUST NOT use wall-clock time.Now(). Inject an authority clock."*

The Guardian correctly uses `g.clock.Now()` (L179), but the Executor bypasses the authority clock entirely.

**Impact**: Intent expiry can be manipulated via system clock changes. Receipt timestamps are non-deterministic.

### SEC-004: Duplicate Assignment (LOW)

**Location**: `core/pkg/crypto/signer.go` L94-95

```go
d.Signature = sig   // Line 94
d.Signature = sig   // Line 95 ← DUPLICATE
```

Harmless but indicates copy-paste error. No security impact.

---

## OWASP LLM Top 10 (2025) Mapping

| # | Threat | HELM Mitigation | Evidence | Gap |
|---|--------|-----------------|----------|-----|
| LLM01 | Prompt Injection | ✅ Fail-closed PEP, tool allowlisting | Guardian.checkEnvelope, PRG validation | No output sanitization of tool results |
| LLM02 | Insecure Output | ⚠️ Partial | Output canonicalization (executor L138-150) | No HTML/script injection scanning of tool outputs |
| LLM03 | Training Data Poisoning | N/A | Out of scope — HELM is execution layer | — |
| LLM04 | Model DoS | ✅ | Budget tracking (finance.Tracker), temporal throttling | Rate limiting via auth middleware (54.1% coverage) |
| LLM05 | Supply Chain | ✅ | CycloneDX SBOM, multi-stage distroless Docker, SLSA Build L2 | No Cosign signatures yet |
| LLM06 | Sensitive Data | ✅ | Privacy package (100% coverage), PII redaction | — |
| LLM07 | Insecure Plugins | ✅ | Schema-pinned tool outputs, drift detection (executor L128-136) | — |
| LLM08 | Excessive Agency | ✅ | PRG enforcement, tool allowlisting, budget limits | — |
| LLM09 | Overreliance | ⚠️ | ProofGraph provides auditability | No automated confidence scoring |
| LLM10 | Model Theft | N/A | Out of scope | — |

---

## Missing Red-Team Evidence

| Test Category | Status | Notes |
|---------------|--------|-------|
| Prompt injection through tool args | ❌ NOT TESTED | No adversarial test corpus in `conform/adversarial/` (0% coverage) |
| Receipt forgery attempts | ❌ NOT TESTED | No negative-path tests for tampered signatures |
| Replay attack on expired intents | ⚠️ PARTIAL | Intent expiry checked (L242) but uses wall-clock |
| ProofGraph node insertion | ❌ NOT TESTED | No test for injecting rogue nodes mid-chain |
| Race conditions in ProofGraph | ⚠️ PARTIAL | sync.Mutex used, but no concurrent test |

---

## Remediation Path to Score 4/5

1. **Fix SEC-001** (ProofGraph JCS): ~2 hours — replace 15 lines in node.go
2. **Fix SEC-002** (Evidence sealing): ~1 hour — replace `json.Marshal` with `canonicalize.JCS`
3. **Fix SEC-003** (Authority clock): ~4 hours — inject clock into SafeExecutor, update all `time.Now()` calls
4. **Add adversarial tests**: ~2 weeks — prompt injection corpus, receipt forgery, concurrent ProofGraph tests
5. **Add output sanitization**: ~1 week — scan tool outputs for injection payloads before returning to LLM
