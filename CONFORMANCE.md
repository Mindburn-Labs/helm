# HELM Conformance Standards

**H**euristic **E**xecution **L**ayer **M**odel (HELM) OSS v0.1 targets **Level 2 (L2) Conformance**.

## Conformance Levels

| Level | Name | Description | Status |
|---|---|---|---|
| **L1** | **Structure** | Schema validity, interface existence, compilation, basic tool execution | ✅ **PASS** |
| **L2** | **Behavior** | Fail-closed security, TCB isolation, determinism, idempotency, basic budgets | ✅ **PASS** |
| **L3** | **Hardened** | Advanced side-channel resistance, formal verification, hardware enclave support | 🚧 *Planned v0.2* |

## Verification Suite

To run the full conformance suite locally:

```bash
# Requires Docker and Go 1.24+
bash scripts/usecases/run_all.sh
```

## Mandatory Use Cases (L2)

Files located in `docs/use_cases/`:

1.  **[UC-001] PEP Allow** — Valid signed intent with known schema must execute.
2.  **[UC-002] PEP Fail-Closed** — Unknown fields, extra params, or schema violations must be rejected (`DENY_SCHEMA_MISMATCH`).
3.  **[UC-003] Approval Ceremony** — High-stakes actions must block for multi-party approval/timelock (`DENY_APPROVAL_REQUIRED`).
4.  **[UC-004] WASM Transform** — Pure logic must execute in WASI sandbox with no network/fs.
5.  **[UC-005] WASM Exhaustion** — Infinite loops or memory leaks must be terminated deterministically (`DENY_GAS_EXHAUSTION`).
6.  **[UC-006] Idempotency** — Re-submitting the same signed intent (nonce) must return the original receipt without re-execution.
7.  **[UC-007] Export** — EvidencePack export must be a deterministic, content-addressed archive.
8.  **[UC-008] Replay** — EvidencePack must be replayable offline to reproduce the exact Ledger state.
9.  **[UC-009] Contract Drift** — Connector outputs violating schema must fail execution (`ERR_CONNECTOR_CONTRACT_DRIFT`).
10. **[UC-010] Trust Rotation** — Key rotation events must be verifiable via the ProofGraph.
11. **[UC-011] Island Mode** — Kernel must function (with reduced availability) during network partition.
12. **[UC-012] OpenAI Proxy** — (Optional) Proxy mode must correctly wrap standard LLM API calls.
