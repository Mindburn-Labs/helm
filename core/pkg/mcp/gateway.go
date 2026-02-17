package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mindburn-Labs/helm/core/pkg/bridge"
	"github.com/Mindburn-Labs/helm/core/pkg/manifest"
)

// GatewayConfig configures the MCP gateway server.
type GatewayConfig struct {
	ListenAddr string `json:"listen_addr"`
}

// Gateway is an MCP server that exposes tool execution with governance.
type Gateway struct {
	catalog Catalog
	config  GatewayConfig
	bridge  *bridge.KernelBridge // governance bridge (optional)
}

// GatewayOption configures optional Gateway settings.
type GatewayOption func(*Gateway)

// WithBridge sets the KernelBridge for governance.
func WithBridge(kb *bridge.KernelBridge) GatewayOption {
	return func(g *Gateway) {
		g.bridge = kb
	}
}

// NewGateway creates a new MCP gateway.
func NewGateway(catalog Catalog, config GatewayConfig, opts ...GatewayOption) *Gateway {
	gw := &Gateway{
		catalog: catalog,
		config:  config,
	}
	for _, opt := range opts {
		opt(gw)
	}
	return gw
}

// MCPToolCallRequest is the wire format for an MCP tool call.
type MCPToolCallRequest struct {
	Method string         `json:"method"`
	Params map[string]any `json:"params,omitempty"`
}

// MCPToolCallResponse is the wire format for an MCP tool result.
type MCPToolCallResponse struct {
	Result     any    `json:"result,omitempty"`
	Error      string `json:"error,omitempty"`
	Decision   string `json:"decision,omitempty"`
	ReasonCode string `json:"reason_code,omitempty"`
	ArgsHash   string `json:"args_hash,omitempty"`
	PGNode     string `json:"proofgraph_node,omitempty"`
}

// MCPCapabilityManifest describes the capabilities this server exposes.
type MCPCapabilityManifest struct {
	ServerName   string    `json:"server_name"`
	Version      string    `json:"version"`
	Capabilities []ToolRef `json:"capabilities"`
	Governance   string    `json:"governance"` // "helm:pep:v1"
}

// RegisterRoutes registers MCP gateway HTTP routes.
func (g *Gateway) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/mcp/v1/capabilities", g.handleCapabilities)
	mux.HandleFunc("/mcp/v1/execute", g.handleExecute)
}

func (g *Gateway) handleCapabilities(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	tools, err := g.catalog.Search(ctx, "")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(MCPToolCallResponse{Error: err.Error()})
		return
	}

	m := MCPCapabilityManifest{
		ServerName:   "helm-mcp-gateway",
		Version:      "1.0.0",
		Capabilities: tools,
		Governance:   "helm:pep:v1",
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(m)
}

func (g *Gateway) handleExecute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req MCPToolCallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(MCPToolCallResponse{Error: "invalid request body"})
		return
	}

	// 1. Validate and canonicalize args via PEP boundary
	var argsHash string
	if req.Params != nil {
		result, err := manifest.ValidateAndCanonicalizeToolArgs(nil, req.Params)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(MCPToolCallResponse{
				Error:      fmt.Sprintf("PEP validation failed: %v", err),
				ReasonCode: "SCHEMA_VALIDATION_FAILED",
			})
			return
		}
		argsHash = result.ArgsHash
	}

	// 2. Governance via KernelBridge (if configured)
	resp := MCPToolCallResponse{ArgsHash: argsHash}

	if g.bridge != nil {
		govResult, govErr := g.bridge.Govern(context.Background(), req.Method, argsHash)
		if govErr != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(MCPToolCallResponse{
				Error:      fmt.Sprintf("governance error: %v", govErr),
				ReasonCode: "POLICY_DECISION_MISSING",
			})
			return
		}

		resp.ReasonCode = govResult.ReasonCode
		resp.PGNode = govResult.NodeID
		if govResult.Decision != nil {
			resp.Decision = govResult.Decision.ID
		}

		if !govResult.Allowed {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			resp.Error = fmt.Sprintf("tool %q denied by governance: %s", req.Method, govResult.ReasonCode)
			_ = json.NewEncoder(w).Encode(resp)
			return
		}

		resp.Result = map[string]any{
			"status":  "governed_allow",
			"tool":    req.Method,
			"message": fmt.Sprintf("tool %q approved by Guardian governance", req.Method),
		}
	} else {
		// No bridge: return governed stub response
		resp.Result = map[string]any{
			"status":  "stub",
			"tool":    req.Method,
			"message": fmt.Sprintf("tool %q requires governance — configure KernelBridge for full governance", req.Method),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
