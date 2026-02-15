package contracts

import "time"

// OrgGenome represents the organization genome structure.
//
//nolint:govet // fieldalignment: struct layout is human-readable
type OrgGenome struct {
	Phenotype         PhenotypeMeta       `json:"phenotype"`
	Graph             OrgGraph            `json:"graph"`
	PhenotypeContract PhenotypeContract   `json:"phenotype_contract"` // Struct
	Morphogenesis     []MorphogenesisRule `json:"morphogenesis"`

	// Legacy / Boot Compatibility
	Meta        GenomeMeta         `json:"meta"`
	Environment EnvironmentProfile `json:"environment"`
	Regulation  map[string]any     `json:"regulation"`
}

// PhenotypeContract defines the phenotype contract.
type PhenotypeContract struct {
	Determinism Determinism `json:"determinism"`
	MustProduce []string    `json:"must_produce"` // Added for tests
}

// DeterminismContract is an alias for Determinism.
type DeterminismContract struct { // Alias for Determinism
	RequiresRandomSeed bool `json:"requires_random_seed"`
}

// Determinism represents determinism requirements.
type Determinism struct {
	RequiresRandomSeed bool `json:"requires_random_seed"`
}

// GenomeMeta contains genome metadata.
type GenomeMeta struct {
	GenomeID  string    `json:"genome_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// MorphogenesisRule defines a morphogenesis rule.
type MorphogenesisRule struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Condition string         `json:"condition"`
	Effect    string         `json:"effect"`
	When      string         `json:"when"`     // Alias
	Generate  map[string]any `json:"generate"` // Field, not method
}

// OrgGraph represents the structural view of the organization.
type OrgGraph struct {
	Nodes []OrgNode `json:"nodes"`
	Edges []OrgEdge `json:"edges"`
}

// OrgNode represents a node in the organization graph.
type OrgNode struct {
	ID         string            `json:"id"`
	NodeID     string            `json:"node_id"` // Alias
	Type       string            `json:"type"`
	NodeType   string            `json:"node_type"` // Alias
	Meta       map[string]string `json:"meta"`
	Properties map[string]any    `json:"properties"`
}

// OrgEdge represents an edge in the organization graph.
type OrgEdge struct {
	From       string         `json:"from"`
	To         string         `json:"to"`
	Rel        string         `json:"rel"`
	EdgeID     string         `json:"edge_id"`
	RelType    string         `json:"rel_type"` // Alias
	Properties map[string]any `json:"properties"`
}
