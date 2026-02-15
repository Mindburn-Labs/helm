package swarm

import (
	"time"
)

// SwarmResultEvent is the payload for "TITAN_EVT.swarm.completed"
type SwarmResultEvent struct {
	// JobID matches the Command ID
	JobID string `json:"job_id"`
	// Success indicates if the swarm job completed successfully
	Success bool `json:"success"`
	// Artifacts is a map of generated/modified paths to their new git hash or diff
	Artifacts map[string]string `json:"artifacts"`
	// Report is a human-readable summary of what the swarm did
	Report string `json:"report"`
	// Error contains details if Success is false
	Error string `json:"error,omitempty"`
	// Timestamp is when the event was emitted
	Timestamp time.Time `json:"timestamp"`
}

const (
	// SubjectResult is the NATS subject where results are published
	SubjectResult = "TITAN_EVT.swarm.completed"
)
