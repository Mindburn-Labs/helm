package kernel

import (
	"context"
	"errors"
	"testing"
)

func TestScheduleLoopTicks(t *testing.T) {
	scheduler := NewInMemoryScheduler()
	runtime := NewCyberneticsRuntime(nil, scheduler)
	ctx := context.Background()

	// Register a loop
	loop := &ControlLoop{
		LoopID:  "test-loop",
		Name:    "Test Loop",
		Cadence: "100ms",
	}
	_ = runtime.RegisterLoop(loop)

	// Schedule 3 ticks
	err := runtime.ScheduleLoopTicks(ctx, "test-loop", 3)
	if err != nil {
		t.Fatalf("ScheduleLoopTicks failed: %v", err)
	}
}

func TestScheduleLoopTicksLoopNotFound(t *testing.T) {
	scheduler := NewInMemoryScheduler()
	runtime := NewCyberneticsRuntime(nil, scheduler)

	err := runtime.ScheduleLoopTicks(context.Background(), "nonexistent", 1)
	if !errors.Is(err, ErrLoopNotFound) {
		t.Errorf("Expected ErrLoopNotFound, got %v", err)
	}
}

func TestScheduleLoopTicksInvalidCadence(t *testing.T) {
	scheduler := NewInMemoryScheduler()
	runtime := NewCyberneticsRuntime(nil, scheduler)

	// Register a loop with invalid cadence
	loop := &ControlLoop{
		LoopID:  "bad-loop",
		Cadence: "invalid",
	}
	_ = runtime.RegisterLoop(loop)

	err := runtime.ScheduleLoopTicks(context.Background(), "bad-loop", 1)
	if err == nil {
		t.Error("Expected error for invalid cadence")
	}
}

func TestParseCadence(t *testing.T) {
	tests := []struct {
		cadence   string
		expectErr bool
	}{
		{"1h", false},
		{"30m", false},
		{"500ms", false},
		{"1s", false},
		{"invalid", true},
		{"", true},
	}

	for _, tt := range tests {
		duration, err := parseCadence(tt.cadence)
		if tt.expectErr && err == nil {
			t.Errorf("parseCadence(%q) should error", tt.cadence)
		}
		if !tt.expectErr && err != nil {
			t.Errorf("parseCadence(%q) error: %v", tt.cadence, err)
		}
		if !tt.expectErr && duration <= 0 {
			t.Errorf("parseCadence(%q) duration should be positive", tt.cadence)
		}
	}
}
