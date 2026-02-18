package pdp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultCedarTimeout = 5 * time.Second
	defaultCedarPath    = "/decide"
)

// CedarConfig configures the Cedar adapter.
type CedarConfig struct {
	// URL is the base URL of the Cedar PDP sidecar (e.g., "http://localhost:8182").
	URL string `json:"url"`
	// DecidePath overrides the default decision path. Default: "/decide"
	DecidePath string `json:"decide_path,omitempty"`
	// Timeout sets the HTTP call timeout. Default: 5s.
	Timeout time.Duration `json:"timeout,omitempty"`
	// PolicyVersion is a human-readable identifier for the Cedar policy set.
	PolicyVersion string `json:"policy_version,omitempty"`
}

// CedarPDP implements PolicyDecisionPoint using a Cedar PDP sidecar.
// Cedar policies are evaluated by a separate process (Node.js/Rust/Java)
// because no production-grade Go Cedar evaluator exists.
//
// Strict fail-closed: any communication failure results in DENY.
type CedarPDP struct {
	config     CedarConfig
	client     *http.Client
	policyHash string
}

// NewCedarPDP creates a Cedar-backed PDP.
func NewCedarPDP(cfg CedarConfig) *CedarPDP {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = defaultCedarTimeout
	}
	if cfg.DecidePath == "" {
		cfg.DecidePath = defaultCedarPath
	}
	return &CedarPDP{
		config:     cfg,
		client:     &http.Client{Timeout: timeout},
		policyHash: fmt.Sprintf("sha256:cedar:%s", cfg.PolicyVersion),
	}
}

// cedarRequest is the Cedar PDP sidecar input.
type cedarRequest struct {
	Principal string            `json:"principal"`
	Action    string            `json:"action"`
	Resource  string            `json:"resource"`
	Context   map[string]any    `json:"context,omitempty"`
	Entities  map[string]string `json:"entities,omitempty"`
}

// cedarResponse is the Cedar PDP sidecar output.
type cedarResponse struct {
	Decision    string `json:"decision"` // "Allow" or "Deny"
	Diagnostics struct {
		Reason []string `json:"reason,omitempty"`
		Errors []string `json:"errors,omitempty"`
	} `json:"diagnostics,omitempty"`
}

// Evaluate implements PolicyDecisionPoint. Fail-closed on all errors.
func (c *CedarPDP) Evaluate(ctx context.Context, req *DecisionRequest) (*DecisionResponse, error) {
	if req == nil {
		return c.deny("DENY_NIL_REQUEST"), nil
	}

	body := cedarRequest{
		Principal: req.Principal,
		Action:    req.Action,
		Resource:  req.Resource,
		Context:   req.Context,
		Entities:  req.Environment,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return c.deny("DENY_MARSHAL_ERROR"), nil
	}

	url := c.config.URL + c.config.DecidePath
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return c.deny("DENY_REQUEST_ERROR"), nil
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return c.deny("DENY_CEDAR_UNREACHABLE"), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.deny(fmt.Sprintf("DENY_CEDAR_HTTP_%d", resp.StatusCode)), nil
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.deny("DENY_CEDAR_READ_ERROR"), nil
	}

	var cedarResp cedarResponse
	if err := json.Unmarshal(respBody, &cedarResp); err != nil {
		return c.deny("DENY_CEDAR_PARSE_ERROR"), nil
	}

	allowed := cedarResp.Decision == "Allow"
	reasonCode := "ALLOW"
	if !allowed {
		reasonCode = "DENY_POLICY"
		if len(cedarResp.Diagnostics.Reason) > 0 {
			reasonCode = cedarResp.Diagnostics.Reason[0]
		}
	}

	decision := &DecisionResponse{
		Allow:      allowed,
		ReasonCode: reasonCode,
		PolicyRef:  fmt.Sprintf("cedar:%s", c.config.PolicyVersion),
	}

	hash, err := ComputeDecisionHash(decision)
	if err != nil {
		return c.deny("DENY_HASH_FAILURE"), nil
	}
	decision.DecisionHash = hash

	return decision, nil
}

// Backend implements PolicyDecisionPoint.
func (c *CedarPDP) Backend() Backend { return BackendCedar }

// PolicyHash implements PolicyDecisionPoint.
func (c *CedarPDP) PolicyHash() string { return c.policyHash }

func (c *CedarPDP) deny(reason string) *DecisionResponse {
	resp := &DecisionResponse{
		Allow:      false,
		ReasonCode: reason,
		PolicyRef:  fmt.Sprintf("cedar:%s", c.config.PolicyVersion),
	}
	hash, _ := ComputeDecisionHash(resp)
	resp.DecisionHash = hash
	return resp
}
