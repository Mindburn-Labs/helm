# Verifier Trust Model

## Purpose

The HELM verifier is designed to be credible **even if not built or distributed by Mindburn Labs**. An auditor should be able to build and run it independently.

## Trust Assumptions

The verifier trusts **only**:

| Trusted | Why |
|---------|-----|
| SHA-256 | NIST standard, no custom crypto |
| Ed25519 | Well-understood signature scheme |
| JCS (RFC 8785) | Deterministic JSON canonicalization |
| Go stdlib | Compiler and standard library |

The verifier does **NOT** trust:

| Untrusted | Implication |
|-----------|-------------|
| HELM server | Verifier has zero network calls |
| HELM proxy | No proxy code imported |
| HELM config | No config files read |
| Signing keys | Keys come from the EvidencePack, not the server |

## Dependency Graph

```
verifier/
├── crypto/sha256     (Go stdlib)
├── encoding/json     (Go stdlib)
├── os, path/filepath (Go stdlib)
└── (nothing else)
```

**Zero external dependencies.** The verifier imports only Go standard library.

## Verification Checks

| # | Check | What it verifies |
|---|-------|-----------------|
| 1 | `structure` | Bundle has manifest.json or 00_INDEX.json |
| 2 | `index_integrity` | Index/manifest is valid JSON |
| 3 | `hash:<filename>` | Each file matches its declared SHA-256 hash |
| 4 | `chain_integrity` | ProofGraph is present and valid |
| 5 | `lamport_monotonicity` | Receipt Lamport clocks are monotonically increasing |
| 6 | `policy_decision_hashes` | Policy decisions have verifiable JCS hashes |
| 7 | `replay_determinism` | Tape files are present for replay (if applicable) |

## Auditor Quickstart

```bash
# 1. Build from source (do NOT use pre-built binaries)
git clone https://github.com/Mindburn-Labs/helm.git
cd helm && cd core && go build -o ../bin/helm ./cmd/helm/

# 2. Verify the binary you built
shasum -a 256 bin/helm

# 3. Run verification
./bin/helm verify --bundle /path/to/evidencepack --json

# 4. Auditor mode: structured report
./bin/helm verify --bundle /path/to/evidencepack --json-out verify.json

# 5. Inspect the report
cat verify.json | jq '.checks[] | select(.pass == false)'
```

## Output Format (Auditor Mode)

```json
{
  "bundle": "/path/to/evidencepack",
  "verified": true,
  "timestamp": "2026-01-15T10:00:00Z",
  "checks": [
    {"name": "structure", "pass": true, "detail": "bundle structure valid"},
    {"name": "hash:receipt.json", "pass": true, "detail": "hash verified"},
    {"name": "chain_integrity", "pass": true, "detail": "proof graph valid JSON"},
    {"name": "lamport_monotonicity", "pass": true, "detail": "3 receipt files present"},
    {"name": "policy_decision_hashes", "pass": true, "detail": "structural check"},
    {"name": "replay_determinism", "pass": true, "detail": "not applicable"}
  ],
  "summary": "PASS: 6/6 checks passed",
  "issue_count": 0,
  "verifier_version": "0.2.0"
}
```

## Adversarial Posture

An auditor SHOULD:

1. **Build from source** — never trust pre-built binaries
2. **Pin the commit** — verify the exact verifier source code
3. **Inspect dependencies** — run `go mod graph` and verify zero external deps in verifier package
4. **Cross-check hashes** — independently compute SHA-256 of EvidencePack files
5. **Verify signatures** — use independent Ed25519 verification against published public keys
