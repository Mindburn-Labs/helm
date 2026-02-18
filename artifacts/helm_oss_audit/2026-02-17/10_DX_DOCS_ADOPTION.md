# 10 — DX, Docs & Adoption

**Score: 4/5** · Gate ≥3 · **✅ PASS**

---

## One-Command Quickstart

```bash
docker compose up -d && curl -s localhost:8080/v1/chat/completions \
  -H 'Content-Type: application/json' \
  -d '{"model":"gpt-4","messages":[{"role":"user","content":"hello"}]}' | jq .id
```

**Assessment:** ✅ Genuine one-command quickstart. But doesn't exercise the receipt/proof surface — the "5-Minute Proof Loop" in `QUICKSTART.md` does, but is 8 steps.

---

## Documentation Structure (Verified: 23 items in `docs/`)

```
docs/
├── CONFORMANCE.md          ✅ Comprehensive
├── DEMO.md                 ✅ Copy-paste demo
├── DEPENDENCIES.md         ✅ Dependency docs
├── HELM_Unified_Canonical_Standard.md  ✅ Normative spec
├── INTEGRATE_IN_5_MIN.md   ✅ Quick guide
├── INTEGRATIONS/           ✅ Integration guides
├── OSS_CUTLINE.md          ✅ SHIP/QUARANTINE/REMOVE
├── OSS_SCOPE.md            ✅ Scope definition
├── PORTS.md                ✅ Port mapping
├── PUBLISHING.md           ✅ SDK publishing guide
├── QUICKSTART.md           ✅ 5-minute proof loop
├── REGIONAL_COMPAT.md      ✅ Regional compliance
├── ROADMAP.md              ✅ 10 items, no dates
├── SECURITY_MODEL.md       ✅ TCB and crypto chain
├── START_HERE.md           ✅ Entry point
├── TCB_POLICY.md           ✅ TCB boundary rules
├── THREAT_MODEL.md         ✅ 9 adversary classes (T1-T9)
├── UNKNOWNs.md             ✅ Known unknowns
├── api/                    ✅ API documentation
├── sdks/                   ✅ SDK documentation
├── standard/               ✅ Canonical standard ref
├── use-cases/              ✅ Use case guides
└── use_cases/              ⚠️ Duplicate directory
```

**23 doc items** — one of the most comprehensive doc sets in the agent framework space.

### Document Discrepancies (Verified)

| Discrepancy | Location | Issue |
|-------------|----------|-------|
| "11 adversary classes" | README | THREAT_MODEL.md has 9 (T1-T9) |
| EvidencePack determinism | QUICKSTART | CI test is a no-op |
| "58 packages" | README | Actually 112 packages now |
| `use-cases/` vs `use_cases/` | docs/ | Duplicate directories |

---

## Copy-Paste Examples (Verified: 7 directories)

| Language | Directory | Status |
|----------|-----------|--------|
| Python (OpenAI SDK) | `examples/python_openai_baseurl/` | ✅ |
| JavaScript (fetch) | `examples/js_openai_baseurl/` | ✅ |
| TypeScript (Vercel AI) | `examples/ts_vercel_baseurl/` | ✅ |
| Go | `examples/go_client/` | ✅ |
| Java | `examples/java_client/` | ✅ |
| Rust | `examples/rust_client/` | ✅ |
| MCP | `examples/mcp_client/` | ✅ |

**7 examples across all major ecosystems.** Verified via `find`.

### Example Quality

| Criterion | Status |
|-----------|--------|
| Copy-pasteable | ✅ |
| Comments explain HELM behavior | ⚠️ Minimal |
| Error handling with reason_codes | ✅ Go example |
| Receipt verification example | ❌ Missing |
| CI smoke test | ✅ `scripts/ci/examples_smoke.sh` |

---

## SDK Quality Matrix (Verified)

| SDK | Generated Types | Tests | CI Gate | Publish | Coverage |
|-----|-----------------|-------|---------|---------|----------|
| TypeScript | ✅ `types.gen.ts` | ✅ vitest | ✅ | ✅ npm --provenance | Strong |
| Python | ✅ `types_gen.py` | ✅ pytest | ✅ | ✅ PyPI OIDC | Strong |
| Go | ✅ `types_gen.go` | ✅ go test | ✅ | N/A | Good |
| Rust | ✅ `types_gen.rs` | ✅ cargo test | ✅ | ✅ crates.io | Good |
| Java | ✅ `TypesGen.java` | ⚠️ Compile only | ✅ | ✅ Maven | Weak |

**Java gap**: No functional tests, only compile check.

---

## Docs Lint Infrastructure

| Tool | Status | Evidence |
|------|--------|----------|
| `tools/doccheck/` | ✅ CI gate | `helm_core_gates.yml` L258-263 |
| `scripts/ci/doc_hash.sh` | ✅ Canonical doc hash | `helm_core_gates.yml` L219-225 |
| Zero-TODO policy | ✅ CI gate | `helm_core_gates.yml` L155-172 |
| Link checker | ❌ Not found | — |
| Claim-evidence linker | ❌ Not found | — |

---

## Operator Console (`console/`, 4,914 LOC) — **PREVIOUSLY UNAUDITED**

| Feature | File | LOC |
|---------|------|-----|
| Operator API | `operator_api.go` | 928 |
| Mission control | `server_mission_control.go` | 403 |
| Runs API | `runs_api.go` | 293 |
| User onboarding | `onboarding.go` | 196 |
| Portal API | `portal_api.go` | 190 |
| Safe mode | `safemode.go` | 176 |
| UI adapter | `adapter.go` + `ui/adapter.go` | 263 |
| Metrics dashboard | `server_metrics.go` | 840 |
| Why inspector | `why_inspector.go` | 123 |
| Chaos injection | `server.go handleChaosInjectAPI` | — |
| Compliance report | `server.go handleComplianceReportAPI` | — |
| Builder API | `server.go handleBuilderAPI` | — |

**Assessment**: Full operator dashboard with 26 files. Well-tested (`operator_api_test.go` 712 LOC, `runs_api_test.go` 242 LOC, `safemode_test.go` 100 LOC). Not documented publicly — operators won't discover this surface without explicit docs.

---

## Adoption Surface Analysis

### Integration Path (Verified)
```diff
- client = openai.OpenAI()
+ client = openai.OpenAI(base_url="http://localhost:8080/v1")
```
One-line change. No SDK wrapping, no decorator patterns.

### MCP Support
`mcp/` package exists (confirmed in codebase). Enables HELM as MCP-compatible tool execution layer.

### Multi-Tenant
`auth/middleware.go` enforces `TenantID` in JWT claims (L22). Multi-tenant ready at the auth layer.

---

## Score: 4/5

**Justification:**
- ✅ One-command quickstart works
- ✅ 7 copy-paste examples across all ecosystems
- ✅ 5 SDKs with CI gates and publishing
- ✅ 23-item docs structure (comprehensive)
- ✅ Doccheck + doc hash CI gates
- ✅ MCP support
- ✅ Full operator console with onboarding, runs, metrics, chaos injection
- ❌ Document discrepancies (adversary count, package count)
- ❌ No receipt verification example
- ❌ Java SDK has no functional tests
- ❌ No claim-evidence linker
- ❌ Console/operator API undocumented in public docs

### To reach 5/5:
1. Fix document discrepancies
2. Add receipt verification example (Python + TypeScript)
3. Add Java SDK tests
4. Implement claims ledger
5. Deduplicate `use-cases/` vs `use_cases/`
6. Document console/operator API in public docs
