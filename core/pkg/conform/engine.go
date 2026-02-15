package conform

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Engine is the conformance engine per §11.1.
// It enumerates gates, runs them deterministically, emits an EvidencePack,
// and signs the final report.
type Engine struct {
	gates   map[string]Gate
	ordered []string // gate execution order
	clock   func() time.Time
}

// NewEngine creates a new conformance engine.
func NewEngine() *Engine {
	return &Engine{
		gates:   make(map[string]Gate),
		ordered: make([]string, 0),
		clock:   time.Now,
	}
}

// WithClock overrides the clock for deterministic testing.
func (e *Engine) WithClock(clock func() time.Time) *Engine {
	e.clock = clock
	return e
}

// RegisterGate adds a gate to the engine.
// Gates are run in registration order.
func (e *Engine) RegisterGate(g Gate) {
	id := g.ID()
	if _, exists := e.gates[id]; !exists {
		e.ordered = append(e.ordered, id)
	}
	e.gates[id] = g
}

// ConformanceReport is the top-level result of a conformance run.
type ConformanceReport struct {
	RunID       string         `json:"run_id"`
	Profile     ProfileID      `json:"profile"`
	Timestamp   time.Time      `json:"timestamp"`
	Pass        bool           `json:"pass"`
	GateResults []*GateResult  `json:"gate_results"`
	Duration    time.Duration  `json:"duration"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

// RunOptions configures a conformance run.
type RunOptions struct {
	Profile      ProfileID
	Jurisdiction string
	GateFilter   []string // if non-empty, only run these gates
	ProjectRoot  string
	OutputDir    string // base dir for EvidencePack (default: artifacts/conformance)
}

// Run executes a full conformance run for the given profile.
// It returns a ConformanceReport and writes the EvidencePack directory.
func (e *Engine) Run(opts *RunOptions) (*ConformanceReport, error) {
	start := e.clock()

	// Pre-flight: check for receipt emission panic
	if opts.OutputDir != "" {
		if panicRec, _ := CheckPanicRecord(opts.OutputDir); panicRec != nil {
			return nil, fmt.Errorf("%s: receipt emission system failed at seq %d: %s",
				ReasonReceiptEmissionPanic, panicRec.LastGoodSeq, panicRec.Reason)
		}
	}
	runID := fmt.Sprintf("run-%s-%d", opts.Profile, start.UnixNano())

	// Determine which gates to run
	requiredGates := e.resolveGates(opts)
	if len(requiredGates) == 0 {
		return nil, fmt.Errorf("no gates to run for profile %s", opts.Profile)
	}

	// Create EvidencePack directory
	dateStr := start.Format("2006-01-02")
	outputBase := opts.OutputDir
	if outputBase == "" {
		outputBase = filepath.Join(opts.ProjectRoot, "artifacts", "conformance")
	}
	evidenceDir := filepath.Join(outputBase, dateStr, runID)
	if err := CreateEvidencePackDirs(evidenceDir); err != nil {
		return nil, fmt.Errorf("failed to create EvidencePack dirs: %w", err)
	}

	// Build run context
	ctx := &RunContext{
		RunID:        runID,
		Profile:      opts.Profile,
		Jurisdiction: opts.Jurisdiction,
		EvidenceDir:  evidenceDir,
		ArtifactsDir: filepath.Join(opts.ProjectRoot, "artifacts"),
		ProjectRoot:  opts.ProjectRoot,
		Clock:        e.clock,
		ExtraConfig:  make(map[string]any),
	}

	// Run gates deterministically in order
	results := make([]*GateResult, 0, len(requiredGates))
	allPass := true
	for _, gateID := range requiredGates {
		g, ok := e.gates[gateID]
		if !ok {
			// Missing gate is a hard fail
			results = append(results, &GateResult{
				GateID:  gateID,
				Pass:    false,
				Reasons: []string{fmt.Sprintf("gate %s not registered", gateID)},
			})
			allPass = false
			continue
		}

		gateStart := e.clock()
		result := g.Run(ctx)
		gateDuration := e.clock().Sub(gateStart)
		result.Metrics.DurationMs = gateDuration.Milliseconds()

		results = append(results, result)
		if !result.Pass {
			allPass = false
		}
	}

	report := &ConformanceReport{
		RunID:       runID,
		Profile:     opts.Profile,
		Timestamp:   start,
		Pass:        allPass,
		GateResults: results,
		Duration:    e.clock().Sub(start),
	}

	// Write 01_SCORE.json
	if err := writeScore(evidenceDir, report); err != nil {
		return report, fmt.Errorf("failed to write score: %w", err)
	}

	// Write 00_INDEX.json
	if err := writeIndex(evidenceDir, runID, opts.Profile, e.clock); err != nil {
		return report, fmt.Errorf("failed to write index: %w", err)
	}

	return report, nil
}

// resolveGates returns the gate IDs to run based on options.
func (e *Engine) resolveGates(opts *RunOptions) []string {
	if len(opts.GateFilter) > 0 {
		return opts.GateFilter
	}

	profileGates := GatesForProfile(opts.Profile)
	if profileGates == nil {
		return e.ordered
	}

	// Intersect profile gates with registered gates, preserving registration order
	required := make(map[string]bool, len(profileGates))
	for _, g := range profileGates {
		required[g] = true
	}

	result := make([]string, 0)
	for _, g := range e.ordered {
		if required[g] {
			result = append(result, g)
		}
	}
	return result
}

// writeScore writes the 01_SCORE.json file.
func writeScore(evidenceDir string, report *ConformanceReport) error {
	scoreData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(evidenceDir, "01_SCORE.json"), scoreData, 0600)
}

// IndexEntry is a single artifact reference in 00_INDEX.json per §3.2.
type IndexEntry struct {
	Path          string `json:"path"`
	SHA256        string `json:"sha256"`
	SizeBytes     int64  `json:"size_bytes"`
	SchemaVersion string `json:"schema_version,omitempty"`
	ContentType   string `json:"content_type"`
}

// IndexManifest is the 00_INDEX.json structure per §3.1.
type IndexManifest struct {
	RunID         string       `json:"run_id"`
	Profile       ProfileID    `json:"profile"`
	CreatedAt     time.Time    `json:"created_at"`
	TopoOrderRule string       `json:"topo_order_rule"` // DAG linearization rule
	Entries       []IndexEntry `json:"entries"`
}

// writeIndex walks the EvidencePack directory and creates 00_INDEX.json.
func writeIndex(evidenceDir, runID string, profile ProfileID, clock func() time.Time) error {
	var entries []IndexEntry

	err := filepath.Walk(evidenceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		relPath, _ := filepath.Rel(evidenceDir, path)
		if relPath == "00_INDEX.json" {
			return nil // skip self
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		hash := sha256.Sum256(data)

		entries = append(entries, IndexEntry{
			Path:        relPath,
			SHA256:      hex.EncodeToString(hash[:]),
			SizeBytes:   info.Size(),
			ContentType: inferContentType(relPath),
		})
		return nil
	})
	if err != nil {
		return err
	}

	manifest := IndexManifest{
		RunID:         runID,
		Profile:       profile,
		CreatedAt:     clock().UTC(),
		TopoOrderRule: "seq_monotonic_dag",
		Entries:       entries,
	}

	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(evidenceDir, "00_INDEX.json"), data, 0600)
}

// inferContentType infers content type from file extension.
func inferContentType(path string) string {
	switch filepath.Ext(path) {
	case ".json":
		return "application/json"
	case ".jsonl":
		return "application/x-ndjson"
	case ".sig":
		return "application/pgp-signature"
	default:
		return "application/octet-stream"
	}
}
