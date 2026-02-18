# SOTA Baseline - Feb 2026

## What "Ahead" Means

Being ahead of governance wrappers (Guardrails AI, Salus, WitnessAI, Cisco AI Defense) means:

1. **Proof survives the database**: Receipts are hash-chained in a DAG, not rows in a vendor database. An auditor can verify without trusting HELM's infrastructure.
2. **Deterministic authority**: Every deny produces the same reason code, the same receipt structure, the same ProofGraph node. There is no "it depends on the deployment."
3. **Offline verification**: An EvidencePack can be validated on a disconnected machine. No API calls, no license checks, no cloud dependencies.
4. **Schema-first governance**: Tool calls are validated against pinned schemas before execution, not just logged after the fact.

## Competitive Comparison

| Capability | HELM OSS (Today) | HELM OSS (Target) | Guardrails AI | WitnessAI | Cisco AI Defense |
|------------|-------------------|--------------------|---------------|-----------|-----------------|
| **Proof chain** | JSONL receipts with LamportClock + PrevHash [CODE] | ProofGraph DAG with signed nodes | No proof chain; policy logs | Session logs | Alert logs |
| **Deterministic deny** | 23 stable reason codes [CODE] | Reason codes in proxy + MCP | Non-deterministic guardrail failures | N/A | Alert severity levels |
| **Offline verify** | EvidencePack export with hash validation [CODE] | Full offline replay | No offline mode | No offline mode | No offline mode |
| **Schema validation** | Manifest tool schemas with JCS canonicalization [CODE] | Pinned schema enforcement in proxy | Input/output rails | No schema enforcement | ML-based detection |
| **Bounded compute** | WASI sandbox with gas/time/memory limits [CODE] | Budget enforcement in proxy | No bounded compute | No bounded compute | No bounded compute |
| **Approval ceremonies** | Timelock + challenge + domain separation [CODE] | Integrated into proxy/UI | No approval workflow | No approval workflow | No approval workflow |
| **Cryptographic signing** | Ed25519 receipt signing [CODE] | ProofGraph node signing | No crypto signing | No crypto signing | No crypto signing |
| **Region profiles** | US/EU/RU/CN config [CODE] | Enforced region constraints | No region support | US-only | No region support |
| **1-line integration** | `OPENAI_BASE_URL` change [CODE] | Full agentic loop | SDK-specific wrappers | Network interposition | Network appliance |

## Where HELM Is Already Ahead

1. **Cryptographic proof chain**: No competitor has hash-chained receipts with Lamport clocks. This is structural differentiation.
2. **WASI sandboxing**: No competitor offers bounded compute for tool execution. HELM has working gas/time/memory limits.
3. **Deterministic reason codes**: No competitor commits to stable, enumerated deny reasons. This is auditor-grade.
4. **Offline verification**: No competitor offers disconnected EvidencePack verification.
5. **Schema-first**: No competitor validates tool_call arguments against pinned schemas before execution.

## Where HELM Must Close Gaps to Credibly Ship

1. **End-to-end proxy loop**: Today's proxy intercepts but does not govern. Competitors at least log everything.
2. **Visual proof**: No UI means no demo conversion. Competitors have dashboards.
3. **npm-install simplicity**: `helm proxy --upstream` works, but the full agentic loop (tool_call -> govern -> execute -> re-invoke) is missing.

## Moat Assessment

```
Structural moat (hard to replicate):  ProofGraph, WASI sandbox, JCS canonicalization
Operational moat (easy to replicate): Dashboards, alerting, logging
```

HELM's structural moat is strong. The gap is purely operational: wire the structural primitives into the user-facing proxy path.
