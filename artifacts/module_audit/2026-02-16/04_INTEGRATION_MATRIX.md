# Integration Matrix

This matrix shows which kernel components are correctly wired together and which connections are missing.

## Component Wiring Status

```
Legend:
  [Y] = Wired and tested
  [P] = Partially wired
  [N] = Not wired
  [-] = N/A
```

| Source | Guardian | Executor | ProofGraph | Budget | Receipts | Trust | Schemas | Conform | Replay | MCP | Proxy |
|--------|----------|----------|------------|--------|----------|-------|---------|---------|--------|-----|-------|
| **Proxy CLI** | [N] | [N] | [N] | [N] | [P] JSONL | [-] | [P] JCS only | [-] | [-] | [-] | - |
| **Proxy Handler** | [N] | [N] | [N] | [N] | [N] | [-] | [N] | [-] | [-] | [-] | - |
| **Server (runServer)** | [Y] | [Y] | [P] | [Y] | [Y] | [Y] | [Y] | [-] | [-] | [P] | [P] |
| **Guardian** | - | [Y] Intent | [N] | [Y] Tracker | [Y] | [Y] PRG | [Y] Envelope | [-] | [-] | [-] | [-] |
| **Executor** | [Y] Gating | - | [N] | [-] | [Y] | [Y] Artifact | [Y] Output | [-] | [-] | [-] | [-] |
| **Conformance** | [-] | [-] | [P] G1 gate | [P] G3a gate | [Y] G1 gate | [Y] G0 gate | [Y] G2a gate | - | [P] G2 gate | [-] | [-] |
| **MCP Gateway** | [N] | [N] | [N] | [N] | [N] | [-] | [-] | [-] | [-] | - | [-] |
| **EvidencePack** | [-] | [Y] | [-] | [-] | [Y] Merkle | [-] | [-] | [Y] | [-] | [-] | [-] |

## Critical Missing Wiring

1. **Proxy -> Guardian**: The proxy CLI does not call `guardian.EvaluateDecision()` for tool calls. This is the #1 wedge gap.
2. **Proxy -> Executor**: The proxy does not call `executor.SafeExecutor.Execute()`. Tool calls are not governed.
3. **Proxy -> ProofGraph**: The proxy uses flat JSONL receipts, not `proofgraph.Graph.Append()`.
4. **Proxy -> Budget**: The proxy does not call `budget.SimpleEnforcer.Check()`.
5. **MCP -> Guardian/Executor**: The MCP gateway is completely disconnected from the kernel pipeline.

## Correctly Wired Paths

The `runServer` path in `main.go:132-299` correctly wires:
- Guardian with PRG, crypto signer, artifact registry, budget tracker
- SafeExecutor with crypto verifier, tool driver, receipt store, artifact store
- Trust registry with TUF, SLSA, Rekor integrations
- Conformance engine with all 24 registered gates

**Key insight**: The server path has correct wiring. The proxy and MCP paths are disconnected. The fix is to reuse the server-path wiring in the proxy path.
