# OSS Wedge Requirements Checklist

## Step 1: OpenAI-Compatible Proxy + Execution Loop

### Endpoint
- [ ] `POST /v1/chat/completions` with full OpenAI wire compatibility
- [ ] Feature flag: `HELM_ENABLE_OPENAI_PROXY=1` (OFF by default)
- [x] OpenAI request/response types defined (`OpenAIChatRequest`, `OpenAIChatResponse`)

### Execution Loop
- [x] Upstream forwarding via `httputil.ReverseProxy` (in `proxy_cmd.go`)
- [ ] Tool call interception routed through `KernelBridge -> Guardian -> Executor`
- [ ] Pinned schema validation for tool args (via `manifest.ValidateToolArgs`)
- [ ] Pinned schema validation for tool outputs (via `executor.OutputSchemaRegistry`)
- [ ] LLM re-invocation loop (tool_call -> execute -> re-invoke until completion)
- [ ] Strict limits: max tool iterations, max wallclock, bounded compute budgets
- [x] Receipt creation with LamportClock, PrevHash (in `proxy_cmd.go`)
- [ ] Receipt stored in ProofGraph (currently JSONL)

### Deterministic Deny
- [x] 23 stable reason codes enumerated (`conform/reason_codes.go`)
- [ ] Proxy uses canonical reason codes (currently ad-hoc strings)
- [ ] Stable error payload structure (JSON with `code`, `reason`, `receipt_id`)
- [ ] HTTP 403 with HELM_DENY reason code for violations

### Developer Integration
- [x] `OPENAI_BASE_URL` is the only config change needed
- [x] `helm proxy --upstream <url> --port <port>` CLI exists
- [ ] Works with LangChain/LangGraph/Vercel AI SDK (untested)

---

## Step 2: Demo That Converts Skeptics

### Scripts and Use Cases
- [x] `scripts/usecases/run_all.sh` exists (89 lines)
- [x] UC-001 through UC-012 documented
- [ ] UC docs contain reproducible input/output assertions (currently pointer files)
- [x] Artifacts output to `artifacts/usecases/`

### Control-Room UI
- [ ] P0 dashboard
- [ ] Inbox approvals (RFC-005 ceremony)
- [ ] ProofGraph timeline (LamportClock, PrevHash, verdicts, reason codes)

### Rogue Agent Demo
- [ ] Dedicated script: sample client exceeds budget
- [ ] HELM denies with 403 + HELM_DENY reason code
- [ ] UI shows deny receipt and chained ProofGraph nodes
- [ ] EvidencePack export matches offline replay verification

---

## Step 3: Conformance Runner

### CLI
- [x] `helm conform` CLI command exists
- [x] `--profile` flag with 6 profiles (SMB, CORE, ENTERPRISE, ...)
- [ ] `--level L1|L2` alias for wedge compatibility
- [x] `--json` output flag
- [ ] Deterministic output bytes (same input -> same output)
- [ ] Git commit + environment fingerprint in output

### Gates
- [x] 24 gates registered (G0-G12, Gx_envelope, Gx_tenant)
- [x] Gates run in deterministic order
- [x] EvidencePack directory created per run

### CI Integration
- [x] CI workflow with 15 jobs
- [ ] Dedicated `helm conform` CI job
- [ ] Fail-on-regression gate

---

## Step 4: Region Profiles

### Profile Files
- [ ] `profile_us.yaml` (separate file)
- [ ] `profile_eu.yaml` (separate file)
- [ ] `profile_ru.yaml` (separate file)
- [ ] `profile_cn.yaml` (separate file)
- [x] All four regions defined (in single `regional.yaml`)

### Profile Controls
- [ ] Outbound networking defaults
- [ ] Allowlists for LLM upstreams
- [ ] Island mode posture
- [ ] Crypto policy allowlists
- [ ] Retention defaults
- [x] Ceremony settings (timelock, hold, challenge, domain separation)
- [x] Data residency
- [x] Compliance framework references
- [x] Encryption algorithm

### Enforcement
- [ ] Config loader selects profile at startup
- [ ] Profile constraints enforced at runtime (not just documented)

---

## Step 5: EvidencePack + Offline Replay

### Export
- [x] `ExportPack` function with deterministic tar.gz
- [x] Sorted paths, fixed mtime(0), stable uid/gid(0)
- [ ] CLI `helm export pack <session_id>` wired to session state
- [ ] Pack includes ProofGraph nodes
- [ ] Pack includes trust roots

### Verify
- [x] `VerifyPack` function validates hash integrity
- [x] `helm verify --bundle` CLI command
- [x] `conform.ValidateEvidencePackStructure` called

### Replay
- [x] `helm replay --evidence --verify` CLI command
- [x] Tape-based verification (sequence, data_class, manifest integrity)
- [ ] Offline effect re-execution with hash comparison
- [ ] Disconnected machine verification (no network required)

---

## OSS Always Includes

- [x] Guardian (PEP): 347 lines, PRG evaluation, budget tracking, intervention
- [x] Executor: 341 lines, SafeExecutor with gating, receipts, schema validation
- [x] ProofGraph: DAG with Lamport clock, JCS hashing, chain validation
- [x] Trust Registry: TUF, SLSA, Rekor, pack loader, signature verifier
- [x] Manifest + Schema validation: JCS canonicalization, tool arg validation
- [x] Canonicalization: JCS (RFC 8785) implementation
- [x] Crypto: Ed25519 signing, verification, audit log
- [x] RFC-004 bounded compute: WASI sandbox with gas/time/memory limits
- [x] RFC-005 approval ceremonies: Timelock, challenge, domain separation
- [x] EvidencePack export + verify
- [x] Conformance runner with 24 gates + profiles
- [x] OpenAI-compatible proxy (CLI version, feature-flagged)
- [x] MCP gateway (stub)
- [ ] Control-room UI
