package ui

import (
	"time"
)

// UIInteraction represents an operator's action on the UI.
type UIInteraction struct {
	InteractionID      string         `json:"interaction_id"`
	SessionID          string         `json:"session_id"`
	ComponentID        string         `json:"component_id"`
	ActionType         string         `json:"action_type"`
	Payload            map[string]any `json:"payload,omitempty"`
	Timestamp          time.Time      `json:"timestamp"`
	PreviousRenderHash string         `json:"previous_render_hash,omitempty"`
}

// UIRenderReceipt confirms what was shown to the user.
type UIRenderReceipt struct {
	RenderID    string    `json:"render_id"`
	SpecHash    string    `json:"spec_hash"`
	RenderedAt  time.Time `json:"rendered_at"`
	ContextHash string    `json:"context_hash,omitempty"`
}

// UIComponentCall defines a component to be rendered.
type UIComponentCall struct {
	ID            string         `json:"id"`
	ComponentName string         `json:"component_name"`
	Props         map[string]any `json:"props"`
}

// UISpec is the full declarative specification for a view.
type UISpec struct {
	Version    string            `json:"version"`
	Layout     map[string]any    `json:"layout,omitempty"`
	Components []UIComponentCall `json:"components"`
	Theme      string            `json:"theme,omitempty"`
}

// UIUpdate is a streamed delta.
type UIUpdate struct {
	SequenceID int64          `json:"sequence_id"`
	Operation  string         `json:"operation"` // e.g. "APPEND", "REPLACE"
	Path       string         `json:"path"`
	Value      map[string]any `json:"value"`
}
