package proofgraph

import (
	"fmt"
	"sync"
)

// Graph is an in-memory ProofGraph DAG.
type Graph struct {
	mu      sync.RWMutex
	nodes   map[string]*Node
	heads   []string // Current head node IDs (tips of the DAG)
	lamport uint64
}

// NewGraph creates a new empty ProofGraph.
func NewGraph() *Graph {
	return &Graph{
		nodes: make(map[string]*Node),
	}
}

// Append adds a node to the graph, linking it to the current heads.
// Returns the finalized node with computed hash.
func (g *Graph) Append(nodeType NodeType, payload []byte) (*Node, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.lamport++

	prevHash := ""
	if len(g.heads) > 0 {
		// Use the first head's hash as prev for linear chains
		if head, ok := g.nodes[g.heads[0]]; ok {
			prevHash = head.NodeHash
		}
	}

	node := NewNode(nodeType, g.heads, payload, prevHash, g.lamport)
	g.nodes[node.ID] = node
	g.heads = []string{node.ID}

	return node, nil
}

// AppendSigned adds a signed node to the graph.
func (g *Graph) AppendSigned(nodeType NodeType, payload []byte, signature, signerKeyID string) (*Node, error) {
	node, err := g.Append(nodeType, payload)
	if err != nil {
		return nil, err
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	node.Signature = signature
	node.SignerKeyID = signerKeyID
	node.NodeHash = node.ComputeNodeHash() // Recompute with signature
	return node, nil
}

// Get retrieves a node by ID.
func (g *Graph) Get(id string) (*Node, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	n, ok := g.nodes[id]
	return n, ok
}

// Heads returns the current head node IDs.
func (g *Graph) Heads() []string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	result := make([]string, len(g.heads))
	copy(result, g.heads)
	return result
}

// LamportClock returns the current Lamport clock value.
func (g *Graph) LamportClock() uint64 {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.lamport
}

// ValidateChain walks from a node back through parents and validates hashes.
func (g *Graph) ValidateChain(nodeID string) error {
	g.mu.RLock()
	defer g.mu.RUnlock()

	visited := make(map[string]bool)
	return g.walkValidate(nodeID, visited)
}

func (g *Graph) walkValidate(nodeID string, visited map[string]bool) error {
	if visited[nodeID] {
		return nil
	}
	visited[nodeID] = true

	node, ok := g.nodes[nodeID]
	if !ok {
		return fmt.Errorf("node %s not found", nodeID)
	}

	if err := node.Validate(); err != nil {
		return err
	}

	for _, pid := range node.ParentIDs {
		if err := g.walkValidate(pid, visited); err != nil {
			return err
		}
	}

	return nil
}

// Len returns the number of nodes in the graph.
func (g *Graph) Len() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return len(g.nodes)
}

// AllNodes returns all nodes (for serialization/export).
func (g *Graph) AllNodes() []*Node {
	g.mu.RLock()
	defer g.mu.RUnlock()
	result := make([]*Node, 0, len(g.nodes))
	for _, n := range g.nodes {
		result = append(result, n)
	}
	return result
}
