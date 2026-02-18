# HELM — Fail-Closed Tool Calling for AI Agents

[![Build](https://github.com/Mindburn-Labs/helm/actions/workflows/helm_core_gates.yml/badge.svg)](https://github.com/Mindburn-Labs/helm/actions/workflows/helm_core_gates.yml)
[![Conformance L1](https://img.shields.io/badge/conformance-L1%20pass-brightgreen)](docs/use_cases/UC-012_openai_proxy.sh)
[![Conformance L2](https://img.shields.io/badge/conformance-L2%20pass-brightgreen)](docs/use_cases/UC-012_openai_proxy.sh)
[![SBOM](https://img.shields.io/badge/SBOM-CycloneDX%201.5-blue)](sbom.json)
[![Provenance](https://img.shields.io/badge/provenance-SLSA-blue)](https://github.com/Mindburn-Labs/helm/releases)

OpenAI-compatible proxy that enforces tool execution and emits verifiable cryptographic receipts.

Spec is broader than v0.1 by design — see [docs/OSS_CUTLINE.md](docs/OSS_CUTLINE.md) for exact shipped guarantees.

- **1-line integration** — swap `base_url`, keep everything else
- **EvidencePack export** — deterministic `.tar.gz`, verify offline, sue-grade
- **Bounded compute** — WASI sandbox with gas/time/memory caps, approval ceremonies with timelocks
- **Canonical Standard** — [HELM Unified Canonical Standard v1.2](HELM_Unified_Canonical_Standard_FINAL_2026-02-15_FINAL_SOTA_v1.2.md)

---

### Quickest path to a receipt

```bash
docker compose up -d && curl -s localhost:8080/v1/chat/completions \
  -H 'Content-Type: application/json' \
  -d '{"model":"gpt-4","messages":[{"role":"user","content":"hello"}]}' | jq .id
```

### Start from source

```bash
git clone https://github.com/Mindburn-Labs/helm.git && cd helm
docker compose up -d
curl -s http://localhost:8080/health   # → OK
```

### Run the proof loop

```bash
make build && make crucible            # 12 use cases + conformance L1/L2
./bin/helm export --evidence ./data/evidence --out pack.tar.gz
./bin/helm verify --bundle pack.tar.gz # offline — no network
```

### One-line integration

```diff
- client = openai.OpenAI()
+ client = openai.OpenAI(base_url="http://localhost:8080/v1")
```

That's it. Your app doesn't change. Every tool call now produces a signed receipt in an append-only DAG.

---

## 5-Minute Proof Loop

**Goal: prove it works without trusting us.** You can verify the EvidencePack and replay without network access.

```bash
# 1. Start
docker compose up -d

# 2. Trigger a deny (schema mismatch → fail-closed)
curl -s http://localhost:8080/v1/tools/execute \
  -H 'Content-Type: application/json' \
  -d '{"tool":"unknown_tool","args":{"bad_field":true}}' | jq .reason_code
# → "ERR_TOOL_NOT_FOUND"

# 3. View receipt
curl -s http://localhost:8080/api/v1/receipts?limit=1 | jq '.[0].receipt_hash'

# 4. Export EvidencePack
./bin/helm export --evidence ./data/evidence --out pack.tar.gz

# 5. Offline replay verify — no network required
./bin/helm verify --bundle pack.tar.gz
# → "verification: PASS"  (air-gapped safe)

# 6. Run conformance L1/L2
./bin/helm conform --profile L2 --json
# → {"profile":"L2","verdict":"PASS","gates":12}
```

Full walkthrough: [docs/QUICKSTART.md](docs/QUICKSTART.md) · [docs/POLICY_BACKENDS.md](docs/POLICY_BACKENDS.md) · [docs/VERIFIER_TRUST_MODEL.md](docs/VERIFIER_TRUST_MODEL.md) · [docs/PROCUREMENT.md](docs/PROCUREMENT.md)

---

## Why Devs Should Care

| Pain (postmortem you're preventing) | HELM behavior | Receipt reason code | Proof |
|------------------------------------|---------------|--------------------|---------|
| Tool-call overspend blows budget | ACID budget locks, fail-closed on ceiling breach | `DENY_BUDGET_EXCEEDED` | [UC-005](docs/use_cases/UC-005_wasi_gas_exhaustion.sh) |
| Schema drift breaks prod silently | Fail-closed on input AND output schema mismatch | `DENY_SCHEMA_MISMATCH` | [UC-002](docs/use_cases/UC-002_schema_mismatch.sh), [UC-009](docs/use_cases/UC-009_connector_drift.sh) |
| Untrusted WASM runs wild | Sandbox: gas + time + memory budgets, deterministic traps | `DENY_GAS_EXHAUSTION` | [UC-004](docs/use_cases/UC-004_wasi_transform.sh) |
| "Who approved that?" disputes | Timelock + challenge/response ceremony, Ed25519 signed | `DENY_APPROVAL_REQUIRED` | [UC-003](docs/use_cases/UC-003_approval_ceremony.sh) |
| No audit trail for regulators | Deterministic EvidencePack, offline verifiable, replay from genesis | — | [UC-008](docs/use_cases/UC-008_replay_verify.sh) |
| Can't prove compliance to auditors | Conformance L1 + L2 gates, 12 runnable use cases | — | [UC-012](docs/use_cases/UC-012_openai_proxy.sh) |

---

## Integrations

### Python — OpenAI SDK

The only change:
```diff
- client = openai.OpenAI()
+ client = openai.OpenAI(base_url="http://localhost:8080/v1")
```

Full snippet:
```python
import openai

client = openai.OpenAI(base_url="http://localhost:8080/v1")

response = client.chat.completions.create(
    model="gpt-4",
    messages=[{"role": "user", "content": "List files in /tmp"}]
)
print(response.choices[0].message.content)
# Response headers include:
#   X-Helm-Receipt-ID: rec_a1b2c3...
#   X-Helm-Output-Hash: sha256:7f83b1...
#   X-Helm-Lamport-Clock: 42
```

→ Full example: [examples/python_openai_baseurl/main.py](examples/python_openai_baseurl/main.py)

### TypeScript — Vercel AI SDK / fetch

The only change:
```diff
- const BASE = "https://api.openai.com/v1";
+ const BASE = "http://localhost:8080/v1";
```

Full snippet:
```typescript
const response = await fetch("http://localhost:8080/v1/chat/completions", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({
    model: "gpt-4",
    messages: [{ role: "user", content: "What time is it?" }],
  }),
});
const data = await response.json();
console.log(data.choices[0].message.content);
// X-Helm-Receipt-ID: rec_d4e5f6...
```

→ Full example: [examples/js_openai_baseurl/main.js](examples/js_openai_baseurl/main.js)

### MCP Gateway

```bash
# List governed capabilities
curl -s http://localhost:8080/mcp/v1/capabilities | jq '.tools[].name'

# Execute a governed tool call
curl -s -X POST http://localhost:8080/mcp/v1/execute \
  -H 'Content-Type: application/json' \
  -d '{"method":"file_read","params":{"path":"/tmp/test.txt"}}' | jq .
# → { "result": ..., "receipt_id": "rec_...", "reason_code": "ALLOW" }
→ Full example: [examples/mcp_client/main.sh](examples/mcp_client/main.sh)

---

## SDKs

Typed clients for 5 languages. All generated from [api/openapi/helm.openapi.yaml](api/openapi/helm.openapi.yaml).

| Language | Installation Command | Package Link |
| :--- | :--- | :--- |
| **TypeScript** | `npm install @mindburn/helm-sdk` | [npm/@mindburn/helm-sdk](https://www.npmjs.com/package/@mindburn/helm-sdk) |
| **Python** | `pip install helm-sdk` | [pypi/helm-sdk](https://pypi.org/project/helm-sdk/) |
| **Go** | `go get github.com/Mindburn-Labs/helm/sdk/go` | [pkg.go.dev](https://pkg.go.dev/github.com/Mindburn-Labs/helm/sdk/go) |
| **Rust** | `cargo add helm-sdk` | [crates.io/helm-sdk](https://crates.io/crates/helm-sdk) |
| **Java** | `implementation 'ai.mindburn.helm:helm-sdk:0.1.0'` | [Maven Central](https://central.sonatype.com/) |

Every SDK exposes the same primitives: `chatCompletions`, `approveIntent`, `listSessions`, `getReceipts`, `exportEvidence`, `verifyEvidence`, `conformanceRun`.

Every error includes a typed `reason_code` (e.g. `DENY_TOOL_NOT_FOUND`).

**Go — 10-line denial-handling example:**

```go
c := helm.New("http://localhost:8080")
res, err := c.ChatCompletions(helm.ChatCompletionRequest{
    Model:    "gpt-4",
    Messages: []helm.ChatMessage{{Role: "user", Content: "List /tmp"}},
})
if apiErr, ok := err.(*helm.HelmApiError); ok {
    fmt.Println("Denied:", apiErr.ReasonCode) // DENY_TOOL_NOT_FOUND
}
```

**Rust:**

```rust
let c = HelmClient::new("http://localhost:8080");
match c.chat_completions(&req) {
    Ok(res) => println!("{:?}", res.choices[0].message.content),
    Err(e) => println!("Denied: {:?}", e.reason_code),
}
```

**Java:**

```java
var helm = new HelmClient("http://localhost:8080");
try { helm.chatCompletions(req); }
catch (HelmApiException e) { System.out.println(e.reasonCode); }
```

Full examples: [examples/](examples/) · SDK docs: [docs/sdks/00_INDEX.md](docs/sdks/00_INDEX.md)

---

## OpenAPI Contract

[api/openapi/helm.openapi.yaml](api/openapi/helm.openapi.yaml) — OpenAPI 3.1 spec.

Single source of truth. SDKs are generated from it. CI prevents drift.

→ [Contract versioning](docs/sdks/contract_versioning.md)

---

## How It Works

```
Your App (OpenAI SDK)
       │
       │ base_url = localhost:8080
       ▼
   HELM Proxy ──→ Guardian (policy: allow/deny)
       │                │
       │           PEP Boundary (JCS canonicalize → SHA-256)
       │                │
       ▼                ▼
   Executor ──→ Tool ──→ Receipt (Ed25519 signed)
       │                        │
       ▼                        ▼
  ProofGraph DAG          EvidencePack (.tar.gz)
  (append-only)           (offline verifiable)
       │
       ▼
  Replay Verify
  (air-gapped safe)
```

---

## What Ships vs What's Spec

| Shipped in OSS v0.1 | Spec (future / enterprise) |
|---------------------|---------------------------|
| ✅ OpenAI-compatible proxy | 🔮 Multi-model gateway |
| ✅ Schema PEP (input + output) | 🔮 ZK-CPI (zero-knowledge proofs) |
| ✅ ProofGraph DAG (Lamport + Ed25519) | 🔮 Hardware TEE attestation |
| ✅ WASI sandbox (gas/time/memory) | 🔮 Post-quantum cryptography |
| ✅ Approval ceremonies (timelock + challenge) | 🔮 Multi-org federation |
| ✅ Trust registry (event-sourced) | 🔮 Formal verification (SMT/LTL) |
| ✅ EvidencePack export + offline replay | 🔮 Cross-tenant ProofGraph merge |
| ✅ Conformance L1 + L2 | 🔮 Conformance L3 (enterprise) |
| ✅ 11 CLI commands | 🔮 Production key management (HSM) |

Full cutline: [docs/OSS_CUTLINE.md](docs/OSS_CUTLINE.md)

---

## Verification

```bash
make test       # 112 packages, 0 failures
make crucible   # 12 use cases + conformance L1/L2
make lint       # go vet, clean
```

---

## Deploy

```bash
# Local demo
docker compose up -d

# Production (DigitalOcean / any Docker host)
docker compose -f docker-compose.demo.yml up -d
```

→ [deploy/README.md](deploy/README.md) — deploy your own in 3 minutes

---

## Project Structure

```
helm/
├── api/openapi/         # OpenAPI 3.1 spec (single source of truth)
├── core/               # Go kernel (8-package TCB + executor + ProofGraph)
│   ├── cmd/helm/       # CLI: proxy, export, verify, replay, conform, ...
│   └── cmd/helm-node/  # Kernel API server
├── sdk/                # Multi-language SDKs (TS, Python, Go, Rust, Java)
│   ├── ts/             #   npm @mindburn/helm-sdk
│   ├── python/         #   pip helm-sdk
│   ├── go/             #   go get .../sdk/go
│   ├── rust/           #   cargo add helm-sdk
│   └── java/           #   mvn ai.mindburn.helm:helm-sdk
├── examples/           # Runnable examples per language + MCP
├── scripts/sdk/        # Type generator (gen.sh)
├── scripts/ci/         # SDK drift + build gates
├── deploy/             # Caddy config, demo compose, deploy guide
├── docs/               # Threat model, quickstart, demo, SDK docs
└── Makefile            # build, test, crucible, demo, release-binaries
```
---

## Scope and Guarantees

OSS v0.1 targets L1/L2 core conformance. Spec contains L2/L3 and enterprise/2030 extensions — see [docs/OSS_CUTLINE.md](docs/OSS_CUTLINE.md) for the exact shipped-vs-spec boundary.

---

## Security Posture

- **TCB isolation gate** — 8-package kernel boundary, CI-enforced forbidden imports ([TCB Policy](TCB_POLICY.md))
- **Bounded compute gate** — WASI sandbox with gas/time/memory caps, deterministic traps on breach ([UC-005](docs/use_cases/UC-005_wasi_gas_exhaustion.sh))
- **Schema drift fail-closed** — JCS canonicalization + SHA-256 on every tool call, both input and output ([UC-002](docs/use_cases/UC-002_schema_mismatch.sh))

See also: [SECURITY.md](SECURITY.md) (vulnerability reporting) · [Threat Model](docs/THREAT_MODEL.md) (9 adversary classes)

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Good first issues: conformance improvements, SDK enhancements, docs truth fixes.

## Roadmap

See [docs/ROADMAP.md](docs/ROADMAP.md). 10 items, no dates, each tied to a conformance level.

## License

[Business Source License 1.1](LICENSE) — converts to Apache 2.0 on 2030-02-15.

---

Built by **[Mindburn Labs](https://mindburn.org)**.