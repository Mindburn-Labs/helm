package firewall

import (
	"context"
	"testing"
)

type testDispatcher struct{}

func (d *testDispatcher) Dispatch(ctx context.Context, toolName string, params map[string]any) (any, error) {
	return map[string]any{"tool": toolName, "params": params}, nil
}

func TestPolicyFirewall_BlockUnknown(t *testing.T) {
	fw := NewPolicyFirewall(&testDispatcher{})
	if err := fw.AllowTool("other_tool", ""); err != nil {
		t.Fatalf("allow tool failed: %v", err)
	}

	_, err := fw.CallTool(context.Background(), PolicyInputBundle{}, "unknown_tool", nil)
	if err == nil {
		t.Error("Expected error for unknown tool, got nil")
	}
}

func TestPolicyFirewall_AllowKnown(t *testing.T) {
	fw := NewPolicyFirewall(&testDispatcher{})
	if err := fw.AllowTool("known_tool", "{}"); err != nil {
		t.Fatalf("allow tool failed: %v", err)
	}

	res, err := fw.CallTool(context.Background(), PolicyInputBundle{}, "known_tool", map[string]any{"foo": "bar"})
	if err != nil {
		t.Errorf("Expected success, got error: %v", err)
	}
	out, ok := res.(map[string]any)
	if !ok {
		t.Fatalf("unexpected result type: %T", res)
	}
	if out["tool"] != "known_tool" {
		t.Errorf("unexpected tool: %v", out["tool"])
	}
}
