package kernel

import (
	"context"
	"testing"
	"time"
)

//nolint:gocognit // test complexity is acceptable
func TestCyberneticsRuntime(t *testing.T) {
	t.Run("Register and tick loop", func(t *testing.T) {
		runtime := NewCyberneticsRuntime(nil, nil)

		loop := &ControlLoop{
			LoopID:  "loop-1",
			Name:    "Test Loop",
			Cadence: "1m",
			State:   "active",
			Variables: []EssentialVariable{
				{
					VariableID:   "var-1",
					Name:         "CPU Usage",
					CurrentValue: 50.0,
					LowerBound:   0.0,
					UpperBound:   100.0,
				},
			},
		}

		err := runtime.RegisterLoop(loop)
		if err != nil {
			t.Fatalf("RegisterLoop failed: %v", err)
		}

		tick, err := runtime.Tick(context.Background(), "loop-1")
		if err != nil {
			t.Fatalf("Tick failed: %v", err)
		}
		if tick.LoopID != "loop-1" {
			t.Errorf("LoopID = %q, want 'loop-1'", tick.LoopID)
		}
		if tick.TickHash == "" {
			t.Error("TickHash should be set")
		}
	})

	t.Run("Register mode and check effects", func(t *testing.T) {
		runtime := NewCyberneticsRuntime(nil, nil)

		mode := &OperationalMode{
			ModeID:         "normal",
			Name:           "Normal Mode",
			AllowedEffects: []string{"DATA_WRITE", "NOTIFY"},
			BlockedEffects: []string{"FUNDS_TRANSFER"},
		}

		err := runtime.RegisterMode(mode)
		if err != nil {
			t.Fatalf("RegisterMode failed: %v", err)
		}

		err = runtime.ActivateMode(context.Background(), "normal")
		if err != nil {
			t.Fatalf("ActivateMode failed: %v", err)
		}

		// Check allowed effect
		allowed, reason := runtime.CheckEffectAllowed("DATA_WRITE")
		if !allowed {
			t.Errorf("DATA_WRITE should be allowed: %s", reason)
		}

		// Check blocked effect
		allowed, _ = runtime.CheckEffectAllowed("FUNDS_TRANSFER")
		if allowed {
			t.Error("FUNDS_TRANSFER should be blocked")
		}
	})

	t.Run("GetActiveMode", func(t *testing.T) {
		runtime := NewCyberneticsRuntime(nil, nil)

		mode := &OperationalMode{
			ModeID: "test-mode",
			Name:   "Test Mode",
		}
		_ = runtime.RegisterMode(mode)
		_ = runtime.ActivateMode(context.Background(), "test-mode")

		active := runtime.GetActiveMode()
		if active == nil {
			t.Fatal("Active mode should not be nil")
		}
		if active.ModeID != "test-mode" {
			t.Errorf("ModeID = %q, want 'test-mode'", active.ModeID)
		}
	})

	t.Run("GetActiveMode no active", func(t *testing.T) {
		runtime := NewCyberneticsRuntime(nil, nil)
		active := runtime.GetActiveMode()
		if active != nil {
			t.Error("Active mode should be nil when none set")
		}
	})

	t.Run("UpdateVariable", func(t *testing.T) {
		runtime := NewCyberneticsRuntime(nil, nil)

		loop := &ControlLoop{
			LoopID:  "loop-1",
			Name:    "Test Loop",
			Cadence: "1m",
			Variables: []EssentialVariable{
				{
					VariableID:   "var-1",
					Name:         "CPU Usage",
					CurrentValue: 50.0,
				},
			},
		}
		_ = runtime.RegisterLoop(loop)

		err := runtime.UpdateVariable("loop-1", "var-1", 75.0)
		if err != nil {
			t.Fatalf("UpdateVariable failed: %v", err)
		}
	})

	t.Run("UpdateVariable loop not found", func(t *testing.T) {
		runtime := NewCyberneticsRuntime(nil, nil)

		err := runtime.UpdateVariable("nonexistent", "var-1", 75.0)
		if err == nil {
			t.Error("Should fail for nonexistent loop")
		}
	})

	t.Run("UpdateVariable variable not found", func(t *testing.T) {
		runtime := NewCyberneticsRuntime(nil, nil)

		loop := &ControlLoop{
			LoopID:    "loop-1",
			Variables: []EssentialVariable{},
		}
		_ = runtime.RegisterLoop(loop)

		err := runtime.UpdateVariable("loop-1", "nonexistent", 75.0)
		if err == nil {
			t.Error("Should fail for nonexistent variable")
		}
	})

	t.Run("Tick loop not found", func(t *testing.T) {
		runtime := NewCyberneticsRuntime(nil, nil)

		_, err := runtime.Tick(context.Background(), "nonexistent")
		if err == nil {
			t.Error("Should fail for nonexistent loop")
		}
	})

	t.Run("ActivateMode not found", func(t *testing.T) {
		runtime := NewCyberneticsRuntime(nil, nil)

		err := runtime.ActivateMode(context.Background(), "nonexistent")
		if err == nil {
			t.Error("Should fail for nonexistent mode")
		}
	})

	t.Run("CheckEffectAllowed no mode", func(t *testing.T) {
		runtime := NewCyberneticsRuntime(nil, nil)

		allowed, _ := runtime.CheckEffectAllowed("DATA_WRITE")
		if !allowed {
			t.Error("Should allow effects when no mode is active")
		}
	})
}

func TestEssentialVariable(t *testing.T) {
	t.Run("Variable with bounds", func(t *testing.T) {
		v := EssentialVariable{
			VariableID:   "test-var",
			Name:         "Test Variable",
			CurrentValue: 50.0,
			LowerBound:   0.0,
			UpperBound:   100.0,
			LastUpdated:  time.Now(),
		}

		if v.CurrentValue < v.LowerBound || v.CurrentValue > v.UpperBound {
			t.Error("CurrentValue should be within bounds")
		}
	})
}

func TestControlLoop(t *testing.T) {
	t.Run("Loop struct", func(t *testing.T) {
		loop := ControlLoop{
			LoopID:   "test-loop",
			Name:     "Test Loop",
			Cadence:  "5m",
			State:    "active",
			LastTick: time.Now(),
			NextTick: time.Now().Add(5 * time.Minute),
		}

		if loop.LoopID == "" {
			t.Error("LoopID should be set")
		}
	})
}

func TestLoopTick(t *testing.T) {
	t.Run("Tick struct", func(t *testing.T) {
		tick := LoopTick{
			TickID:    "tick-1",
			LoopID:    "loop-1",
			Timestamp: time.Now(),
		}

		if tick.TickID == "" {
			t.Error("TickID should be set")
		}
	})
}
