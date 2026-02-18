# Standard RFP Answers: HELM
This document contains standard responses for Risk, Compliance, and Technical RFP sections.

## 1. Security & Compliance
**Q: How does the system handle unauthorized agent actions?**
A: HELM uses a Policy Enforcement Point (PEP). Every side-effectful action must be authorized by a signed Kernel Verdict. Without this verdict, the Executor (SafeExecutor) cannot dispatch the action.

**Q: Can logs be tampered with?**
A: HELM produces a ProofGraph (Merkle-DAG). Any attempt to alter history invalidates the subsequent node hashes and signatures, making tampering mathematically detectable during verification.

## 2. Technical Architecture
**Q: What is the "7-Plane Model"?**
A: It is the HELM architectural standard that separates Identity, Policy, Truth, Record, Tools, Knowledge, and Surface into explicit trust boundaries.

**Q: Does it support OPA?**
A: Yes, HELM includes a native OPA adapter (`HELM_POLICY_BACKEND=opa`) that delegates decisions to a central OPA server while maintaining local proof binding.

## 3. Reliability
**Q: What happens if the policy engine is down?**
A: HELM implements fail-closed semantics. If the PDP is unreachable, all actions are denied until the connection is restored.