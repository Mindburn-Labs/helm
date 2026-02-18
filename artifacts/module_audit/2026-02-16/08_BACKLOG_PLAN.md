# Backlog Plan

## Phase 0: P0 Wedge Blockers (Ship Blockers)

### P0-1: Wire Proxy Execution Loop Through Kernel

**[MODIFY] `core/cmd/helm/proxy_cmd.go`**

Current state: `runProxyCmd` uses `httputil.ReverseProxy` for upstream forwarding and intercepts tool_calls in the response, but does not route them through Guardian -> Executor.

Implementation:
1. Import `guardian`, `executor`, `proofgraph`, `budget`, `contracts` packages
2. In `runProxyCmd`, initialize kernel components (or accept them as dependencies):
   - `guardian.NewGuardian(signer, prg, registry)`
   - `executor.NewSafeExecutor(verifier, signer, driver, store, ...)`
   - `proofgraph.NewGraph()`
   - `budget.NewSimpleEnforcer(memoryStore)`
3. When tool_calls are intercepted from LLM response:
   ```go
   // For each tool_call in response:
   req := guardian.DecisionRequest{
       Principal: tenantID,
       Action:    toolCall.Function.Name,
       Resource:  toolCall.Function.Arguments,
   }
   decision, err := g.EvaluateDecision(ctx, req)
   if decision.Verdict != "PASS" {
       // deterministic deny: use conform.ReasonSchemaValidationFailed etc.
       writeProxyDeny(w, decision, receipt)
       return
   }
   intent, _ := g.IssueExecutionIntent(ctx, decision, effect)
   receipt, artifact, _ := exec.Execute(ctx, effect, decision, intent)
   ```
4. After execution, append receipt to ProofGraph:
   ```go
   payload, _ := proofgraph.EncodePayload(receipt)
   node, _ := pg.AppendSigned(proofgraph.NodeTypeEffect, payload, sig, principal, seq)
   ```

**Definition of Done**: Tool calls in proxied responses go through Guardian evaluation and SafeExecutor execution. Deny decisions produce HTTP 403 with stable reason code.
**Verification**: `go test ./cmd/helm -run TestProxyGovernedToolCall -v`

---

### P0-2: Implement Agentic Re-Invocation Loop

**[MODIFY] `core/cmd/helm/proxy_cmd.go`**

Add iteration loop:
```go
const maxIterations = 10
for iteration := 0; iteration < maxIterations; iteration++ {
    resp := forwardToUpstream(req)
    if !hasToolCalls(resp) {
        return resp // final response
    }
    toolResults := []ToolResult{}
    for _, tc := range resp.ToolCalls {
        result := governAndExecute(tc) // Guardian -> Executor
        toolResults = append(toolResults, result)
    }
    req = buildFollowupRequest(req, resp, toolResults) // re-invoke LLM
}
return deterministicDeny("MAX_ITERATIONS_EXCEEDED")
```

**Definition of Done**: Multi-turn agent conversations complete through proxy with governance at each tool call.
**Verification**: Integration test with mock upstream returning 3 sequential tool_calls.

---

### P0-3: Integrate ProofGraph Into Proxy

**[MODIFY] `core/cmd/helm/proxy_cmd.go`**

Replace `receiptStore` (JSONL) with `proofgraph.Graph` + persistence:
1. Initialize `pg := proofgraph.NewGraph()` at proxy startup
2. For each receipt, call `pg.AppendSigned(...)` instead of `store.Append(...)`
3. Add filesystem serialization: write ProofGraph JSON after each append
4. In receipt response headers, include `X-Helm-Receipt-NodeHash` and `X-Helm-Receipt-Lamport`

**Definition of Done**: All proxy receipts are ProofGraph nodes. Chain validation passes. `pg.ValidateChain(headID)` returns nil.
**Verification**: `go test -run TestProxyProofGraphChain`

---

### P0-4: Map Proxy Outcomes to Stable Reason Codes

**[MODIFY] `core/cmd/helm/proxy_cmd.go`**
**[MODIFY] `core/pkg/conform/reason_codes.go`**

1. Add proxy-specific reason codes:
   ```go
   ReasonProxyBudgetExceeded      = "PROXY_BUDGET_EXCEEDED"
   ReasonProxySchemaViolation     = "PROXY_SCHEMA_VIOLATION"  
   ReasonProxyMaxIterations       = "PROXY_MAX_ITERATIONS"
   ReasonProxyUpstreamError       = "PROXY_UPSTREAM_ERROR"
   ReasonProxyToolCallDenied      = "PROXY_TOOL_CALL_DENIED"
   ```
2. Use these codes in all proxy deny responses
3. Return JSON error body: `{"error": {"code": "PROXY_BUDGET_EXCEEDED", "receipt_id": "...", "message": "..."}}`

**Definition of Done**: Every proxy deny uses a code from `conform.AllReasonCodes()`. Error payload structure is stable.
**Verification**: `go test -run TestProxyReasonCodes`

---

### P0-5: Wire Budget Enforcement Into Proxy

**[MODIFY] `core/cmd/helm/proxy_cmd.go`**

1. Initialize `budgetEnforcer := budget.NewSimpleEnforcer(budget.NewMemoryStore())`
2. Before each tool execution:
   ```go
   decision, _ := budgetEnforcer.Check(ctx, tenantID, budget.Cost{Amount: estimatedCost})
   if !decision.Allowed {
       return deterministicDeny(conform.ReasonBudgetExhausted, decision.Receipt)
   }
   ```
3. Add `--daily-limit` and `--monthly-limit` flags to proxy CLI

**Definition of Done**: Budget violations produce 403 with `BUDGET_EXHAUSTED` reason code. Receipts are chained.
**Verification**: `go test -run TestProxyBudgetDeny`

---

### P0-6: Wire EvidencePack Export to CLI

**[MODIFY] `core/cmd/helm/export_cmd.go`**

1. Implement `helm export pack <session_id>` subcommand
2. Look up session receipts from ProofGraph or receipt JSONL
3. Call `ExportPack(sessionID, files, outputPath)`
4. Print pack hash to stdout

**Definition of Done**: `helm export pack session123 -o pack.tar.gz` produces deterministic tarball.
**Verification**: `helm export pack test-session -o /tmp/pack.tar.gz && helm verify --bundle /tmp/pack.tar.gz`

---

## Phase 1: Demo Conversion Surface

### P1-1: Rogue Agent Demo Script

**[NEW] `scripts/demo/rogue_agent.sh`**

Script that:
1. Starts `helm proxy` with `--daily-limit 100`
2. Sends curl requests that exercise budget
3. Sends one more request that exceeds budget
4. Asserts HTTP 403 response with `BUDGET_EXHAUSTED` reason code
5. Extracts receipt and ProofGraph node hash from response
6. Runs `helm export pack` and `helm verify`

**Definition of Done**: Script runs end-to-end and produces passing output.
**Verification**: `bash scripts/demo/rogue_agent.sh`

---

### P1-2: Minimal Control-Room UI

**[NEW] `apps/control-room/`**

Minimal Next.js or Vite app with three views:
1. **P0 Dashboard**: Active sessions, recent receipts, deny count, budget utilization
2. **Approvals Inbox**: Pending approval requests (from RFC-005 ceremonies)
3. **ProofGraph Timeline**: Scrollable list of ProofGraph nodes with LamportClock, NodeHash, Kind, Verdict

API endpoints needed (in existing server or new):
- `GET /api/receipts?limit=N` - recent receipts
- `GET /api/proofgraph/nodes?limit=N` - recent ProofGraph nodes
- `GET /api/approvals/pending` - pending approval requests
- `POST /api/approvals/:id/approve` - approve a request

**Definition of Done**: UI starts, shows ProofGraph nodes from a running proxy session.
**Verification**: Start proxy, send requests, open UI, see receipts in timeline.

---

### P1-3: Wire MCP Gateway to Kernel

**[MODIFY] `core/pkg/mcp/gateway.go`**

Replace stub `handleExecute` with:
1. Look up tool in catalog
2. Build `guardian.DecisionRequest` from MCP request
3. Call `guardian.EvaluateDecision`
4. If allowed, call `executor.Execute`
5. Return tool result or deny

**Definition of Done**: MCP tool calls are governed by Guardian/Executor.
**Verification**: `go test ./pkg/mcp -run TestGatewayGoverned`

---

### P1-4: Add `--level L1|L2` to Conformance CLI

**[MODIFY] `core/cmd/helm/conform.go`**

Add `--level` flag that maps:
- L1 -> gates: G0, G1, G2a (deterministic bytes, ProofGraph signing, EvidencePack)
- L2 -> gates: all L1 + G3a (budget), G8 (HITL), G2 (replay), Gx_tenant, Gx_envelope

**Definition of Done**: `helm conform --level L1` runs correct gate subset.
**Verification**: `helm conform --level L1 --json | jq .gate_results[].gate_id`

---

## Phase 2: Conformance and CI Hardening

### P2-1: Deterministic Conformance Output

**[MODIFY] `core/pkg/conform/engine.go`**

1. Use JCS canonicalization for all JSON output
2. Add fields to `ConformanceReport`:
   ```go
   GitCommit    string `json:"git_commit"`
   EnvFingerprint string `json:"env_fingerprint"` // OS + arch + go version
   ```
3. Replace `json.MarshalIndent` with canonical JSON writer

**Definition of Done**: Two runs with same input produce byte-identical SHA-256.
**Verification**: `helm conform --profile CORE --json > /tmp/a.json && helm conform --profile CORE --json > /tmp/b.json && sha256sum /tmp/a.json /tmp/b.json`

---

### P2-2: Conformance CI Gate

**[MODIFY] `.github/workflows/helm_core_gates.yml`**

Add job:
```yaml
conformance:
  runs-on: ubuntu-latest
  needs: build
  steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-go@v5
      with:
        go-version: '1.24'
    - name: Build
      run: cd core && go build -o ../bin/helm ./cmd/helm/
    - name: Conformance L1
      run: bin/helm conform --level L1 --json > artifacts/conformance/l1.json
    - name: Conformance L2
      run: bin/helm conform --level L2 --json > artifacts/conformance/l2.json
    - name: Verify determinism
      run: |
        bin/helm conform --level L1 --json > /tmp/l1_check.json
        diff artifacts/conformance/l1.json /tmp/l1_check.json
```

**Definition of Done**: CI runs conformance and fails on regression or non-determinism.

---

### P2-3: ProofGraph Persistence

**[NEW] `core/pkg/proofgraph/fs_store.go`**

Implement filesystem-backed ProofGraph store:
```go
type FSStore struct {
    dir string
}
func (s *FSStore) Save(g *Graph) error
func (s *FSStore) Load() (*Graph, error)
```

**Definition of Done**: ProofGraph survives proxy restart. Load-save-load produces same graph.
**Verification**: `go test -run TestFSStoreRoundtrip`

---

## Phase 3: Region Profiles + MCP

### P3-1: Split Regional Profiles

**[NEW] `core/pkg/config/profiles/profile_us.yaml`**
**[NEW] `core/pkg/config/profiles/profile_eu.yaml`**
**[NEW] `core/pkg/config/profiles/profile_ru.yaml`**
**[NEW] `core/pkg/config/profiles/profile_cn.yaml`**
**[DELETE] `core/pkg/config/profiles/regional.yaml`**

Each profile file includes:
```yaml
name: "United States"
outbound_networking:
  default_policy: allow
  blocked_domains: []
upstream_allowlist:
  - "api.openai.com"
  - "api.anthropic.com"
island_mode: false
crypto_policy:
  allowed_algorithms: ["Ed25519", "AES-256-GCM", "SHA-256"]
retention:
  default_days: 365
  max_days: 2555
ceremony:
  min_timelock_ms: 2000
  min_hold_ms: 1000
data_residency: "us-east-1"
compliance: ["SOC2", "NIST-800-53"]
```

---

### P3-2: Profile Config Loader

**[MODIFY] `core/pkg/config/config.go`**

Implement:
```go
func LoadProfile(region string) (*RegionProfile, error)
func (p *RegionProfile) EnforceNetworking(req *http.Request) error
func (p *RegionProfile) EnforceCrypto(algorithm string) error
```

---

## Phase 4: Polish and Publish

### P4-1: Expand Use Case Docs

**[MODIFY] `docs/use-cases/UC-001.md` through `UC-012.md`**

Each UC should contain:
- Input: exact curl command or SDK code
- Expected output: exact JSON response
- Assertion: deterministic hash or field check
- Artifact: path to generated evidence

---

### P4-2: Reason Codes Reference Doc

**[NEW] `docs/REASON_CODES.md`**

Generated from `conform.AllReasonCodes()` with descriptions for each code.

---

### P4-3: Conformance Levels Reference Doc

**[NEW] `docs/CONFORMANCE_LEVELS.md`**

Documents L1 and L2 gate mappings and what each proves.

---

### P4-4: EvidencePack Verify Tutorial

**[NEW] `docs/EVIDENCEPACK_VERIFY.md`**

Step-by-step guide for auditors to verify an EvidencePack offline.

---

### P4-5: Proxy Threat Model

**[MODIFY] `docs/THREAT_MODEL.md`**

Add section covering proxy-specific attack vectors:
- Prompt injection via tool_call arguments
- Budget exhaustion attacks
- Receipt chain tampering
- Upstream impersonation
