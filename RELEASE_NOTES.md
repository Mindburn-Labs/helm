# HELM OSS Release Notes
Version: 0.1.0
Date: 2026-02-18

## New Features
- **Hardened PDP Engine**: Native CEL (Common Expression Language) support for context-aware policy enforcement.
- **Pluggable Policy Backends**: Production-grade adapters for **OPA (Open Policy Agent)** and **Cedar**.
- **Verified Genesis Loop (VGL)**: New `synthesize` command for compiling OrgGenome with Deterministic Semantic Mirroring.
- **Deterministic Archiving**: EvidencePack exports (`--tar`) are byte-identical across platforms per UCS Appendix A.3.
- **Auditor-Grade Verification**: Standalone `verify` command with `--json` output for automated compliance reporting.
- **Signed Conformance**: Conformance reports can now be cryptographically signed via `--signed`.

## Architectural Improvements
- Extracted standalone `verifier` library for third-party auditing.
- Registered `GXSDKDrift` CI gate to prevent OpenAPI/SDK desynchronization.
- Fixed namespace collisions in conformance gates.

## Security
- Pervasive use of JCS (RFC 8785) for all cryptographic pre-images.
- Fail-closed semantics reinforced across all external PDP integrations.