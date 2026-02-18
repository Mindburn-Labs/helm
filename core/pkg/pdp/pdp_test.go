package pdp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// --- Interface conformance tests ---

// allBackends returns PDP instances for each backend.
// OPA and Cedar use httptest servers.
func allBackends(t *testing.T) map[Backend]PolicyDecisionPoint {
	t.Helper()

	// 1. HELM PDP
	helmPDP := NewHelmPDP("v0.1.0", map[string]bool{
		"allowed_resource": true,
		"denied_resource":  false,
	})

	// 2. OPA mock server
	opaMux := http.NewServeMux()
	opaMux.HandleFunc("/v1/data/helm/authz", func(w http.ResponseWriter, r *http.Request) {
		var req opaRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		allow := req.Input.Resource != "denied_resource"
		reason := "ALLOW"
		if !allow {
			reason = "DENY_POLICY"
		}
		resp := opaResponse{
			Result: &opaResult{Allow: allow, ReasonCode: reason},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
	opaServer := httptest.NewServer(opaMux)
	t.Cleanup(opaServer.Close)

	opaPDP := NewOPAPDP(OPAConfig{
		URL:           opaServer.URL,
		PolicyVersion: "test-v1",
	})

	// 3. Cedar mock server
	cedarMux := http.NewServeMux()
	cedarMux.HandleFunc("/decide", func(w http.ResponseWriter, r *http.Request) {
		var req cedarRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		decision := "Allow"
		if req.Resource == "denied_resource" {
			decision = "Deny"
		}
		resp := cedarResponse{Decision: decision}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
	cedarServer := httptest.NewServer(cedarMux)
	t.Cleanup(cedarServer.Close)

	cedarPDP := NewCedarPDP(CedarConfig{
		URL:           cedarServer.URL,
		PolicyVersion: "test-v1",
	})

	return map[Backend]PolicyDecisionPoint{
		BackendHELM:  helmPDP,
		BackendOPA:   opaPDP,
		BackendCedar: cedarPDP,
	}
}

func TestPDPInterfaceConformance_Allow(t *testing.T) {
	backends := allBackends(t)
	for name, pdp := range backends {
		t.Run(string(name)+"_allow", func(t *testing.T) {
			resp, err := pdp.Evaluate(context.Background(), &DecisionRequest{
				Principal: "user:alice",
				Action:    "read",
				Resource:  "allowed_resource",
				Timestamp: time.Now(),
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !resp.Allow {
				t.Errorf("expected allow, got deny (reason=%s)", resp.ReasonCode)
			}
			if resp.DecisionHash == "" {
				t.Error("decision hash must not be empty")
			}
			if !strings.HasPrefix(resp.DecisionHash, "sha256:") {
				t.Errorf("decision hash must start with sha256:, got %s", resp.DecisionHash)
			}
			if resp.PolicyRef == "" {
				t.Error("policy ref must not be empty")
			}
		})
	}
}

func TestPDPInterfaceConformance_Deny(t *testing.T) {
	backends := allBackends(t)
	for name, pdp := range backends {
		t.Run(string(name)+"_deny", func(t *testing.T) {
			resp, err := pdp.Evaluate(context.Background(), &DecisionRequest{
				Principal: "user:alice",
				Action:    "write",
				Resource:  "denied_resource",
				Timestamp: time.Now(),
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp.Allow {
				t.Error("expected deny, got allow")
			}
			if resp.ReasonCode == "" || resp.ReasonCode == "ALLOW" {
				t.Errorf("expected deny reason code, got %q", resp.ReasonCode)
			}
			if resp.DecisionHash == "" {
				t.Error("decision hash must not be empty on deny")
			}
		})
	}
}

func TestPDPInterfaceConformance_NilRequest(t *testing.T) {
	backends := allBackends(t)
	for name, pdp := range backends {
		t.Run(string(name)+"_nil", func(t *testing.T) {
			resp, err := pdp.Evaluate(context.Background(), nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp.Allow {
				t.Error("nil request must be denied (fail-closed)")
			}
		})
	}
}

func TestPDPInterfaceConformance_BackendIdentity(t *testing.T) {
	backends := allBackends(t)
	expected := map[Backend]bool{BackendHELM: true, BackendOPA: true, BackendCedar: true}
	for name, pdp := range backends {
		t.Run(string(name)+"_backend", func(t *testing.T) {
			b := pdp.Backend()
			if !expected[b] {
				t.Errorf("unexpected backend: %v", b)
			}
			if pdp.PolicyHash() == "" {
				t.Error("policy hash must not be empty")
			}
		})
	}
}

func TestPDPInterfaceConformance_DecisionHashDeterminism(t *testing.T) {
	backends := allBackends(t)
	for name, pdp := range backends {
		t.Run(string(name)+"_determinism", func(t *testing.T) {
			req := &DecisionRequest{
				Principal: "user:bob",
				Action:    "execute",
				Resource:  "allowed_resource",
				Timestamp: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			}
			resp1, _ := pdp.Evaluate(context.Background(), req)
			resp2, _ := pdp.Evaluate(context.Background(), req)

			if resp1.DecisionHash != resp2.DecisionHash {
				t.Errorf("decision hash not deterministic: %s vs %s",
					resp1.DecisionHash, resp2.DecisionHash)
			}
		})
	}
}

// --- Fail-closed tests ---

func TestOPA_FailClosed_Unreachable(t *testing.T) {
	pdp := NewOPAPDP(OPAConfig{
		URL:           "http://127.0.0.1:1", // unreachable port
		PolicyVersion: "test-v1",
		Timeout:       100 * time.Millisecond,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	resp, err := pdp.Evaluate(ctx, &DecisionRequest{
		Principal: "user:alice",
		Action:    "read",
		Resource:  "test",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Allow {
		t.Error("unreachable OPA must deny (fail-closed)")
	}
	if resp.ReasonCode != "DENY_OPA_UNREACHABLE" {
		t.Errorf("expected DENY_OPA_UNREACHABLE, got %s", resp.ReasonCode)
	}
}

func TestCedar_FailClosed_Unreachable(t *testing.T) {
	pdp := NewCedarPDP(CedarConfig{
		URL:           "http://127.0.0.1:1", // unreachable port
		PolicyVersion: "test-v1",
		Timeout:       100 * time.Millisecond,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	resp, err := pdp.Evaluate(ctx, &DecisionRequest{
		Principal: "user:alice",
		Action:    "read",
		Resource:  "test",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Allow {
		t.Error("unreachable Cedar must deny (fail-closed)")
	}
	if resp.ReasonCode != "DENY_CEDAR_UNREACHABLE" {
		t.Errorf("expected DENY_CEDAR_UNREACHABLE, got %s", resp.ReasonCode)
	}
}

func TestOPA_FailClosed_BadResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	pdp := NewOPAPDP(OPAConfig{URL: srv.URL, PolicyVersion: "err-test"})
	resp, _ := pdp.Evaluate(context.Background(), &DecisionRequest{
		Principal: "user:alice",
		Action:    "read",
		Resource:  "test",
	})
	if resp.Allow {
		t.Error("500 response must deny")
	}
	if resp.ReasonCode != "DENY_OPA_HTTP_500" {
		t.Errorf("expected DENY_OPA_HTTP_500, got %s", resp.ReasonCode)
	}
}

func TestComputeDecisionHash(t *testing.T) {
	resp := &DecisionResponse{
		Allow:      true,
		ReasonCode: "ALLOW",
		PolicyRef:  "helm:v1",
	}
	hash, err := ComputeDecisionHash(resp)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(hash, "sha256:") {
		t.Errorf("expected sha256: prefix, got %s", hash)
	}

	// Deterministic
	hash2, _ := ComputeDecisionHash(resp)
	if hash != hash2 {
		t.Errorf("hash not deterministic: %s vs %s", hash, hash2)
	}
}
