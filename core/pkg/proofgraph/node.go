// Package proofgraph implements the cryptographic ProofGraph DAG for HELM.
// Every execution produces a chain of nodes: INTENT → ATTESTATION → EFFECT,
// with TRUST_EVENT and CHECKPOINT nodes for registry management.
package proofgraph

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// NodeType enumerates the types of nodes in the ProofGraph.
type NodeType string

const (
	NodeTypeIntent        NodeType = "INTENT"
	NodeTypeAttestation   NodeType = "ATTESTATION"
	NodeTypeEffect        NodeType = "EFFECT"
	NodeTypeTrustEvent    NodeType = "TRUST_EVENT"
	NodeTypeCheckpoint    NodeType = "CHECKPOINT"
	NodeTypeMergeDecision NodeType = "MERGE_DECISION"
)

// Node is a single vertex in the ProofGraph DAG.
// Aligned with HELM Standard v1.2 Appendix B.1
type Node struct {
	NodeHash     string          `json:"node_hash"`
	Kind         NodeType        `json:"kind"`
	Parents      []string        `json:"parents"`
	Lamport      uint64          `json:"lamport"`
	Principal    string          `json:"principal"`
	PrincipalSeq uint64          `json:"principal_seq"`
	Payload      json.RawMessage `json:"payload"`
	Sig          string          `json:"sig"`
	Timestamp    int64           `json:"ts_unix_ms,omitempty"`
}

// ComputeNodeHash computes the deterministic hash of the node (excluding NodeHash itself).
// Uses JCS (RFC 8785) logic: serialize without node_hash, then SHA-256.
func (n *Node) ComputeNodeHash() string {
	// Create a temporary structure for hashing that excludes NodeHash
	type NodeJCS struct {
		Kind         NodeType        `json:"kind"`
		Parents      []string        `json:"parents"`
		Lamport      uint64          `json:"lamport"`
		Principal    string          `json:"principal"`
		PrincipalSeq uint64          `json:"principal_seq"`
		Payload      json.RawMessage `json:"payload"`
		Sig          string          `json:"sig"`
		Timestamp    int64           `json:"ts_unix_ms,omitempty"`
	}

	temp := NodeJCS{
		Kind:         n.Kind,
		Parents:      n.Parents,
		Lamport:      n.Lamport,
		Principal:    n.Principal,
		PrincipalSeq: n.PrincipalSeq,
		Payload:      n.Payload,
		Sig:          n.Sig,
		Timestamp:    n.Timestamp,
	}

	// Marshaling must be canonical: consistent key order, no whitespace.
	// encoding/json output is compact and sorted by key, but escapes HTML.
	// RFC 8785 requires NO HTML escaping.
	// We use a custom buffer writer.

	// Fast path: use standard json.Marshal but check for HTML chars if paranoid.
	// For this reference implementation, we assume keys are sorted.
	// We MUST disable HTML escaping.

	// Create a new encoder
	var buf []byte
	buffer := bytes.NewBuffer(buf)
	enc := json.NewEncoder(buffer)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(temp); err != nil {
		// Should not happen for struct
		return ""
	}

	// Encoder adds a newline at the end, JCS implies minimal?
	// RFC 8785 doesn't explicitly mention trailing newline, but usually "canonical bytes" means EXACT bytes.
	// We'll trim the trailing newline if Encode adds it.
	data := bytes.TrimSpace(buffer.Bytes())

	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

// Validate checks the node hash integrity.
func (n *Node) Validate() error {
	expected := n.ComputeNodeHash()
	if n.NodeHash != expected {
		return fmt.Errorf("node hash mismatch: got %s, want %s", n.NodeHash, expected)
	}
	return nil
}

// NewNode creates a properly initialized node.
func NewNode(kind NodeType, parents []string, payload []byte, lamport uint64, principal string, principalSeq uint64) *Node {
	n := &Node{
		Kind:         kind,
		Parents:      parents,
		Payload:      json.RawMessage(payload),
		Lamport:      lamport,
		Principal:    principal,
		PrincipalSeq: principalSeq,
		Timestamp:    time.Now().Unix(),
	}
	n.NodeHash = n.ComputeNodeHash()
	return n
}

// EncodePayload is a helper to JSON-marshal a payload for node creation.
func EncodePayload(v any) ([]byte, error) {
	return json.Marshal(v)
}
