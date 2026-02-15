package executor

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Mindburn-Labs/helm/core/pkg/contracts"
	"github.com/Mindburn-Labs/helm/core/pkg/crypto"
	"github.com/gowebpki/jcs"
)

// PackExporter defines the interface for generating and signing proof packs.
type PackExporter interface {
	ExportChangePack(ctx context.Context, input *contracts.ChangePack) (*contracts.ChangePack, error)
	ExportIncidentPack(ctx context.Context, input *contracts.IncidentPack) (*contracts.IncidentPack, error)
	ExportAccessReviewPack(ctx context.Context, input *contracts.AccessReviewPack) (*contracts.AccessReviewPack, error)
	ExportVendorDueDiligencePack(ctx context.Context, input *contracts.VendorDueDiligencePack) (*contracts.VendorDueDiligencePack, error)
}

type packExporter struct {
	signerID string
	signer   crypto.Signer
}

// NewPackExporter creates a new PackExporter with a real cryptographic signer.
// The signer MUST implement crypto.Signer (e.g. Ed25519Signer).
func NewPackExporter(signerID string, signer crypto.Signer) PackExporter {
	return &packExporter{
		signerID: signerID,
		signer:   signer,
	}
}

func (e *packExporter) ExportChangePack(ctx context.Context, input *contracts.ChangePack) (*contracts.ChangePack, error) {
	if input.Attestation.GeneratedAt.IsZero() {
		input.Attestation.GeneratedAt = time.Now().UTC()
	}
	input.Attestation.SignerID = e.signerID

	hash, err := e.computePackHash(input)
	if err != nil {
		return nil, fmt.Errorf("failed to compute pack hash: %w", err)
	}
	input.Attestation.PackHash = hash

	sig, err := e.signer.Sign([]byte(hash))
	if err != nil {
		return nil, fmt.Errorf("failed to sign pack: %w", err)
	}
	input.Attestation.Signature = sig

	input.Attestation.Signature = sig

	return input, nil
}

func (e *packExporter) ExportIncidentPack(ctx context.Context, input *contracts.IncidentPack) (*contracts.IncidentPack, error) {
	if input.Attestation.GeneratedAt.IsZero() {
		input.Attestation.GeneratedAt = time.Now().UTC()
	}

	hash, err := e.computePackHash(input)
	if err != nil {
		return nil, fmt.Errorf("failed to compute pack hash: %w", err)
	}
	input.Attestation.PackHash = hash

	sig, err := e.signer.Sign([]byte(hash))
	if err != nil {
		return nil, fmt.Errorf("failed to sign pack: %w", err)
	}
	input.Attestation.Signature = sig

	input.Attestation.Signature = sig

	return input, nil
}

func (e *packExporter) ExportAccessReviewPack(ctx context.Context, input *contracts.AccessReviewPack) (*contracts.AccessReviewPack, error) {
	if input.Attestation.GeneratedAt.IsZero() {
		input.Attestation.GeneratedAt = time.Now().UTC()
	}

	hash, err := e.computePackHash(input)
	if err != nil {
		return nil, fmt.Errorf("failed to compute pack hash: %w", err)
	}
	input.Attestation.PackHash = hash

	sig, err := e.signer.Sign([]byte(hash))
	if err != nil {
		return nil, fmt.Errorf("failed to sign pack: %w", err)
	}
	input.Attestation.Signature = sig

	input.Attestation.Signature = sig

	return input, nil
}

func (e *packExporter) ExportVendorDueDiligencePack(ctx context.Context, input *contracts.VendorDueDiligencePack) (*contracts.VendorDueDiligencePack, error) {
	if input.Attestation.GeneratedAt.IsZero() {
		input.Attestation.GeneratedAt = time.Now().UTC()
	}

	hash, err := e.computePackHash(input)
	if err != nil {
		return nil, fmt.Errorf("failed to compute pack hash: %w", err)
	}
	input.Attestation.PackHash = hash

	sig, err := e.signer.Sign([]byte(hash))
	if err != nil {
		return nil, fmt.Errorf("failed to sign pack: %w", err)
	}
	input.Attestation.Signature = sig

	input.Attestation.Signature = sig

	return input, nil
}

// computePackHash computes the SHA-256 hash of the JCS (RFC 8785) canonicalized pack data.
// It strips the attestation hash/signature fields before hashing to avoid circular references.
func (e *packExporter) computePackHash(data interface{}) (string, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	var flatMap map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &flatMap); err != nil {
		return "", err
	}

	// Remove attestation.pack_hash and attestation.signature before hashing.
	if attestation, ok := flatMap["attestation"].(map[string]interface{}); ok {
		delete(attestation, "pack_hash")
		delete(attestation, "signature")
	}

	modifiedJSON, err := json.Marshal(flatMap)
	if err != nil {
		return "", err
	}

	canonicalJSON, err := jcs.Transform(modifiedJSON)
	if err != nil {
		return "", err
	}

	hasher := sha256.New()
	hasher.Write(canonicalJSON)
	return fmt.Sprintf("sha256:%s", hex.EncodeToString(hasher.Sum(nil))), nil
}
