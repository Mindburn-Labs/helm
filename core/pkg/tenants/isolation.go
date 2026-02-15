// Package tenants — Tenant Isolation Proofs.
//
// Per HELM 2030 Spec:
//   - Formal cross-tenant prevention with runtime assertions
//   - IsolationChecker verifies operations don't leak across tenant boundaries
//   - Produces IsolationReceipt as proof
package tenants

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// IsolationReceipt proves no cross-tenant leakage occurred.
type IsolationReceipt struct {
	ReceiptID    string    `json:"receipt_id"`
	TenantID     string    `json:"tenant_id"`
	OperationID  string    `json:"operation_id"`
	ChecksPassed int       `json:"checks_passed"`
	ChecksFailed int       `json:"checks_failed"`
	Violations   []string  `json:"violations,omitempty"`
	Isolated     bool      `json:"isolated"`
	ContentHash  string    `json:"content_hash"`
	Timestamp    time.Time `json:"timestamp"`
}

// IsolationChecker performs cross-tenant boundary checks.
type IsolationChecker struct {
	mu         sync.RWMutex
	tenantData map[string]map[string]bool // tenantID → set of resource IDs
	seq        int64
	clock      func() time.Time
}

// NewIsolationChecker creates a new checker.
func NewIsolationChecker() *IsolationChecker {
	return &IsolationChecker{
		tenantData: make(map[string]map[string]bool),
		clock:      time.Now,
	}
}

// WithClock overrides clock for testing.
func (c *IsolationChecker) WithClock(clock func() time.Time) *IsolationChecker {
	c.clock = clock
	return c
}

// RegisterResource associates a resource with a tenant.
func (c *IsolationChecker) RegisterResource(tenantID, resourceID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.tenantData[tenantID] == nil {
		c.tenantData[tenantID] = make(map[string]bool)
	}
	c.tenantData[tenantID][resourceID] = true
}

// CheckAccess verifies a tenant can only access its own resources.
func (c *IsolationChecker) CheckAccess(tenantID string, resourceIDs []string) *IsolationReceipt {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.seq++
	receipt := &IsolationReceipt{
		ReceiptID:   fmt.Sprintf("iso-%d", c.seq),
		TenantID:    tenantID,
		OperationID: fmt.Sprintf("op-%d", c.seq),
		Isolated:    true,
		Timestamp:   c.clock(),
	}

	tenantResources := c.tenantData[tenantID]

	for _, resourceID := range resourceIDs {
		if tenantResources != nil && tenantResources[resourceID] {
			receipt.ChecksPassed++
			continue
		}

		// Check if resource belongs to another tenant
		crossTenant := false
		for otherTenant, resources := range c.tenantData {
			if otherTenant != tenantID && resources[resourceID] {
				crossTenant = true
				receipt.Violations = append(receipt.Violations,
					fmt.Sprintf("tenant %s attempted to access resource %s owned by %s", tenantID, resourceID, otherTenant))
				break
			}
		}

		if crossTenant {
			receipt.ChecksFailed++
			receipt.Isolated = false
		} else {
			// Resource not registered to any tenant — could be unregistered
			receipt.ChecksPassed++
		}
	}

	// Compute content hash
	hashInput := fmt.Sprintf("%s:%s:%d:%d", receipt.TenantID, receipt.OperationID, receipt.ChecksPassed, receipt.ChecksFailed)
	h := sha256.Sum256([]byte(hashInput))
	receipt.ContentHash = "sha256:" + hex.EncodeToString(h[:])

	return receipt
}

// VerifyIsolation does a comprehensive cross-tenant check for all tenants.
func (c *IsolationChecker) VerifyIsolation() (bool, []string) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var violations []string
	// Check that no resource is claimed by multiple tenants
	resourceOwners := make(map[string]string)
	for tenantID, resources := range c.tenantData {
		for resourceID := range resources {
			if owner, exists := resourceOwners[resourceID]; exists {
				violations = append(violations,
					fmt.Sprintf("resource %s claimed by both %s and %s", resourceID, owner, tenantID))
			} else {
				resourceOwners[resourceID] = tenantID
			}
		}
	}

	return len(violations) == 0, violations
}
