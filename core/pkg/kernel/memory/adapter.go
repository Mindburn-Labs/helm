package memory

import "time"

// MemoryAdapter acts as a bridge between the kernel and the semantic memory engine.
type MemoryAdapter interface {
	// IngestSourceBundle ingests raw documents.
	IngestSourceBundle(bundle SourceBundle) (MemoryBuildManifest, error)

	// QueryMemory performs evidence-based retrieval.
	QueryMemory(query string) ([]QueryResult, error)

	// Promote marks a memory graph as trusted.
	Promote(ref DocumentGraphRef) error
}

type SourceBundle struct {
	BundleID   string   `json:"bundle_id"`
	SourceURIs []string `json:"source_uris"`
}

type MemoryBuildManifest struct {
	ManifestID string           `json:"manifest_id"`
	BundleID   string           `json:"bundle_id"`
	BuiltAt    time.Time        `json:"built_at"`
	GraphRef   DocumentGraphRef `json:"graph_ref"`
}

type DocumentGraphRef struct {
	GraphID string `json:"graph_id"`
	Version string `json:"version"`
	Hash    string `json:"hash"`
}

type QueryResult struct {
	Content   string  `json:"content"`
	SourceURI string  `json:"source_uri"`
	Score     float64 `json:"score"`
}
