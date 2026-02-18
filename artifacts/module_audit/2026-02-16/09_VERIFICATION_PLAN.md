# Verification Plan

## Unit Tests

### Proxy Governance Loop
```bash
cd core && go test ./cmd/helm -run TestProxyGovernedToolCall -v
# Expected: tool_call intercepted, Guardian.EvaluateDecision called, receipt created
```

### Proxy Budget Deny
```bash
cd core && go test ./cmd/helm -run TestProxyBudgetDeny -v
# Expected: HTTP 403, reason code BUDGET_EXHAUSTED, receipt in ProofGraph
```

### Proxy Reason Codes
```bash
cd core && go test ./cmd/helm -run TestProxyReasonCodes -v
# Expected: all proxy deny paths produce codes from conform.AllReasonCodes()
```

### ProofGraph Chain Validation
```bash
cd core && go test ./pkg/proofgraph -v
# Expected: all tests pass, chain validation succeeds
```

### Guardian Evaluation
```bash
cd core && go test ./pkg/guardian -v
# Expected: all tests pass including intervention and persistence tests
```

### Executor Gating
```bash
cd core && go test ./pkg/executor -v -count=1
# Expected: SafeExecutor validates gating, creates receipts, checks idempotency
```

### Budget Enforcer
```bash
cd core && go test ./pkg/budget -v
# Expected: fail-closed on errors, daily/monthly limits enforced
```

### Conformance Gates
```bash
cd core && go test ./pkg/conform/... -v
# Expected: all 24 gates registered and tested
```

### EvidencePack Export/Verify
```bash
cd core && go test ./cmd/helm -run TestExportPack -v
cd core && go test ./cmd/helm -run TestVerifyPack -v
# Expected: deterministic tar with sorted paths, hash verification passes
```

---

## Integration Tests

### Full Proxy End-to-End
```bash
# Start proxy with budget limits
bin/helm proxy --upstream https://api.openai.com/v1 --port 9090 --daily-limit 1000 --sign &
PROXY_PID=$!
sleep 2

# Happy path: tool call within budget
curl -X POST http://localhost:9090/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{"model":"gpt-4","messages":[{"role":"user","content":"test"}]}'
# Expected: 200 with governed response, receipt header present

# Budget exceed: send many requests
for i in $(seq 1 100); do
  curl -s -o /dev/null -w "%{http_code}" \
    -X POST http://localhost:9090/v1/chat/completions \
    -H "Content-Type: application/json" \
    -d '{"model":"gpt-4","messages":[{"role":"user","content":"test"}]}'
done
# Expected: eventually 403 with BUDGET_EXHAUSTED

kill $PROXY_PID
```

### Rogue Agent Demo
```bash
bash scripts/demo/rogue_agent.sh
# Expected: exit 0, output shows:
#   - Requests served within budget
#   - Budget exceeded -> 403 BUDGET_EXHAUSTED
#   - Receipt chain verified
#   - EvidencePack exported and verified
```

---

## Repeated-Run Determinism Checks

### Conformance Determinism
```bash
bin/helm conform --profile CORE --json > /tmp/run1.json
bin/helm conform --profile CORE --json > /tmp/run2.json
sha256sum /tmp/run1.json /tmp/run2.json
# Expected: identical SHA-256 hashes
```

### EvidencePack Determinism
```bash
SESSION="test-determinism"
bin/helm export pack $SESSION -o /tmp/pack1.tar.gz
bin/helm export pack $SESSION -o /tmp/pack2.tar.gz
sha256sum /tmp/pack1.tar.gz /tmp/pack2.tar.gz
# Expected: identical SHA-256 hashes (fixed mtime, uid/gid, sorted paths)
```

### ProofGraph Hash Stability
```bash
cd core && go test ./pkg/proofgraph -run TestNodeHashDeterminism -v
# Expected: same payload -> same NodeHash across runs
```

---

## Conformance Runner Outputs

### Level L1
```bash
bin/helm conform --level L1 --json
# Expected output fields:
# {
#   "run_id": "run-L1-...",
#   "profile": "L1",
#   "pass": true/false,
#   "gate_results": [
#     {"gate_id": "G0_BUILD_IDENTITY", "pass": true},
#     {"gate_id": "G1_PROOF_RECEIPTS", "pass": true},
#     {"gate_id": "G2A_SCHEMA_FIRST", "pass": true}
#   ],
#   "git_commit": "abc123...",
#   "env_fingerprint": "linux/amd64/go1.24"
# }
```

### Level L2
```bash
bin/helm conform --level L2 --json
# Expected: all L1 gates + G3A_BUDGET, G8_HITL, G2_REPLAY, GX_TENANT, GX_ENVELOPE
```

---

## EvidencePack Export/Verify/Replay

### Export
```bash
bin/helm export pack session-123 -o /tmp/evidence.tar.gz
# Expected: tar.gz file created
# Verify contents:
tar tzf /tmp/evidence.tar.gz
# Expected: manifest.json, receipt files, sorted alphabetically
```

### Verify
```bash
bin/helm verify --bundle /tmp/evidence-dir
# Expected output: "EvidencePack verification PASSED"
```

### Replay
```bash
bin/helm replay --evidence /path/to/evidence-dir --verify --json
# Expected output:
# {
#   "replay_status": "COMPLETE",
#   "entry_count": N,
#   "issues": []
# }
```

---

## Rogue Agent Deny Reproduction

```bash
# 1. Start proxy with tight budget
bin/helm proxy --upstream mock://localhost:8888 --port 9090 --daily-limit 10 &

# 2. Send request exceeding budget
RESPONSE=$(curl -s -w "\n%{http_code}" \
  -X POST http://localhost:9090/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{"model":"gpt-4","messages":[{"role":"user","content":"do expensive thing"}]}')

HTTP_CODE=$(echo "$RESPONSE" | tail -1)
BODY=$(echo "$RESPONSE" | head -n -1)

# 3. Assert deterministic deny
test "$HTTP_CODE" = "403" || echo "FAIL: expected 403 got $HTTP_CODE"
echo "$BODY" | jq -e '.error.code == "BUDGET_EXHAUSTED"' || echo "FAIL: wrong reason code"
echo "$BODY" | jq -e '.error.receipt_id != null' || echo "FAIL: missing receipt_id"

# 4. Export evidence and verify
bin/helm export pack current-session -o /tmp/rogue-evidence.tar.gz
bin/helm verify --bundle /tmp/rogue-evidence.tar.gz

echo "Rogue agent deny reproduction: COMPLETE"
```

---

## Artifact Locations

| Artifact | Path |
|----------|------|
| Conformance L1 output | `artifacts/conformance/YYYY-MM-DD/run-L1-*/01_SCORE.json` |
| Conformance L2 output | `artifacts/conformance/YYYY-MM-DD/run-L2-*/01_SCORE.json` |
| Use case logs | `artifacts/usecases/UC-001.log` through `UC-012.log` |
| EvidencePack exports | `artifacts/evidence-packs/` |
| Proxy receipts | `artifacts/proxy/receipts.jsonl` (legacy) or ProofGraph JSON |
| CI workflow artifacts | GitHub Actions artifact uploads |
