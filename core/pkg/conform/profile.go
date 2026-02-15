package conform

// ProfileID identifies a conformance profile per §2.2.
type ProfileID string

const (
	ProfileSMB              ProfileID = "SMB"
	ProfileCore             ProfileID = "CORE"
	ProfileEnterprise       ProfileID = "ENTERPRISE"
	ProfileRegulatedFinance ProfileID = "REGULATED_FINANCE"
	ProfileRegulatedHealth  ProfileID = "REGULATED_HEALTH"
	ProfileAgenticWebRouter ProfileID = "AGENTIC_WEB_ROUTER"
)

// ProfileDefinition describes which gates a profile requires
// and any per-gate threshold overrides.
type ProfileDefinition struct {
	ID            ProfileID      `json:"id"`
	Description   string         `json:"description"`
	RequiredGates []string       `json:"required_gates"`
	Inherits      ProfileID      `json:"inherits,omitempty"`
	Overrides     map[string]any `json:"overrides,omitempty"`
}

// Profiles returns the built-in profile definitions per §9.
//
// Compliance is profile-scoped: a runtime is compliant for a declared
// profile iff it passes all mandatory gates for that profile and produces
// a signed conformance report.
func Profiles() map[ProfileID]*ProfileDefinition {
	return map[ProfileID]*ProfileDefinition{
		ProfileSMB: {
			ID:          ProfileSMB,
			Description: "Minimum autonomy runtime safety for startups and small businesses",
			RequiredGates: []string{
				"G0", "G1", "G2", "G3", "G3A",
				"G5", "G7", "G8", "GX_ENVELOPE",
			},
		},
		ProfileCore: {
			ID:          ProfileCore,
			Description: "Baseline enterprise autonomy — minimum 'autonomy runtime safety' bar",
			RequiredGates: []string{
				"G0", "G1", "G2", "G2A", "G3", "G3A",
				"G5", "G6", "G7", "G8", "G12",
				"GX_ENVELOPE",
			},
		},
		ProfileEnterprise: {
			ID:          ProfileEnterprise,
			Description: "CORE plus operability, identity hardening, and tenant isolation",
			Inherits:    ProfileCore,
			RequiredGates: []string{
				"G0", "G1", "G2", "G2A", "G3", "G3A",
				"G4", "G5", "G6", "G7", "G8", "G9", "G11", "G12",
				"GX_ENVELOPE", "GX_TENANT",
			},
		},
		ProfileRegulatedFinance: {
			ID:          ProfileRegulatedFinance,
			Description: "ENTERPRISE plus strict replay, audit, and formal verification",
			Inherits:    ProfileEnterprise,
			RequiredGates: []string{
				"G0", "G1", "G2", "G2A", "G3", "G3A",
				"G4", "G5", "G5A", "G6", "G7", "G8", "G9", "G10", "G11", "G12",
				"GX_ENVELOPE", "GX_TENANT",
			},
			Overrides: map[string]any{
				"replay_bit_identical":   true,
				"schema_first_hard_fail": true,
				"a2a_proof_required":     true,
			},
		},
		ProfileRegulatedHealth: {
			ID:          ProfileRegulatedHealth,
			Description: "ENTERPRISE plus strict privacy, retention, and tape residency",
			Inherits:    ProfileEnterprise,
			RequiredGates: []string{
				"G0", "G1", "G2", "G2A", "G3", "G3A",
				"G4", "G5", "G6", "G7", "G8", "G9", "G10", "G11", "G12",
				"GX_ENVELOPE", "GX_TENANT",
			},
			Overrides: map[string]any{
				"privacy_erasure_strict":    true,
				"retention_policy_required": true,
				"tape_residency_enforced":   true,
			},
		},
		ProfileAgenticWebRouter: {
			ID:          ProfileAgenticWebRouter,
			Description: "ENTERPRISE plus A2A proof routing and semantic quarantine",
			Inherits:    ProfileEnterprise,
			RequiredGates: []string{
				"G0", "G1", "G2", "G2A", "G3", "G3A",
				"G4", "G5", "G5A", "G6", "G7", "G8", "G9", "G11", "G12",
				"GX_ENVELOPE", "GX_TENANT",
			},
			Overrides: map[string]any{
				"proof_capsule_routing": true,
			},
		},
	}
}

// GatesForProfile returns the gate IDs required by a profile.
func GatesForProfile(id ProfileID) []string {
	p := Profiles()
	def, ok := p[id]
	if !ok {
		return nil
	}
	return def.RequiredGates
}
