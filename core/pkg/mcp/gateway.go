package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GatewayConfig configures the MCP gateway server.
type GatewayConfig struct {
	ListenAddr string `json:"listen_addr"`
}

// Gateway is an MCP server that exposes tool execution with governance.
type Gateway struct {
	catalog Catalog
	config  GatewayConfig
}

// NewGateway creates a new MCP gateway.
func NewGateway(catalog Catalog, config GatewayConfig) *Gateway {
	return &Gateway{
		catalog: catalog,
		config:  config,
	}
}

// MCPToolCallRequest is the wire format for an MCP tool call.
type MCPToolCallRequest struct {
	Method string         `json:"method"`
	Params map[string]any `json:"params,omitempty"`
}

// MCPToolCallResponse is the wire format for an MCP tool result.
type MCPToolCallResponse struct {
	Result any    `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
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

	manifest := MCPCapabilityManifest{
		ServerName:   "helm-mcp-gateway",
		Version:      "1.0.0",
		Capabilities: tools,
		Governance:   "helm:pep:v1",
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(manifest)
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

	// In production, this would:
	// 1. Validate args via PEP boundary
	// 2. Request decision from Guardian
	// 3. Execute via SafeExecutor
	// 4. Validate output for drift
	// For now, return a governed response
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(MCPToolCallResponse{
		Error: fmt.Sprintf("tool %q requires governance approval — use HELM console to approve", req.Method),
	})
}
