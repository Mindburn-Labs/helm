package arc

import (
	"context"
	"time"
)

// SourceArtifact represents a raw, content-addressed authoritative source.
type SourceArtifact struct {
	ArtifactID      string            `json:"artifact_id"`  // Unique ID (UUID)
	ContentHash     string            `json:"content_hash"` // SHA-256 of RawContent (CAS key)
	SourceID        string            `json:"source_id"`    // ID of the connector/source (e.g., "eu-eurlex")
	ExternalID      string            `json:"external_id"`  // Stable ID in external system (e.g., "ELI/...")
	IngestedAt      time.Time         `json:"ingested_at"`
	MimeType        string            `json:"mime_type"`
	ConnectorConfig map[string]string `json:"connector_config"` // Config used during ingestion
	Provenance      *SourceProvenance `json:"provenance"`
}

// SourceProvenance captures the chain of custody for the source.
type SourceProvenance struct {
	ConnectorName    string    `json:"connector_name"`
	ConnectorVersion string    `json:"connector_version"`
	RetrievalMethod  string    `json:"retrieval_method"` // "api", "crawl", "manual"
	RetrievedAt      time.Time `json:"retrieved_at"`
	Signature        string    `json:"signature,omitempty"` // Optional signature of the artifact
}

// IngestionReceipt records the outcome of an ingestion attempt.
// This is the "Proof of Work" for Phase 1.
type IngestionReceipt struct {
	ReceiptID     string    `json:"receipt_id"`
	SourceID      string    `json:"source_id"`
	ArtifactID    string    `json:"artifact_id,omitempty"`
	Status        string    `json:"status"` // "SUCCESS", "ERROR", "SKIPPED"
	BytesIngested int64     `json:"bytes_ingested"`
	CostUSD       float64   `json:"cost_usd"`
	Timestamp     time.Time `json:"timestamp"`
	Error         string    `json:"error,omitempty"`
}

// SourceConnector defines the interface for fetching external regulations.
// Implementations must be deterministic and economically bounded.
type SourceConnector interface {
	// ID returns the unique identifier of the connector (e.g., "eu-eurlex").
	ID() string

	// Fetch retrieves the content for a given external ID.
	// It returns the raw bytes, mime type, and any specific metadata.
	Fetch(ctx context.Context, externalID string) ([]byte, string, map[string]string, error)

	// TrustClass returns the trust level of this source.
	TrustClass() TrustClass
}

type TrustClass string

const (
	TrustClassOfficial  TrustClass = "official"  // Official Government API/Portal
	TrustClassPartner   TrustClass = "partner"   // Paid Partner Feed
	TrustClassCommunity TrustClass = "community" // Scraped/Community maintained
)
