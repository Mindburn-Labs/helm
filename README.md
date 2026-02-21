# HELM — Fail-Closed Tool Calling for AI Agents

[![Build](https://github.com/Mindburn-Labs/helm/actions/workflows/helm_core_gates.yml/badge.svg)](https://github.com/Mindburn-Labs/helm/actions/workflows/helm_core_gates.yml)
[![Conformance](https://img.shields.io/badge/conformance-L1%20%2B%20L2-brightgreen)](docs/CONFORMANCE.md)
[![Verify](https://img.shields.io/badge/verify-npx%20%40mindburn%2Fhelm-blue)](docs/verify.md)
[![Provenance](https://img.shields.io/badge/provenance-SLSA-blue)](https://github.com/Mindburn-Labs/helm/releases)

**Models propose. The kernel disposes.**

HELM is a high-performance, deterministic proxy for LLM tool calling. It enforces mathematical and legal boundaries on AI agents in real-time, generating a tamper-proof **ProofGraph** of every decision.

**[Read HELM for Humans 🧠](docs/FOR_HUMANS.md)** — A non-technical overview of why this exists.

---

## 🚀 SOTA 2026 Quickstart

Install the HELM CLI and start governing in 60 seconds:

```bash
curl -fsSL https://raw.githubusercontent.com/Mindburn-Labs/helm/main/install.sh | bash
helm server
```

*No Postgres required! HELM auto-provisions a local SQLite database and persistent trust root by default.*

---

## 📊 Performance

HELM is built for high-stakes, low-latency environments. To measure the overhead on your machine:

```bash
./scripts/bench/latency.sh
```

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

## 🔍 Verify Any Release

```bash
npx @mindburn/helm
```

One command, progressive disclosure, cryptographic proof. Supports interactive and CI modes:

```bash
# CI mode — JSON on stdout, exit code 0/1
npx @mindburn/helm --ci --bundle ./evidence 2>/dev/null | jq .verdict
```

→ Full guide: [docs/verify.md](docs/verify.md)

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

Full scope details in [docs/OSS_SCOPE.md](docs/OSS_SCOPE.md)

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
├── packages/
│   └── mindburn-helm-cli/  # @mindburn/helm v3 (npm CLI)
├── sdk/                # Multi-language SDKs (TS, Python, Go, Rust, Java)
├── examples/           # Runnable examples per language + MCP
├── scripts/            # Release, CI, SDK generation
├── deploy/             # Caddy config, demo compose, deploy guide
├── docs/               # Threat model, quickstart, verify, conformance
└── Makefile            # build, test, crucible, demo, release-binaries
```
---

## Scope and Guarantees

OSS targets L1/L2 core conformance. The spec contains L3 and enterprise extensions — see [docs/OSS_SCOPE.md](docs/OSS_SCOPE.md) for the shipped-vs-spec boundary.

---

## Security Posture

- **TCB isolation gate** — 8-package kernel boundary, CI-enforced forbidden imports ([TCB Policy](docs/TCB_POLICY.md))
- **Bounded compute gate** — WASI sandbox with gas/time/memory caps, deterministic traps on breach ([UC-005](docs/use_cases/UC-005_wasi_gas_exhaustion.sh))
- **Schema drift fail-closed** — JCS canonicalization + SHA-256 on every tool call, both input and output ([UC-002](docs/use_cases/UC-002_schema_mismatch.sh))

See also: [SECURITY.md](SECURITY.md) (vulnerability reporting) · [Threat Model](docs/THREAT_MODEL.md) (9 adversary classes)

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Good first issues: conformance improvements, SDK enhancements, docs truth fixes.

## Roadmap

See [docs/ROADMAP.md](docs/ROADMAP.md). 10 items, no dates, each tied to a conformance level.

## License

[Apache License 2.0](LICENSE)

---

Built by **[Mindburn Labs](https://mindburn.org)**.