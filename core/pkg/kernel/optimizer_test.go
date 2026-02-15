package kernel

import (
	"context"
	"testing"
)

func TestOptimizer_Attenuation_I36(t *testing.T) {
	// Allow the call to pass the hard rate-limit gate so we can test
	// attenuation behavior (mode selection) as a soft constraint.
	policy := BackpressurePolicy{RPM: 60, Burst: 1}
	store := NewInMemoryLimiterStore()
	opt := NewOptimizer(policy, StrategyCostAware)

	// 1. Initial State
	if opt.currentMode != ModeSmart {
		t.Errorf("expected Smart mode, got %s", opt.currentMode)
	}

	// 2. Trigger Attenuation (low budget => Fast)
	mode, err := opt.CheckAndAttenuate(
		context.Background(),
		"actor-1",
		store,
		map[string]float64{"budget_remaining": 0.10},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 3. Verify Downgrade
	if mode != ModeFast {
		t.Errorf("expected Fast mode after pressure, got %s", mode)
	}
}
