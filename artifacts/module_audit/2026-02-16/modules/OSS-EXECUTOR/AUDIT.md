# OSS-EXECUTOR Audit

## Scope
Run effects if valid decision + intent exist. Signed receipts. Output schema validation. Idempotency. Fail-closed.

## Reality
- **executor.go** (341 lines): SafeExecutor with gating, idempotency, signed receipts with causal links. [CODE]
- OutputSchemaRegistry interface. Evidence Pack + Merkle tree proofs.
- `validateGating`: Requires signed decision + intent + matching IDs. [CODE]
- `createReceipt`: Signed receipt with causal links. [CODE]
- Tests: executor_test.go, evidence_pack_test.go, merkle_test.go.
- **Quality: SOTA** — Fail-closed, idempotent, cryptographic receipts.

## Conformance
| Check | Status |
|-------|--------|
| Gating (valid decision + intent) | ✅ |
| Idempotency via receipt store | ✅ |
| Signed receipts with causal links | ✅ |
| Output schema registry interface | ✅ |
| Merkle tree evidence proofs | ✅ |
| Called from proxy | ❌ |
| Called from MCP gateway | ❌ |

## Gaps
| # | Gap | Severity |
|---|-----|----------|
| 1 | Not invoked from proxy tool_call handling | P0 |
| 2 | Not invoked from MCP gateway | P1 |
| 3 | OutputSchemaRegistry not populated in proxy context | P1 |

## Recommendations
Executor code is complete. Gaps are in callers only:
1. Wire via KernelBridge: tool_call -> Guardian -> Executor.Execute
2. Populate OutputSchemaRegistry from manifest tool definitions.
