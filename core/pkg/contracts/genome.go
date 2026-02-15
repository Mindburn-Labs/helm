package contracts

import "time"

// OrgPhenotype represents the runtime configuration of the organization.
// OrgPhenotype represents the runtime configuration of the organization.
type OrgPhenotype struct {
	Meta       PhenotypeMeta  `json:"meta"`
	Regulation map[string]any `json:"regulation"`

	// Extended Attributes
	OrgGraph  OrgGraph      `json:"graph"`
	Workflows []WorkflowDef `json:"workflows"`
	Policies  []PolicyRef   `json:"policies"`
	Resources []ResourceDef `json:"resources"`
}

// PhenotypeMeta contains metadata about the phenotype.
type PhenotypeMeta struct {
	GenomeID        string    `json:"genome_id"` // Added for compatibility
	PhenotypeID     string    `json:"phenotype_id"`
	CanonicalHash   string    `json:"canonical_hash"`
	SpecVersion     string    `json:"spec_version"`
	CompilerVersion string    `json:"compiler_version,omitempty"`
	Iterations      int       `json:"iterations,omitempty"`
	CompiledAt      time.Time `json:"compiled_at,omitempty"`
	BuildInfo       BuildInfo `json:"build_info,omitempty"`
}
