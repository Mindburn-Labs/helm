# HELM Architectural Mapping (Baseline Snapshot)
Date: 2026-02-18
Standard: UCS v1.2

## Enforcement Flow
1. **Agent Proposal**: Agent submits an Intent via SDK or CLI.
2. **PEP Boundary**: `Guardian` intercepts the intent.
3. **Policy Evaluation**: `Guardian` delegates to `PDP` (CEL/OPA/Cedar).
4. **Kernel Verdict**: `SignDecision` produces a signed `DecisionRecord`.
5. **Execution**: `SafeExecutor` dispatches to `Connector` only if `PASS`.
6. **ProofGraph**: All nodes (INTENT → ATTESTATION → EFFECT) are hashed via JCS and appended to the DAG.

## Directory Structure
- `core/pkg/kernel`: Truth Plane orchestration.
- `core/pkg/guardian`: PEP implementation.
- `core/pkg/proofgraph`: Merkle-DAG storage.
- `core/pkg/pdp`: Pluggable policy engines.
- `apps/helm-node`: Primary daemon and CLI bridge.
- `sdk/`: Multi-language wire-format implementations.