// Package kernel provides cybernetics runtime for control loops.
// Per Section 4 - Cybernetics Schemas
package kernel

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
	"sync"
	"time"
)

// EssentialVariable represents a critical system variable that must be regulated.
// Per Section 4.1 - Essential Variable Schema
type EssentialVariable struct {
	VariableID       string    `json:"variable_id"`
	Name             string    `json:"name"`
	CurrentValue     float64   `json:"current_value"`
	TargetValue      float64   `json:"target_value"`
	LowerBound       float64   `json:"lower_bound"`
	UpperBound       float64   `json:"upper_bound"`
	EvaluationWindow string    `json:"evaluation_window"` // e.g., "5m", "1h"
	ViolationAction  string    `json:"violation_action"`  // alert, throttle, shutdown
	LastUpdated      time.Time `json:"last_updated"`
}

// ControlLoop represents a regulatory feedback loop.
// Per Section 4.2 - Control Loop Schema
type ControlLoop struct {
	LoopID    string              `json:"loop_id"`
	Name      string              `json:"name"`
	Variables []EssentialVariable `json:"variables"`
	Cadence   string              `json:"cadence"` // e.g., "1m", "5m"
	LastTick  time.Time           `json:"last_tick"`
	NextTick  time.Time           `json:"next_tick"`
	State     string              `json:"state"` // running, paused, error
}

// LoopTick represents a single evaluation tick of a control loop.
type LoopTick struct {
	TickID           string             `json:"tick_id"`
	LoopID           string             `json:"loop_id"`
	Timestamp        time.Time          `json:"timestamp"`
	VariableStates   map[string]float64 `json:"variable_states"`
	Violations       []string           `json:"violations,omitempty"`
	ActionsTriggered []string           `json:"actions_triggered,omitempty"`
	TickHash         string             `json:"tick_hash"`
}

// OperationalMode represents an operational mode with effect constraints.
// Per Section 4.4 - Mode Schema
type OperationalMode struct {
	ModeID         string   `json:"mode_id"`
	Name           string   `json:"name"`
	EntryPredicate string   `json:"entry_predicate"` // CEL expression
	ExitPredicate  string   `json:"exit_predicate"`  // CEL expression
	AllowedEffects []string `json:"allowed_effects"`
	BlockedEffects []string `json:"blocked_effects"`
	IsActive       bool     `json:"is_active"`
}

// CyberneticsRuntime manages control loops and modes.
type CyberneticsRuntime struct {
	mu         sync.RWMutex
	loops      map[string]*ControlLoop
	modes      map[string]*OperationalMode
	activeMode string
	eventLog   EventLog
	scheduler  DeterministicScheduler
}

// NewCyberneticsRuntime creates a new cybernetics runtime.
func NewCyberneticsRuntime(eventLog EventLog, scheduler DeterministicScheduler) *CyberneticsRuntime {
	return &CyberneticsRuntime{
		loops:     make(map[string]*ControlLoop),
		modes:     make(map[string]*OperationalMode),
		eventLog:  eventLog,
		scheduler: scheduler,
	}
}

// RegisterLoop registers a control loop.
func (r *CyberneticsRuntime) RegisterLoop(loop *ControlLoop) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.loops[loop.LoopID] = loop
	return nil
}

// RegisterMode registers an operational mode.
func (r *CyberneticsRuntime) RegisterMode(mode *OperationalMode) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.modes[mode.ModeID] = mode
	return nil
}

// Tick evaluates a control loop and produces a LoopTick event.
// Per Section 4.2 - LoopTick event production
func (r *CyberneticsRuntime) Tick(ctx context.Context, loopID string) (*LoopTick, error) {
	r.mu.Lock()
	loop, ok := r.loops[loopID]
	if !ok {
		r.mu.Unlock()
		return nil, ErrLoopNotFound
	}
	r.mu.Unlock()

	now := time.Now().UTC()
	tick := &LoopTick{
		TickID:           generateTickID(loopID, now),
		LoopID:           loopID,
		Timestamp:        now,
		VariableStates:   make(map[string]float64),
		Violations:       []string{},
		ActionsTriggered: []string{},
	}

	// Evaluate each variable
	for _, v := range loop.Variables {
		tick.VariableStates[v.VariableID] = v.CurrentValue

		// Check for violations
		if v.CurrentValue < v.LowerBound || v.CurrentValue > v.UpperBound {
			tick.Violations = append(tick.Violations, v.VariableID)
			tick.ActionsTriggered = append(tick.ActionsTriggered, v.ViolationAction)
		}
	}

	// Compute tick hash for determinism
	tick.TickHash = computeTickHash(tick)

	// Log to event log
	if r.eventLog != nil {
		_, _ = r.eventLog.Append(ctx, &EventEnvelope{
			EventID:   tick.TickID,
			EventType: "loop.tick",
			Payload: map[string]interface{}{
				"loop_id":    loopID,
				"violations": tick.Violations,
				"tick_hash":  tick.TickHash,
			},
		})
	}

	// Update loop state
	r.mu.Lock()
	loop.LastTick = now
	r.mu.Unlock()

	return tick, nil
}

// CheckEffectAllowed checks if an effect is allowed in the current mode.
// Per Section 4.4 - Mode envelope effect type constraints
func (r *CyberneticsRuntime) CheckEffectAllowed(effectType string) (bool, string) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.activeMode == "" {
		return true, ""
	}

	mode, ok := r.modes[r.activeMode]
	if !ok {
		return true, ""
	}

	// Check blocked effects
	for _, blocked := range mode.BlockedEffects {
		if blocked == effectType || blocked == "*" {
			return false, "effect blocked by mode: " + mode.Name
		}
	}

	// If allowed effects specified, check whitelist
	if len(mode.AllowedEffects) > 0 {
		for _, allowed := range mode.AllowedEffects {
			if allowed == effectType || allowed == "*" {
				return true, ""
			}
		}
		return false, "effect not in allowed list for mode: " + mode.Name
	}

	return true, ""
}

// ActivateMode activates an operational mode.
func (r *CyberneticsRuntime) ActivateMode(ctx context.Context, modeID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	mode, ok := r.modes[modeID]
	if !ok {
		return ErrModeNotFound
	}

	// Deactivate previous mode
	if r.activeMode != "" {
		if prevMode, ok := r.modes[r.activeMode]; ok {
			prevMode.IsActive = false
		}
	}

	mode.IsActive = true
	r.activeMode = modeID

	// Log mode transition
	if r.eventLog != nil {
		_, _ = r.eventLog.Append(ctx, &EventEnvelope{
			EventID:   "mode-" + modeID + "-" + time.Now().UTC().Format(time.RFC3339Nano),
			EventType: "mode.activated",
			Payload: map[string]interface{}{
				"mode_id": modeID,
				"name":    mode.Name,
			},
		})
	}

	return nil
}

// GetActiveMode returns the currently active mode.
func (r *CyberneticsRuntime) GetActiveMode() *OperationalMode {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.activeMode == "" {
		return nil
	}
	return r.modes[r.activeMode]
}

// UpdateVariable updates an essential variable's current value.
func (r *CyberneticsRuntime) UpdateVariable(loopID, variableID string, value float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	loop, ok := r.loops[loopID]
	if !ok {
		return ErrLoopNotFound
	}

	for i := range loop.Variables {
		if loop.Variables[i].VariableID == variableID {
			loop.Variables[i].CurrentValue = value
			loop.Variables[i].LastUpdated = time.Now().UTC()
			return nil
		}
	}

	return ErrVariableNotFound
}

// ScheduleLoopTicks schedules periodic ticks for a control loop.
// Per Section 4.2 - Cadence enforcement
func (r *CyberneticsRuntime) ScheduleLoopTicks(ctx context.Context, loopID string, count int) error {
	r.mu.RLock()
	loop, ok := r.loops[loopID]
	if !ok {
		r.mu.RUnlock()
		return ErrLoopNotFound
	}
	cadence := loop.Cadence
	r.mu.RUnlock()

	duration, err := parseCadence(cadence)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	for i := 0; i < count; i++ {
		tickTime := now.Add(duration * time.Duration(i+1))
		event := &SchedulerEvent{
			EventID:     generateTickID(loopID, tickTime),
			EventType:   "loop.tick",
			ScheduledAt: tickTime,
			Priority:    1, // High priority for control loops
			Payload: map[string]interface{}{
				"loop_id": loopID,
			},
			LoopID: loopID,
		}
		if err := r.scheduler.Schedule(ctx, event); err != nil {
			return err
		}
	}

	return nil
}

func generateTickID(loopID string, t time.Time) string {
	data := loopID + t.Format(time.RFC3339Nano)
	h := sha256.Sum256([]byte(data))
	return "tick-" + hex.EncodeToString(h[:8])
}

func computeTickHash(tick *LoopTick) string {
	// Sort variable states for determinism
	keys := make([]string, 0, len(tick.VariableStates))
	for k := range tick.VariableStates {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	data := map[string]interface{}{
		"tick_id":   tick.TickID,
		"loop_id":   tick.LoopID,
		"timestamp": tick.Timestamp.Format(time.RFC3339Nano),
		"states":    tick.VariableStates,
	}
	jsonData, _ := json.Marshal(data)
	h := sha256.Sum256(jsonData)
	return hex.EncodeToString(h[:])
}

func parseCadence(cadence string) (time.Duration, error) {
	return time.ParseDuration(cadence)
}

// Error types
var (
	ErrLoopNotFound     = errorString("control loop not found")
	ErrModeNotFound     = errorString("mode not found")
	ErrVariableNotFound = errorString("essential variable not found")
)
