# 02 — Scope & Threat Model

**Score: 5/5** · Gate ≥3 · **✅ PASS**

---

## Scope Definition (Verified)

### What HELM OSS Is (Code-Verified)
| Claim | Package | LOC | Status |
|-------|---------|-----|--------|
| Fail-closed execution kernel | `executor/executor.go` | 341 | ✅ Every stage returns error, never silently passes |
| OpenAI-compatible proxy | `api/`, OpenAPI spec | — | ✅ `/v1/chat/completions` |
| Ed25519 signed receipts | `crypto/signer.go` | 146 | ✅ `ed25519.Sign()` |
| JCS-canonical serialization | `canonicalize/jcs.go` | 139 | ✅ RFC 8785 compliant |
| Lamport-ordered ProofGraph DAG | `proofgraph/graph.go` | 130 | ✅ `sync.Mutex`, parent-linking |
| WASI sandbox | `runtime/sandbox/sandbox.go` | 196 | ✅ wazero-based, memory/time/output limits |
| Conformance harness | `conform/` | 74 files | ✅ L1/L2 gates |

### Explicit Non-Goals (Unchanged)
- ❌ Agent framework / orchestration / planning
- ❌ Content safety / prompt injection in text domain
- ❌ Self-improving or autonomous system
- ❌ Upstream LLM provider security
- ❌ Host OS / hardware side-channel defense

---

## Trust Boundaries (Verified from Source)

```
┌─────────────────────────────────────────────────────┐
│                    UNTRUSTED                        │
│  LLM Provider · User Prompts · Connector Outputs    │
│  WASM Modules · External APIs · MCP Messages        │
└───────────────────────┬─────────────────────────────┘
                        │
                  ┌─────▼─────┐
                  │   PEP     │  JCS + SHA-256 schema validation
                  │ Boundary  │  Guardian default deny (PRG + CEL PDP)
                  │           │  SafeExecutor (Ed25519 signed receipts)
                  │           │  PolicyFirewall (JSON Schema Draft 2020)
                  │           │  Perimeter enforcer (JSON-config)
                  │           │  JWT Auth middleware (fail-closed)
                  │           │  ReBAC authorization (Zanzibar-style)
                  └─────┬─────┘
                        │
┌───────────────────────▼─────────────────────────────┐
│                    TRUSTED                          │
│  Signed Receipt Store · ProofGraph DAG              │
│  Trust Registry (event-sourced) · EvidencePack      │
└─────────────────────────────────────────────────────┘
```

### Trust Assumptions (Source-Verified)

| Actor | Trust Level | Code Evidence |
|-------|-------------|---------------|
| LLM Provider | **Untrusted** | Guardian blocks all undeclared tools (default deny) |
| Tool Arguments | **Untrusted** | JCS canonical + SHA-256 ArgsHash binding (`canonical.go` L53) |
| Connector Outputs | **Untrusted** | Output drift detection (`executor.go` L127-136) |
| WASM Modules | **Untrusted** | wazero sandbox: MemoryLimitBytes, CPUTimeLimit, OutputMaxBytes 1MB (`sandbox.go` L27-31) |
| Operator | **Semi-trusted** | Controls config via PerimeterPolicy JSON; cannot forge receipts |
| Receipt Store | **Trusted** | Ed25519 signed, PrevHash-chained, Lamport-ordered |
| Human Approver | **Trusted** | Timelock + ceremony (`escalation/` package) |

---

## Fail-Closed Conditions (Source-Verified)

| Condition | Reason Code | Package | Evidence |
|-----------|-------------|---------|----------|
| Unknown tool URN | `DENY_TOOL_NOT_FOUND` | firewall | `firewall.go` L58: `"not in allowlist"` |
| Schema validation (input) | `DENY_SCHEMA_MISMATCH` | firewall | `firewall.go` L67: `schema.Validate(params)` |
| Schema validation (output) | `ERR_CONNECTOR_CONTRACT_DRIFT` | executor | `executor.go` L127-136 |
| Budget exceeded | `DENY_BUDGET_EXCEEDED` | guardian | `guardian.go` L120-130 |
| Approval required | `DENY_APPROVAL_REQUIRED` | guardian | `guardian.go` L131 |
| Approval timeout | `DENY_APPROVAL_TIMEOUT` | escalation | escalation package |
| WASI gas exhaustion | `ERR_COMPUTE_TIME_EXHAUSTED` | sandbox | `sandbox.go` L126-129 |
| WASI memory exceeded | `ERR_COMPUTE_MEMORY_EXHAUSTED` | sandbox | `sandbox.go` L131-134 |
| WASI output exceeded | `ERR_COMPUTE_OUTPUT_EXHAUSTED` | sandbox | `sandbox.go` L146-149 |
| Policy violation (PDP) | `DENY` (default) | governance | `pdp.go` L242-244: "fail-closed (DENY) if no policy matches" |
| Decision expired | Gating check | executor | `executor.go` L236-247 |
| Idempotency duplicate | Return existing receipt | executor | `executor.go` L206-213 |
| Auth not configured | `401 Unauthorized` | auth | `middleware.go` L103: nil validator → reject all |
| Perimeter violation | `ErrPolicyDenied` | boundary | `perimeter.go` L35 |

**14 fail-closed conditions verified** — all with deterministic error responses.

---

## WASI Sandbox (Verified: `runtime/sandbox/sandbox.go`, 196 LOC)

### Implementation: wazero (pure Go WebAssembly runtime)

```go
type SandboxConfig struct {
    MemoryLimitBytes int64         // Converted to WASM pages (64KB each)
    CPUTimeLimit     time.Duration // context.WithTimeout
    AllowedSyscalls  []string
    NetworkEnabled   bool          // Always false (WASI deny-by-default)
}
```

### Execution Pipeline
1. Fetch WASM binary from Artifact Store by hash
2. Set context deadline from CPUTimeLimit
3. Configure stdin/stdout/stderr capture (1MB limit)
4. Compile module
5. Instantiate with enforced limits
6. Return stdout bytes (or deterministic error)

### Error Codes
| Code | Trigger |
|------|---------|
| `ERR_COMPUTE_TIME_EXHAUSTED` | context.DeadlineExceeded |
| `ERR_COMPUTE_MEMORY_EXHAUSTED` | wazero memory.grow failure |
| `ERR_COMPUTE_OUTPUT_EXHAUSTED` | stdout+stderr > 1MB |

### Security Notes
- `InProcessSandbox` exists for dev — clearly marked `WARNING: NOT SECURE. DO NOT USE IN PRODUCTION` (L33)
- WASI Preview 1 only (no filesystem, no network by default)
- Uses `wasi_snapshot_preview1.Instantiate()` for WASI syscall emulation (L77)

---

## Threat Categories (OWASP LLM Top 10 + Code Evidence)

| OWASP | HELM Relevance | Code Defense | Package |
|-------|---------------|--------------|---------|
| LLM01: Prompt Injection | IN SCOPE (tool outputs) | Output drift detection | `executor` |
| LLM02: Insecure Output | IN SCOPE | Schema validation + drift check | `executor`, `firewall` |
| LLM03: Training Data | OUT OF SCOPE | — | — |
| LLM04: Model DoS | PARTIAL | Budget governors (tokens, cost, time) | `budget`, `guardian` |
| LLM05: Supply Chain | IN SCOPE | SBOM, SLSA L2, signed releases | `release.yml` |
| LLM06: Sensitive Info | PARTIAL | Data classification constraints | `boundary` (RedactPatterns) |
| LLM07: Insecure Plugin | IN SCOPE | Schema pinning, allowlist, WASI | `firewall`, `sandbox` |
| LLM08: Excessive Agency | IN SCOPE | Budget, approval ceremonies, perimeter | `guardian`, `boundary` |
| LLM09: Overreliance | OUT OF SCOPE | — | — |
| LLM10: Model Theft | OUT OF SCOPE | — | — |

---

## Score: 5/5

**Justification:**
- ✅ All 14 fail-closed conditions verified with code line references
- ✅ 7 trust assumptions mapped to specific enforcement code
- ✅ WASI sandbox verified with wazero implementation, 3 deterministic error codes
- ✅ OWASP LLM Top 10 mapping with package-level evidence
- ✅ Clear scope boundaries with non-goals documented
- ✅ Multi-layer PEP boundary (Guardian + PDP + Firewall + Perimeter + Auth + ReBAC)
