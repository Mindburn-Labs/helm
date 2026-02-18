# Gateboard - OSS Wedge Go/No-Go

## Step 1: OpenAI-Compatible Proxy + Execution Loop

| Gate | Status | Evidence |
|------|--------|----------|
| POST /v1/chat/completions exists | **PARTIAL** | [CODE] `core/pkg/api/openai_proxy.go:50` handler exists but returns static response |
| Feature flag (OFF by default) | **YES** | [CODE] `HELM_ENABLE_OPENAI_PROXY` env var in `core/cmd/helm/main.go` |
| Upstream forwarding | **PARTIAL** | [CODE] `proxy_cmd.go:140-250` uses httputil.ReverseProxy for upstream, but handler in `openai_proxy.go` does not |
| Tool call interception | **PARTIAL** | [CODE] `proxy_cmd.go:200-280` intercepts tool_calls from response, validates args with JCS |
| KernelBridge -> Guardian -> Executor routing | **NO** | Proxy tool interception does not invoke Guardian.EvaluateDecision or SafeExecutor.Execute |
| Pinned schema validation (tool args) | **PARTIAL** | [CODE] `proxy_cmd.go:86-104` validates JSON + JCS hash, but does not use `manifest.ToolSchema` pinned schemas |
| Pinned schema validation (tool outputs) | **NO** | No output schema validation in proxy loop |
| Receipt signing + ProofGraph | **PARTIAL** | [CODE] `proxy_cmd.go:27-42` proxyReceipt with LamportClock + PrevHash, Ed25519 signing, but uses flat JSONL not ProofGraph DAG |
| Deterministic deny with stable reason codes | **PARTIAL** | [CODE] `conform/reason_codes.go` has 23 stable reason codes, but proxy uses ad-hoc status strings |
| Strict iteration/wallclock/budget limits | **PARTIAL** | [CODE] `proxy_cmd.go` has `--max-tokens` flag; `budget/enforcer.go` has fail-closed budget check, but not wired into proxy |
| LLM re-invocation loop | **NO** | Proxy does single passthrough, no tool_call -> result -> re-invoke loop |

**Verdict: NO-GO** - Core execution loop wiring missing.

---

## Step 2: Demo Readiness

| Gate | Status | Evidence |
|------|--------|----------|
| scripts/usecases/run_all.sh | **YES** | [CODE] `scripts/usecases/run_all.sh` (89 lines, runs UC-001..UC-012) |
| UC-001..UC-012 exist | **YES** | [CODE] `docs/use-cases/UC-001.md` through `UC-012.md` (pointer docs, ~120B each) |
| Artifact output structure | **YES** | [CODE] `artifacts/usecases/` created by run_all.sh |
| Deterministic re-runs | **PARTIAL** | Use cases run `go test` commands; test determinism depends on individual test design |
| Control-room UI | **NO** | No UI directory exists, no React/web UI code in repo |
| Rogue agent demo | **NO** | No dedicated rogue agent script; budget deny path exists in tests but no turnkey demo |

**Verdict: PARTIAL** - Script infrastructure exists, but no UI and no rogue agent reproduction.

---

## Step 3: Conformance Runner

| Gate | Status | Evidence |
|------|--------|----------|
| helm conformance run --level L1/L2 | **YES** | [CODE] `core/cmd/helm/conform.go` implements `helm conform --profile` |
| Deterministic output JSON | **PARTIAL** | [CODE] `conform/engine.go:143-152` writes JSON, but uses `json.MarshalIndent` (not canonicalized), no environment fingerprint |
| CI gates that run conformance | **PARTIAL** | [CODE] `.github/workflows/helm_core_gates.yml:123-136` runs use-cases, but does not run `helm conform` directly |
| Profile-to-gate mapping | **YES** | [CODE] `conform/profile.go` maps profiles to gate sets |
| 24 gates registered | **YES** | [CODE] `conform/gates/registry.go` + G0-G12, Gx gates |
| Versioning + git commit in output | **NO** | ConformanceReport struct lacks git commit, env fingerprint fields |

**Verdict: PARTIAL** - Engine works, but output is not byte-deterministic and CI does not run `helm conform`.

---

## Step 4: Region Profiles

| Gate | Status | Evidence |
|------|--------|----------|
| profile_{us,eu,ru,cn}.yaml | **PARTIAL** | [CODE] `core/pkg/config/profiles/regional.yaml` has all 4 profiles in a single file, not separate per-region files |
| Config loader profile selection | **NO** | [CODE] `core/pkg/config/config.go` (945B) is minimal, no profile loading by region |
| Outbound networking defaults | **NO** | Not in profile yaml |
| LLM upstream allowlists | **NO** | Not in profile yaml |
| Island mode posture | **NO** | Concept mentioned in docs but not in profile config |
| Crypto policy allowlists | **PARTIAL** | Profile yaml has `encryption` field but no allowlist enforcement |
| Retention defaults | **NO** | Not in profile yaml |

**Verdict: NO-GO** - Profile structure is a document, not enforced configuration.

---

## Step 5: EvidencePack + Offline Replay

| Gate | Status | Evidence |
|------|--------|----------|
| helm export pack <session_id> | **PARTIAL** | [CODE] `core/cmd/helm/export_pack.go:25-79` ExportPack function with deterministic tar, but CLI wiring is incomplete (function not exposed as subcommand with session_id lookup) |
| Deterministic tarball bytes | **YES** | [CODE] `export_pack.go:81-97` fixed mtime(0), uid/gid(0), sorted paths |
| helm verify pack | **PARTIAL** | [CODE] `export_pack.go:99-160` VerifyPack reads tar + validates hashes; CLI in `verify_cmd.go` delegates to `conform.ValidateEvidencePackStructure` |
| helm replay verify | **YES** | [CODE] `replay_cmd.go` implements `helm replay --evidence --verify` with tape integrity checks |
| Offline replay without DB | **PARTIAL** | Replay loads from filesystem tapes, no DB needed, but does not re-execute effects |
| EvidencePack self-contained | **PARTIAL** | Export includes file hashes in manifest, but no ProofGraph nodes or trust roots bundled |

**Verdict: PARTIAL** - Primitives exist but not wired end-to-end.
