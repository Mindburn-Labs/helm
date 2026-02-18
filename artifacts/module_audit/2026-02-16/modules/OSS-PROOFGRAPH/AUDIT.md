# OSS-PROOFGRAPH Audit

## Scope
Cryptographic DAG recording all decisions, effects, and denials. Hash-chained nodes with Lamport clock. JCS canonicalization. Chain validation.

## Reality
- **graph.go** (130 lines): In-memory DAG with Append, AppendSigned, ValidateChain, LamportClock. [CODE]
- **node.go** (123 lines): 6 node types, JCS-style hash (SHA-256), Validate. [CODE]
- **store.go** (2840B): Store interface defined but not wired. [CODE]
- **Quality: OK** — DAG logic correct, no persistence, no proxy integration.

## Conformance
| Check | Status |
|-------|--------|
| DAG with parent linking + Lamport clock | ✅ |
| JCS-style hash computation (SHA-256) | ✅ |
| Chain validation (recursive parent walk) | ✅ |
| Signed node support, 6 node types | ✅ |
| Filesystem persistence | ❌ |
| Integrated into proxy receipts | ❌ |
| Authority clock (not time.Now()) | ❌ |

## Gaps
| # | Gap | Severity |
|---|-----|----------|
| 1 | Not integrated into proxy receipt chain | P0 |
| 2 | No filesystem persistence | P1 |
| 3 | Uses time.Now() not authority clock | P1 |

## Recommendations
1. **[NEW]** `fs_store.go` — filesystem-backed store
2. **[MODIFY]** `proxy_cmd.go` — replace JSONL receiptStore with ProofGraph
3. **[MODIFY]** `node.go:112` — inject clock
4. **[NEW]** `export.go` — JSON serialization for EvidencePack
