package replay

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"testing"
)

func TestReplay_ValidChain(t *testing.T) {
	r1 := Receipt{ID: "r1", ToolName: "calc", Timestamp: "2026-01-01T00:00:00Z", ReasonCode: "TOOL_ALLOWED"}
	data1, _ := json.Marshal(r1)
	h1 := sha256.Sum256(data1)

	r2 := Receipt{ID: "r2", ToolName: "calc", PrevHash: hex.EncodeToString(h1[:]), Timestamp: "2026-01-01T00:00:01Z", ReasonCode: "TOOL_ALLOWED"}

	result, err := Replay([]Receipt{r1, r2})
	if err != nil {
		t.Fatalf("Replay: %v", err)
	}
	if !result.ValidChain {
		t.Errorf("expected valid chain, got breaks: %v", result.ChainBreaks)
	}
	if result.TotalReceipts != 2 {
		t.Errorf("expected 2 receipts, got %d", result.TotalReceipts)
	}
	if result.Summary["TOOL_ALLOWED"] != 2 {
		t.Errorf("expected 2 TOOL_ALLOWED, got %d", result.Summary["TOOL_ALLOWED"])
	}
}

func TestReplay_BrokenChain(t *testing.T) {
	r1 := Receipt{ID: "r1", ToolName: "calc", Timestamp: "2026-01-01T00:00:00Z"}
	r2 := Receipt{ID: "r2", ToolName: "calc", PrevHash: "wrong_hash", Timestamp: "2026-01-01T00:00:01Z"}

	result, err := Replay([]Receipt{r1, r2})
	if err != nil {
		t.Fatalf("Replay: %v", err)
	}
	if result.ValidChain {
		t.Error("expected broken chain")
	}
	if len(result.ChainBreaks) != 1 {
		t.Errorf("expected 1 chain break, got %d", len(result.ChainBreaks))
	}
}

func TestReplay_DuplicateIDs(t *testing.T) {
	r1 := Receipt{ID: "dup", ToolName: "calc", Timestamp: "2026-01-01T00:00:00Z"}
	r2 := Receipt{ID: "dup", ToolName: "calc", Timestamp: "2026-01-01T00:00:01Z"}

	result, err := Replay([]Receipt{r1, r2})
	if err != nil {
		t.Fatalf("Replay: %v", err)
	}
	if len(result.DuplicateIDs) != 1 {
		t.Errorf("expected 1 duplicate, got %d", len(result.DuplicateIDs))
	}
}

func TestReplayFromReader_JSONL(t *testing.T) {
	r1 := Receipt{ID: "r1", ToolName: "calc", Timestamp: "2026-01-01T00:00:00Z"}
	data1, _ := json.Marshal(r1)
	h1 := sha256.Sum256(data1)

	r2 := Receipt{ID: "r2", ToolName: "calc", PrevHash: hex.EncodeToString(h1[:]), Timestamp: "2026-01-01T00:00:01Z"}

	line1, _ := json.Marshal(r1)
	line2, _ := json.Marshal(r2)
	jsonl := append(line1, '\n')
	jsonl = append(jsonl, line2...)

	result, err := ReplayFromReader(bytes.NewReader(jsonl))
	if err != nil {
		t.Fatalf("ReplayFromReader: %v", err)
	}
	if !result.ValidChain {
		t.Errorf("expected valid chain from JSONL reader")
	}
}

func TestReplay_Empty(t *testing.T) {
	result, err := Replay(nil)
	if err != nil {
		t.Fatalf("Replay(nil): %v", err)
	}
	if result.TotalReceipts != 0 {
		t.Errorf("expected 0 receipts, got %d", result.TotalReceipts)
	}
	if !result.ValidChain {
		t.Error("empty chain should be valid")
	}
}
