package console

import (
	"testing"
)

func TestSafeModeStartsAutonomous(t *testing.T) {
	c := NewSafeModeController()
	if !c.IsAutonomyAllowed() {
		t.Fatal("expected autonomous by default")
	}
	s := c.Status()
	if s.State != ModeAutonomous {
		t.Fatalf("expected AUTONOMOUS, got %s", s.State)
	}
}

func TestSafeModePause(t *testing.T) {
	c := NewSafeModeController()
	receipt := c.Pause("admin", "suspicious activity")
	if receipt.ToState != ModePaused {
		t.Fatal("expected PAUSED state")
	}
	if c.IsAutonomyAllowed() {
		t.Fatal("autonomy should not be allowed when paused")
	}
}

func TestSafeModeQuarantine(t *testing.T) {
	c := NewSafeModeController()
	c.Quarantine("admin", []string{"run-1", "run-2"})
	if !c.IsRunQuarantined("run-1") {
		t.Fatal("run-1 should be quarantined")
	}
	if c.IsRunQuarantined("run-3") {
		t.Fatal("run-3 should not be quarantined")
	}
}

func TestSafeModeRequireManual(t *testing.T) {
	c := NewSafeModeController()
	c.RequireManual("admin", "security review")
	if c.IsAutonomyAllowed() {
		t.Fatal("autonomy should not be allowed in manual mode")
	}
	s := c.Status()
	if s.State != ModeManualRequired {
		t.Fatalf("expected MANUAL_REQUIRED, got %s", s.State)
	}
}

func TestSafeModeResume(t *testing.T) {
	c := NewSafeModeController()
	c.Pause("admin", "test")
	receipt, err := c.Resume("senior-admin", "all clear after investigation")
	if err != nil {
		t.Fatal(err)
	}
	if receipt.ToState != ModeAutonomous {
		t.Fatal("expected AUTONOMOUS after resume")
	}
	if !c.IsAutonomyAllowed() {
		t.Fatal("autonomy should be allowed after resume")
	}
}

func TestSafeModeResumeRequiresJustification(t *testing.T) {
	c := NewSafeModeController()
	c.Pause("admin", "test")
	_, err := c.Resume("admin", "")
	if err == nil {
		t.Fatal("expected error for empty justification")
	}
}

func TestSafeModeResumeAlreadyAutonomous(t *testing.T) {
	c := NewSafeModeController()
	_, err := c.Resume("admin", "nothing to resume")
	if err == nil {
		t.Fatal("expected error when already autonomous")
	}
}

func TestSafeModeTransitionReceipts(t *testing.T) {
	c := NewSafeModeController()
	c.Pause("a", "r1")
	c.Resume("b", "r2")
	s := c.Status()
	if s.TransitionCount != 2 {
		t.Fatalf("expected 2 transitions, got %d", s.TransitionCount)
	}
}

func TestSafeModeReceiptHash(t *testing.T) {
	c := NewSafeModeController()
	r := c.Pause("admin", "test")
	if r.ContentHash == "" {
		t.Fatal("expected content hash")
	}
}
