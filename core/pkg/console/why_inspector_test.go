package console

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// testWhyStore is an in-memory store for testing.
type testWhyStore struct {
	decisions map[string][]ProvenanceLink
	effects   map[string][]ProvenanceLink
	runs      map[string][]ProvenanceLink
}

func (s *testWhyStore) GetDecisionLinks(ctx context.Context, decisionID string) ([]ProvenanceLink, error) {
	links, ok := s.decisions[decisionID]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return links, nil
}
func (s *testWhyStore) GetEffectLinks(ctx context.Context, effectID string) ([]ProvenanceLink, error) {
	links, ok := s.effects[effectID]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return links, nil
}
func (s *testWhyStore) GetRunLinks(ctx context.Context, runID string) ([]ProvenanceLink, error) {
	links, ok := s.runs[runID]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return links, nil
}

func testLinks() []ProvenanceLink {
	now := time.Now()
	return []ProvenanceLink{
		{LinkType: "policy", LinkID: "pol-1", Description: "PDP evaluated CEL policy", Timestamp: now, ContentHash: "sha256:aaa"},
		{LinkType: "judgment", LinkID: "jdg-1", Description: "Classified as AUTONOMOUS", Timestamp: now, ContentHash: "sha256:bbb"},
		{LinkType: "envelope", LinkID: "env-1", Description: "Within envelope bounds", Timestamp: now, ContentHash: "sha256:ccc"},
		{LinkType: "evidence", LinkID: "evi-1", Description: "Evidence satisfied", Timestamp: now, ContentHash: "sha256:ddd"},
	}
}

func TestWhyInspectorDecision(t *testing.T) {
	store := &testWhyStore{decisions: map[string][]ProvenanceLink{"dec-1": testLinks()}}
	inspector := NewWhyInspector(store)

	chain, err := inspector.Explain(context.Background(), WhyQuery{DecisionID: "dec-1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(chain.Links) != 4 {
		t.Fatalf("expected 4 links, got %d", len(chain.Links))
	}
	if chain.ChainHash == "" {
		t.Fatal("expected chain hash")
	}
}

func TestWhyInspectorEffect(t *testing.T) {
	store := &testWhyStore{effects: map[string][]ProvenanceLink{"eff-1": testLinks()[:2]}}
	inspector := NewWhyInspector(store)

	chain, err := inspector.Explain(context.Background(), WhyQuery{EffectID: "eff-1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(chain.Links) != 2 {
		t.Fatalf("expected 2 links, got %d", len(chain.Links))
	}
}

func TestWhyInspectorNotFound(t *testing.T) {
	store := &testWhyStore{decisions: map[string][]ProvenanceLink{}}
	inspector := NewWhyInspector(store)

	_, err := inspector.Explain(context.Background(), WhyQuery{DecisionID: "nonexistent"})
	if err == nil {
		t.Fatal("expected error for missing decision")
	}
}

func TestWhyInspectorEmptyQuery(t *testing.T) {
	store := &testWhyStore{}
	inspector := NewWhyInspector(store)

	_, err := inspector.Explain(context.Background(), WhyQuery{})
	if err == nil {
		t.Fatal("expected error for empty query")
	}
}

func TestVerifyChainValid(t *testing.T) {
	chain := &ProvenanceChain{Links: testLinks()}
	ok, _ := VerifyChain(chain)
	if !ok {
		t.Fatal("expected valid chain")
	}
}

func TestVerifyChainMissingHash(t *testing.T) {
	links := testLinks()
	links[1].ContentHash = "" // Remove hash from one link
	chain := &ProvenanceChain{Links: links}
	ok, _ := VerifyChain(chain)
	if ok {
		t.Fatal("expected invalid chain with missing hash")
	}
}

// Portal API tests

func TestPortalRegisterPolicy(t *testing.T) {
	portal := NewPortalAPI()
	portal.RegisterPolicy(&PolicyRegistryEntry{PolicyID: "P1", PolicyName: "Deploy", Version: "1.0", Active: true})

	p, err := portal.GetPolicy("P1")
	if err != nil {
		t.Fatal(err)
	}
	if p.PolicyName != "Deploy" {
		t.Fatalf("expected Deploy, got %s", p.PolicyName)
	}
}

func TestPortalPolicyChangeLedger(t *testing.T) {
	portal := NewPortalAPI()
	portal.RegisterPolicy(&PolicyRegistryEntry{PolicyID: "P1", PolicyName: "Deploy", Version: "1.0", Active: true})
	portal.RegisterPolicy(&PolicyRegistryEntry{PolicyID: "P1", PolicyName: "Deploy", Version: "2.0", Active: true})

	changes := portal.ListChanges()
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].OldVersion != "1.0" || changes[0].NewVersion != "2.0" {
		t.Fatal("version mismatch in change")
	}
}

func TestPortalCompliancePosture(t *testing.T) {
	portal := NewPortalAPI()
	portal.SetCompliancePosture(&CompliancePosture{
		TenantID: "t1",
		Frameworks: map[string]FrameworkState{
			"gdpr": {FrameworkID: "gdpr", Name: "GDPR", Score: 85.0, ControlsMet: 17, ControlsTotal: 20},
		},
		OverallScore: 85.0,
	})

	p, err := portal.GetCompliancePosture("t1")
	if err != nil {
		t.Fatal(err)
	}
	if p.OverallScore != 85.0 {
		t.Fatalf("expected 85.0, got %.1f", p.OverallScore)
	}
}

func TestPortalExportEvidence(t *testing.T) {
	portal := NewPortalAPI()
	entries := []EvidenceExportEntry{
		{EntryType: "receipt", EntryID: "r1", ContentHash: "sha256:aaa"},
		{EntryType: "decision", EntryID: "d1", ContentHash: "sha256:bbb"},
	}

	pack, err := portal.ExportEvidence(context.Background(), "t1", "run-1", entries)
	if err != nil {
		t.Fatal(err)
	}
	if pack.ContentHash == "" {
		t.Fatal("expected content hash")
	}
	if len(pack.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(pack.Entries))
	}
}
