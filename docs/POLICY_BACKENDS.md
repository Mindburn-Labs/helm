# Policy Backend Configuration

HELM supports pluggable policy backends. The enforcement kernel (Guardian) delegates policy evaluation to a configurable Policy Decision Point (PDP) while retaining signing, proof binding, and receipt chain integrity.

## Supported Backends

| Backend | Description | Config Value |
|---------|-------------|-------------|
| **HELM** (default) | Built-in PRG evaluation | `helm` |
| **OPA** | Open Policy Agent (Rego) | `opa` |
| **Cedar** | AWS Cedar (via sidecar) | `cedar` |

## Configuration

### Environment Variable

```bash
export HELM_POLICY_BACKEND=opa   # or "cedar" or "helm" (default)
```

### OPA Backend

```bash
# 1. Start OPA with your policies
docker run -d -p 8181:8181 openpolicyagent/opa:latest run --server

# 2. Load a HELM authorization policy
cat <<'EOF' > /tmp/helm_authz.rego
package helm.authz

default allow = false

allow {
    input.action == "read"
}

reason_code = "ALLOW" { allow }
reason_code = "DENY_POLICY" { not allow }
EOF

curl -X PUT http://localhost:8181/v1/policies/helm \
  --data-binary @/tmp/helm_authz.rego

# 3. Configure HELM
export HELM_POLICY_BACKEND=opa
export HELM_OPA_URL=http://localhost:8181
export HELM_OPA_POLICY_PATH=/v1/data/helm/authz  # default
export HELM_OPA_TIMEOUT=5s                        # default
export HELM_OPA_POLICY_VERSION=v1.0.0             # for receipt binding
```

### Cedar Backend

Cedar requires a sidecar PDP (no production Go evaluator exists).

```bash
# 1. Start the Cedar PDP sidecar
cd tools/cedar-pdp
docker compose up -d
# Sidecar listens on http://localhost:8182

# 2. Configure HELM
export HELM_POLICY_BACKEND=cedar
export HELM_CEDAR_URL=http://localhost:8182
export HELM_CEDAR_DECIDE_PATH=/decide            # default
export HELM_CEDAR_TIMEOUT=5s                      # default
export HELM_CEDAR_POLICY_VERSION=v1.0.0           # for receipt binding
```

## How It Works

```
Request → Guardian.EvaluateDecision()
              │
              ├── PDP configured? ──→ PDP.Evaluate() ──→ allow/deny
              │                            │
              │                    Bind into DecisionRecord:
              │                    - policy_backend: "opa"
              │                    - policy_content_hash: sha256:...
              │                    - policy_decision_hash: sha256:...
              │
              ├── PRG + Temporal checks (always run)
              │
              └── Sign DecisionRecord → Receipt → ProofGraph
```

## Receipt Binding

Every receipt includes:

| Field | Description |
|-------|-------------|
| `policy_backend` | Which PDP made the decision |
| `policy_content_hash` | Content-addressed hash of the active policy set |
| `policy_decision_hash` | SHA-256 of the JCS-canonical decision (deterministic) |

An auditor can verify: *"this deny happened because policy X in backend Y produced decision hash Z"*.

## Fail-Closed Semantics

All backends enforce fail-closed:

- **Timeout** → DENY (`DENY_OPA_UNREACHABLE` / `DENY_CEDAR_UNREACHABLE`)
- **HTTP error** → DENY (`DENY_OPA_HTTP_500` etc.)
- **Parse error** → DENY (`DENY_OPA_PARSE_ERROR` etc.)
- **Nil request** → DENY (`DENY_NIL_REQUEST`)

## Writing Custom Backends

Implement the `PolicyDecisionPoint` interface:

```go
type PolicyDecisionPoint interface {
    Evaluate(ctx context.Context, req *DecisionRequest) (*DecisionResponse, error)
    Backend() Backend
    PolicyHash() string
}
```

See `core/pkg/pdp/` for reference implementations.
