package proofgraph

import (
	"testing"
)

func TestGraph_AppendAndValidate(t *testing.T) {
	g := NewGraph()

	n1, err := g.Append(NodeTypeIntent, []byte(`{"intent":"create_file"}`))
	if err != nil {
		t.Fatal(err)
	}
	if n1.LamportClock != 1 {
		t.Errorf("lamport = %d, want 1", n1.LamportClock)
	}

	n2, err := g.Append(NodeTypeAttestation, []byte(`{"decision":"PASS"}`))
	if err != nil {
		t.Fatal(err)
	}
	if n2.LamportClock != 2 {
		t.Errorf("lamport = %d, want 2", n2.LamportClock)
	}
	// n2 should have n1 as parent
	if len(n2.ParentIDs) != 1 || n2.ParentIDs[0] != n1.ID {
		t.Errorf("n2 parents = %v, want [%s]", n2.ParentIDs, n1.ID)
	}

	n3, err := g.Append(NodeTypeEffect, []byte(`{"effect":"file_created"}`))
	if err != nil {
		t.Fatal(err)
	}

	// Validate full chain
	if err := g.ValidateChain(n3.ID); err != nil {
		t.Fatalf("chain validation failed: %v", err)
	}

	if g.Len() != 3 {
		t.Errorf("graph len = %d, want 3", g.Len())
	}
}

func TestGraph_LamportMonotonicity(t *testing.T) {
	g := NewGraph()

	var prevClock uint64
	for i := 0; i < 100; i++ {
		n, err := g.Append(NodeTypeEffect, []byte(`{}`))
		if err != nil {
			t.Fatal(err)
		}
		if n.LamportClock <= prevClock {
			t.Fatalf("lamport not monotonic: %d <= %d at step %d", n.LamportClock, prevClock, i)
		}
		prevClock = n.LamportClock
	}
}

func TestNode_HashIntegrity(t *testing.T) {
	n := NewNode(NodeTypeIntent, nil, []byte(`test`), "", 1)
	if err := n.Validate(); err != nil {
		t.Fatalf("fresh node should validate: %v", err)
	}

	// Tamper with payload hash
	n.PayloadHash = "tampered"
	if err := n.Validate(); err == nil {
		t.Fatal("tampered node should fail validation")
	}
}

func TestGraph_TrustEvent(t *testing.T) {
	g := NewGraph()

	payload := []byte(`{"event":"KEY_ROTATED","key_id":"k-1","public_key":"abc123"}`)
	n, err := g.Append(NodeTypeTrustEvent, payload)
	if err != nil {
		t.Fatal(err)
	}
	if n.Type != NodeTypeTrustEvent {
		t.Errorf("type = %s, want TRUST_EVENT", n.Type)
	}
}
