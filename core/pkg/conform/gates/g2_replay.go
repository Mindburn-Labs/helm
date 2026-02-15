package gates

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Mindburn-Labs/helm/core/pkg/conform"
)

// G2Replay validates deterministic replay with VCR Tape per §G2.
type G2Replay struct{}

func (g *G2Replay) ID() string   { return "G2" }
func (g *G2Replay) Name() string { return "Deterministic Replay with VCR Tape" }

func (g *G2Replay) Run(ctx *conform.RunContext) *conform.GateResult {
	result := &conform.GateResult{
		GateID:        g.ID(),
		Pass:          true,
		Reasons:       []string{},
		EvidencePaths: []string{},
		Metrics:       conform.GateMetrics{Counts: make(map[string]int)},
	}

	tapesDir := filepath.Join(ctx.EvidenceDir, "08_TAPES")
	diffsDir := filepath.Join(ctx.EvidenceDir, "05_DIFFS")

	// 1. Check tape_manifest.json exists
	manifestPath := filepath.Join(tapesDir, "tape_manifest.json")
	if !fileExists(manifestPath) {
		result.Pass = false
		result.Reasons = append(result.Reasons, conform.ReasonReplayTapeMiss)
		return result
	}
	result.EvidencePaths = append(result.EvidencePaths, "08_TAPES/tape_manifest.json")

	// 2. Check diffs directory — empty on PASS
	diffEntries, err := os.ReadDir(diffsDir)
	if err == nil && len(diffEntries) > 0 {
		result.Pass = false
		result.Reasons = append(result.Reasons, conform.ReasonReplayHashDivergence)
		result.EvidencePaths = append(result.EvidencePaths, "05_DIFFS/")
	}

	// 3. Check determinism_manifest.json
	detManifestPath := filepath.Join(ctx.EvidenceDir, "02_PROOFGRAPH", "determinism_manifest.json")
	if fileExists(detManifestPath) {
		result.EvidencePaths = append(result.EvidencePaths, "02_PROOFGRAPH/determinism_manifest.json")

		data, err := os.ReadFile(detManifestPath)
		if err == nil {
			var dm map[string]any
			if json.Unmarshal(data, &dm) == nil {
				if liveHash, ok := dm["live_hash"].(string); ok {
					if replayHash, ok := dm["replay_hash"].(string); ok {
						if liveHash != replayHash {
							result.Pass = false
							result.Reasons = append(result.Reasons, conform.ReasonReplayHashDivergence)
						}
						result.Metrics.Counts["hash_comparisons"] = 1
					}
				}
			}
		}
	}

	return result
}
