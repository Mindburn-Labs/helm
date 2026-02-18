# 11 — Competition Benchmark

## Methodology

Web research conducted 2026-02-17 covering the latest public documentation, release notes, CVE databases, and analyst reports for competing frameworks. All claims are sourced.

---

## Landscape: AI Agent Security & Execution Frameworks (Feb 2026)

### Category 1: Agent Orchestration Frameworks

| Framework | Type | Key Feature | Audit Trail | Deterministic Exec | Crypto Receipts |
|-----------|------|-------------|-------------|-------------------|-----------------|
| **HELM** | Execution Kernel | Fail-closed PEP | ✅ ProofGraph DAG + signed receipts | ✅ JCS + replay engine | ✅ Ed25519 chain |
| LangChain | Orchestrator | Ecosystem breadth | ⚠️ LangSmith tracing (SaaS) | ❌ No | ❌ No |
| CrewAI | Multi-agent | Role-based crews | ⚠️ Audit logs (enterprise) | ⚠️ Task guardrails (validation only) | ❌ No |
| AutoGen (MS) | Multi-agent | Code execution | ⚠️ v0.4 debugging support | ❌ No | ❌ No |
| Semantic Kernel (MS) | SDK | Enterprise integration | ⚠️ OTel tracing (planned Q3 2026) | ❌ No | ❌ No |

### Category 2: Guardrail / Safety Systems

| Framework | Type | Key Feature | Agent-Native | Fail-Closed | Policy Engine |
|-----------|------|-------------|-------------|-------------|---------------|
| **HELM** | Execution Kernel | Crypto-verified execution | ✅ | ✅ | ✅ PRG + Temporal |
| NeMo Guardrails | Input/Output Rails | 5 rail types (Thoughtworks Adopt) | ✅ | ⚠️ Configurable | ✅ Colang |
| Guardrails AI | Output Validation | Pydantic-style validation | ⚠️ Partial | ❌ No | ⚠️ RAIL spec |
| LLM Guard | Scanning | PII redaction, prompt injection | ❌ Scanning only | ❌ No | ❌ No |

### Category 3: Tool Routing / Integration

| Platform | Tools | Key Feature | Auth | MCP | Crypto |
|----------|-------|-------------|------|-----|--------|
| **HELM** | Schema-pinned | Drift detection, receipts | N/A (kernel) | ✅ MCP support | ✅ |
| Composio | 850+ | Managed auth, OAuth | ✅ | ✅ | ❌ |
| Toolhouse | MCP-native | Visual builder, memory | ✅ | ✅ | ❌ |

---

## Head-to-Head Analysis

### HELM vs LangChain

| Dimension | HELM | LangChain |
|-----------|------|-----------|
| **Purpose** | Execution firewall | Orchestration framework |
| **Security Posture** | Cryptographic proofs | SOC 2 Type II (LangSmith SaaS) |
| **Known CVEs** | None public | CVE-2025-68664 (LangGrinch, RCE via serialization injection in `langchain-core`), CVE-2023-46229 (SSRF), CVE-2023-44467 (prompt injection) |
| **Audit Trail** | Ed25519-signed receipt chain | LangSmith runs API (proprietary) |
| **DX** | `base_url` swap | pip install + extensive ecosystem |
| **Adoption** | Pre-launch | 2M+ weekly downloads |
| **Weakness** | Small ecosystem | No deterministic execution, no crypto proofs |

### HELM vs NeMo Guardrails

| Dimension | HELM | NeMo Guardrails |
|-----------|------|-----------------|
| **Scope** | Execution + proof | Input/output filtering |
| **Policy** | PRG + CEL PDP + Perimeter (multi-layer) | Colang (declarative DSL) |
| **Rail Types** | 1 (PEP boundary) | 5 (input, dialog, retrieval, execution, output) |
| **Hot-Reload** | ✅ Yes (`UpdatePolicyBundle` + `LoadPolicy`) | ✅ Yes |
| **Integration** | OpenAI proxy | LangChain, LlamaIndex, NVIDIA NIM |
| **Maturity** | New | Thoughtworks Adopt (Nov 2025) |
| **Unique Value** | Cryptographic proofs, causal ordering | Dialog flow control, RAG validation |

### HELM vs CrewAI

| Dimension | HELM | CrewAI |
|-----------|------|--------|
| **Multi-Agent** | Single-agent kernel | Multi-agent orchestration |
| **Determinism** | JCS + replay engine | Task guardrails (retry until valid) |
| **Governance** | PRG + budget + temporal | Enterprise audit logs |
| **DX** | Go SDK primary | Python-first |
| **Weakness** | No multi-agent support | No crypto proofs |

### HELM vs Microsoft Agent Framework (AutoGen + Semantic Kernel)

| Dimension | HELM | MS Agent Framework |
|-----------|------|-------------------|
| **Stack** | Go kernel + multi-SDK | .NET/Python |
| **Observability** | OTel (362 LOC) | OTel tracing (planned) |
| **Governance** | PRG + Guardian | Agent 365 (enterprise, conditional access) |
| **Policy** | CEL PDP + PRG + Perimeter (multi-layer) | Azure AD-integrated |
| **Unique** | Receipt chain | Enterprise identity (treat agents as service principals) |
| **Risk** | Small team | Vendor lock-in to Azure |

---

## OpenAI API Compatibility — Responses API Impact

### Current State (Feb 2026)
- **HELM implements**: Chat Completions API (OpenAI-compatible proxy)
- **OpenAI deprecations**:
  - GPT-4o API support ended **Feb 16, 2026** (yesterday!)
  - Assistants API sunset: **Aug 26, 2026**
  - Realtime API Beta removal: **Feb 27, 2026**
- **New standard**: Responses API (launched March 2025)
  - Uses `Items` instead of `messages` array
  - Built-in tools: web search, file search, computer use
  - Stateful context management

### Impact on HELM
1. **Immediate**: GPT-4o deprecation means users must update model references
2. **Strategic**: Chat Completions API continues to be supported "indefinitely" per OpenAI
3. **Opportunity**: Adding Responses API support would expand compatibility
4. **Risk**: If OpenAI eventually deprecates Chat Completions, HELM needs migration path

---

## Unique HELM Differentiators (No Equivalent in Market)

| Feature | Competitor Equivalent | HELM Status |
|---------|----------------------|-------------|
| Signed receipt per tool call | ❌ None | ✅ Implemented |
| Causal chain (PrevHash + Lamport) | ❌ None | ✅ Implemented |
| EvidencePack deterministic export | ❌ None | ✅ Implemented |
| Fail-closed PEP boundary | ⚠️ NeMo Guardrails (configurable) | ✅ By default |
| ArgsHash binding (PEP boundary) | ❌ None | ✅ Implemented |
| Output drift detection | ❌ None | ✅ Implemented |
| Authority clock (non-wall-clock) | ❌ None | ⚠️ Guardian only, not Executor |

---

## Competitive Position Summary

```
                    ┌─────────────────────────┐
            Crypto  │                    HELM │
           Proofs   │                    ●    │
                    │                         │
                    │                         │
                    │   NeMo ●                │
            Safety  │                         │
           Rails    │ Guardrails AI ●         │
                    │   LLM Guard ●           │
                    │                         │
                    │ CrewAI ●    SK/AutoGen ● │
         Orchestr.  │            LangChain ●  │
                    │                         │
                    └─────────────────────────┘
                    Simple ←── DX ──→ Rich
```

**HELM occupies a unique niche**: no other framework provides cryptographically-verified, causally-ordered execution receipts with deterministic replay. The trade-off is DX maturity and ecosystem breadth.

**Biggest competitive threat**: NeMo Guardrails adding execution rails + crypto proofs (currently has execution rails but no crypto). If NVIDIA adds receipt signing to NeMo, HELM's differentiation narrows significantly.

**Biggest opportunity**: Agent governance is predicted to become its own product category in 2026 (per Vectra AI). HELM is positioned to define the standard.
