# Architecture Decisions

## ADR-001: Proxy Has Two Implementations - Must Converge

**Status**: Gap identified

**Context**: The repo has two OpenAI proxy codepaths:
1. `core/pkg/api/openai_proxy.go` (93 lines) - HTTP handler that returns static "governance proxy active" response. Registered in the server mux. Does not forward to upstream.
2. `core/cmd/helm/proxy_cmd.go` (388 lines) - Standalone CLI command with `httputil.ReverseProxy` upstream forwarding, tool_call interception, JSONL receipt store with LamportClock and PrevHash, Ed25519 signing.

**Decision**: Converge on `proxy_cmd.go` as the canonical proxy path. The handler in `openai_proxy.go` should either be removed or replaced with a thin wrapper that delegates to the same logic. The CLI proxy already has the right architecture (upstream forwarding, tool interception, receipt chaining).

**Consequence**: The handler used when starting via `helm serve` (runServer path) needs the same capabilities as `helm proxy`.

---

## ADR-002: ProofGraph Is In-Memory Only

**Status**: Architectural limitation

**Context**: `proofgraph.Graph` is an in-memory DAG. It has mutex-protected operations, Lamport clock, and chain validation. But it has no persistence layer. The `proofgraph/store.go` (2840B) defines a store interface, but it is not used by the proxy.

**Decision**: Accept in-memory for v0.1 OSS. Add filesystem-backed store for single-node deployments. The receipt chain in `proxy_cmd.go` already persists to JSONL, which is a pragmatic alternative, but must be migrated to use ProofGraph.

---

## ADR-003: Conformance Profiles vs. Levels

**Context**: The wedge spec asks for `--level L1|L2`. The implementation uses `--profile SMB|CORE|ENTERPRISE|...`.

**Decision**: Add `--level L1|L2` as an alias that maps to a subset of gates. L1 = deterministic bytes + ProofGraph signing + EvidencePack export/verify. L2 = L1 + bounded compute + approvals + offline replay + ACID concurrency.

---

## ADR-004: Regional Profiles - Single YAML vs. Separate Files

**Context**: Wedge requires `profile_{us,eu,ru,cn}.yaml` at `core/pkg/config/profiles/`. Current repo has `regional.yaml` with all four embedded.

**Decision**: Split into separate files for clarity and git-blame isolation. Each file should be a complete, self-contained profile with all required fields (outbound_networking, upstream_allowlist, island_mode, crypto_policy, retention_defaults).

---

## ADR-005: Budget Enforcement Uses SimpleEnforcer - Adequate for OSS

**Context**: `budget.SimpleEnforcer` is fail-closed, has daily/monthly limits, and creates enforcement receipts. It uses a `Storage` interface with memory and Postgres implementations.

**Decision**: SimpleEnforcer is adequate for OSS v0.1. Must be wired into the proxy execution loop. The fail-closed behavior matches the wedge requirement.

---

## ADR-006: Stable Reason Code Taxonomy Exists but Not Used Everywhere

**Context**: `conform/reason_codes.go` defines 23 stable reason codes. But the proxy (`proxy_cmd.go`) and budget enforcer (`budget/enforcer.go`) use ad-hoc string reasons.

**Decision**: All components must use and reference the canonical reason codes from `conform/reason_codes.go`. Add HELM_DENY_* prefix for proxy-specific codes. Extend the reason code taxonomy to cover all deny scenarios in the proxy path.
