# OSS-EVIDENCEPACK Audit

## Scope
Deterministic tar.gz export. Manifest with file hashes. Verify integrity offline. Self-contained for auditor review.

## Reality
- **export_pack.go** (161 lines): Deterministic tar (sorted paths, fixed mtime(0), uid/gid(0)), manifest with SHA-256 hashes. [CODE]
- **verify_cmd.go** (89 lines): `helm verify --bundle`. Validates structure + signature. [CODE]
- **executor/evidence_pack.go** (8067B): EvidencePack struct with Merkle tree. [CODE]
- **executor/merkle.go** (10653B): Merkle tree for evidence proofs. [CODE]
- **Quality: Partial** — Primitives exist, not wired as CLI with session state.

## Gaps
| # | Gap | Severity |
|---|-----|----------|
| 1 | CLI not wired as `helm export pack <session_id>` | P0 |
| 2 | ProofGraph not included in pack | P1 |
| 3 | Trust roots not included in pack | P2 |

## Recommendations
1. Wire CLI: `helm export pack <session_id>` -> lookup receipts -> call ExportPack.
2. Include ProofGraph JSON. 3. Bundle trust anchors.
