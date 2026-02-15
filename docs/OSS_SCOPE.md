# HELM OSS Scope

## Kernel TCB (Trusted Computing Base)

The following packages are part of the kernel TCB — the minimal trusted core:

| Package | Purpose | Status |
|---------|---------|--------|
| `guardian/` | Policy Resolution Graph (PRG) enforcement | ✅ Active |
| `executor/` | SafeExecutor with receipt generation | ✅ Active |
| `contracts/` | DecisionRecord, Effect, Receipt, Approval | ✅ Active |
| `crypto/` | Ed25519 signing, verification | ✅ Active |
| `canonicalize/` | RFC 8785 JCS implementation | ✅ Active |
| `manifest/` | Tool args/output validation (PEP boundary) | ✅ Active |
| `proofgraph/` | Cryptographic ProofGraph DAG | ✅ Active |
| `trust/registry/` | Event-sourced trust registry | ✅ Active |
| `agent/adapter.go` | KernelBridge choke point | ✅ Active |
| `runtime/sandbox/` | WASI sandbox (wazero, deny-by-default) | ✅ Active |
| `runtime/budget/` | Compute budget enforcement | ✅ Active |
| `escalation/ceremony/` | RFC-005 Approval Ceremony | ✅ Active |
| `evidence/` | Evidence pack export/verify | ✅ Active |
| `replay/` | Replay engine for verification | ✅ Active |
| `mcp/` | Tool catalog + MCP gateway | ✅ Active |
| `kernel/` | Rate limiting, backpressure | ✅ Active |

## Removed from TCB (Enterprise)

The following packages were removed to minimize the attack surface:

| Package | Reason |
|---------|--------|
| `access/` | Enterprise access control |
| `ingestion/` | Brain subsystem data pipeline |
| `verification/refinement/` | Enterprise verification |
| `cockpit/` | UI dashboard |
| `ops/` | Operations tooling |
| `multiregion/` | Multi-region orchestration |
| `hierarchy/` | Enterprise hierarchy |
| `heuristic/` | Heuristic analysis |
| `perimeter/` | Network perimeter |
