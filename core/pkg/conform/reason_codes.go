package conform

// Reason codes are stable identifiers per §6.2.
// They MUST NOT change between releases.
const (
	// --- Build & Identity ---
	ReasonBuildIdentityMissing = "BUILD_IDENTITY_MISSING"
	ReasonTrustRootsMissing    = "TRUST_ROOTS_MISSING" // G0: signing keys not verifiable

	// --- Receipt chain ---
	ReasonReceiptChainBroken        = "RECEIPT_CHAIN_BROKEN"
	ReasonReceiptDAGBroken          = "RECEIPT_DAG_BROKEN" // DAG parent hash unresolvable
	ReasonSignatureInvalid          = "SIGNATURE_INVALID"
	ReasonPayloadCommitmentMismatch = "PAYLOAD_COMMITMENT_MISMATCH"

	// --- Receipt emission ---
	ReasonReceiptEmissionPanic = "RECEIPT_EMISSION_PANIC" // kernel panic: receipts cannot be emitted

	// --- Replay ---
	ReasonReplayHashDivergence = "REPLAY_HASH_DIVERGENCE"
	ReasonReplayTapeMiss       = "REPLAY_TAPE_MISS"

	// --- Tape ---
	ReasonTapeResidencyViolation = "TAPE_RESIDENCY_VIOLATION" // taped payload violates jurisdiction/data handling

	// --- Policy ---
	ReasonPolicyDecisionMissing  = "POLICY_DECISION_MISSING"
	ReasonSchemaValidationFailed = "SCHEMA_VALIDATION_FAILED"

	// --- Budget ---
	ReasonBudgetExhausted = "BUDGET_EXHAUSTED"

	// --- Containment ---
	ReasonContainmentNotTriggered = "CONTAINMENT_NOT_TRIGGERED"

	// --- A2A ---
	ReasonA2AProofMissing    = "A2A_PROOF_MISSING"
	ReasonA2APassportInvalid = "A2A_PASSPORT_INVALID"

	// --- Taint ---
	ReasonTaintFlowViolation = "TAINT_FLOW_VIOLATION"

	// --- Formal ---
	ReasonFormalExportInvalid = "FORMAL_EXPORT_INVALID"

	// --- Tenant Isolation ---
	ReasonTenantIsolationViolation = "TENANT_ISOLATION_VIOLATION" // cross-tenant access detected
	ReasonTenantIDMissing          = "TENANT_ID_MISSING"          // receipt/evidence lacks tenant_id

	// --- Envelope Binding ---
	ReasonEnvelopeNotBound        = "ENVELOPE_NOT_BOUND"         // effect without active envelope
	ReasonEnvelopeNotEnforced     = "ENVELOPE_NOT_ENFORCED"      // envelope constraints not checked
	ReasonEnvelopeDenialNoReceipt = "ENVELOPE_DENIAL_NO_RECEIPT" // denial without receipt

	// --- Proxy Governance ---
	ReasonProxyToolAllowed     = "PROXY_TOOL_ALLOWED"     // tool call passed governance
	ReasonProxyToolDenied      = "PROXY_TOOL_DENIED"      // tool call failed governance
	ReasonProxyBudgetExhausted = "PROXY_BUDGET_EXHAUSTED" // budget limit hit via proxy
	ReasonProxyIterationLimit  = "PROXY_ITERATION_LIMIT"  // max iterations reached
	ReasonProxyWallclockLimit  = "PROXY_WALLCLOCK_LIMIT"  // session wallclock exceeded
)

// AllReasonCodes returns the full set of normative reason codes.
func AllReasonCodes() []string {
	return []string{
		ReasonBuildIdentityMissing,
		ReasonTrustRootsMissing,
		ReasonReceiptChainBroken,
		ReasonReceiptDAGBroken,
		ReasonSignatureInvalid,
		ReasonPayloadCommitmentMismatch,
		ReasonReceiptEmissionPanic,
		ReasonReplayHashDivergence,
		ReasonReplayTapeMiss,
		ReasonTapeResidencyViolation,
		ReasonPolicyDecisionMissing,
		ReasonSchemaValidationFailed,
		ReasonBudgetExhausted,
		ReasonContainmentNotTriggered,
		ReasonA2AProofMissing,
		ReasonA2APassportInvalid,
		ReasonTaintFlowViolation,
		ReasonFormalExportInvalid,
		ReasonTenantIsolationViolation,
		ReasonTenantIDMissing,
		ReasonEnvelopeNotBound,
		ReasonEnvelopeNotEnforced,
		ReasonEnvelopeDenialNoReceipt,
		ReasonProxyToolAllowed,
		ReasonProxyToolDenied,
		ReasonProxyBudgetExhausted,
		ReasonProxyIterationLimit,
		ReasonProxyWallclockLimit,
	}
}
