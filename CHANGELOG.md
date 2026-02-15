# Changelog

All notable changes to HELM Core OSS are documented here.

## [0.1.0] — 2026-02-15

### Added

- **Proxy sidecar** (`helm proxy`) — OpenAI-compatible reverse proxy. One line changed, every tool call gets a receipt.
- **SafeExecutor** — single execution boundary with schema validation, hash binding, and signed receipts.
- **Guardian** — policy engine with configurable tool allowlists and deny-by-default.
- **ProofGraph DAG** — signed nodes (INTENT, ATTESTATION, EFFECT, TRUST_EVENT, CHECKPOINT) with Lamport clocks and causal `PrevHash` chains.
- **Trust Registry** — event-sourced key lifecycle (add/revoke/rotate), replayable at any height.
- **WASI Sandbox** — deny-by-default (no FS, no net) with gas/time/memory budgets and deterministic trap codes.
- **Approval Ceremonies** — timelock + deliberate confirmation + challenge/response, suitable for disputes.
- **EvidencePack Export** — deterministic `.tar.gz` with sorted paths, epoch mtime, root uid/gid.
- **Replay Verify** — offline session replay with full signature and schema re-validation.
- **CLI** — 11 commands: `proxy`, `export`, `verify`, `replay`, `conform`, `doctor`, `init`, `trust add/revoke`, `version`, `serve`.
- **SDK Stubs** — TypeScript and Python client libraries.
- **Regional Profiles** — US, EU, RU, CN with Island Mode for network partitions.
- **12 executable use cases** with scripted validation.
- **Conformance gates** — L1 (kernel invariants) and L2 (profile-specific).

### Security

- Fail-closed execution: undeclared tools are blocked, schema drift is a hard error.
- Ed25519 signatures on all decisions, intents, and receipts.
- ArgsHash (PEP boundary) cryptographically bound into signed receipt chain.
- 8-package TCB with forbidden-import linter.
