# HELM RFP Answers

Standard responses for enterprise procurement questionnaires.

---

## 1. How do you ensure AI decisions are auditable?

Every AI decision produces a cryptographically signed receipt containing:
- **Decision hash** (SHA-256 of JCS-canonical decision)
- **Policy backend** (which engine evaluated the policy)
- **Policy content hash** (exact version of policies used)
- **Effect digest** (what the AI was authorized to do)

These receipts are chained into an EvidencePack with Lamport clock ordering and can be verified offline without any HELM server access.

## 2. Can we use our existing policy engine?

Yes. HELM supports pluggable Policy Decision Points:
- **OPA** (Open Policy Agent) — Rego policies via HTTP
- **Cedar** (AWS) — Cedar policies via sidecar
- **HELM** (built-in) — Proof Requirement Graph

Policy backend is selected via `HELM_POLICY_BACKEND` environment variable. The Guardian kernel delegates evaluation while retaining enforcement and signing.

## 3. How do we verify compliance independently?

Build the verifier from source (`core/pkg/verifier/`) — it has **zero external dependencies**.

```bash
./bin/helm verify --bundle /path/to/pack --json-out audit.json
```

The auditor report includes structure checks, file hash integrity, chain integrity, Lamport monotonicity, and policy decision hash verification.

See `docs/VERIFIER_TRUST_MODEL.md` for the full adversarial posture.

## 4. What conformance levels do you support?

| Level | Description |
|-------|-------------|
| **L1** | Trust roots, deterministic bytes, EvidencePack structure |
| **L2** | L1 + ProofGraph, budget enforcement, replay, HITL, envelope, tenant isolation |

Run `helm conform --level L2 --signed` to get a signed conformance report with SHA-256 hash and signature.

## 5. How does HELM handle policy engine failures?

All policy backends enforce **fail-closed semantics**:
- Network timeout → DENY
- HTTP error → DENY
- Parse error → DENY
- Nil request → DENY

The specific deny reason is recorded in the receipt for auditability.

## 6. What is the deployment topology?

```
Your App → HELM Proxy → LLM Provider
               │
               ├── Guardian (enforcement kernel)
               ├── PDP (policy evaluation — HELM/OPA/Cedar)
               ├── ProofGraph (receipt chain)
               └── EvidencePack (exportable audit bundle)
```

HELM runs as a sidecar or reverse proxy. No changes to your application code required for the proxy path.

## 7. What frameworks does HELM map to?

| Framework | Coverage |
|-----------|---------|
| SOC 2 Type II | CC6.1, CC7.2 via EvidencePack |
| ISO 27001 | A.12, A.14 via conformance gates |
| EU AI Act | Article 13 via proof chain transparency |
| NIST AI RMF | GOVERN, MAP, MEASURE via decision receipts |

## 8. Is there vendor lock-in?

No. HELM's design prevents lock-in at every layer:
- **Open source** (Apache 2.0 / MIT)
- **Pluggable policy engines** (OPA, Cedar, or custom)
- **Offline verification** (no server dependency)
- **Standard formats** (JSON, SHA-256, Ed25519, JCS)
- **OpenAPI contract** for integration
