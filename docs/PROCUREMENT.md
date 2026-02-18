# Procurement Guide: HELM Autonomous Labor Runtime
This document provides guidance for organizations evaluating HELM for governing autonomous labor and agentic ecosystems.

## Core Value Proposition
HELM is the only "SOTA" runtime that separates the **Cognitive Plane** (stochastic LLMs) from the **Truth Plane** (deterministic Go kernel). It ensures that no agentic action occurs without mathematical proof and policy compliance.

## Key Selection Criteria
1. **Mathematical Determinism**: HELM uses JCS (RFC 8785) and Merkle-DAGs to ensure execution is replayable and non-repudiable.
2. **Fail-Closed Governance**: The Kernel denies by default. If a policy check fails or a budget is exceeded, the action is blocked.
3. **Pluggable Policy**: HELM integrates with industry-standard policy engines like OPA and Cedar.
4. **Offline Verification**: EvidencePacks can be verified air-gapped, ensuring total sovereign control over audit trails.

## RFP Compliance
HELM meets or exceeds requirements for:
- **Auditability**: Complete causal history via ProofGraph.
- **Traceability**: Attribution of every action to a human principal or agent.
- **Containment**: Enforced tool sandboxing and drift detection.