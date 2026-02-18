# Pluggable Policy Backends (PDP)
HELM supports multiple Policy Decision Point (PDP) backends to integrate with existing enterprise policy infrastructure.

## Configuration
Use the following environment variables to configure the active PDP:

- `HELM_POLICY_BACKEND`: `helm` (default CEL), `opa`, or `cedar`.
- `HELM_POLICY_VERSION`: A human-readable identifier for the policy set (e.g., git hash).

### OPA (Open Policy Agent)
- `OPA_URL`: Base URL of the OPA server (e.g., `http://localhost:8181`).
- HELM dispatches requests to `/v1/data/helm/authz`.

### Cedar
- `CEDAR_URL`: Base URL of the Cedar sidecar (e.g., `http://localhost:8182`).
- HELM dispatches requests to `/decide`.

## Fail-Closed Semantics
All external PDP adapters implement strict fail-closed behavior. If the PDP is unreachable, returns a non-200 status, or provides a malformed response, the Kernel will DENY the action.

## Proof Binding
Every decision produced by a PDP is hashed via JCS and bound into the `DecisionRecord` and `ProofGraph`.