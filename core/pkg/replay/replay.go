package replay

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
)

// Receipt is the minimal structure needed from a proxy receipt for replay.
type Receipt struct {
	ID         string `json:"receipt_id"`
	ToolName   string `json:"tool_name"`
	ArgsHash   string `json:"args_hash"`
	PrevHash   string `json:"prev_hash"`
	Timestamp  string `json:"timestamp"`
	Status     string `json:"status"`
	ReasonCode string `json:"reason_code"`
	Signature  string `json:"signature,omitempty"`
}

// ReplayResult holds the outcome of replaying a receipt chain.
type ReplayResult struct {
	TotalReceipts  int            `json:"total_receipts"`
	ValidChain     bool           `json:"valid_chain"`
	ChainBreaks    []string       `json:"chain_breaks,omitempty"`
	DuplicateIDs   []string       `json:"duplicate_ids,omitempty"`
	OrderValid     bool           `json:"order_valid"`
	HashesVerified int            `json:"hashes_verified"`
	HashMismatches []string       `json:"hash_mismatches,omitempty"`
	Summary        map[string]int `json:"summary"` // reason_code -> count
}

// ReplayFromFile reads a JSONL receipt file and replays the chain offline.
func ReplayFromFile(path string) (*ReplayResult, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open receipts: %w", err)
	}
	defer f.Close()

	return ReplayFromReader(f)
}

// ReplayFromReader replays a receipt chain from an io.Reader (JSONL format).
func ReplayFromReader(r io.Reader) (*ReplayResult, error) {
	dec := json.NewDecoder(r)

	var receipts []Receipt
	for dec.More() {
		var receipt Receipt
		if err := dec.Decode(&receipt); err != nil {
			return nil, fmt.Errorf("decode receipt: %w", err)
		}
		receipts = append(receipts, receipt)
	}

	return Replay(receipts)
}

// Replay verifies a receipt chain for causal integrity.
func Replay(receipts []Receipt) (*ReplayResult, error) {
	result := &ReplayResult{
		TotalReceipts: len(receipts),
		ValidChain:    true,
		OrderValid:    true,
		Summary:       make(map[string]int),
	}

	if len(receipts) == 0 {
		return result, nil
	}

	// Check for duplicate IDs
	idSeen := make(map[string]bool, len(receipts))
	for _, r := range receipts {
		if idSeen[r.ID] {
			result.DuplicateIDs = append(result.DuplicateIDs, r.ID)
		}
		idSeen[r.ID] = true
	}

	// Verify causal chain via prevHash linking
	prevHash := ""
	for i, r := range receipts {
		// Count reason codes
		if r.ReasonCode != "" {
			result.Summary[r.ReasonCode]++
		}

		// Verify prevHash chain
		if i == 0 {
			// First receipt should have empty or zero prevHash
			if r.PrevHash != "" && r.PrevHash != "0000000000000000000000000000000000000000000000000000000000000000" {
				result.ChainBreaks = append(result.ChainBreaks,
					fmt.Sprintf("receipt[0] %s: expected empty prevHash, got %s", r.ID, r.PrevHash))
				result.ValidChain = false
			}
		} else {
			if r.PrevHash != prevHash {
				result.ChainBreaks = append(result.ChainBreaks,
					fmt.Sprintf("receipt[%d] %s: prevHash mismatch (expected %s, got %s)",
						i, r.ID, prevHash, r.PrevHash))
				result.ValidChain = false
			}
		}

		// Compute this receipt's hash for next iteration
		receiptBytes, _ := json.Marshal(r)
		h := sha256.Sum256(receiptBytes)
		prevHash = hex.EncodeToString(h[:])
		result.HashesVerified++
	}

	// Verify ordering (timestamps should be monotonic)
	timestamps := make([]string, len(receipts))
	for i, r := range receipts {
		timestamps[i] = r.Timestamp
	}
	if !sort.StringsAreSorted(timestamps) {
		result.OrderValid = false
	}

	return result, nil
}
