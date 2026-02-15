// Package console — Trust Portal backend API.
//
// Per HELM 2030 Spec — Trust Portal:
//   - Exportable EvidencePacks for external auditors
//   - Policy registry with change ledger
//   - Compliance posture mapping
package console

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// PolicyRegistryEntry represents a policy in the registry.
type PolicyRegistryEntry struct {
	PolicyID    string    `json:"policy_id"`
	PolicyName  string    `json:"policy_name"`
	Version     string    `json:"version"`
	ContentHash string    `json:"content_hash"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Active      bool      `json:"active"`
}

// PolicyChange records a change to a policy.
type PolicyChange struct {
	ChangeID    string    `json:"change_id"`
	PolicyID    string    `json:"policy_id"`
	OldVersion  string    `json:"old_version"`
	NewVersion  string    `json:"new_version"`
	ChangedBy   string    `json:"changed_by"`
	ChangedAt   time.Time `json:"changed_at"`
	Description string    `json:"description"`
}

// CompliancePosture represents the compliance state of a tenant.
type CompliancePosture struct {
	TenantID     string                    `json:"tenant_id"`
	Frameworks   map[string]FrameworkState `json:"frameworks"`
	OverallScore float64                   `json:"overall_score"` // 0-100
	AssessedAt   time.Time                 `json:"assessed_at"`
}

// FrameworkState is the compliance state for one framework.
type FrameworkState struct {
	FrameworkID   string  `json:"framework_id"`
	Name          string  `json:"name"`
	Score         float64 `json:"score"` // 0-100
	ControlsMet   int     `json:"controls_met"`
	ControlsTotal int     `json:"controls_total"`
}

// ExportableEvidencePack is a self-contained evidence bundle for external auditors.
type ExportableEvidencePack struct {
	ExportID    string                `json:"export_id"`
	TenantID    string                `json:"tenant_id"`
	RunID       string                `json:"run_id,omitempty"`
	ExportedAt  time.Time             `json:"exported_at"`
	ContentHash string                `json:"content_hash"`
	Entries     []EvidenceExportEntry `json:"entries"`
}

// EvidenceExportEntry is one item in an exported evidence pack.
type EvidenceExportEntry struct {
	EntryType   string                 `json:"entry_type"` // receipt, decision, event, attestation
	EntryID     string                 `json:"entry_id"`
	ContentHash string                 `json:"content_hash"`
	Data        map[string]interface{} `json:"data"`
}

// PortalAPI is the Trust Portal backend.
type PortalAPI struct {
	mu        sync.RWMutex
	policies  map[string]*PolicyRegistryEntry
	changes   []PolicyChange
	postures  map[string]*CompliancePosture
	exports   map[string]*ExportableEvidencePack
	clock     func() time.Time
	changeSeq int64
}

// NewPortalAPI creates a new trust portal backend.
func NewPortalAPI() *PortalAPI {
	return &PortalAPI{
		policies: make(map[string]*PolicyRegistryEntry),
		changes:  make([]PolicyChange, 0),
		postures: make(map[string]*CompliancePosture),
		exports:  make(map[string]*ExportableEvidencePack),
		clock:    time.Now,
	}
}

// WithClock overrides clock for testing.
func (p *PortalAPI) WithClock(clock func() time.Time) *PortalAPI {
	p.clock = clock
	return p
}

// RegisterPolicy adds or updates a policy in the registry.
func (p *PortalAPI) RegisterPolicy(entry *PolicyRegistryEntry) {
	p.mu.Lock()
	defer p.mu.Unlock()

	old, exists := p.policies[entry.PolicyID]
	p.policies[entry.PolicyID] = entry

	if exists {
		p.changeSeq++
		p.changes = append(p.changes, PolicyChange{
			ChangeID:    fmt.Sprintf("chg-%d", p.changeSeq),
			PolicyID:    entry.PolicyID,
			OldVersion:  old.Version,
			NewVersion:  entry.Version,
			ChangedAt:   p.clock(),
			Description: "policy updated",
		})
	}
}

// GetPolicy retrieves a policy by ID.
func (p *PortalAPI) GetPolicy(policyID string) (*PolicyRegistryEntry, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	entry, ok := p.policies[policyID]
	if !ok {
		return nil, fmt.Errorf("policy %q not found", policyID)
	}
	return entry, nil
}

// ListChanges returns policy change history.
func (p *PortalAPI) ListChanges() []PolicyChange {
	p.mu.RLock()
	defer p.mu.RUnlock()
	result := make([]PolicyChange, len(p.changes))
	copy(result, p.changes)
	return result
}

// SetCompliancePosture records compliance state for a tenant.
func (p *PortalAPI) SetCompliancePosture(posture *CompliancePosture) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.postures[posture.TenantID] = posture
}

// GetCompliancePosture retrieves compliance state.
func (p *PortalAPI) GetCompliancePosture(tenantID string) (*CompliancePosture, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	posture, ok := p.postures[tenantID]
	if !ok {
		return nil, fmt.Errorf("no compliance posture for tenant %q", tenantID)
	}
	return posture, nil
}

// ExportEvidence creates an exportable evidence pack.
func (p *PortalAPI) ExportEvidence(ctx context.Context, tenantID, runID string, entries []EvidenceExportEntry) (*ExportableEvidencePack, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := p.clock()
	exportID := fmt.Sprintf("export-%s-%d", tenantID, now.UnixNano())

	data, err := json.Marshal(entries)
	if err != nil {
		return nil, err
	}
	h := sha256.Sum256(data)

	pack := &ExportableEvidencePack{
		ExportID:    exportID,
		TenantID:    tenantID,
		RunID:       runID,
		ExportedAt:  now,
		ContentHash: "sha256:" + hex.EncodeToString(h[:]),
		Entries:     entries,
	}

	p.exports[exportID] = pack
	return pack, nil
}
