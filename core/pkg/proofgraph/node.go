// Package proofgraph implements the cryptographic ProofGraph DAG for HELM.
// Every execution produces a chain of nodes: INTENT → ATTESTATION → EFFECT,
// with TRUST_EVENT and CHECKPOINT nodes for registry management.
package proofgraph

import (
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
type Node struct {
	ID           string   `json:"id"`
	Type         NodeType `json:"type"`
	ParentIDs    []string `json:"parent_ids"`
	PayloadHash  string   `json:"payload_hash"`
	Payload      []byte   `json:"payload,omitempty"`
	Signature    string   `json:"signature,omitempty"`
	SignerKeyID  string   `json:"signer_key_id,omitempty"`
	LamportClock uint64   `json:"lamport_clock"`
	Timestamp    int64    `json:"timestamp_unix"`
	PrevNodeHash string   `json:"prev_node_hash"`
	NodeHash     string   `json:"node_hash"`
}

// ComputeNodeHash computes the deterministic hash of the node (excluding NodeHash itself).
func (n *Node) ComputeNodeHash() string {
	h := sha256.New()
	h.Write([]byte(n.ID))
	h.Write([]byte(string(n.Type)))
	for _, p := range n.ParentIDs {
		h.Write([]byte(p))
	}
	h.Write([]byte(n.PayloadHash))
	h.Write([]byte(n.Signature))
	h.Write([]byte(n.PrevNodeHash))
	h.Write([]byte(fmt.Sprintf("%d", n.LamportClock)))
	h.Write([]byte(fmt.Sprintf("%d", n.Timestamp)))
	return hex.EncodeToString(h.Sum(nil))
}

// Validate checks the node hash integrity.
func (n *Node) Validate() error {
	expected := n.ComputeNodeHash()
	if n.NodeHash != expected {
		return fmt.Errorf("node %s hash mismatch: got %s, want %s", n.ID, n.NodeHash, expected)
	}
	return nil
}

// PayloadHashOf computes the SHA-256 of arbitrary payload bytes.
func PayloadHashOf(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

// NewNode creates a properly initialized node.
func NewNode(nodeType NodeType, parentIDs []string, payload []byte, prevHash string, lamport uint64) *Node {
	n := &Node{
		ID:           fmt.Sprintf("pg-%s-%d", nodeType, time.Now().UnixNano()),
		Type:         nodeType,
		ParentIDs:    parentIDs,
		PayloadHash:  PayloadHashOf(payload),
		Payload:      payload,
		LamportClock: lamport,
		Timestamp:    time.Now().Unix(),
		PrevNodeHash: prevHash,
	}
	n.NodeHash = n.ComputeNodeHash()
	return n
}

// EncodePayload is a helper to JSON-marshal a payload for node creation.
func EncodePayload(v any) ([]byte, error) {
	return json.Marshal(v)
}
