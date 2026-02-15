package pqc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPQCMerkleBuilder(t *testing.T) {
	builder := NewPQCMerkleBuilder(SHA256PQC)
	require.NotNil(t, builder)
	require.Equal(t, SHA256PQC, builder.algorithm)
}

func TestPQCMerkleBuilderDefaultAlgorithm(t *testing.T) {
	builder := NewPQCMerkleBuilder("")
	require.Equal(t, SHA256PQC, builder.algorithm)
}

func TestPQCMerkleAddLeaf(t *testing.T) {
	builder := NewPQCMerkleBuilder(SHA256PQC)

	err := builder.AddLeaf("/intent", map[string]string{"action": "execute"}, false)
	require.NoError(t, err)
	require.Len(t, builder.leaves, 1)
}

func TestPQCMerkleBuild(t *testing.T) {
	builder := NewPQCMerkleBuilder(SHA256PQC)
	builder.AddLeafBytes("/a", []byte("data-a"), false)
	builder.AddLeafBytes("/b", []byte("data-b"), false)

	tree, err := builder.Build()
	require.NoError(t, err)
	require.NotNil(t, tree.Root)
	require.Equal(t, SHA256PQC, tree.Algorithm)
	require.Len(t, tree.Leaves, 2)
}

func TestPQCMerkleBuildEmpty(t *testing.T) {
	builder := NewPQCMerkleBuilder(SHA256PQC)

	_, err := builder.Build()
	require.Error(t, err)
}

func TestPQCMerkleRootHex(t *testing.T) {
	builder := NewPQCMerkleBuilder(SHA256PQC)
	builder.AddLeafBytes("/test", []byte("test-data"), false)

	tree, err := builder.Build()
	require.NoError(t, err)

	rootHex := tree.RootHex()
	require.Len(t, rootHex, 64) // 32 bytes = 64 hex chars
}

func TestPQCMerkleProofGeneration(t *testing.T) {
	builder := NewPQCMerkleBuilder(SHA256PQC)
	builder.AddLeafBytes("/a", []byte("a"), false)
	builder.AddLeafBytes("/b", []byte("b"), false)
	builder.AddLeafBytes("/c", []byte("c"), false)

	tree, err := builder.Build()
	require.NoError(t, err)

	proof, err := tree.GenerateProof(1)
	require.NoError(t, err)
	require.Equal(t, 1, proof.LeafIndex)
	require.Equal(t, "/b", proof.LeafPath)
	require.Equal(t, SHA256PQC, proof.Algorithm)
}

func TestPQCMerkleProofVerification(t *testing.T) {
	builder := NewPQCMerkleBuilder(SHA256PQC)
	builder.AddLeafBytes("/a", []byte("data-a"), false)
	builder.AddLeafBytes("/b", []byte("data-b"), false)

	tree, err := builder.Build()
	require.NoError(t, err)

	proof, err := tree.GenerateProof(0)
	require.NoError(t, err)

	valid, err := VerifyPQCProof(proof)
	require.NoError(t, err)
	require.True(t, valid)
}

func TestPQCMerkleEnhanced(t *testing.T) {
	builder := NewPQCMerkleBuilder(SHA256Enhanced)
	builder.AddLeafBytes("/test", []byte("enhanced-test"), false)

	tree, err := builder.Build()
	require.NoError(t, err)
	require.Equal(t, SHA256Enhanced, tree.Algorithm)
	require.Len(t, tree.Root, 32)
}

func TestPQCMerkleEnhancedVerify(t *testing.T) {
	builder := NewPQCMerkleBuilder(SHA256Enhanced)
	builder.AddLeafBytes("/x", []byte("x"), false)
	builder.AddLeafBytes("/y", []byte("y"), false)

	tree, err := builder.Build()
	require.NoError(t, err)

	proof, err := tree.GenerateProof(1)
	require.NoError(t, err)

	valid, err := VerifyPQCProof(proof)
	require.NoError(t, err)
	require.True(t, valid)
}

func TestPQCMerkleDeterminism(t *testing.T) {
	builder1 := NewPQCMerkleBuilder(SHA256PQC)
	builder1.AddLeafBytes("/a", []byte("same"), false)
	builder1.AddLeafBytes("/b", []byte("data"), false)

	builder2 := NewPQCMerkleBuilder(SHA256PQC)
	builder2.AddLeafBytes("/a", []byte("same"), false)
	builder2.AddLeafBytes("/b", []byte("data"), false)

	tree1, _ := builder1.Build()
	tree2, _ := builder2.Build()

	require.Equal(t, tree1.RootHex(), tree2.RootHex())
}

func TestPQCMerkleSignRoot(t *testing.T) {
	signer, err := NewPQCSigner(DefaultPQCConfig())
	require.NoError(t, err)

	builder := NewPQCMerkleBuilder(SHA256PQC)
	builder.AddLeafBytes("/test", []byte("signed-data"), false)

	tree, err := builder.Build()
	require.NoError(t, err)

	err = tree.SignRoot(signer)
	require.NoError(t, err)
	require.NotNil(t, tree.Signature)
	require.NotEmpty(t, tree.SignerID)
}

func TestPQCMerkleVerifySignature(t *testing.T) {
	signer, err := NewPQCSigner(&PQCConfig{EnablePQC: true})
	require.NoError(t, err)

	builder := NewPQCMerkleBuilder(SHA256PQC)
	builder.AddLeafBytes("/evidence", []byte("critical-data"), false)

	tree, err := builder.Build()
	require.NoError(t, err)

	err = tree.SignRoot(signer)
	require.NoError(t, err)

	valid, err := tree.VerifyRootSignature(signer)
	require.NoError(t, err)
	require.True(t, valid)
}

func TestPQCMerkleNoSignatureError(t *testing.T) {
	signer, _ := NewPQCSigner(&PQCConfig{EnablePQC: true})

	builder := NewPQCMerkleBuilder(SHA256PQC)
	builder.AddLeafBytes("/test", []byte("unsigned"), false)

	tree, _ := builder.Build()

	_, err := tree.VerifyRootSignature(signer)
	require.Error(t, err)
}

func TestPQCMerkleProofIndexOutOfRange(t *testing.T) {
	builder := NewPQCMerkleBuilder(SHA256PQC)
	builder.AddLeafBytes("/only", []byte("one"), false)

	tree, _ := builder.Build()

	_, err := tree.GenerateProof(-1)
	require.Error(t, err)

	_, err = tree.GenerateProof(5)
	require.Error(t, err)
}

func TestHashAlgorithmConstants(t *testing.T) {
	require.Equal(t, HashAlgorithm("SHA256-PQC"), SHA256PQC)
	require.Equal(t, HashAlgorithm("SHA256-ENHANCED"), SHA256Enhanced)
}

func TestPQCMerkleProfileID(t *testing.T) {
	require.Equal(t, "merkle-pqc-v1", PQCMerkleProfileID)
}
