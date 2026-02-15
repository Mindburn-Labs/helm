package executor

import (
	"context"
	"testing"

	"github.com/Mindburn-Labs/helm/core/pkg/contracts"
	"github.com/Mindburn-Labs/helm/core/pkg/crypto"
)

func mustSigner(t *testing.T) crypto.Signer {
	t.Helper()
	s, err := crypto.NewEd25519Signer("test-key")
	if err != nil {
		t.Fatalf("failed to create signer: %v", err)
	}
	return s
}

func TestPackExporter_Determinism(t *testing.T) {
	signer := mustSigner(t)
	exporter := NewPackExporter("test-signer", signer)
	ctx := context.Background()

	// 1. ChangePack
	changeInput := &contracts.ChangePack{
		PackID:       "pack-1",
		PackType:     "CHANGE_PACK",
		TargetSystem: "prod-api",
		ChangeContext: contracts.ChangeContext{
			Repo:      "helm/core",
			CommitSHA: "abc123456",
			Branch:    "main",
		},
		EvidenceRefs: contracts.ChangeEvidenceRefs{
			ApprovalReceiptID: "rcpt-approve-1",
			BuildReceiptID:    "rcpt-build-1",
		},
	}

	pack1, err := exporter.ExportChangePack(ctx, changeInput)
	if err != nil {
		t.Fatalf("ExportChangePack failed: %v", err)
	}

	if pack1.Attestation.Signature == "" {
		t.Fatal("Signature is empty — signing failed silently")
	}

	// Verify signature is a real hex-encoded Ed25519 signature (128 hex chars = 64 bytes)
	if len(pack1.Attestation.Signature) != 128 {
		t.Errorf("Expected 128-char hex signature, got %d chars", len(pack1.Attestation.Signature))
	}

	// Determinism check: same inputs => same hash
	changeInput3 := &contracts.ChangePack{
		PackID:       "pack-1",
		PackType:     "CHANGE_PACK",
		TargetSystem: "prod-api",
		ChangeContext: contracts.ChangeContext{
			Repo:      "helm/core",
			CommitSHA: "abc123456",
			Branch:    "main",
		},
		EvidenceRefs: contracts.ChangeEvidenceRefs{
			ApprovalReceiptID: "rcpt-approve-1",
			BuildReceiptID:    "rcpt-build-1",
		},
		Attestation: contracts.ChangePackAttestation{
			GeneratedAt: pack1.Attestation.GeneratedAt, // SYNC TIME
		},
	}

	pack2, err := exporter.ExportChangePack(ctx, changeInput3)
	if err != nil {
		t.Fatalf("ExportChangePack 2 failed: %v", err)
	}

	if pack1.Attestation.PackHash != pack2.Attestation.PackHash {
		t.Errorf("ChangePack hash mismatch (indeterministic): %s vs %s", pack1.Attestation.PackHash, pack2.Attestation.PackHash)
	}
}

func TestPackExporter_TamperCheck(t *testing.T) {
	signer := mustSigner(t)
	exporter := NewPackExporter("test-signer", signer)
	ctx := context.Background()

	input := &contracts.IncidentPack{
		PackID:     "inc-pack-1",
		PackType:   "INCIDENT_PACK",
		IncidentID: "inc-123",
		Timeline: []contracts.IncidentEvent{
			{Description: "Detected"},
		},
	}

	pack, err := exporter.ExportIncidentPack(ctx, input)
	if err != nil {
		t.Fatalf("ExportIncidentPack failed: %v", err)
	}
	originalHash := pack.Attestation.PackHash

	if pack.Attestation.Signature == "" {
		t.Fatal("Signature is empty")
	}

	// Tamper
	pack.Timeline[0].Description = "Covered Up"

	// Re-export (which re-hashes)
	packTampered, err := exporter.ExportIncidentPack(ctx, pack)
	if err != nil {
		t.Fatalf("ExportIncidentPack tampered failed: %v", err)
	}

	if packTampered.Attestation.PackHash == originalHash {
		t.Error("PackHash did NOT change after tampering with body!")
	}
}

func TestPackExporter_SignatureVerifiable(t *testing.T) {
	signer := mustSigner(t)
	exporter := NewPackExporter("test-signer", signer)
	ctx := context.Background()

	input := &contracts.ChangePack{
		PackID:       "verifiable-1",
		PackType:     "CHANGE_PACK",
		TargetSystem: "staging",
		ChangeContext: contracts.ChangeContext{
			Repo:      "helm/test",
			CommitSHA: "def456",
			Branch:    "main",
		},
	}

	pack, err := exporter.ExportChangePack(ctx, input)
	if err != nil {
		t.Fatalf("ExportChangePack failed: %v", err)
	}

	// Verify the signature using the signer's public key
	valid, err := crypto.Verify(signer.PublicKey(), pack.Attestation.Signature, []byte(pack.Attestation.PackHash))
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}
	if !valid {
		t.Error("Signature verification FAILED — pack cannot be trusted")
	}
}
