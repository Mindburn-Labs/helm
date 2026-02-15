package api

import (
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

// TestOpenAPISpec_Integrity verifies the OpenAPI spec loads and has required endpoints.
func TestOpenAPISpec_Integrity(t *testing.T) {
	// Find openapi.yaml relative to repo root
	paths := []string{
		"../../docs/api/openapi.yaml",
		"../../../docs/api/openapi.yaml",
	}

	var data []byte
	var err error
	for _, p := range paths {
		data, err = os.ReadFile(p)
		if err == nil {
			break
		}
	}
	if err != nil {
		t.Skip("openapi.yaml not found (run from repo root)")
		return
	}

	var spec map[string]interface{}
	if err := yaml.Unmarshal(data, &spec); err != nil {
		t.Fatalf("openapi.yaml parse error: %v", err)
	}

	// Verify required paths exist
	pathsMap, ok := spec["paths"].(map[string]interface{})
	if !ok {
		t.Fatal("openapi.yaml missing paths section")
	}

	required := []string{
		"/health",
		"/api/v1/kernel/dispatch",
		"/api/v1/kernel/approve",
		"/api/v1/trust/keys/add",
		"/api/v1/trust/keys/revoke",
		"/v1/chat/completions",
		"/mcp/v1/capabilities",
		"/mcp/v1/execute",
	}

	for _, path := range required {
		if _, exists := pathsMap[path]; !exists {
			t.Errorf("openapi.yaml missing required path: %s", path)
		}
	}
}
