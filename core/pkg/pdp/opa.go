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
	defaultOPATimeout = 5 * time.Second
	defaultOPAPath    = "/v1/data/helm/authz"
)

// OPAConfig configures the OPA adapter.
type OPAConfig struct {
	// URL is the base URL of the OPA server (e.g., "http://localhost:8181").
	URL string `json:"url"`
	// PolicyPath overrides the default OPA decision path.
	// Default: "/v1/data/helm/authz"
	PolicyPath string `json:"policy_path,omitempty"`
	// Timeout sets the HTTP call timeout. Default: 5s.
	Timeout time.Duration `json:"timeout,omitempty"`
	// PolicyVersion is a human-readable identifier for the policy bundle.
	PolicyVersion string `json:"policy_version,omitempty"`
}

// OPAPDP implements PolicyDecisionPoint using a remote OPA HTTP API.
// Strict fail-closed semantics: any error, timeout, or non-200 response
// results in a DENY.
type OPAPDP struct {
	config     OPAConfig
	client     *http.Client
	policyHash string
}

// NewOPAPDP creates an OPA-backed PDP.
func NewOPAPDP(cfg OPAConfig) *OPAPDP {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = defaultOPATimeout
	}
	if cfg.PolicyPath == "" {
		cfg.PolicyPath = defaultOPAPath
	}
	return &OPAPDP{
		config: cfg,
		client: &http.Client{Timeout: timeout},
		// Policy hash is fetched lazily or set from bundle revision
		policyHash: fmt.Sprintf("sha256:opa:%s", cfg.PolicyVersion),
	}
}

// opaRequest is the OPA input envelope.
type opaRequest struct {
	Input *opaInput `json:"input"`
}

type opaInput struct {
	Principal   string            `json:"principal"`
	Action      string            `json:"action"`
	Resource    string            `json:"resource"`
	Context     map[string]any    `json:"context,omitempty"`
	SchemaHash  string            `json:"schema_hash,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
}

// opaResponse is the OPA result envelope.
type opaResponse struct {
	Result *opaResult `json:"result"`
}

type opaResult struct {
	Allow      bool   `json:"allow"`
	ReasonCode string `json:"reason_code,omitempty"`
}

// Evaluate implements PolicyDecisionPoint. Fail-closed on all errors.
func (o *OPAPDP) Evaluate(ctx context.Context, req *DecisionRequest) (*DecisionResponse, error) {
	if req == nil {
		return o.deny("DENY_NIL_REQUEST"), nil
	}

	body := opaRequest{
		Input: &opaInput{
			Principal:   req.Principal,
			Action:      req.Action,
			Resource:    req.Resource,
			Context:     req.Context,
			SchemaHash:  req.SchemaHash,
			Environment: req.Environment,
		},
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return o.deny("DENY_MARSHAL_ERROR"), nil
	}

	url := o.config.URL + o.config.PolicyPath
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return o.deny("DENY_REQUEST_ERROR"), nil
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(httpReq)
	if err != nil {
		// Fail-closed: timeout, connection refused, etc.
		return o.deny("DENY_OPA_UNREACHABLE"), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return o.deny(fmt.Sprintf("DENY_OPA_HTTP_%d", resp.StatusCode)), nil
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return o.deny("DENY_OPA_READ_ERROR"), nil
	}

	var opaResp opaResponse
	if err := json.Unmarshal(respBody, &opaResp); err != nil {
		return o.deny("DENY_OPA_PARSE_ERROR"), nil
	}

	if opaResp.Result == nil {
		return o.deny("DENY_OPA_NO_RESULT"), nil
	}

	reasonCode := opaResp.Result.ReasonCode
	if reasonCode == "" {
		if opaResp.Result.Allow {
			reasonCode = "ALLOW"
		} else {
			reasonCode = "DENY_POLICY"
		}
	}

	decision := &DecisionResponse{
		Allow:      opaResp.Result.Allow,
		ReasonCode: reasonCode,
		PolicyRef:  fmt.Sprintf("opa:%s:%s", o.config.PolicyVersion, o.config.PolicyPath),
	}

	hash, err := ComputeDecisionHash(decision)
	if err != nil {
		return o.deny("DENY_HASH_FAILURE"), nil
	}
	decision.DecisionHash = hash

	return decision, nil
}

// Backend implements PolicyDecisionPoint.
func (o *OPAPDP) Backend() Backend { return BackendOPA }

// PolicyHash implements PolicyDecisionPoint.
func (o *OPAPDP) PolicyHash() string { return o.policyHash }

func (o *OPAPDP) deny(reason string) *DecisionResponse {
	resp := &DecisionResponse{
		Allow:      false,
		ReasonCode: reason,
		PolicyRef:  fmt.Sprintf("opa:%s", o.config.PolicyVersion),
	}
	hash, _ := ComputeDecisionHash(resp)
	resp.DecisionHash = hash
	return resp
}
