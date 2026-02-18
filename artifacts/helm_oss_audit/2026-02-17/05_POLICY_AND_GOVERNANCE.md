# 05 — Policy & Governance

**Score: 4/5** · Gate ≥3 · **✅ PASS**

---

## Policy Architecture (Multi-Layer, NOT Hardcoded)

> [!IMPORTANT]
> Contrary to prior audit claims, HELM **does** have an externalized policy engine using CEL (Common Expression Language) and a configurable JSON-based perimeter policy. The architecture is significantly more advanced than previously documented.

### Layer 1: Guardian (Legacy PEP) — `guardian/guardian.go`, 347 LOC
- Hardcoded Go logic for PRG validation, budget, temporal, envelope checks
- Sets `PolicyVersion: "v1.0.0"` (L296) — static string
- **This is the legacy path**, still active in the execution pipeline

### Layer 2: CEL Policy Decision Point — `governance/pdp.go`, 324 LOC
- Full `PolicyDecisionPoint` interface (L20-27):
  - `Evaluate(ctx, PDPRequest) → PDPResponse` — deterministic
  - `PolicyVersion() → string` — content-addressed policy bundle hash
- `CELPolicyDecisionPoint` implementation (L182-309):
  - **Fail-closed**: Default decision = `DENY` (L244)
  - **Deterministic decision IDs**: `uuid.NewSHA1(OID, hash(request))` (L221)
  - **Decision expiry**: 5-minute default (L250)
  - **Effect type classification**: Allow/Deny/RequireApproval (L255-291)
  - **FactOracle integration**: Knowledge graph queries for contextual decisions (L269-278)
  - **Hot-reload**: `UpdatePolicyBundle(newVersionHash)` (L304-309)
- Request hashing: Uses JCS (`compliance/jcs.Marshal`) for canonical hashing (L314)

### Layer 3: CEL Policy Evaluator — `governance/policy_evaluator_cel.go`, 158 LOC
- Uses `google/cel-go` library (real CEL runtime)
- **System rules** ("Constitution"): namespaced module names, semantic versioning (L34-39)
- **Module self-policy**: Each module can carry its own CEL expression (L82-92)
- **Program caching**: Double-checked locking with `sync.RWMutex` (L113-148)
- **Computational limits**: `cel.CostLimit(10000)`, `cel.InterruptCheckFrequency(100)` (L134-135)
- **Fail-closed**: System policy errors → deny (L72-76)

### Layer 4: Swarm PDP — `governance/swarm_pdp.go`, 410 LOC
- **Parallel policy evaluation** with K2.5-style PARL architecture
- Domain decomposition: `authorization`, `compliance`, `risk`, `audit`, `general` (L90-96)
- Batch evaluation with goroutine fan-out (L149-238)
- PARL metrics tracking: critical path, parallelism, batch sizes
- Strict merge mode for conflict resolution (L290-328)
- `PolicyVersion()` delegates to base PDP (L142-147)

### Layer 5: Boundary Perimeter — `boundary/perimeter.go`, 337 LOC
- **JSON-configurable policy** with 4 constraint domains:

```go
type PerimeterPolicy struct {
    Version     string      `json:"version"`
    PolicyID    string      `json:"policy_id"`
    Constraints Constraints `json:"constraints"`
    Enforcement Enforcement `json:"enforcement"`
}
```

| Constraint Domain | Capabilities |
|-------------------|-------------|
| **Network** | AllowedHosts, DeniedHosts, AllowedPorts, MaxRequestsPerMin, MaxBandwidthBytes, RequireTLS |
| **Tools** | AllowedTools, DeniedTools, RequireAttestation, MaxConcurrentCalls, TimeoutSeconds |
| **Data** | AllowedClasses, DeniedClasses, MaxContextTokens, MaxResponseTokens, RedactPatterns |
| **Temporal** | AllowedHours, AllowedDays, MaxExecutionSeconds, CooldownSeconds |

- **Hot-reload**: `LoadPolicy(policy)` with `sync.RWMutex` protection (L137-164)
- **Enforcement modes**: `enforce`, `audit`, `disabled` (L18-22)
- **Fail-closed**: `ErrPolicyDenied` on constraint violation (L35)

### Layer 6: PolicyFirewall — `firewall/firewall.go`, 84 LOC
- Tool allowlist + JSON Schema validation (Draft 2020-12)
- Fail-closed: `nil` dispatcher → error (L74-76)
- Schema compilation and caching per tool (L32-52)

### Layer 7: ReBAC Authorization — `authz/engine.go`, 120 LOC
- Relationship-Based Access Control (Zanzibar-inspired)
- Transitive group membership checks with cycle detection
- RelationTuple: `(Object#Relation@Subject)` (L11-15)

---

## Policy Version Binding (EXISTS — Prior Audit Was Wrong)

| Location | Field | Evidence |
|----------|-------|----------|
| `contracts/decision.go` L24 | `PolicyVersion string` | ✅ In DecisionRecord |
| `governance/pdp.go` L103 | `PolicyVersion string` | ✅ In PDPResponse |
| `governance/pdp.go` L246 | Sets PolicyVersion in response | ✅ |
| `guardian/guardian.go` L296 | `PolicyVersion: "v1.0.0"` | ⚠️ Hardcoded string (legacy) |
| `contracts/evidence.go` L54 | `PolicyVersion string` | ✅ In EvidenceMetadata |
| `executor/evidence_pack.go` L38 | `PolicyVersion string` | ✅ In EvidencePackInput |
| `conform/gates/g1_proof_receipts.go` L49 | `PolicyVersion string` | ✅ REQUIRED in conformance |

**Verdict**: PolicyVersion IS bound to decisions and evidence. The Guardian legacy path hardcodes it, but the governance PDP uses content-addressed hashes.

---

## Decision Trace (PDPResponse)

```go
type PDPResponse struct {
    Decision      Decision            // ALLOW | DENY | REQUIRE_APPROVAL | REQUIRE_EVIDENCE | DEFER
    DecisionID    string              // Deterministic: uuid.NewSHA1(OID, JCS(request))
    PolicyVersion string              // Content-addressed policy bundle hash
    Constraints   DecisionConstraints // Envelope, approvals, evidence, compensations
    Trace         DecisionTrace       // Evaluation graph hash, rules fired, input hashes
    IssuedAt      time.Time
    ExpiresAt     time.Time           // Default: 5 minutes
}
```

The `DecisionTrace` includes:
- `EvaluationGraphHash`: SHA-256 of all input hashes (deterministic)
- `RulesFired`: Ordered list of matched rules (audit trail)
- `InputsHashes`: Map of all input component hashes
- `EngineSubtraces`: For multi-engine PDP setups

---

## Coverage

| Package | Coverage | Status |
|---------|----------|--------|
| `governance` | 60.8% | ✅ Good |
| `authz` (authorization) | 93.3% | ✅ Excellent |
| `firewall` | 36.2% | ⚠️ Low for security component |
| `boundary` (perimeter) | 60.2% | ✅ OK |
| `kernel/pdp` | 100% | ✅ |
| `kernel/celdp` | 71.8% | ✅ Good |

---

## Score: 4/5 (upgraded from prior 3/5)

**Justification:**
- ✅ CEL-based policy engine exists with deterministic evaluation
- ✅ PolicyVersion bound to decisions and evidence
- ✅ JSON-configurable perimeter policy with hot-reload
- ✅ Multi-layer architecture (Guardian + PDP + Perimeter + Firewall + ReBAC)
- ✅ Fail-closed at every layer
- ⚠️ Guardian still hardcodes `"v1.0.0"` instead of content-addressed hash
- ⚠️ CEL rules are system-defined strings, not loaded from external policy bundle files
- ❌ No OPA/Rego support (CEL only)

### To reach 5/5:
1. Replace Guardian's hardcoded `"v1.0.0"` with PDP's content-addressed PolicyVersion
2. Load CEL rules from external policy bundle files (not hardcoded strings)
3. Add policy bundle signing and verification
4. Document the multi-layer policy architecture in a public spec
