//go:build !go1.24

package pqc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestNewPQCSigner tests the legacy signer (PQC disabled).
func TestNewPQCSigner(t *testing.T) {
	signer, err := NewPQCSigner(nil)
	require.NoError(t, err)
	require.NotNil(t, signer)
	require.False(t, signer.IsPQCEnabled(), "PQC should be disabled in legacy mode")
}

func TestNewPQCSignerWithConfig(t *testing.T) {
	config := &PQCConfig{
		KeyID:     "test-key-id",
		EnablePQC: true, // Ignored in legacy
		KeyExpiry: 24 * time.Hour,
	}

	signer, err := NewPQCSigner(config)
	require.NoError(t, err)
	require.Equal(t, "test-key-id", signer.KeyID())
	require.False(t, signer.IsPQCEnabled(), "PQC should be disabled in legacy mode")
}

// TestRealMLKEMEncapsulateDecapsulate skips on legacy builds.
func TestRealMLKEMEncapsulateDecapsulate(t *testing.T) {
	t.Skip("ML-KEM requires Go 1.24+ with helmpqc build tag")
}

func TestEncapsulateToSelf(t *testing.T) {
	t.Skip("ML-KEM requires Go 1.24+ with helmpqc build tag")
}

func TestEncapsulateDisabledPQC(t *testing.T) {
	config := &PQCConfig{
		KeyID:     "no-pqc",
		EnablePQC: false,
	}

	signer, err := NewPQCSigner(config)
	require.NoError(t, err)

	_, err = signer.Encapsulate([]byte("dummy"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "not supported")
}

func TestPublicKeys(t *testing.T) {
	signer, err := NewPQCSigner(nil)
	require.NoError(t, err)

	keys := signer.PublicKeys()
	require.Contains(t, keys, "ed25519")
	// ML-KEM not present in legacy
	require.NotContains(t, keys, "ml-kem-768")

	// Verify key lengths (hex encoded)
	require.Len(t, keys["ed25519"], Ed25519PublicKeySize*2)
}

func TestVerifyWrongAlgorithm(t *testing.T) {
	signer, err := NewPQCSigner(nil)
	require.NoError(t, err)

	sig := &Signature{
		Value:     []byte("fake"),
		Algorithm: "unknown-algo",
	}

	// Legacy verifier returns false, not error for unknown algorithm
	valid, _ := signer.Verify([]byte("data"), sig)
	require.False(t, valid)
}

func TestDefaultPQCConfig(t *testing.T) {
	config := DefaultPQCConfig()
	require.False(t, config.EnablePQC, "PQC should be disabled in legacy config")
	require.Equal(t, 365*24*time.Hour, config.KeyExpiry)
	require.NotEmpty(t, config.KeyID)
}

func TestKeyPairFromSigner(t *testing.T) {
	signer, err := NewPQCSigner(nil)
	require.NoError(t, err)

	// Get Ed25519 key pair
	edKP := signer.KeyPairFromSigner(AlgorithmEd25519)
	require.NotNil(t, edKP)
	require.Equal(t, AlgorithmEd25519, edKP.Algorithm)

	// Get ML-KEM key pair (should be nil in legacy)
	mlKP := signer.KeyPairFromSigner(AlgorithmMLKEM768)
	require.Nil(t, mlKP)

	// Invalid algorithm
	nilKP := signer.KeyPairFromSigner("invalid")
	require.Nil(t, nilKP)
}

// TestFIPS203Compliance is skipped in legacy mode.
func TestFIPS203Compliance(t *testing.T) {
	t.Skip("ML-KEM FIPS 203 requires Go 1.24+ with helmpqc build tag")
}

// Note: FuzzMLKEMDecapsulation is in pqc_fuzz_test.go (helmpqc tag only)
