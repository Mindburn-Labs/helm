# OSS-TRUST-REGISTRY Audit

## Scope
Verify build identity, artifact provenance, signing keys via TUF, SLSA, Rekor. Offline trust pack. Signature verifier.

## Reality
- **tuf.go** (2893B), **slsa.go** (3165B), **rekor.go** (2960B): Clients. [CODE]
- **pack.go** (3268B): Offline trust pack. **verifier.go** (3142B): Sig verifier. [CODE]
- Integrated in server wiring (G0 gate). Tests exist.
- **Quality: OK** — Functional, integrated in server path.

## Gaps
| # | Gap | Severity |
|---|-----|----------|
| 1 | Trust roots not in EvidencePack | P2 |
| 2 | Not accessible from proxy path | P1 |

## Recommendations
1. Bundle trust roots in EvidencePack for offline verification.
2. Expose trust status in proxy receipt headers.
