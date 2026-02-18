# 06 — ProofGraph & Receipts

**Score: 4/5** · Gate ≥4 · **✅ PASS**

---

## Receipt Structure (Verified: `contracts/receipt.go`, 62 LOC)

```go
type Receipt struct {
    ReceiptID           string         `json:"receipt_id"`
    DecisionID          string         `json:"decision_id"`
    EffectID            string         `json:"effect_id"`
    ExternalReferenceID string         `json:"external_reference_id"`
    Status              string         `json:"status"`
    BlobHash            string         `json:"blob_hash,omitempty"`    // Input Snapshot CAS
    OutputHash          string         `json:"output_hash,omitempty"`  // Tool Output CAS
    Timestamp           time.Time      `json:"timestamp"`
    Signature           string         `json:"signature,omitempty"`
    // V2: Tamper-Evidence
    MerkleRoot        string             `json:"merkle_root,omitempty"`
    WitnessSignatures []WitnessSignature `json:"witness_signatures,omitempty"`
    // V3: Causal chain
    PrevHash     string `json:"prev_hash"`            // Previous receipt's signature
    LamportClock uint64 `json:"lamport_clock"`        // Monotonic logical clock
    ArgsHash     string `json:"args_hash,omitempty"`   // SHA-256 of JCS(tool_args)
    // Extensions
    ReplayScript     *ReplayScriptRef   `json:"replay_script,omitempty"`
    Provenance       *ReceiptProvenance `json:"provenance,omitempty"`
    BundledArtifacts []ParsedArtifact   `json:"bundled_artifacts,omitempty"`
}
```

### Dual-Hash Design
- `BlobHash`: SHA-256 of **input snapshot** (stored in CAS)
- `OutputHash`: SHA-256 of **canonicalized tool output**
- `ArgsHash`: SHA-256 of **JCS-canonicalized tool arguments** (PEP boundary binding)
- Enables independent verification: verify inputs, outputs, and arguments separately

### Signature Payload (Verified: `crypto/canonical.go` L53-54)
```go
CanonicalizeReceipt(receiptID, decisionID, effectID, status, outputHash, prevHash, lamportClock, argsHash)
// Format: "receipt_id:decision_id:effect_id:status:output_hash:prev_hash:lamport_clock:args_hash"
```

**All 8 fields** are bound into the signed payload — comprehensive coverage.

---

## ProofGraph DAG (Verified: `proofgraph/`, 253 LOC total)

### Node Types (`proofgraph/node.go` L17-24)
```
INTENT → ATTESTATION → EFFECT → TRUST_EVENT → CHECKPOINT → MERGE_DECISION
```

### Node Structure
| Field | Purpose |
|-------|---------|
| `NodeHash` | SHA-256 of JCS-like serialization (excluding self) |
| `Kind` | Node type enum |
| `Parents` | List of parent node hashes (DAG edges) |
| `Lamport` | Monotonic counter across all nodes |
| `Principal` | Identity of the actor |
| `PrincipalSeq` | Per-principal sequence number |
| `Payload` | Arbitrary JSON data |
| `Sig` | Optional signature |
| `Timestamp` | Unix milliseconds |

### Graph Operations (`proofgraph/graph.go`, 130 LOC)
| Operation | Description | Thread Safety |
|-----------|-------------|---------------|
| `Append()` | Add node, link to current heads, increment Lamport | `sync.Mutex` |
| `AppendSigned()` | Add node with signature, recompute hash | `sync.Mutex` |
| `Get()` | Retrieve node by hash | `sync.RWMutex` |
| `Heads()` | Current DAG tips (copied) | `sync.RWMutex` |
| `ValidateChain()` | Walk parents recursively, verify all hashes | `sync.RWMutex` |
| `AllNodes()` | Export all nodes | `sync.RWMutex` |

### Causal Properties
- **Lamport clock**: Incremented on every `Append()` (L29)
- **Multi-head**: Supports concurrent branches (DAG, not chain)
- **PrevHash link**: In Receipt, `PrevHash = prev.Signature` (executor.go L284) — cryptographic link to previous receipt
- **Hash integrity**: `Validate()` recomputes hash and compares (node.go L95-101)

---

## Receipt Creation Pipeline (Verified: `executor.go` L270-313)

```
createReceipt():
  1. Query previous receipt for session (GetLastForSession)
  2. Set prevHash = prev.Signature (GENESIS if first)
  3. Set lamportClock = prev.LamportClock + 1
  4. Construct Receipt with BlobHash + OutputHash + ArgsHash
  5. Sign receipt (fail-closed)
  6. Store receipt
```

### Idempotency (`executor.go` L206-213)
- Before execution, checks if receipt already exists for DecisionID
- If found, returns existing receipt (no re-execution)

---

## Evidence Exporter (Verified: `evidence/exporter.go`, 174 LOC)

### Bundle Types
| Type | Purpose |
|------|---------|
| `SOC2_AUDIT` | Compliance evidence packaging |
| `INCIDENT_REPORT` | Incident diagnostics packaging |

### Sealing Process (`sealBundle()`, L134-172)
1. Sort artifacts by name (deterministic)
2. Create payload with ID, Type, TraceID, Timestamp, artifact hashes
3. Marshal payload → compute SHA-256 → set `BundleHash` and `SignatureMessageHash`
4. Sign the marshaled bytes → set `Signature`

###The evidence bundle sealing was using `json.Marshal`. This has been **Fixed (SEC-002)** by using `canonicalize.JCS`. (see Security doc)
- Fail-closed: returns error if signing key not configured (L66-68, L99-101)

---

## Replay Engine (Verified: `replay/engine.go`, 256 LOC)

### Session Lifecycle
```
RUNNING → COMPLETE (all steps match)
RUNNING → DIVERGED (output hash mismatch at step N)
RUNNING → FAILED (execution error at step N)
```

### Hash Verification
- `computeRunHash()`: SHA-256 of array of all event PayloadHashes
- `computeReplayHash()`: SHA-256 of array of all step OutputHashes
- `VerifyReplayIntegrity()`: Compares original vs replay hashes

### PRNG Support
- `RunEvent.PRNGSeed` field present — enables deterministic replay of random-dependent operations

---

## Coverage

| Package | Coverage | Status |
|---------|----------|--------|
| `proofgraph` | 35.3% | ⚠️ Under-tested for core trust component |
| `receipts` | [no statements] | Conformance tests exist in separate file |
| `receipts/policies` | 67.6% | OK |
| `evidence` | Tested via `overall_test.go` and `registry_test.go` | OK |
| `replay` | 59.1% | ⚠️ Should be higher for audit component |
| `merkle` | 94.1% | ✅ |

---

## Gaps for Score 5/5

| Gap | Impact | Effort |
|-----|--------|--------|
| ProofGraph node hashing used json.Encoder (Fixed) | Fixed (SEC-001) | — |
| No transparency log anchoring (RFC 9162) | Receipts not independently verifiable without HELM | 2 weeks |
| No EvidencePack-level signature | Bundle integrity relies on individual receipt sigs | 1 week |
| ProofGraph coverage 35.3% | Core trust component under-tested | 1 week |
| No witness quorum verification | WitnessSignatures field exists but no quorum logic | 1 week |
