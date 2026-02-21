package swarm

import (
	"time"
)

// SwarmDeployCommand is the payload for "TITAN_CMD.swarm.deploy"
type SwarmDeployCommand struct {
	// ID is the unique Correlation ID for this swarm job (UUID)
	ID string `json:"id"`
	// Intent is the high-level goal (e.g., "Migrate service X")
	Intent string `json:"intent"`
	// Strategy defines the parallelization strategy
	Strategy string `json:"strategy"` // "file_sharding", "service_sharding", "monolithic"
	// TargetScope defines the files or directories to operate on
	TargetScope []string `json:"target_scope"`
	// ReplyTo is the NATS subject to publish the result to (optional)
	ReplyTo string `json:"reply_to,omitempty"`
	// Governance
	PhenotypeHash string `json:"phenotype_hash"` // The law under which this runs
	StateCursor   string `json:"state_cursor"`   // The state revision
	// Timestamp is when the command was issued
	Timestamp time.Time `json:"timestamp"`
}

const (
	// SubjectDeploy is the NATS subject for deploying a swarm
	SubjectDeploy = "TITAN_CMD.swarm.deploy"
)
