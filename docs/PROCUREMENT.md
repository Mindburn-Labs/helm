# HELM Procurement & Enterprise Adoption Guide

## For Enterprise Evaluators

HELM provides **cryptographic proof of AI governance** — not just logs, but independently verifiable evidence chains.

### What Makes HELM Different

| Feature | HELM | Traditional AI Governance |
|---------|------|--------------------------|
| **Proof Model** | Cryptographic receipts (SHA-256 + Ed25519) | Logs (mutable) |
| **Policy Engine** | Pluggable (HELM, OPA, Cedar) | Vendor lock-in |
| **Verification** | Offline, standalone, zero-trust | Server-dependent |
| **Conformance** | Machine-verifiable gates | Self-attestation |
| **Decision Binding** | JCS-canonical hashes in receipts | None |

### Conformance Levels

| Level | Gates | Suitable For |
|-------|-------|-------------|
| **L1** | G0 (trust roots), G1 (deterministic bytes), G2A (EvidencePack structure) | SMB, initial adoption |
| **L2** | L1 + G2 (ProofGraph), G3A (budget), G5 (replay), G8 (HITL), GX_ENVELOPE, GX_TENANT | Enterprise, regulated |

### Quick Verification

```bash
# Build verifier from source (zero external deps)
git clone https://github.com/Mindburn-Labs/helm.git
cd helm && make build

# Verify any EvidencePack
./bin/helm verify --bundle /path/to/evidencepack --json

# Run conformance suite
./bin/helm conform --level L2 --signed

# Output: conform_report.json + .sha256 + .sig
```

## Integration Paths

### 1. Drop-in Proxy (Fastest)

```bash
# Start HELM proxy in front of any OpenAI-compatible API
./bin/helm proxy --target https://api.openai.com/v1 --port 8080
# All requests/responses get EvidencePack receipts automatically
```

### 2. SDK Integration

- **Python SDK**: `pip install helm-sdk`
- **TypeScript SDK**: `npm install @mindburn/helm-sdk`

### 3. Policy Engine Integration

Bring your own policies:
- **OPA (Rego)**: `export HELM_POLICY_BACKEND=opa`
- **Cedar**: `export HELM_POLICY_BACKEND=cedar`

See [POLICY_BACKENDS.md](POLICY_BACKENDS.md) for full setup.

## Compliance Mapping

| Framework | HELM Coverage |
|-----------|--------------|
| **SOC 2 Type II** | EvidencePack export covers CC6.1, CC7.2 |
| **ISO 27001** | Conformance gates map to A.12, A.14 |
| **EU AI Act** | Proof chain satisfies Article 13 transparency |
| **NIST AI RMF** | Decision receipts cover GOVERN, MAP, MEASURE |

## Procurement Checklist

- [ ] Build HELM from source and verify SHA
- [ ] Run `helm conform --level L2 --signed` and review report
- [ ] Run `helm verify --bundle <pack> --json-out audit.json`
- [ ] Review `docs/VERIFIER_TRUST_MODEL.md` with security team
- [ ] Evaluate policy backend integration (OPA/Cedar)
- [ ] Review `docs/POLICY_BACKENDS.md` for deployment topology
