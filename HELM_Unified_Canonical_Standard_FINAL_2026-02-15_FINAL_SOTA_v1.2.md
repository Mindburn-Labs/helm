# HELM Unified Canonical Standard (Unified Core)
Version: 1.2
Status: Canonical
Date: 2026-02-15
Audience: CTO, CISO, Compliance Officers, Systems Engineering Leadership, Platform Engineers

Canonical-Doc-Hash: sha256:cd2129389a68353722f76f0920858045c82c0fa5e791efeef7220f880140f6bd
Hash-Scope: SHA-256 over canonical Markdown UTF-8 bytes of this document, with the 64-hex digest value above treated as all-zero during hashing (Appendix A.5). This Markdown is the canonical artifact. Any PDF is a non-normative rendering.

## 0. Normative posture

### 0.1 Non-goals
This document does not:
- Specify individual model architectures or weights.
- Define proprietary application-layer business logic for any specific company.
- Provide legal advice for specific regulatory frameworks.
- Mandate any single vendor cloud, hardware, or model provider.

This document does:
- Define a deterministic autonomy runtime architecture.
- Define normative wire contracts and hashing rules for verifiable execution.
- Define conformance levels and an executable harness profile.

### 0.2 Normative language
The key words MUST, MUST NOT, REQUIRED, SHALL, SHALL NOT, SHOULD, SHOULD NOT, RECOMMENDED, MAY, and OPTIONAL are to be interpreted as described in RFC 2119 and RFC 8174.

### 0.3 Canonical terminology
To prevent drift, the following canonical terms are normative and supersede prior nomenclature:

- **OrgGenome**: the compiled, versioned specification of an organization (policies, budgets, roles, tool access, risk classes, connector bindings).
- **OrgPhenotype**: the deterministic runtime state compiled from the OrgGenome and ProofGraph checkpoints.
- **ProofGraph**: an immutable causal DAG of intents, verifications, receipts, and effects.
- **Policy Enforcement Point (PEP)**: the deterministic boundary that verifies cryptographic intent and policy compliance before side effects.
- **Verified Planning Loop (VPL)**: the protocol that turns stochastic proposals into deterministic execution.
- **Constraint and Proof Interface (CPI)**: the deterministic validator for plans and policy constraints.
- **EvidencePack**: a deterministic, signed evidence bundle proving what happened and why.
- **HELM A2A Envelope**: the minimal deterministic wire format enabling agent-to-agent interoperability.
- **Verified Genesis Loop (VGL)**: the protocol for compiling OrgGenome from human intent and context with deterministic semantic reflection and explicit approval.
- **Proof Condensation**: cryptographic compaction of low-risk activity into checkpoints to preserve unit economics without losing auditability.

### 0.4 Single Source of Truth rule
If a fact is not derivable from (a) OrgGenome, (b) ProofGraph, or (c) signed EvidencePacks referenced by ProofGraph, it is NOT canonical. UIs and external APIs MUST NOT invent parallel truth.

### 0.5 Policy namespace and precedence (terminology hardening)
This standard uses these policy levels:
Note: policy levels P0/P1/P2 are distinct from the Plane 1-7 architecture model.
- **P0 Platform Ceilings**: platform-level immutable ceilings and deny lists set via explicit UI controls. P0 is evaluated first and cannot be bypassed by any other policy.
- **P1 OrgGenome Policies**: the compiled MAPL Policy Bundle that defines capability sets, constraints, budgets, and risk classes.
- **P2 Transient Overlays**: time-bounded, signed overlays (for example emergency freezes, incident throttles, or temporary exceptions). P2 MAY further restrict behavior and MUST NOT expand authority beyond P0 and P1 unless an explicit break-glass policy exists requiring stronger signatures and time locks.

Precedence:
1) P0 evaluated first (hard ceilings, hard denies).
2) P1 determines effective permissions via intersection semantics (Section 7).
3) P2 is applied last as an overlay (restrict-by-default). Any break-glass expansion, if supported, MUST be explicit, time-bounded, and auditable with dedicated reason codes.



### 0.6 Autonomy trilemma and the Universal Deployment Framework (HUDF)

HELM acknowledges a fundamental constraint of autonomous systems engineering:

- Agentic fluidity (dynamic reasoning loops, improvisation, rapid adaptation)
- Cryptographic determinism (bit-stable receipts, replayable proofs, non-repudiation)
- Low operational overhead (minimal latency, minimal storage, minimal human intervention)

These three cannot all be maximized at the exact moment of runtime execution. Any attempt to make the deterministic runtime "partially fuzzy" (for example: runtime auto-healing of connector payloads, unlogged execution paths, or side-effect sandboxes that bypass CPI) destroys HELM's mathematical and legal guarantees and fractures ecosystem trust.

Therefore:

- The Kernel MUST NOT be forked by market segment.
- The ProofGraph, Policy Enforcement Point (PEP), and Constraint and Proof Interface (CPI) MUST remain uncompromising across all deployments.
- Flexibility MUST be shifted left into asynchronous compilation, UI, policy scoping, and risk delegation.

The HELM Universal Deployment Framework (HUDF) is the canonical mechanism for supporting both SMB and enterprise without forking the math:

- Fork capability scope via P1 policy (what an agent is allowed to do).
- Fork toolchain contracts via Plane 5 (how side effects are executed and contained).
- Fork economics via Plane 4 and Section 8 (what is retained, condensed, and for how long).
- Never fork enforcement semantics (PEP/CPI) or proof semantics (ProofGraph).


## 1. Executive vision

### 1.1 B2C and SMB
HELM is a business compiler and operating system. A user describes an outcome, grants access (accounts, docs, inbox, calendar, payments, repositories), and HELM compiles an OrgGenome that can run a real business as governed autonomous labor. The user delegates outcomes, not steps. HELM executes end-to-end workflows while producing proof-grade receipts and exportable evidence.

Users do not author LTL or policy code. HELM compiles formal constraints from context, then forces deterministic semantic reflection and explicit approval before those constraints become law (Section 4).


Additional HUDF requirements for B2C and SMB:

- B2C and SMB deployments MUST use HUDF profiles rather than kernel forks. The enforcement loop (PEP/CPI/ProofGraph) remains identical to enterprise deployments.
- SMB deployments MAY instantiate OrgGenome from a vendor-signed template (Franchised genome). In this model:
  - the template vendor signs the P1 policy bundle and impact report
  - the SMB principal signs only P0 ceilings (sovereign blast radius)
  - the kernel binds both signatures into a single ORG_GENESIS_APPROVAL attestation



### 1.2 B2B and enterprise
HELM is an enterprise autonomy runtime that turns autonomous work into a verifiable system of record. Every decision and action is attributable, replayable, and policy-bound by default. HELM enables multi-vendor agent ecosystems through a shared trust envelope and cryptographic accounting, while remaining merge-auditable under replication and partitions.

### 1.3 The core thesis
Stochastic models are useful for proposing plans. They are not acceptable as authorities for execution. HELM isolates stochasticity to proposal and forces execution through deterministic governance, deterministic replay, and cryptographic receipts.

## 2. Architecture model

### 2.1 Seven-plane model
HELM is defined by explicit trust boundaries:

- **Plane 1 Identity and Trust**: principals, keys, trust roots, optional hardware attestations.
- **Plane 2 OrgGenome**: compiled org spec (policies, budgets, risk classes, tool permissions).
- **Plane 3 Deterministic Kernel (Truth Plane)**: arbitration, CPI validation, PEP enforcement, deterministic scheduling.
- **Plane 4 ProofGraph**: immutable causal record and replay substrate.
- **Plane 5 Tools and Connectors**: sandboxed side-effect surface and connector contracts.
- **Plane 6 Knowledge Plane**: learned vs curated knowledge with promotion protocol.
- **Plane 7 Surfaces and Ecosystem**: UI and external APIs that render from canonical truth and issue signed intents.

### 2.2 Deterministic boundary: the PEP
All side-effectful actions MUST cross the PEP. The PEP MUST:
- Verify the caller principal and authorization.
- Verify the action is permitted by OrgGenome policies (MAPL) and P0 ceilings (if present).
- Verify CPI verdict for the proposed plan slice and action.
- Enforce tool sandbox requirements and connector preconditions.
- Emit receipts into ProofGraph for every allow or deny decision.

Fail-closed invariant: if verification cannot be completed, execution MUST be denied with a deterministic reason code.



### 2.3 Knowledge Plane contract (LKS vs CKS)
HELM separates memory into two classes:

- Learned Knowledge Store (LKS): untrusted, probabilistic, and mutable. It may contain summaries, embeddings, heuristics, and model-produced notes. LKS MAY influence plan proposals but MUST NOT directly authorize side effects.
- Curated Knowledge Store (CKS): trusted, signed, versioned facts and policies. CKS entries MUST have provenance and MUST be referenced by hash from ProofGraph when used to justify actions.

Promotion protocol:
- Any LKS-derived claim that is required for a side effect MUST be promoted to CKS by attaching provenance (source, snapshot, signature) and emitting an OBSERVATION node referencing the promoted artifact.
- Policies MAY require dual-source corroboration before promotion for specific data classes.

This split prevents "memory drift" from becoming execution authority.



### 2.4 Universal Deployment Framework (HUDF)

HUDF defines how a single canonical Kernel supports different market realities without non-conformant bypasses.

HUDF is implemented as deployment profiles. Profiles MUST be expressed as policy and toolchain configuration, never as changes to ProofGraph semantics or CPI verdict logic.

A profile MUST define:

- max_risk_class (T0-T3) allowed for the tenant
- CPI ladder enabled (Tier 1, Tier 2, Tier 2.5 HITL, Tier 3 formal)
- toolchain constraints (for example: dynamic code = WASI-only, no raw network I/O)
- evidence retention policy (default TTLs, condensation cadence, dispute replay minimums)
- approval routing policy (who can sign, by what mechanism, with what time windows)

Two reference profiles are canonical:

- SMB profile: max_risk_class <= T2, Tier 1 and Tier 2, Tier 2.5 approvals for escalations, aggressive condensation, and franchised genomes.
- Enterprise profile: Tier 3 enabled for selected intents, longer evidence retention, stronger sandbox attestations, and CI/CD-based connector patching.

Implementations MAY define additional profiles, but MUST state profile parameters in machine-readable policy and MUST remain conformant to PEP/CPI semantics.

#
### 2.5 Reality Anchor Extensions (RFC-002)

HUDF keeps the Kernel uncompromising by forking capability scope and toolchains, not cryptography. RFC-002 hardens HUDF for a 2030-scale autonomous economy by closing five systemic blind spots:

1. Privacy-Verification Paradox - Cross-org verification MUST NOT force disclosure of raw ProofGraph or EvidencePack contents to counterparties.
2. Silicon Void - A deterministic kernel running on an untrusted host can be a "fake judge" that signs perfectly formatted lies.
3. Algorithmic Denial of Wallet - Sender can externalize cost by forcing receiver to execute NP-hard checks (SMT zip-bombs).
4. Cognitive Drift - Wargaming is invalidated if the cognitive engine changes under the same OrgGenome and policies.
5. Cryptographic Mortality - Zero-trust systems need deterministic succession and freeze semantics when the sovereign principal disappears.

Normative resolution:

- ZK-CPI (Section 6.4 and Section 11.2): When a counterparty demands verification, the sender MUST produce a zero-knowledge proof of compliance and private claims. The receiver MUST verify in bounded time. Compute burden is inverted.
- Hardware Attested Kernel (Section 3.3.2 and Section 11.3): KERNEL_VERDICT attestations consumed outside a tenant trust boundary MUST carry a silicon quote proving the measured kernel binary and execution profile.
- Cognitive Engine Pinning (Section 4.3): OrgGenome approvals MUST bind to a pinned cognitive engine attestation. Silent engine swaps MUST force deprecation mode and restrict to low-risk reads until re-wargamed and re-approved.
- Succession and Dead Man's Clock (Section 4.5): Every OrgGenome MUST define heartbeat and recovery quorum semantics. Heartbeat expiry MUST deterministically freeze external effects to prevent "ghost ship" autonomy.


#
### 2.6 Entropy and bounded compute controls (RFC-004)

HELM is deterministic, not fragile. Determinism does not require unbounded computation. Any untrusted computation surface MUST be bounded in a deterministic way to prevent asymmetric algorithmic denial-of-service (DoS) and systemic synchronization shocks.

#### 2.6.1 Deterministic compute budgets (normative)

Any of the following MUST execute under deterministic compute budgets, recorded in receipts:

- Tier 3 solvers (SMT, LTL, or equivalent)
- Zero-knowledge verification (Profile-gated)
- Dynamic code execution (Section 10.1.1)
- Any verification routine whose input can be influenced by an external counterparty

Minimum budget fields (recorded):
- `gas_limit_steps` and `gas_used_steps` (deterministic instruction or tick budget)
- `time_limit_ms` and `time_used_ms` (deterministic wallclock budget, used only for termination, never for ordering)

Fail-closed:
- If a budget is exhausted, CPI/PEP MUST NOT fail open. It MUST return DENY or DEFER with reason_code = ERR_COMPUTE_GAS_EXHAUSTED or ERR_COMPUTE_TIME_EXHAUSTED.

#### 2.6.2 Velocity ceilings and dispersion windows (systemic-risk optional)

For intents that can create systemic correlated actions (for example: markets, global logistics, mass procurement), policies MAY require velocity ceilings and deterministic dispersion windows.

Recommended P0 ceiling additions:
- MAX_SPEND_USD_PER_HOUR
- MAX_SIDE_EFFECT_ACTIONS_PER_MIN
- MAX_NEW_COUNTERPARTIES_PER_DAY
- MAX_EFFECT_BURST_WINDOW_MS (controls time-local clustering)

Dispersion windows:
- If enabled by policy, the kernel MAY apply a deterministic delay in the range [0, MAX_EFFECT_BURST_WINDOW_MS] derived from:
  - tenant root public key, and
  - intent_hash, and
  - effect_index
- The derived delay MUST be recorded in the ToolReceipt as `dispatch_delay_ms`.

This creates heterogeneity across tenants even when they share a franchised genome, without introducing nondeterminism.

#### 2.6.3 Analog proxy risk controls (normative minimum)

Any connector operation that can materially change the physical world via third-party human labor or non-attested off-platform execution MUST declare executor_kind = ANALOG (Section 10.2).

Default behavior:
- executor_kind = ANALOG MUST require explicit approval (Tier 2.5) unless a policy allowlist explicitly permits the destination and constrains task scope. If not approved, the system MUST DENY with reason_code = ERR_ANALOG_EXECUTOR_REQUIRES_APPROVAL.

This prevents "clean digital payments" from becoming an ungoverned analog actuation channel by default.

### 2.7 Human governance ceremony and anti-atrophy controls (RFC-005)

Cryptographic approvals are not automatically meaningful governance. At scale, humans habituate. If Tier 2.5 approvals can be rubber-stamped, non-repudiation collapses and the system becomes "autonomy theater" again.

HELM therefore defines a minimal, deterministic approval ceremony protocol that preserves legal and operational validity without weakening fail-closed semantics.

#### 2.7.1 Approval ceremony v1 (normative minimum when enabled)

For any approval of an intent where:
- risk_class >= T2, and
- effect_class = IRREVERSIBLE, or policy marks the action as HIGH_EXTERNALITY,

the runtime MUST require an approval ceremony, unless an emergency override policy explicitly permits bypass.

Minimum ceremony properties:

- Binding: the ceremony MUST bind to the exact intent_hash and include a ui_summary_hash over the human-visible summary that was displayed.
- Timelock: the runtime MUST enforce a not-before delay between request issuance and acceptance. The delay is policy-configured; default 30 seconds. The enforced value MUST be recorded as approval_timelock_ms.
- Deliberate action: the approver MUST perform at least one non-reflexive confirmation step in addition to biometric or passkey presence. Permitted minimum mechanisms:
  - hold-to-approve for >= min_hold_ms (policy-configured), or
  - re-enter amount and destination last-4 (or equivalent high-salience fields), or
  - select a reason code from a bounded set.
- Transcript hashing: the UI MUST hash the challenge and response into challenge_hash and response_hash. Raw challenge content MAY be retained locally as Sensitive evidence (Section 3.4) if required by policy.

The approval MUST be recorded as an ATTESTATION of type APPROVAL_CEREMONY_V1 and MUST be included in the EvidencePack slice for any dispute.

#### 2.7.2 Approval rate limits and quorum escalation (recommended, profile-gated)

Policies SHOULD constrain approval throughput per principal to prevent "approval fatigue" and credential abuse:

- MAX_APPROVALS_PER_HOUR
- MAX_HIGH_RISK_APPROVALS_PER_DAY
- REQUIRED_APPROVAL_QUORUM for specified action_urn patterns

If a rate limit is exceeded, CPI MUST return REQUIRE_APPROVAL with a quorum requirement, or DENY with reason_code = ERR_APPROVAL_RATE_LIMIT_EXCEEDED, depending on policy.

#### 2.7.3 Emergency override (optional, profile-gated)

Emergency overrides are permitted only when the policy explicitly enables them. An override MUST:

- be recorded as an ATTESTATION of type EMERGENCY_OVERRIDE
- bind to the intent_hash and the specific bypassed requirement (timelock, ceremony, or solver)
- include an emergency_reason string from a bounded policy-defined enum
- tighten post-action monitoring (for example: mandatory CHECKPOINT after the effect)

Emergency overrides MUST be treated as auditable exceptions, not a steady-state path.

### 2.8 Threat model (normative minimum)

HELM implementations MUST explicitly defend against the following adversary classes. Conformance claims are invalid unless the implementation can demonstrate mitigations for each class.

- A1 Prompt injection and instruction smuggling: hostile inputs in tools, documents, emails, web pages, and connector payloads.
- A2 Compromised connector or tool: a tool returns malicious outputs, exfiltrates data, or mutates state unexpectedly.
- A3 Replay and reordering attacks: captured intents, envelopes, or receipts are replayed or reordered to induce duplicate side effects.
- A4 Cross-tenant contamination: one tenant attempts to read or influence another tenant’s policies, receipts, or memory.
- A5 Insider and privileged operator risk: administrative access attempts to alter receipts, bypass PEP, or rewrite history.
- A6 Supply chain compromise: modified dependency, build pipeline, or signing keys attempt to counterfeit EvidencePacks or binaries.
- A7 Partial network partitions and clock skew: replicated systems diverge; attackers exploit nondeterminism caused by time or ordering.
- A8 Model substitution: upstream swaps models or weights to change behavior without proper attestation.
- A9 Systemic synchronization and monoculture: correlated agents act simultaneously, creating flash-crash dynamics.
- A10 Analog proxy abuse: an agent causes physical harm by delegating to human labor markets or off-platform actors.
- A11 Approval fatigue and rubber-stamping: repeated approvals degrade human governance into reflex, undermining intent review.

Normative security goals:
- G1 Fail-closed side effects: no external mutation occurs without a kernel-issued KERNEL_VERDICT bound to plan_hash and policy_hash.
- G2 Deterministic replay: any allowed effect can be replay-verified from EvidencePack or condensation inclusion proof.
- G3 Minimal disclosure: sensitive payloads are off-graph; proof nodes store hashes and ciphertext references only (Section 3.4).
- G4 Drift intolerance: connector schema drift causes fail-closed behavior, never probabilistic parsing (Section 10.5).


## 3. ProofGraph v1

### 3.1 Node envelope
Each ProofGraph node MUST be canonical JSON (RFC 8785) and MUST include:

- `node_hash`: `sha256:<hex>` over the canonical bytes of the node payload excluding `node_hash` itself.
- `kind`: one of the node kinds in Section 3.2.
- `parents`: array of parent node hashes (DAG).
- `lamport`: monotonic Lamport clock (uint64).
- `principal`: stable principal identifier (string, for example DID, SPIFFE, or X.509 subject).
- `principal_seq`: monotonic sequence number per principal (anti-replay).
- `ts_unix_ms`: OPTIONAL and MUST NOT be used for determinism. If present, it MUST be excluded from any deterministic tie-breakers.
- `payload`: kind-specific payload (Section 3.3).
- `sig`: detached signature over `HELM/ProofGraphNode/v1` domain-separated canonical bytes (Appendix A.2).

### 3.2 Node kinds
The `kind` enum is stable:

- `INTENT` - a signed request to act.
- `POLICY` - a policy publication or policy change receipt.
- `OBSERVATION` - observed facts (read-only) with provenance.
- `ATTESTATION` - signed proofs or approvals (includes kernel verdicts).
- `TRUST_EVENT` - trust registry events (trust roots add/rotate/revoke/pin) anchored for replayable verification.
- `EFFECT` - executed side effects with receipts.
- `CHECKPOINT` - state snapshots and condensation checkpoints.
- `MERGE_DECISION` - deterministic merge resolution under replication.

Kernel verdict representation:
- `KERNEL_VERDICT` MUST be represented as `kind: "ATTESTATION"` with `payload.attestation.type = "KERNEL_VERDICT"`.

### 3.3 Minimal payload contracts

#### 3.3.1 INTENT payload
- `intent_id` (uuid)
- `action_urn` (string, tool or connector action URN)
- `args_c14n_hash` (sha256 of canonical args JSON)
- `plan_hash` (sha256 of canonical Plan IR)
- `policy_hash` (sha256 of effective policy bundle)
- `risk_class` (enum: T0, T1, T2, T3)
- `ttl_ms` (uint32)
- `nonce` (string)

#### 3.3.2 ATTESTATION payload
- `type` (enum string, registry in Appendix F)
- `subject_hash` (sha256:<hex> of the object being attested)
- `evidence` (array of referenced hashes)
- `claims` (canonical JSON map)
- `signer` (principal)
- `sig_profile` (string, for example "EdDSA-Ed25519-JWS-Compact")

RFC-002 - Reality anchor claims (normative):

- If attestation_type == KERNEL_VERDICT and the attestation is intended for cross-tenant consumption (A2A Profile 3+), claims MUST include:
  - kernel_measurement_sha256: string (hex, 32 bytes) - SHA-256 of the exact kernel binary measured at load
  - execution_profile: string enum {TEE_ATTESTED, SOFTWARE_ONLY}
  - silicon_quote_b64: string - vendor quote/attestation blob proving kernel_measurement_sha256 inside the execution profile
  - tee_vendor: string enum {AWS_NITRO, INTEL_TDX, AMD_SEV_SNP, APPLE_SEP, OTHER}
  - tee_quote_format: string - opaque identifier for decoding/verification
- Any A2A verifier MUST reject KERNEL_VERDICT attestations with execution_profile != TEE_ATTESTED unless the envelope explicitly declares a weaker trust profile (Profile 2 only) and the receiver policy permits it.

- If attestation_type == MODEL_VERSION_ATTESTATION, claims MUST include:
  - cognitive_engine_id: string (provider-scoped identifier)
  - cognitive_engine_hash: string (hex, 32 bytes) if weights are locally available
  - provider_attestation_ref: string (hash ref) if weights are not locally available (closed API), binding to a provider-signed statement

- If attestation_type == HEARTBEAT_ATTESTATION, claims MUST include:
  - heartbeat_lamport: int64 (monotonic, per-org)
  - heartbeat_scope: string enum {SOVEREIGN, RECOVERY_QUORUM}

#### 3.3.3 EFFECT payload
- `tool_urn`
- `request_hash`
- `response_hash`
- `sandbox_profile`
- `connector_receipt` (optional)
- `effect_type` (optional; reserved strings for specialized kernel-controlled effects)
- `stdout_hash` / `stderr_hash` (optional)
- `status` ("ok" | "error")
- `reason_code` (if denied or failed)

#### 3.3.4 CHECKPOINT payload
- `checkpoint_type` ("STATE_SNAPSHOT" | "CONDENSATION_CHECKPOINT")
- `snapshot_hash` (for state snapshots)
- `condensation` (for condensation checkpoints, Section 8)

#### 3.3.5 MERGE_DECISION payload
- `conflict_class` (enum)
- `candidates` (array of node_hash)
- `winner` (node_hash)
- `tie_breaker` ("LEX_NODE_HASH")
- `reason_code`
- `policy_hash`

### 3.4 Sensitive payload rules (crypto-shredding and right-to-forget)
Problem: ProofGraph is immutable. Some jurisdictions require deletion or irrecoverability of sensitive data (for example PII).

Normative requirements:
- Permanent nodes (all T2 and T3, and any node retained beyond TTL) MUST NOT store raw sensitive payloads in cleartext.
- Each node payload MUST declare `payload_class`:
  - `Public`: safe to persist in cleartext.
  - `Sensitive`: MUST be stored off-graph as encrypted content.
- For `Sensitive` payloads, the node MUST store only:
  - `ciphertext_hash`: sha256 of the encrypted blob bytes.
  - `blob_ref`: opaque location reference (content-addressed where possible).
  - `payload_schema_ref`: schema or type identifier for decoding (non-sensitive).
  - `kek_ref`: key-encryption-key reference (tenant-scoped).
  - `redaction_policy_ref`: reference to the redaction policy in OrgGenome.

Crypto-shredding protocol (right-to-forget):
- Sensitive content MUST be encrypted with a unique data-encryption key (DEK) per blob.
- DEKs MUST be wrapped by a tenant key-encryption key (KEK) referenced by `kek_ref`.
- Key management determinism (required):
  - Implementations MUST treat key policy as part of the proof surface.
  - Any receipt or EvidencePack that references `kek_ref` MUST also reference `key_policy_hash` (sha256) for the applicable tenant key policy (rotation, revocation, legal hold behavior).
  - Key rotation MUST preserve replay verification of historical ciphertext_hash values; rotation cannot mutate past ciphertext.
- "Forget" MUST be implemented by deleting or revoking the KEK (or the DEK wrapper) so the ciphertext is permanently unreadable.
- ProofGraph MUST remain intact. Replay MUST still verify hashes, ordering, signatures, policy verdicts, and inclusion proofs. Human-readable sensitive content MAY be rendered as REDACTED.

Legal hold (optional, policy-gated):
- OrgGenome MAY define a legal hold policy that prevents shredding for specific classes for a bounded window.
- Any hold MUST emit an ATTESTATION of type `LEGAL_HOLD_APPLIED` with scope, reason, and expiry.

### 3.5 Versioning and extension rules (normative)

ProofGraph and EvidencePack schemas MUST be evolvable without silent forks.

- Every node MUST include:
  - schema_version (string, for example "1.0")
  - kind (stable enum)
  - payload_schema (string identifier, for example "proofgraph.intent.v1")
  - payload_schema_hash (sha256 of the canonical schema text or compiled schema bundle)

- Backward compatibility:
  - Implementations MUST accept unknown optional fields and ignore them.
  - Implementations MUST fail-closed on unknown required fields, unknown required kinds, or unknown payload_schema values referenced by policy.

- Extension mechanism:
  - New kinds MUST be introduced as OPTIONAL until a conformance profile requires them.
  - A feature MUST be declared via a FEATURE node that includes feature_id, min_schema_version, and required_conformance_level.
  - Reason codes MUST be stable strings. New reason codes MAY be added, but existing codes MUST NOT change semantics.

- Negotiation and downgrade:
  - For A2A, parties MUST exchange supported schema_versions and feature sets during session establishment.
  - If negotiation fails, the system MUST deny side effects and emit reason_code = ERR_VERSION_NEGOTIATION_FAILED.



### 3.6 Key management, rotation, and cryptographic agility (normative minimum)

Implementations MUST support cryptographic agility while preserving replay verification.

Key classes:
- identity_key: binds principal identity (human, org, agent).
- signing_key: signs ProofGraph nodes, A2A envelopes, and EvidencePack manifests.
- encryption_key: encrypts off-graph sensitive payloads (Section 3.4).

Requirements:
- Each signature MUST carry kid (key id) and the Trust Registry reference used to resolve it.
- Key rotation MUST be supported without invalidating historical replay:
  - Historical verification MUST use the key material valid at the time the node was signed.
  - Revocation MUST be represented as an immutable Trust Registry event with a lamport cutoff.
- Tenant isolation:
  - Each tenant MUST have a unique root signing namespace. Cross-tenant key reuse is non-conformant.
- Agility:
  - Ed25519 MUST be supported.
  - Implementations SHOULD support at least one additional algorithm suite for migration.
  - If an algorithm is deprecated, conformance harness MUST provide migration vectors and negative tests.

Replay rule:
- A replay verifier MUST evaluate signature validity against (kid, trust_registry_state_at_lamport, node_lamport) and MUST NOT consult wallclock time.



### 3.7 Redaction and access profiles (required for any human-facing surface)
Human-facing receipts and explorers MUST support deterministic redaction.

Normative requirements:
- OrgGenome MUST define redaction policies per payload class and per role.
- Renderers MUST declare `redaction_profile` on each view (for example OPERATOR, AUDITOR, REGULATOR, PUBLIC).
- A renderer MUST NOT infer or "fill in" missing sensitive fields. If data is unavailable (for example shredded), it MUST render REDACTED and include a deterministic reason code.
- A receipt view MUST include hashes of any revealed plaintext fields so that selective disclosure is auditable.

## 4. Verified Genesis Loop (VGL) v1
VGL is mandatory when compiling or materially modifying an OrgGenome from natural language and contextual ingestion (B2C/SMB) and RECOMMENDED for enterprise onboarding.

Goal: prevent compiler hallucination from becoming deterministic law.

### 4.1 Genesis phases
1) Ingest: gather user intent plus authorized context sources.
2) Propose OrgGenome: stochastic compiler proposes a concrete OrgGenome AST and an impact summary.
3) Deterministic Semantic Mirror (mandatory):
   - Kernel renders the OrgGenome AST into deterministic human text (no models).
   - User signs both the AST hash and the rendered text hash.
4) Blast Radius Wargaming (mandatory for B2C/SMB, optional for enterprise):
   - Run bounded simulation and adversarial scenario generation in a sandbox.
   - Produce an Impact Report with worst-case bounds and discovered loopholes.
5) P0 Ceilings binding (mandatory for B2C/SMB):
   - User sets platform-level ceilings via explicit UI controls (not free text).
   - P0 ceilings apply beneath OrgGenome policies and cannot be bypassed.
6) Activate:
   - Kernel emits an ATTESTATION node ORG_GENESIS_APPROVAL referencing the OrgGenome hash, mirror hash, impact report hash, and P0 ceiling hash.
   - OrgGenome becomes active only after this attestation exists.

### 4.2 Deterministic Semantic Mirror (Isomorphic Semantic Mirroring)
Requirements:
- The OrgGenome compiler MUST emit a constrained AST (MAPL AST plus budgets plus risk classes).
- The kernel MUST provide a deterministic AST-to-text renderer.
- The signed approval MUST be over:
  - org_genome_hash
  - mirror_text_hash
  - p0_ceilings_hash (if present)
  - impact_report_hash (if present)

Fail-closed: if semantic mirror cannot be rendered, activation MUST be denied.


### 4.2.1 Delegated semantic mirroring for franchised genomes

If an OrgGenome is instantiated from a franchised genome, the strict requirement for the tenant principal to dual-sign the full P1 MAPL AST is waived, provided all conditions hold:

1. The template vendor principal MUST sign the org_genome_hash and impact_report_hash.
2. The tenant principal MUST sign a P0_CEILINGS_ACTIVE attestation representing their sovereign limits.
3. The UI rendering the P0 signature MUST NOT abstract values. It MUST display exact scalar integers (for example: MAX_DAILY_SPEND_USD: 50).
4. The kernel MUST emit an ORG_GENESIS_APPROVAL attestation binding both signatures and the active profile parameters (Section 2.4).

### 4.3 Blast Radius Wargaming
Purpose: compute blast radius of a compiled OrgGenome before activation.

Minimum requirements:
- Wargaming MUST run in a deterministic sandbox with:
  - `seed` (uint64) recorded.
  - `gas_limit_steps` (max simulation ticks) recorded.
  - `time_limit_ms` recorded.
  - network disabled. All connectors MUST be stubbed or fixture-driven.
- Scenario generation MUST cover at minimum:
  - spend escalation and budget exhaustion
  - data exfil and egress caps
  - privilege escalation across tools and principals
  - runaway loops and retry storms

Deterministic fixtures:
- Any external dependency MUST be represented by recorded fixtures (request/response pairs) or synthetic stubs.
- The fixture set MUST be hashed and referenced by the impact report.

Impact report (required fields):
- `seed`, `gas_limit_steps`, `gas_used_steps`, `time_limit_ms`, `time_used_ms`
- maximum theoretical daily spend given policies and P0 ceilings
- maximum API-call volume given ceilings
- maximum data egress and datasets touched
- all loopholes found with minimal reproductions
- recommended policy patches
- termination: NORMAL or GAS_EXHAUSTED or TIME_EXHAUSTED

Fail-closed:
- If termination is GAS_EXHAUSTED or TIME_EXHAUSTED, the OrgGenome MUST be considered unstable and MUST NOT activate without explicit override attestation.

Packaging:
- The impact report MUST be packaged as an EvidencePack and referenced by hash.

### 4.4 P0 Immutable Ceilings
P0 ceilings are platform invariants evaluated before OrgGenome (P1) and transient overlays (P2). If a velocity ceiling is breached, the system MUST DENY with reason_code = ERR_VELOCITY_CEILING_BREACH.

Minimum ceiling set:
- MAX_DAILY_SPEND_USD
- MAX_SPEND_USD_PER_HOUR (velocity ceiling)
- MAX_VENDOR_SPEND_USD per vendor or connector
- MAX_SIDE_EFFECT_ACTIONS_PER_DAY
- MAX_SIDE_EFFECT_ACTIONS_PER_MIN (velocity ceiling)
- FORBIDDEN_ACTION_URNS (deny list)
- MAX_DATA_EGRESS_BYTES_PER_DAY

Two-phase commit (required):
1) P0_PENDING (grace window)
- Any P0 change MUST first enter P0_PENDING for a deterministic grace period (default 15 minutes).
- During grace, the root principal MAY revert the pending change without time-lock.
- A P0_PENDING record MUST be written as an ATTESTATION of type `P0_CEILINGS_PENDING` containing:
  - p0_ceilings_hash
  - grace_deadline_lamport or grace_deadline_unix_ms (informational only)
  - Enforcement MUST be based on Lamport order. Any unix_ms value MUST NOT affect determinism, ordering, or eligibility for activation.
  - signer and auth profile

2) P0_ACTIVE (time-lock)
- After grace, the change becomes P0_ACTIVE by emitting an ATTESTATION of type `P0_CEILINGS_ACTIVE`.
- Increases to ceilings MUST be time-locked (minimum 24 hours) and require strong authentication.
- Decreases MAY be immediate but MUST still emit `P0_CEILINGS_ACTIVE` and MUST not bypass policy hashing.

Emergency override (required):
- A root principal MUST have an emergency path to unbrick obvious mistakes (for example MAX_DAILY_SPEND_USD=0).
- Emergency override MUST require MFA and MUST emit an ATTESTATION of type `P0_EMERGENCY_OVERRIDE` including:
  - prior_p0_hash, new_p0_hash
  - reason_code
  - auth_profile
- Policies MAY restrict frequency and scope of emergency overrides (rate limit).

Hashing and enforcement:
- P0 MUST be included in policy hashing and VPL decisions.
- Any side effect that would exceed P0 MUST be denied with a deterministic reason code and recorded.


### 4.5 Temporal key decay and succession (Dead Man's Clock)

Goal: prevent "ghost ship" autonomy and avoid permanent bricking under sovereign key loss, without introducing any fail-open backdoors.

Normative requirements:

- Every OrgGenome MUST define:
  - heartbeat_interval_lamport: int64 (minimum cadence at which the sovereign principal must renew control)
  - heartbeat_grace_unix_ms: int64 (informational only; enforcement MUST be lamport-based)
  - recovery_quorum: object:
      - scheme: enum {SHAMIR_K_OF_N}
      - k: int
      - n: int
      - participants: array of {participant_id, public_key_ref, role}
  - resurrection_timelock_unix_ms: int64 (default: 72h)
- The kernel MUST require periodic HEARTBEAT_ATTESTATION from the sovereign principal (or an explicitly configured quorum signer).
- If heartbeat expires, the kernel MUST enter P0_TERMINAL_FREEZE:
  - all external connectors become deny-by-default
  - only read-only and evidence export are permitted
  - any WRITE/EFFECT step MUST return ERR_HEARTBEAT_EXPIRED_FREEZE
- Identity resurrection:
  - a recovery_quorum MAY submit an IDENTITY_RESURRECTION_INTENT referencing the next_root_public_key_ref and a quorum proof
  - the kernel MUST enforce a deterministic time-lock window (resurrection_timelock_unix_ms) before rotation
  - during the time-lock, the current sovereign key MAY veto by emitting a RESURRECTION_VETO attestation
  - if no veto occurs, the kernel MUST rotate the root principal to the new key and emit an ATTESTATION of type IDENTITY_ROTATION


## 5. Verified Planning Loop (VPL) v1

### 5.1 Phases
1) Propose plan (stochastic) as Plan IR v1 (Section 6).
2) Governance pre-check (optional sensors, Section 9).
3) CPI validation (tiered ladder, Section 6.3).
4) Kernel verdict emitted (KERNEL_VERDICT attestation).
5) Execute via PEP to tools/connectors with sandbox enforcement.
6) Emit receipts and effects into ProofGraph.
7) Checkpoint and, if applicable, condense low-risk trails (Section 8).

### 5.2 VPL invariants (MUST)
- No side effect without a matching INTENT and KERNEL_VERDICT attestation binding plan_hash plus policy_hash.
- principal_seq MUST be strictly increasing per principal. Any stale or replayed sequence MUST be denied.
- Plan hashes and policy hashes MUST be computed over canonical bytes (Appendix A).
- Denials MUST be recorded with deterministic reason codes.
- The kernel MUST be able to replay any workflow deterministically using ProofGraph plus EvidencePacks and reach the same verdicts.

## 6. Plan IR v1 and CPI v1

### 6.1 Plan IR v1 (canonical contract)
Plan IR is canonical JSON with stable ordering.

Minimum fields:
- plan_id (uuid)
- org_genome_hash
- created_by_principal
- created_at_lamport
- risk_class (T0-T3)
- nodes: map keyed by node_id
- edges: array of objects with keys from, to, type
- budgets: object with usd_max, api_calls_max, time_ms_max
- constraints: array of typed constraints

Node shape:
- node_id
- op (enum: TOOL_CALL, CONNECTOR_CALL, APPROVAL, READ, WRITE, COMPUTE, CHECKPOINT, REASONING_LOOP)
- action_urn
- args (canonical JSON)
- preconditions (array)
- postconditions (array)
- side_effect (bool)
- effect_class (REQUIRED if side_effect=true; enum: REVERSIBLE, COMPENSABLE, IRREVERSIBLE)
- compensation_urn (OPTIONAL; required for COMPENSABLE when a deterministic compensate path exists)
- irreversible_reason (OPTIONAL; required for IRREVERSIBLE to explain why no rollback exists)
- idempotency_scope (OPTIONAL; hints for connector idempotency domain)
- estimated_cost (optional)

Constraints (minimum enum):
- SEQUENCE (must occur before or after)
- BUDGET (caps)
- DATA_SCOPE (allowed datasets, egress caps)
- AUTH (required approvals)
- TOOL_SANDBOX (required sandbox profile)
- CONNECTOR_IDEMPOTENCY (required idempotency keys)


### 6.1.1 Epistemic reasoning phase (condensation-backed reasoning loops)

To preserve agentic fluidity without bloating the immutable ProofGraph, a Plan IR MAY include a REASONING_LOOP node.

Normative requirements:

- A REASONING_LOOP node MUST have side_effect: false and risk_class: T0.
- Within a REASONING_LOOP, the agent MAY perform READ and COMPUTE operations against the Learned Knowledge Store (LKS) and T0 connectors explicitly permitted by policy.
- The runtime MUST continuously hash all reasoning loop operations (prompts, tool reads, observations, and intermediate outputs) into a local Merkle tree.
- The runtime MUST emit a reasoning_merkle_root at loop exit.
- Any subsequent state-mutating INTENT that depends on the loop MUST include the reasoning_merkle_root in evidence_refs, so causal provenance remains cryptographically unbroken.
- Retention of raw reasoning-loop transcripts MUST follow risk-class routing and condensation policy (Section 8). At minimum, the reasoning_merkle_root MUST be retained.

### 6.2 CPI v1 (Constraint and Proof Interface)
CPI is a deterministic function.

Input:
- plan_ir (canonical)
- effective_policy_bundle (canonical)
- org_phenotype_snapshot_hash (optional)
- p0_ceilings (optional)
- mode (FAST_PATH | BOUNDED_TEMPORAL | SOLVER)

Output:
- verdict (ALLOW | DENY | DEFER | REQUIRE_APPROVAL)
- reason_code (stable string)
- unsat_core (optional, minimal counterexample)
- evidence_refs (hashes of inputs, solver config, sandbox profile)
- validator_hash (sha256 of CPI implementation identity)

All CPI outputs MUST be canonical JSON.

Kernel verdict binding (normative):
- The kernel MUST emit a KERNEL_VERDICT ATTESTATION for every CPI evaluation (ALLOW, DENY, DEFER, REQUIRE_APPROVAL).
- The attestation MUST bind subject_hash = plan_hash and MUST include policy_hash.
- The attestation claims MUST include verdict and reason_code.
- Side effects MUST execute only when verdict=ALLOW and the attestation exists.

### 6.3 Layered validation ladder (mandatory)

CPI defines a layered validation ladder. The ladder is selected by HUDF profile parameters (Section 2.4) and by MAPL policy.

CPI is deterministic: given the same CPIRequest and policy bundle, it MUST return the same verdict.

#### Tier 1 - Structural fast-path (required)

Tier 1 MUST execute for every step.

Tier 1 checks include, at minimum:

- bytes contract: canonical JSON (RFC 8785) for all signed envelopes
- schema validation of Plan IR and tool args
- capability scope check (allowed ops, connectors, and destinations)
- P0 ceilings and budgets (hard-deny)
- connector contract pinning and request/response shape checks (fail-closed on drift)
- idempotency key derivation requirements for side effects

Verdict outcomes: ALLOW, DENY.

#### Tier 2 - Bounded temporal and checkpoint invariants (recommended)

Tier 2 is used for multi-step workflows where local sequencing matters.

Tier 2 requirements:

- bounded horizon verification: verify the next N steps (default N=5) for type and precondition satisfaction
- checkpoint invariants: enforce global constraints at declared checkpoints (for example: "no irreversible effects before approval checkpoint")
- irreversibility ordering: any step tagged IRREVERSIBLE MUST be topologically sorted to the end of the execution DAG, or CPI MUST require an explicit approval checkpoint before crossing the irreversible threshold

Verdict outcomes: ALLOW, DENY, DEFER, REQUIRE_APPROVAL.

#### Tier 2.5 - Cryptographic HITL (human-in-the-loop) substitution (profile-gated)

Tier 2.5 provides a deterministic path to execute high-stakes actions in environments that do not enable Tier 3 formal solving (common for SMB profiles), without failing open.

Trigger conditions:
- If policy returns REQUIRE_APPROVAL, execution MUST halt.
- If Tier 3 is unavailable or times out and the active profile enables TIER_2_5_DOWNGRADE, CPI MAY convert DEFER into REQUIRE_APPROVAL.

Tier 2.5 requirements (normative):

- The kernel MUST emit a signing request as an ATTESTATION node of type HITL_SIGN_REQUEST containing:
  - intent_hash
  - predicted impact summary fields (amounts, destinations, irreversible_reason if applicable)
  - active P0 ceilings and any velocity ceilings applicable to this step
  - required ceremony parameters (if any) and required quorum (if any)
- The approval MUST be recorded as an ATTESTATION node of type HITL_APPROVAL containing:
  - intent_hash
  - signer principal identifiers and signatures (sig or sigset)
  - referenced HITL_SIGN_REQUEST node_hash
  - approval_mode (SIMPLE or CEREMONY_V1)
  - approval_timelock_ms enforced (if required)
  - optional embedded APPROVAL_CEREMONY_V1 fields or a reference to an APPROVAL_CEREMONY_V1 node

High-risk irreversible approvals:
- For risk_class >= T2 and effect_class = IRREVERSIBLE, the runtime MUST enforce:
  - an approval timelock (default 30 seconds) and
  - an approval ceremony (RFC-005 Section 2.7.1)
- If the ceremony is missing, CPI MUST DENY with reason_code = ERR_APPROVAL_CEREMONY_REQUIRED.
- If the timelock has not elapsed, CPI MUST DEFER with reason_code = ERR_APPROVAL_TIMELOCK_ACTIVE.

Rate limits and quorum:
- If policy defines approval rate limits, the kernel MUST enforce them. If exceeded:
  - CPI MUST REQUIRE_APPROVAL with a higher quorum, or
  - CPI MUST DENY with reason_code = ERR_APPROVAL_RATE_LIMIT_EXCEEDED,
  depending on policy.

Emergency override:
- If policy enables emergency overrides, a bypass MUST be recorded as an ATTESTATION of type EMERGENCY_OVERRIDE and referenced by HITL_APPROVAL. Otherwise, bypass attempts MUST be denied with reason_code = ERR_APPROVAL_EMERGENCY_OVERRIDE_REQUIRED.

Verdict outcomes: ALLOW (after valid approval), DENY, DEFER.

#### Tier 3 - Formal proof (enterprise profile)

Tier 3 is reserved for actions whose blast radius justifies formal proof latency (for example: large payments, IAM changes, regulated data exports).

Tier 3 requirements:

- the plan MUST be translated into a formal constraint system (for example: LTL or SMT constraints) whose input symbols are derived deterministically from Plan IR, policy, and the current OrgPhenotype snapshot hash
- the solver MUST run with a deterministic timeout budget
- if the solver returns UNSAT, CPI MUST DENY and SHOULD emit a machine-readable counterexample (UNSAT core) when available
- if the solver times out, CPI MUST NOT fail open. CPI MUST return DEFER with ERR_CPI_SOLVER_TIMEOUT, unless Tier 2.5 downgrade is enabled, in which case CPI MAY return REQUIRE_APPROVAL

Verdict outcomes: ALLOW, DENY, DEFER, REQUIRE_APPROVAL.




## 7. Policy semantics (MAPL v1)

### 7.1 Core model
MAPL policies are compiled into a canonical Policy Bundle:
- Caller policy yields capability set C
- Resource policy yields capability set R
- Effective permission set is E = C intersect R
- Deny overrides allow at rule level within each bundle, and global deny overrides allow.

A capability element is:
- verb (string)
- resource_pattern (string)
- constraints (canonical JSON map, for example budget caps, time windows)
- risk_class_max (T0-T3)

### 7.2 Policy compilation and hashing
- MAPL source MUST compile to a canonical AST.
- The AST MUST compile to a canonical Policy Bundle JSON.
- policy_hash is sha256 over canonical bytes of the Policy Bundle plus P0 ceilings (if present).
- Any policy change MUST emit a POLICY node with prior hash, new hash, diff summary, and signer.

### 7.3 Attestation semantics (minimum)
Attestations referenced by policy MUST declare:
- attestation.type
- required signature profile
- validity window
- trust root source (Trust Registry)

Minimum attestation types:
- ORG_GENESIS_APPROVAL
- KERNEL_VERDICT
- P0_CEILINGS_PENDING
- P0_CEILINGS_ACTIVE
- P0_EMERGENCY_OVERRIDE
- LEGAL_HOLD_APPLIED (if legal hold is supported)
- CODE_REVIEW_SIGNED
- REPRO_BUILD_ATTESTED
- SANDBOX_PROFILE_ATTESTED
- CONNECTOR_CONTRACT_ATTESTED
- MERGE_DECISION_SIGNED
- IRREVERSIBLE_THRESHOLD_APPROVED (for irreversible execution checkpoints)

## 8. Proof Condensation (unit economics without losing truth)

### 8.1 Risk-tiered evidence routing
Every action MUST be assigned a risk class.

- T3 High risk: money movement, identity or IAM changes, data exfil, production deploys, legal commitments.
  - Full EvidencePack required.
  - Permanent ProofGraph nodes for each action.
- T2 Medium risk: customer messaging at scale, data mutations, vendor contract changes.
  - EvidencePack required by policy or thresholds.
  - Permanent nodes, optional condensation.
- T1 Low risk: reads, drafts, local computations, non-sensitive scraping.
  - Hash-only receipts plus periodic checkpoints.
- T0 Routine: high-volume routine operations.
  - Always condensed.

Risk assignment MUST be policy-derived and tool-derived (tool metadata declares default risk; policy can raise it).

### 8.2 Merkle condensation checkpoint
For T0 and optionally T1 and T2, HELM MAY store ephemeral event details off the canonical graph, and MUST commit:
- A CHECKPOINT node of type CONDENSATION_CHECKPOINT containing:
  - range_start_lamport, range_end_lamport
  - event_count
  - merkle_root (sha256)
  - aggregate_stats (costs, counts, violations=0/1)
  - policy_hash
  - gc_ttl_days
  - cold_storage_ref (opaque reference, hashed)

The Merkle tree leaves are canonical hashes of ephemeral events. Inclusion proofs MUST be reproducible.

### 8.3 Asymmetric evidence GC
- Ephemeral payloads for condensed ranges MUST be stored in cheap cold storage with a TTL.
- During TTL, disputes MUST be resolvable by rehydrating payload and proving inclusion under merkle_root.
- After TTL, payloads MAY be deleted. The condensation checkpoint remains permanent.

Sensitive payload handling (required):
- If condensed payloads contain sensitive data, they MUST be stored encrypted, and the keying MUST support crypto-shredding.
- After TTL, deletion MAY be implemented by deleting blobs or by deleting/revoking the wrapping keys (crypto-shred). In both cases, the ProofGraph checkpoint remains valid because it commits only hashes and merkle_root.

Legal hold (optional, policy-gated):
- OrgGenome MAY prevent deletion or shredding for specific condensed ranges for a bounded window.
- Any hold MUST be recorded as an ATTESTATION of type LEGAL_HOLD_APPLIED referencing the condensation checkpoint.

## 9. Governance sensors v1 (baseline plus optional)

Governance sensors are pre-execution checks that classify or gate intents based on deterministic signals. Sensors never execute side effects. Sensors MUST emit canonical OBSERVATION nodes so their influence is auditable and replayable.

### 9.1 Governance Sensor Interface (GSI)
A sensor evaluation MUST have:
- input: INTENT hash, plan_hash, policy_hash, current OrgPhenotype snapshot hash (if available)
- output: one of {PASS, WARN, REQUIRE_APPROVAL, HALT}
- reason_code (Appendix C)
- evidence_refs (hashes of features, thresholds, config)
- sensor_hash (sha256 of the sensor implementation identity)

If a sensor uses any model or statistical component, the model identity and config MUST be hashed and included in evidence_refs. The sensor decision MUST still be deterministic for the same inputs.

### 9.2 Baseline sensors (normative)
Implementations MUST include at least these baseline sensors:

- Budget proximity sensor: escalates when projected spend approaches caps.
- Egress anomaly sensor: escalates when data scope or egress volume deviates from historical norms defined in OrgGenome.
- Permission delta sensor: halts if a plan implies privilege changes without explicit approvals.
- Deny-loop sensor: halts if the system repeats the same denied action pattern without changing plan_hash or policy_hash.
- Connector drift sensor: halts if any connector reports CONTRACT_DRIFT within the plan window.

### 9.3 Optional sensors (non-normative)
Optional sensors MAY include:
- Activation telemetry sensors and interpretability probes.
- Red-team prompt injection detectors.
- Vendor reputation and threat intel feeds.

Optional sensors MUST obey the same deterministic GSI contract and must not become hard dependencies for L1 or L2 conformance.

## 10. Tool execution containment and connector determinism

### 10.1 Tool sandbox minimums (required for L1+)
Any side-effectful tool execution MUST run in an enforced sandbox profile.

Minimum requirements:
- Ephemeral execution environment.
- Explicit network egress allowlist.
- Read-only filesystem by default, explicit mounts only.
- Base image or runtime digest pinned.
- Sandbox config hashed and recorded in receipts.

If sandbox requirements cannot be satisfied, execution MUST be denied.


### 10.1.1 Dynamic code execution capabilities

If a tool executes dynamically generated code (for example: agent-authored data mappers or scripts), general-purpose containers are insufficient because they create unbounded I/O surfaces.

A conformant runtime MUST enforce one of the following I/O isolation models:

1. WASI-native execution (recommended)
   - dynamic code MUST compile to WebAssembly and execute under a WASI profile
   - network and filesystem capabilities MUST default to deny
   - the module MUST act as a pure data transformer: f(inputs) = outputs
   - the wasm_module_sha256 MUST be recorded in the ToolReceipt and referenced by the resulting INTENT or EFFECT
   - the runtime MUST enforce deterministic fuel (gas) and time limits for module execution
   - if fuel or time is exhausted, the sandbox MUST terminate and MUST emit ERR_COMPUTE_GAS_EXHAUSTED or ERR_COMPUTE_TIME_EXHAUSTED

2. eBPF syscall trapping (enterprise profile)
   - if arbitrary runtimes (for example: Python) are permitted, the hypervisor MUST inject an eBPF trap
   - any attempt to open a network socket or unauthorized file descriptor MUST be trapped, serialized into canonical JSON, and routed synchronously through CPI/PEP
   - if the exact trapped payload and destination are not pre-authorized by policy, the sandbox MUST be terminated and MUST emit ERR_SANDBOX_SYSCALL_VIOLATION

### 10.2 Connector contract (required for any external API)
Connectors MUST provide:
- Executor declaration (mandatory): each connector operation MUST declare executor_kind (DIGITAL, KINETIC, ANALOG).
- Deterministic idempotency keys for every side-effectful request.
- Idempotency key derivation (mandatory): key MUST be derived from (intent_hash, connector_urn, effect_index) and MUST remain stable across retries and replays.
- Timeout and retry safety (mandatory): on ambiguous outcomes (for example HTTP 504), the connector MUST treat the request as unknown and MUST resolve state via idempotency lookup or reconciliation before any retry that could duplicate effects.
- Canonical request serialization and hashing.
- Response hashing with schema version.
- Effect metadata (mandatory): each connector operation MUST declare effect_class (REVERSIBLE, COMPENSABLE, IRREVERSIBLE) and executor_kind (DIGITAL, KINETIC, ANALOG) and MUST expose:
  - rollback or compensate entrypoint when applicable
  - irreversible_reason when no rollback exists
- Retry and backoff rules recorded deterministically.
- Saga and compensation semantics for reversible or compensable operations.
- Drift detection: if remote schema or behavior changes, connector MUST fail-closed and emit a CONTRACT_DRIFT receipt.


### 10.5 Connector drift detection and contract test suite (normative)

Connectors bind HELM to nondeterministic external APIs. Drift is guaranteed. HELM implementations MUST detect drift and fail-closed.

Minimum drift detection signals:
- schema_hash_mismatch: response does not match pinned schema hash (OpenAPI/JSON Schema/Pact).
- required_field_missing: any required field absent or type mismatch.
- enum_extension_unrecognized: new enum value received where policy requires closed-set.
- endpoint_fingerprint_change: TLS certificate pin change or endpoint identity change where pinning is required.
- rate_limit_behavior_change: unexpected retry-after semantics that could cause duplication.

Normative connector contract:
- Each connector MUST publish a Connector Contract Bundle containing:
  - connector_id, version, schema artifacts, schema_hash, and allowed endpoints
  - idempotency key strategy and required headers
  - saga and compensation capabilities and limitations
  - irreversibility tags per action

Contract tests:
- A conformance harness MUST include a connector drift test suite that can be executed offline using recorded fixtures.
- On schema_hash_mismatch or required_field_missing, the connector MUST:
  - DENY the action
  - emit reason_code = ERR_CONNECTOR_CONTRACT_DRIFT
  - require explicit operator intervention to update the contract bundle
- Probabilistic parsing or best-effort adaptation to drift is strictly forbidden.



#### 10.5.1 Asynchronous translation shims (mechanic agents)

Probabilistic parsing or runtime auto-healing by the executing agent is strictly forbidden. Drift MUST fail closed.

To reduce mean-time-to-repair without opening a cryptographic backdoor, implementations MAY employ out-of-band mechanic agents that generate deterministic translation shims.

Requirements:

1. On ERR_CONNECTOR_CONTRACT_DRIFT, the execution thread MUST halt safely.
2. A mechanic agent MAY analyze the drift and produce a deterministic shim (recommended: WebAssembly module) that maps the new schema to the pinned canonical schema.
3. The shim MUST be simulated offline against the failed payload and MUST demonstrate policy equivalence (no new capabilities, no widened destinations).
4. The shim MUST NOT be applied automatically. It MUST be explicitly authorized:
   - enterprise: via CI/CD review and platform principal signature
   - SMB: via a tenant principal signature (for example: biometric passkey approval)
5. Authorization MUST be recorded as a POLICY update (P2 transient overlay or P1 epoch change) that binds shim_sha256 to the connector contract hash.
6. Once authorized, the halted intent MAY be replayed with the shim in place, and ProofGraph MUST record the shim_sha256 in the replay receipt.


## 11. HELM A2A Envelope v1
Normative encoding is canonical JSON (RFC 8785). JSON-LD and procurement metadata MAY be attached as non-normative metadata but MUST NOT affect signature bytes.
Unsigned metadata is display-only.
- Any field that influences verification, policy, billing, or replay MUST be inside the signed canonical payload.
- If JSON-LD is attached as optional metadata, all @context data MUST be embedded (no remote fetch) and MUST NOT affect signature bytes.


### 11.1 Required fields
- v ("1")
- msg_id (uuid)
- iss (issuer principal)
- sub (subject principal)
- aud (intended audience)
- nonce
- ttl_ms
- issued_lamport
- capabilities (array of capability grants)
- plan_hash (optional)
- policy_hash
- risk_class
- payload (canonical JSON)
- sig (JWS compact, EdDSA Ed25519, Section 11.3)

### 11.2 Verification ladder (profiles)
- Profile 1 Baseline:
  - signature plus nonce plus TTL plus capability scope
- Profile 2 Reproducible:
  - Profile 1 plus reproducible execution receipt (build digest, container digest) referenced by hash
- Profile 3 Hardware-Assured (optional):
  - Profile 2 plus hardware attestation evidence (vendor-neutral claim set)

### 11.3 JOSE profile (minimum)
- JWS serialization: Compact
- alg: EdDSA (Ed25519)
- Required claims in payload: iss, sub, aud, exp or TTL mapping, jti or msg_id, nonce, policy_hash
- Trust roots are distributed via a Trust Registry anchored in ProofGraph.

### 11.4 Trust Registry v1 (normative minimal)

The Trust Registry distributes, rotates, and revokes trust roots used to verify signatures and attestations. It is anchored in ProofGraph so verification is replayable.

Trust Registry records MUST be represented as ProofGraph nodes of kind TRUST_EVENT with canonical payloads.

Minimum record types:
- TRUST_ROOT_ADD: adds a trust root (tenant root, vendor root, or ecosystem root).
- TRUST_ROOT_ROTATE: introduces a new key and deprecates a prior key.
- TRUST_ROOT_REVOKE: revokes a key id (kid) with a lamport cutoff.
- TRUST_PIN_SET: pins a trust root set for a tenant or policy bundle.

Minimum payload fields:
- trust_event_type (enum)
- subject (tenant_id or vendor_id)
- kid
- public_key (JWK or raw bytes, canonical encoding)
- alg
- valid_from_lamport
- valid_to_lamport (optional)
- revoke_from_lamport (optional, required for REVOKE)
- reason (stable string)
- signer (who authorized the change)
- evidence_refs (optional)

Distribution:
- Implementations MUST be able to export a deterministic Trust Registry snapshot for any lamport height.
- A verifier MUST be able to verify any node or envelope using only:
  - the EvidencePack slice
  - the Trust Registry snapshot at the relevant lamport height
- Wallclock time MUST NOT be used to determine trust validity.

Rotation and pinning:
- Policies MAY require that a given action be signed under a pinned trust root set.
- If a required trust root is missing, the system MUST DENY with reason_code = ERR_TRUST_ROOT_MISSING.

Island mode (profile-gated):
- If the Trust Registry is unreachable, a runtime MAY enter ISLAND_MODE if the active HUDF profile enables it.
- In ISLAND_MODE, the runtime MUST pin to the last verified Trust Registry snapshot and MUST DENY any action that requires a missing trust root.
- In ISLAND_MODE, the runtime MUST DENY:
  - new counterparty introductions (no TRUST_ROOT_ADD) with reason_code = ERR_ISLAND_MODE_RESTRICTED
  - cross-island settlement operations with reason_code = ERR_ISLAND_MODE_RESTRICTED
- In ISLAND_MODE, the runtime MAY ALLOW local operations to already-pinned counterparties within tightened P0 ceilings.




## 12. Metering receipts and dispute replay

### 12.1 Usage receipt (USAGE_NODE)
A usage record MUST be an ATTESTATION of type USAGE_NODE containing:
- metering_period
- principal
- vendor (optional)
- units (tokens_in, tokens_out, tool_calls, api_calls, cpu_ms, egress_bytes)
- merkle_ref (optional, if condensed)
- evidence_refs

### 12.2 Dispute resolution (normative)
Any dispute MUST be resolvable by deterministic replay:
- Tenant provides an EvidencePack slice or condensation checkpoint plus inclusion proof.
- Provider replays with helm replay and compares canonical usage outputs.
- If outputs differ, the discrepancy is resolved in favor of the deterministic replay artifact.

### 12.3 EvidencePack minimum contents (normative)

EvidencePack is the portable, replayable proof bundle for a specific scope (single action, macro-plan slice, or condensation checkpoint). EvidencePack MUST be deterministically packaged (Appendix A) and MUST include a manifest that enumerates all included artifacts and their sha256 digests.

Minimum required artifacts by risk class:

- T3 High risk (full pack REQUIRED):
  - manifest.json (canonical, includes tool class, connector id, and risk class)
  - orggenome_policy_bundle.json (policy bundle + hashes referenced by policy_hash)
  - p0_ceilings.json (if applicable) and p0_attestations (pending/active/override)
  - plan_ir.json (canonical) and plan_hash
  - cpi_inputs.json (canonical) and cpi_output.json (canonical, includes reason_code and validator_hash)
  - kernel_verdict.json (signed KERNEL_VERDICT)
  - tool_receipt.json (stdout/stderr hashes, egress summary, sandbox profile hash)
  - connector_transcript.enc (encrypted request/response transcript, with ciphertext_hash recorded on-graph)
  - redaction_map.json (optional) and encryption_metadata.json (required if any encrypted artifact exists)
  - environment.json (sandbox image digest, OS/arch, runtime versions)
  - replay_instructions.txt (non-normative)

- T2 Medium risk (reduced pack REQUIRED):
  - manifest.json
  - plan_ir.json, cpi_output.json, kernel_verdict.json
  - tool_receipt.json
  - minimal connector transcript (encrypted or redacted)
  - policy_hash references

- T0/T1 Low risk (condensed):
  - individual actions MAY be ephemeral.
  - at each condensation checkpoint, a checkpoint EvidencePack MUST include:
    - manifest.json
    - condensation_checkpoint.json (merkle_root, counts, period, risk mix)
    - inclusion_proof_format.json (schema for inclusion proofs)
    - optional sample transcripts (policy-defined)

Tool class requirements:
- For tools with irreversible effects, EvidencePack MUST include the irreversible gate attestation (IRREVERSIBLE_THRESHOLD_APPROVED) if applicable.
- For money movement, EvidencePack MUST include idempotency key derivation evidence (intent_hash, connector_idempotency_key).
- For data egress, EvidencePack MUST include data scope policy bundle and egress measurements.

Retention:
- Evidence retention is policy-defined. Conformance requires that deletion of raw payloads never invalidates the ability to verify hashes and inclusion proofs.



## 13. Replication and merge semantics (L2)

### 13.1 Merge preconditions
- Both branches MUST be valid DAGs.
- Both branches MUST have consistent policy hashes for the merge range or an explicit policy-change receipt.

### 13.2 Conflict classes
Minimum classes:
- STATE_CONFLICT
- POLICY_CONFLICT
- RECEIPT_CONFLICT
- IDENTITY_CONFLICT

### 13.3 Deterministic tie-breaker (MUST)
Timestamps MUST NOT decide merges.
Absolute tie-breaker is lexicographic order of node_hash values.

### 13.4 MERGE_DECISION node (mandatory)
A merge MUST emit a MERGE_DECISION node signed by the merging authority, recording candidates, winner, conflict class, and reason code.


### 13.5 Island Mode and partition operations (L2+)

Partitions are normal. Conformance requires that partitions do not create nondeterministic side effects.

Requirements:
- A runtime MAY continue local operation during partitions only if ISLAND_MODE is enabled by the active HUDF profile.
- In ISLAND_MODE, the runtime MUST:
  - pin to a specific Trust Registry snapshot hash (Section 11.4)
  - record `island_epoch_id` in any EFFECT receipt emitted during the partition
  - deny any effect that introduces new trust roots or new counterparties with reason_code = ERR_ISLAND_MODE_RESTRICTED

Reconciliation:
- After connectivity is restored, branches MUST be merged using MERGE_DECISION rules.
- Financial and irreversible effects MUST rely on connector idempotency and reconciliation, not CRDT-style optimistic merges.
- Conflicts that cannot be auto-compensated MUST halt with ERR_MERGE_CONFLICT_UNRESOLVED until a signed resolution is provided.

## 14. Conformance levels and harness

### 14.1 L1 Deterministic offline-first core
Requirements:
- VGL supported.
- VPL enforced.
- HUDF profile parameters supported and bound into genesis attestations (Section 2.4).
- CPI Tier 1 and Tier 2 (Tier 2.5 and Tier 3 are profile-gated).
- Tier 2.5 HITL approval ingestion and ATTESTATION binding when REQUIRE_APPROVAL is returned.
- REASONING_LOOP condensation root propagation (reasoning_merkle_root included in subsequent INTENT evidence_refs).
- Connector drift fail-closed behavior (ERR_CONNECTOR_CONTRACT_DRIFT) and shim authorization replay path (10.5.1).
- ProofGraph plus EvidencePack deterministic packaging.
- Tool sandbox minimums, including dynamic code isolation (10.1.1).

### 14.2 L2 Replication and merge-auditable ProofGraph
Adds:
- Replication protocol implementation.
- MERGE_DECISION nodes and deterministic tie-breakers.
- Partition reconciliation tests.

### 14.3 L3 Multi-tenant and ecosystem interop
Adds:
- Tenant isolation boundaries.
- A2A Envelope verification.
- Cryptographic metering receipts.
- Optional profile-3 attestations.

### 14.4 Conformance harness contract
A conformant implementation MUST ship a CLI that supports:
- helm verify proofgraph <path>
- helm verify pack <path>
- helm replay <evidencepack>
- helm diff state <a> <b>
- helm test conformance --level L1|L2|L3

Harness outputs MUST be deterministic JSON with a stable error taxonomy (Appendix C).

## 15. Reference implementation mapping (HELM repo reality)
This section binds the standard to the current HELM reference implementation shape (not required for third-party implementations, but required for HELM itself).

- Kernel entrypoints: core/cmd/helm, core/cmd/helm-migrate, core/cmd/fas-score, core/cmd/fas-verify, core/cmd/helm-loadtest
- Console: served by kernel on port 8080; UI static at /app/ui; health on 8081
- Node daemon: apps/helm-node/main.go on port 9090
- UI: apps/control-room-ui (React plus Vite) consuming kernel console API
- Types: generated from JSON Schema via go generate ./... into apps/control-room-ui/src/types

Any divergence MUST be documented in ProofGraph as a POLICY or ATTESTATION node and gated by conformance tests.

## 16. 2030 UX targets (non-normative, product north star)
These are UX surfaces HELM is architected to support as models and tooling evolve.

### 16.1 B2C/SMB UX
- Outcome Stack (portfolio of outcomes with budgets, risk posture, next irreversible action)
- Decision Inbox (single interrupt surface for approvals)
- Receipt Timeline (exportable proof of what happened and why)
- Pack Library (installable outcome playbooks with verifiable versioning)
- Mode controls (safe, throttle, freeze)

### 16.2 B2B/Enterprise UX
- Org Mission Control (live map of autonomous workflows and blast radius)
- ProofGraph Explorer (timeline plus replay)
- Policy Studio with simulation (policy diff and historical what-if)
- Audit Portal (EvidencePack exports and deterministic verification)
- Vendor mesh (capability grants, A2A envelopes, procurement metadata as optional)

## 17. Optional modules
The following are optional and MUST NOT be required for L1 or L2 conformance:
- Governance sensors based on activation telemetry (research-grade).
- Hardware attestation that claims to bind specific model weights in memory (Profile 3 only).
- Zero-knowledge proofs for usage receipts (optional extension to condensation).

Any optional module MUST:
- Have a deterministic wire contract.
- Be domain-separated for signatures.
- Be gated by policy and conformance tests when enabled.

## References
[1] RFC 2119: Key words for use in RFCs to Indicate Requirement Levels - IETF - 1997 - https://www.rfc-editor.org/rfc/rfc2119
[2] RFC 8174: Ambiguity of Uppercase vs Lowercase in RFC 2119 Key Words - IETF - 2017 - https://www.rfc-editor.org/rfc/rfc8174
[3] RFC 8785: JSON Canonicalization Scheme (JCS) - IETF - 2020 - https://www.rfc-editor.org/rfc/rfc8785
[4] RFC 7515: JSON Web Signature (JWS) - IETF - 2015 - https://www.rfc-editor.org/rfc/rfc7515
[5] RFC 8037: CFRG Elliptic Curve Diffie-Hellman (ECDH) and Signatures in JOSE - IETF - 2017 - https://www.rfc-editor.org/rfc/rfc8037
[6] Trusted Execution Environments - AWS Nitro Enclaves User Guide - AWS - n.d. - https://docs.aws.amazon.com/nitro-enclaves/latest/user/nitro-enclaves.html
[7] Intel Trust Domain Extensions (Intel TDX) - Technical overview - Intel - n.d. - https://www.intel.com/content/www/us/en/developer/articles/technical/intel-trust-domain-extensions.html
[8] Web Authentication: An API for accessing Public Key Credentials Level 2 - W3C - n.d. - https://www.w3.org/TR/webauthn-2/
[9] WebAssembly System Interface (WASI) - Project and specifications - WebAssembly - n.d. - https://github.com/WebAssembly/WASI
[10] Zero Trust Architecture - NIST - 2020 - https://csrc.nist.gov/publications/detail/sp/800-207/final
[11] Secret Sharing (reference overview) - MIT - 2019 - https://web.mit.edu/6.857/OldStuff/Fall03/ref/Shamir/secret_sharing.pdf
[12] A Review of zk-SNARKs (arXiv:2202.06877) - arXiv - 2022 (v4 2023) - https://arxiv.org/abs/2202.06877
[13] SoK: Understanding zk-SNARKs: The Gap Between Research and Practice (arXiv:2502.02387) - arXiv / USENIX Security 2025 - 2025 - https://arxiv.org/abs/2502.02387
[14] Zero-Knowledge Proof Frameworks: A Survey (arXiv:2502.07063) - arXiv - 2025 - https://arxiv.org/abs/2502.07063
[15] NIST - FIPS 203: Module-Lattice-Based Key-Encapsulation Mechanism Standard (ML-KEM) - 2024-08-13 - https://csrc.nist.gov/pubs/fips/203/final
[16] NIST - FIPS 204: Module-Lattice-Based Digital Signature Standard (ML-DSA) - https://csrc.nist.gov/pubs/fips/204/final
[17] NIST - FIPS 205: Stateless Hash-Based Digital Signature Standard (SLH-DSA) - https://csrc.nist.gov/pubs/fips/205/final
[18] European Union - Regulation (EU) 2024/1183 (amending eIDAS, European Digital Identity Framework) - Official Journal - https://eur-lex.europa.eu/eli/reg/2024/1183/oj/eng
[19] C2PA - Content Provenance and Authenticity Specification - https://c2pa.org/specifications/specifications/1.4/specs/C2PA_Specification.html
[20] TLSNotary (TLS session origin proofs) - https://github.com/tlsnotary/tlsn

# Appendix A - Bytes Contract v1.1 (Deterministic hashing, signing, packaging)

## A.1 Canonical bytes

**Signature preimage rule (normative):** For any object that includes detached signatures (for example `sig`, `sigset`, `silicon_quote`, or other attestation fields), the canonical bytes for verification MUST be computed over the object with all signature and attestation fields removed. This prevents self-referential signatures and allows multi-signature sets. The removed fields MUST NOT influence canonicalization or hashing.

All canonical objects in scope (ProofGraph nodes, Plan IR, Policy Bundles, EvidencePack manifests, A2A envelopes) MUST be serialized as:
- UTF-8
- LF newlines for any text artifacts
- No BOM
- JSON MUST use RFC 8785 canonicalization

## A.2 Signature domain separation (mandatory)
Before signing, implementations MUST prepend a context string:
- HELM/ProofGraphNode/v1
- HELM/EvidencePack/v1
- HELM/A2AEnvelope/v1
- HELM/PolicyBundle/v1
- HELM/GenesisApproval/v1

The signed bytes are: context + newline + canonical object bytes.

## A.3 EvidencePack deterministic archive
EvidencePack MUST be a deterministic .tar.gz with:
- Tar format: POSIX pax
- Sorted file paths lexicographically (byte order)
- Normalized paths: forward slashes, no leading ./, no .., no absolute paths
- File metadata normalized: uid=0, gid=0, uname="", gname="", mtime=0
- Directory entries: MUST be included for all directories with normalized metadata
- Symlinks: MUST NOT be present (fail-closed)
- gzip determinism: no original filename, no timestamp (equivalent to gzip -n)

## A.4 EvidencePack manifest
Each EvidencePack MUST contain MANIFEST.json with:
- pack_id
- created_lamport
- org_genome_hash
- range (ProofGraph slice)
- objects (array of path plus sha256)
- signatures (array of detached signatures)

MANIFEST.json MUST be canonical JSON.

## A.5 Document hash computation (this document)
To compute Canonical-Doc-Hash:
- Replace the 64-hex value in Canonical-Doc-Hash with 64 zeros.
- Compute sha256 over the UTF-8 bytes of the full Markdown document.
- Set Canonical-Doc-Hash to the resulting digest.

# Appendix B - Minimal JSON schemas (executable contracts)
This appendix defines minimal schema shapes. Implementations MAY extend fields but MUST NOT change meanings or required fields.

## B.1 ProofGraphNode schema (minimal)
```json
{
  "type": "object",
  "required": ["node_hash","kind","parents","lamport","principal","principal_seq","payload","sig"],
  "properties": {
    "node_hash": {"type":"string"},
    "kind": {"type":"string","enum":["INTENT","POLICY","OBSERVATION","ATTESTATION","TRUST_EVENT","EFFECT","CHECKPOINT","MERGE_DECISION"]},
    "parents": {"type":"array","items":{"type":"string"}},
    "lamport": {"type":"integer","minimum":0},
    "principal": {"type":"string"},
    "principal_seq": {"type":"integer","minimum":0},
    "payload": {"type":"object"},
    "sig": {"type":"string"}
  }
}
```

### B.1.1 Sensitive payload extension (normative if payload_class=Sensitive)
If a node payload is classified as Sensitive (Section 3.4), implementations MUST support these fields inside the node payload:
- payload_class: "Sensitive"
- ciphertext_hash: "sha256:<hex>" of the encrypted blob bytes
- blob_ref: opaque storage reference
- kek_ref: tenant key reference
- key_policy_hash: "sha256:<hex>"
- redaction_policy_ref: reference to OrgGenome redaction policy

Nodes MUST NOT store cleartext sensitive payload fields when payload_class=Sensitive.

## B.2 A2AEnvelope schema (minimal)
```json
{
  "type":"object",
  "required":["v","msg_id","iss","sub","aud","nonce","ttl_ms","issued_lamport","capabilities","policy_hash","risk_class","payload","sig"],
  "properties":{
    "v":{"type":"string","const":"1"},
    "msg_id":{"type":"string"},
    "iss":{"type":"string"},
    "sub":{"type":"string"},
    "aud":{"type":"string"},
    "nonce":{"type":"string"},
    "ttl_ms":{"type":"integer","minimum":1},
    "issued_lamport":{"type":"integer","minimum":0},
    "capabilities":{"type":"array","items":{"type":"object"}},
    "policy_hash":{"type":"string"},
    "risk_class":{"type":"string","enum":["T0","T1","T2","T3"]},
    "payload":{"type":"object"},
    "sig":{"type":"string"}
  }
}
```

# Appendix C - Reason code registry (core)

Reason codes are stable, machine-consumable strings.

This appendix defines the baseline reason codes referenced by this standard. Implementations MAY add additional reason codes, but MUST NOT change the meaning of codes defined here.

## C.1 Core execution and verification
- OK
- ERR_COMPUTE_GAS_EXHAUSTED
- ERR_COMPUTE_TIME_EXHAUSTED
- ERR_STALE_SEQUENCE
- ERR_TTL_EXPIRED
- ERR_SIGNATURE_INVALID
- ERR_TRUST_ROOT_UNKNOWN
- ERR_TRUST_ROOT_MISSING
- ERR_VERSION_NEGOTIATION_FAILED

## C.2 Policy and budgets
- ERR_POLICY_DENY
- ERR_VELOCITY_CEILING_BREACH
- ERR_POLICY_BUDGET
- ERR_P0_CEILING_BREACH
- ERR_BUDGET_EXCEEDED

## C.3 CPI and proof
- ERR_CPI_STRUCTURAL
- ERR_CPI_TEMPORAL
- ERR_CPI_UNSAT_CORE
- ERR_CPI_SOLVER_TIMEOUT

## C.4 Governance approvals
- ERR_APPROVAL_CEREMONY_REQUIRED
- ERR_APPROVAL_TIMELOCK_ACTIVE
- ERR_APPROVAL_RATE_LIMIT_EXCEEDED
- ERR_APPROVAL_EMERGENCY_OVERRIDE_REQUIRED

## C.5 Sandbox and connectors
- ERR_SANDBOX_POLICY_VIOLATION
- ERR_SANDBOX_SYSCALL_VIOLATION
- ERR_CONNECTOR_CONTRACT_DRIFT
- ERR_CONNECTOR_IDEMPOTENCY_MISSING
- ERR_ANALOG_EXECUTOR_REQUIRES_APPROVAL
- ERR_ISLAND_MODE_RESTRICTED

## C.6 Evidence and replication
- ERR_EVIDENCEPACK_INVALID
- ERR_MERGE_CONFLICT_UNRESOLVED

## C.7 Identity lifecycle
- ERR_HEARTBEAT_EXPIRED_FREEZE

## C.8 Zero-knowledge extensions
- ERR_ZK_VERIFY_FAILED
- ERR_ZK_UNSUPPORTED_STATEMENT
- ERR_ZK_VERIFIER_BUDGET

# Appendix D - Conformance harness core suite (minimal)
A conformant implementation MUST include a test suite with:
- Golden pack verification (deterministic archive reproducibility)
- ProofGraph replay determinism
- CPI tier tests (structural, bounded temporal, solver path)
- Policy intersection and deny precedence tests
- Sandbox enforcement tests
- Compute budget exhaustion tests (ERR_COMPUTE_GAS_EXHAUSTED, ERR_COMPUTE_TIME_EXHAUSTED)
- Approval ceremony tests (ERR_APPROVAL_CEREMONY_REQUIRED)
- Approval timelock tests (ERR_APPROVAL_TIMELOCK_ACTIVE)
- Approval rate limit tests (ERR_APPROVAL_RATE_LIMIT_EXCEEDED)
- Island mode restriction tests (ERR_ISLAND_MODE_RESTRICTED)
- Replication merge tests (L2)
- A2A envelope verification tests (L3)
- Metering dispute replay tests (L3)

---

# Appendix E - 2030 Readiness Extensions (Normative)

This appendix defines optional, forward-compatible extensions that harden HELM against cross-jurisdiction trade, adversarial counterparties, cyber-physical actuation, compute exhaustion, and cryptographic transitions. These extensions MUST NOT weaken any invariants in the Truth Plane or ProofGraph.

Implementations MAY claim conformance to any subset of these extensions. Policies MAY require specific extensions for specific intent classes.

## E.1 Privacy-preserving A2A verification (ZK-CPI)

### E.1.1 Prover-pays rule (mandatory when enabled)
To avoid "verification by data leakage" and to invert the compute burden in A2A commerce:
- The party constructing an A2A envelope (the prover) MUST bear the heavy computation required to satisfy the counterparty's constraints.
- The receiving party (the verifier) MUST perform verification in bounded, deterministic time.

### E.1.2 Statement registry (mandatory when enabled)
ZK verification MUST NOT accept arbitrary statements. Implementations MUST define a registry of supported statements, each with:
- statement_id (string, stable)
- public_inputs schema
- witness schema
- verification_key_hash (sha256 of VK bytes)
- circuit_hash (sha256 of circuit IR)

The registry MUST be referenced by hash from the Trust Registry.

### E.1.3 ZK proof bundle schema (normative)
When ZK-CPI is used, the A2A envelope MUST include:

```json
{
  "zk_proofs": [
    {
      "statement_id": "rev_90d_gte_usd_1000000",
      "public_inputs": {},
      "proof_b64": "...",
      "verification_key_hash": "sha256:...",
      "prover_cost_receipt": {
        "metering_root": "sha256:...",
        "max_cost_usd": 5.00
      }
    }
  ]
}
```

### E.1.4 Fail-closed behavior
If ZK verification fails, is unsupported, or exceeds verifier budget, CPI MUST return:
- verdict: DENY or DEFER
- reason_code: ERR_ZK_VERIFY_FAILED, ERR_ZK_UNSUPPORTED_STATEMENT, or ERR_ZK_VERIFIER_BUDGET

No downgrade is permitted unless policy explicitly enables it.

## E.2 Jurisdictional selective disclosure (Regulator View-Key)

Pure privacy-preserving trade can conflict with AML and data-locality requirements. HELM resolves this by supporting selective disclosure to regulators without disclosing raw ledgers to counterparties.

### E.2.1 Disclosure profiles (normative)
A2A envelopes MAY declare:

- disclosure_profile: PRIVATE_ZK
  - counterparties receive only ZK proofs and public fields.
- disclosure_profile: SELECTIVE_REGULATOR_VIEW
  - counterparties receive only ZK proofs and public fields.
  - regulators can decrypt specific disclosures under lawful process.
- disclosure_profile: FULL_DISCLOSURE
  - counterparties receive explicit receipts and/or redacted evidence.

### E.2.2 Regulator view key bundle (normative)
When disclosure_profile is SELECTIVE_REGULATOR_VIEW, the envelope MUST include:

```json
{
  "regulator_view_keys": [
    {
      "jurisdiction": "EU",
      "recipient_key_id": "did:...#key-1",
      "scope": ["KYC_ID", "AML_TRACE", "BENEFICIAL_OWNER"],
      "ciphertext_b64": "..."
    }
  ]
}
```

The ciphertext MUST be generated by the sender and MUST NOT require the receiver to possess any secret for regulator decryption.

## E.3 Multi-sovereign attestation quorums

### E.3.1 Attestation evidence set (normative)
When hardware attestation is required by policy, KERNEL_VERDICT MUST include either:
- a single silicon_quote meeting the receiver's Trust Registry policy, or
- a silicon_quote_set with N-of-M quorum semantics.

```json
{
  "silicon_quote_set": {
    "quorum": { "n": 2, "m": 3 },
    "quotes": [
      { "tee_vendor": "aws_nitro", "quote_b64": "...", "measurement": "sha256:..." },
      { "tee_vendor": "intel_tdx", "quote_b64": "...", "measurement": "sha256:..." },
      { "tee_vendor": "amd_sev_snp", "quote_b64": "...", "measurement": "sha256:..." }
    ]
  }
}
```

### E.3.2 Trust Registry attestation policy hook (normative)
The Trust Registry MUST support an attestation_policy object that can require:
- allowed tee_vendor values
- optional quorum (n, m)
- allowed measurements (kernel builds) by hash
- allowed geos (optional)

## E.4 Cryptographic perception (Oracle hardening)

Deterministic execution of poisoned inputs is still failure. HELM hardens perception by requiring provenance for any external data used to justify Tier 2 or Tier 3 actions.

### E.4.1 Provenance-gated promotion (mandatory when enabled)
If a plan step reads external data that is later referenced by a Tier 2 or Tier 3 intent, the runtime MUST promote that data to an OBSERVATION node that includes a provenance attestation.

### E.4.2 Provenance attestation forms (normative)
OBSERVATION provenance MUST include at least one of:
- content credential proof (for signed media and documents)
- TLS session proof (for web origin claims)
- hardware-signed sensor quote (for IoT and cyber-physical telemetry)

Provenance objects MUST be hash-addressed and referenced from the OBSERVATION payload.

### E.4.3 Policy hook
MAPL MUST support a predicate that requires provenance class and minimum quorum, for example:
- require_provenance(class="content_credentials", quorum=1)
- require_provenance(class="sensor_quote", quorum=2)

## E.5 Cyber-physical actuation profile (Plane 0)

HELM cannot "rollback physics". For cyber-physical systems, HELM MUST control actuation through a safe, deterministic boundary.

### E.5.1 Actuator enclave boundary (normative)
For any connector that triggers kinetic effects, the effect MUST be routed to a dedicated Actuator Enclave that:
- accepts only signed HELM effect envelopes
- enforces an allowlist of actuators and bounded parameters
- supports a safe-state command that is always allowed

### E.5.2 Digital twin pre-simulation (recommended)
Policies MAY require that kinetic intents include a digital_twin_proof_ref, produced by a deterministic simulation pipeline bound to the current policy epoch.

### E.5.3 Kinetic safe state (mandatory)
If CPI returns DENY, DEFER, or any runtime fault occurs while a kinetic workflow is in-flight, the Actuator Enclave MUST transition to a safe state, such as:
- hover and auto-land
- retract and lock
- valve closed

The transition MUST emit an EFFECT node with tool_urn = "helm.actuator.safe_state" and payload.effect_type = "ACTUATION_SAFE_STATE", and it MUST be included in the EvidencePack for the event.

## E.6 Energy and cost governance

Large-scale verification introduces economic and thermodynamic constraints.

### E.6.1 Budget ceilings (normative)
MAPL MUST support optional ceilings for:
- max_verification_cost_usd_per_day
- max_solver_seconds_per_day
- max_energy_kwh_per_day (optional when available)

If a ceiling would be exceeded, CPI MUST return DEFER with ERR_BUDGET_EXCEEDED.

### E.6.2 Carbon-aware scheduling (optional)
If enabled, the scheduler MAY route proof generation to approved compute zones, provided:
- the resulting artifacts are hash-identical
- the execution remains deterministic
- location and metering receipts are recorded in the condensation root

## E.7 Post-quantum cryptographic agility

### E.7.1 Algorithm registry (mandatory)
The standard MUST maintain a registry of approved signature algorithms. Implementations MUST support:
- a classical algorithm for baseline interoperability
- at least one post-quantum algorithm for future safety
- hybrid mode where both signatures are present

### E.7.2 Signature set field (backwards-compatible)
ProofGraph nodes and A2A envelopes MAY include an additional detached signature set:

```json
{
  "sig": "jws:... (baseline)",
  "sigset": [
    { "alg": "Ed25519", "kid": "did:...#k1", "sig_b64": "..." },
    { "alg": "ML-DSA-65", "kid": "did:...#k2", "sig_b64": "..." }
  ]
}
```

Verifiers MUST compute signature preimages over the canonical object with sig and sigset removed.

### E.7.3 PQ epoch sealing (optional)
Organizations MAY define a QUANTUM_EPOCH_SEAL checkpoint that:
- records the ProofGraph root at a given epoch
- is signed with a post-quantum algorithm
- is included in periodic EvidencePacks

Receivers MAY require a recent epoch seal for high-stakes A2A interactions.
# Appendix F - Attestation type registry (core)

Attestation types are stable, machine-consumable strings.

This appendix defines the baseline attestation types referenced by this standard. Implementations MAY add additional attestation types, but MUST NOT change the meaning of types defined here.

## F.1 Kernel and verification
- KERNEL_VERDICT
  - Produced by the kernel after CPI evaluation.
  - MUST bind: intent_hash, verdict, tier, policy_hash, org_phenotype_hash, budgets, and reason_code on deny/defer.
- MODEL_VERSION_ATTESTATION
  - Binds a cognitive engine identifier (model_id, weights hash when available, runtime build id) to a pinned policy epoch.
- CONDENSATION_CHECKPOINT
  - Proves that a condensation operation preserved verifiable reachability for required nodes and produced a new checkpoint hash.

## F.2 Governance and approvals
- ORG_GENESIS_APPROVAL
  - Binds vendor-signed P1 policy bundle and tenant-signed P0 ceilings into a single genesis approval.
- P0_CEILINGS_ACTIVE
  - A tenant principal signature over the exact scalar P0 ceilings that define sovereign blast radius.
- HITL_SIGN_REQUEST
  - Kernel-issued signing request for Tier 2.5 approvals. MUST include intent_hash and required parameters.
- HITL_APPROVAL
  - Human approval signature (sig or sigset) referencing a HITL_SIGN_REQUEST and binding to intent_hash.
- APPROVAL_CEREMONY_V1
  - Proof that a deliberate approval ceremony occurred (RFC-005). MUST bind to intent_hash and ui_summary_hash and include timelock and challenge-response hashes.
- EMERGENCY_OVERRIDE
  - Recorded bypass of a required control (timelock, ceremony, or solver). MUST bind to intent_hash and include bounded emergency_reason.

## F.3 Identity lifecycle
- HEARTBEAT_ATTESTATION
  - Periodic liveness proof for the sovereign principal or configured quorum.
- RESURRECTION_VETO
  - A veto produced by the prior sovereign key during an identity resurrection timelock.
- IDENTITY_ROTATION
  - Receipt that the root identity was rotated to a new key or quorum.

## F.4 Compliance and operations
- LEGAL_HOLD_APPLIED
  - Receipt that a legal hold policy was applied to evidence retention and condensation behavior.
- USAGE_NODE
  - Metering record used for billing and dispute replay. MUST bind to a metering period, units, and evidence references.