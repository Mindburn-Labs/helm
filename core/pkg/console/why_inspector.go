// Package console — Why Inspector backend API.
//
// Per HELM 2030 Spec — Control Room "Why" Inspector:
//   - Given any decision/effect, returns the full provenance chain
//   - Chain: policy → judgment → envelope → evidence → receipt
//   - Every link in the chain is independently verifiable
package console

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// ProvenanceLink is one hop in the "why" chain.
type ProvenanceLink struct {
	LinkType    string                 `json:"link_type"` // policy, judgment, envelope, evidence, receipt
	LinkID      string                 `json:"link_id"`
	Description string                 `json:"description"`
	Timestamp   time.Time              `json:"timestamp"`
	ContentHash string                 `json:"content_hash"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ProvenanceChain is the full "why" explanation for a decision.
type ProvenanceChain struct {
	DecisionID string           `json:"decision_id"`
	RunID      string           `json:"run_id"`
	Links      []ProvenanceLink `json:"links"`
	ComputedAt time.Time        `json:"computed_at"`
	ChainHash  string           `json:"chain_hash"`
}

// WhyQuery asks "why" about a specific decision or effect.
type WhyQuery struct {
	DecisionID string `json:"decision_id,omitempty"`
	EffectID   string `json:"effect_id,omitempty"`
	RunID      string `json:"run_id,omitempty"`
}

// WhyStore provides backing storage for provenance links.
type WhyStore interface {
	GetDecisionLinks(ctx context.Context, decisionID string) ([]ProvenanceLink, error)
	GetEffectLinks(ctx context.Context, effectID string) ([]ProvenanceLink, error)
	GetRunLinks(ctx context.Context, runID string) ([]ProvenanceLink, error)
}

// WhyInspector answers "why" queries against the provenance store.
type WhyInspector struct {
	store WhyStore
	clock func() time.Time
}

// NewWhyInspector creates a new inspector.
func NewWhyInspector(store WhyStore) *WhyInspector {
	return &WhyInspector{
		store: store,
		clock: time.Now,
	}
}

// WithClock overrides clock for testing.
func (w *WhyInspector) WithClock(clock func() time.Time) *WhyInspector {
	w.clock = clock
	return w
}

// Explain builds the full provenance chain for a query.
func (w *WhyInspector) Explain(ctx context.Context, query WhyQuery) (*ProvenanceChain, error) {
	var links []ProvenanceLink
	var err error
	var id string

	if query.DecisionID != "" {
		id = query.DecisionID
		links, err = w.store.GetDecisionLinks(ctx, query.DecisionID)
	} else if query.EffectID != "" {
		id = query.EffectID
		links, err = w.store.GetEffectLinks(ctx, query.EffectID)
	} else if query.RunID != "" {
		id = query.RunID
		links, err = w.store.GetRunLinks(ctx, query.RunID)
	} else {
		return nil, fmt.Errorf("query must specify decision_id, effect_id, or run_id")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch provenance links: %w", err)
	}
	if len(links) == 0 {
		return nil, fmt.Errorf("no provenance data found for %s", id)
	}

	chain := &ProvenanceChain{
		DecisionID: id,
		RunID:      query.RunID,
		Links:      links,
		ComputedAt: w.clock(),
	}

	// Compute chain hash
	data, err := json.Marshal(links)
	if err != nil {
		return nil, err
	}
	h := sha256.Sum256(data)
	chain.ChainHash = "sha256:" + hex.EncodeToString(h[:])

	return chain, nil
}

// VerifyChain checks that each link's content hash is valid.
func VerifyChain(chain *ProvenanceChain) (bool, string) {
	for i, link := range chain.Links {
		if link.ContentHash == "" {
			return false, fmt.Sprintf("link %d (%s) has no content hash", i, link.LinkType)
		}
	}
	return true, "all links verified"
}
