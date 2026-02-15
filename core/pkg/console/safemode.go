// Package console — SafeModeController.
//
// Per HELM 2030 Spec §3A — Workbench Safe-Mode Controls:
//
//	Pause autonomy, quarantine runs, require manual approvals, resume with receipts.
package console

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// SafeModeState represents the current safe mode.
type SafeModeState string

const (
	ModeAutonomous     SafeModeState = "AUTONOMOUS"      // Normal operation
	ModePaused         SafeModeState = "PAUSED"          // All autonomy paused
	ModeManualRequired SafeModeState = "MANUAL_REQUIRED" // Every action needs approval
	ModeQuarantined    SafeModeState = "QUARANTINED"     // Specific runs quarantined
)

// SafeModeReceipt records a state transition.
type SafeModeReceipt struct {
	ReceiptID    string        `json:"receipt_id"`
	Transition   string        `json:"transition"` // e.g. "AUTONOMOUS→PAUSED"
	FromState    SafeModeState `json:"from_state"`
	ToState      SafeModeState `json:"to_state"`
	InitiatedBy  string        `json:"initiated_by"`
	Reason       string        `json:"reason"`
	AffectedRuns []string      `json:"affected_runs,omitempty"`
	Timestamp    time.Time     `json:"timestamp"`
	ContentHash  string        `json:"content_hash"`
}

// SafeModeStatus reports current state.
type SafeModeStatus struct {
	State           SafeModeState `json:"state"`
	Since           time.Time     `json:"since"`
	QuarantinedRuns []string      `json:"quarantined_runs,omitempty"`
	TransitionCount int           `json:"transition_count"`
}

// SafeModeController manages autonomy modes.
type SafeModeController struct {
	mu              sync.Mutex
	state           SafeModeState
	stateSince      time.Time
	quarantinedRuns map[string]bool
	receipts        []SafeModeReceipt
	seq             int64
	clock           func() time.Time
}

// NewSafeModeController creates a controller starting in AUTONOMOUS mode.
func NewSafeModeController() *SafeModeController {
	return &SafeModeController{
		state:           ModeAutonomous,
		stateSince:      time.Now(),
		quarantinedRuns: make(map[string]bool),
		receipts:        make([]SafeModeReceipt, 0),
		clock:           time.Now,
	}
}

// WithClock overrides clock for testing.
func (c *SafeModeController) WithClock(clock func() time.Time) *SafeModeController {
	c.clock = clock
	c.stateSince = clock()
	return c
}

func (c *SafeModeController) transition(toState SafeModeState, initiatedBy, reason string, affectedRuns []string) *SafeModeReceipt {
	c.seq++
	now := c.clock()
	receiptID := fmt.Sprintf("sm-%d", c.seq)
	fromState := c.state
	transition := fmt.Sprintf("%s→%s", fromState, toState)

	hashInput := fmt.Sprintf("%s:%s:%s:%s", receiptID, transition, initiatedBy, reason)
	h := sha256.Sum256([]byte(hashInput))

	receipt := SafeModeReceipt{
		ReceiptID:    receiptID,
		Transition:   transition,
		FromState:    fromState,
		ToState:      toState,
		InitiatedBy:  initiatedBy,
		Reason:       reason,
		AffectedRuns: affectedRuns,
		Timestamp:    now,
		ContentHash:  "sha256:" + hex.EncodeToString(h[:]),
	}

	c.state = toState
	c.stateSince = now
	c.receipts = append(c.receipts, receipt)
	return &receipt
}

// Pause halts all autonomous operations.
func (c *SafeModeController) Pause(initiatedBy, reason string) *SafeModeReceipt {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.transition(ModePaused, initiatedBy, reason, nil)
}

// Quarantine adds runs to quarantine.
func (c *SafeModeController) Quarantine(initiatedBy string, runIDs []string) *SafeModeReceipt {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, id := range runIDs {
		c.quarantinedRuns[id] = true
	}
	return c.transition(ModeQuarantined, initiatedBy, "runs quarantined", runIDs)
}

// RequireManual switches to manual-approval mode.
func (c *SafeModeController) RequireManual(initiatedBy, reason string) *SafeModeReceipt {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.transition(ModeManualRequired, initiatedBy, reason, nil)
}

// Resume returns to autonomous mode. Requires justification.
func (c *SafeModeController) Resume(approverID, justification string) (*SafeModeReceipt, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.state == ModeAutonomous {
		return nil, fmt.Errorf("already in AUTONOMOUS mode")
	}
	if justification == "" {
		return nil, fmt.Errorf("justification required to resume autonomy")
	}

	// Clear quarantined runs on resume
	c.quarantinedRuns = make(map[string]bool)
	return c.transition(ModeAutonomous, approverID, justification, nil), nil
}

// IsRunQuarantined checks if a specific run is quarantined.
func (c *SafeModeController) IsRunQuarantined(runID string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.quarantinedRuns[runID]
}

// IsAutonomyAllowed returns true if autonomy actions are permitted.
func (c *SafeModeController) IsAutonomyAllowed() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.state == ModeAutonomous
}

// Status returns the current safe-mode status.
func (c *SafeModeController) Status() SafeModeStatus {
	c.mu.Lock()
	defer c.mu.Unlock()

	runs := make([]string, 0, len(c.quarantinedRuns))
	for id := range c.quarantinedRuns {
		runs = append(runs, id)
	}

	return SafeModeStatus{
		State:           c.state,
		Since:           c.stateSince,
		QuarantinedRuns: runs,
		TransitionCount: len(c.receipts),
	}
}
