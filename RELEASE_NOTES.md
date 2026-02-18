# HELM v0.9.0-rc.1 — Competitive OSS Release

## Release Highlights

This release transforms HELM from a single-engine governance system into a **pluggable, auditable, enterprise-ready** platform.

### 🔌 Pluggable Policy Engines (P0.1)

HELM now supports three policy backends:

| Backend | Status | Setup |
|---------|--------|-------|
| **HELM** (built-in) | Production | Default — no config needed |
| **OPA** (Open Policy Agent) | Production | `HELM_POLICY_BACKEND=opa` |
| **Cedar** (AWS) | Beta | `HELM_POLICY_BACKEND=cedar` + sidecar |

All backends enforce **fail-closed semantics** and produce **deterministic decision hashes** (SHA-256 of JCS-canonical decisions) that are bound into every receipt.

**New files**: `core/pkg/pdp/` — stable PDP interface with 9/9 conformance tests.

### 🔍 Standalone Offline Verifier (P0.2)

A verifier library with **zero external dependencies** — auditors can build and run it independently:

```bash
./bin/helm verify --bundle /path/to/evidencepack --json-out audit.json
```

7 verification checks: structure, index integrity, file hashes, chain integrity, Lamport monotonicity, policy decision hashes, replay determinism.

**New files**: `core/pkg/verifier/`, `docs/VERIFIER_TRUST_MODEL.md`

### 📊 Signed Conformance Reports (P0.3)

```bash
./bin/helm conform --level L2 --signed
# Outputs: conform_report.json + .sha256 + .sig
```

Enhanced conformance gates:
- **G2**: Lamport clock monotonicity + policy decision hash cross-reference
- **GX_SDK_DRIFT**: OpenAPI spec / SDK directory drift detection

**New docs**: `PROCUREMENT.md`, `RFP_ANSWERS.md` — enterprise buyer enablement.

### 🛡️ Dispute Replay Viewer (P1.1)

Offline HTML/JS viewer for verification reports — drop `verify.json` and see the decision path, check results, and metadata. No server required.

**New files**: `tools/dispute-viewer/`

### 🔗 Orchestration Framework Examples (P1.2)

Integration examples showing how to route agentic pipelines through HELM:
- `examples/langgraph/` — LangGraph agent with HELM proxy
- `examples/openai_agents/` — OpenAI Agents SDK with tool use

### 🧬 OrgDNA Seed (P3)

Organizational policy genome — content-addressed, composable policy packs:

```bash
./bin/helm orgdna validate --pack examples/orgdna/saas_support.json
./bin/helm orgdna hash --pack examples/orgdna/finance_approval.json
```

**New files**: `schemas/orgdna.schema.json`, `docs/ORGDNA_OVERVIEW.md`

---

## Breaking Changes

None. All changes are additive.

## New `DecisionRecord` Fields

```diff
+PolicyBackend      string `json:"policy_backend,omitempty"`
+PolicyContentHash  string `json:"policy_content_hash,omitempty"`
+PolicyDecisionHash string `json:"policy_decision_hash,omitempty"`
```

## New CLI Commands

| Command | Description |
|---------|-------------|
| `helm verify --json-out FILE` | Auditor mode: structured verification report |
| `helm conform --signed` | Signed conformance report artifacts |
| `helm orgdna validate --pack FILE` | Validate OrgDNA pack |
| `helm orgdna hash --pack FILE` | Content-address an OrgDNA pack |

## New Environment Variables

| Variable | Description |
|----------|-------------|
| `HELM_POLICY_BACKEND` | Policy engine: `helm`, `opa`, `cedar` |
| `HELM_OPA_URL` | OPA server URL |
| `HELM_CEDAR_URL` | Cedar sidecar URL |

## Test Results

- **All 60+ packages pass** (zero regressions)
- **9/9 PDP conformance tests** (HELM, OPA, Cedar)
- **5/5 verifier golden tests**
- Build: ✅ | Lint: ✅

## Commits

| SHA | Description |
|-----|-------------|
| `c43846d` | feat(pdp): OPA/Cedar policy adapters |
| `325c6ff` | feat(verifier): standalone offline verifier |
| `6a0706d` | feat(conform): signed conformance reports |
| `4b7af25` | feat(conform): Lamport + policy hash + SDK drift gates |
| `8f210e2` | feat(p1): dispute viewer + orchestration shims |
| `e42b7d7` | feat(orgdna): OrgDNA seed |
