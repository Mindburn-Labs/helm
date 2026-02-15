//go:build go1.24

package pqc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewPQCSigner(t *testing.T) {
	signer, err := NewPQCSigner(nil)
	require.NoError(t, err)
	require.NotNil(t, signer)
	require.True(t, signer.IsPQCEnabled())
}

func TestNewPQCSignerWithConfig(t *testing.T) {
	config := &PQCConfig{
		KeyID:     "test-key-id",
		EnablePQC: true,
		KeyExpiry: 24 * time.Hour,
	}

	signer, err := NewPQCSigner(config)
	require.NoError(t, err)
	require.Equal(t, "test-key-id", signer.KeyID())
	require.True(t, signer.IsPQCEnabled())
}

func TestPQCSignerClassicalOnly(t *testing.T) {
	config := &PQCConfig{
		KeyID:     "classical-only",
		EnablePQC: false,
	}

	signer, err := NewPQCSigner(config)
	require.NoError(t, err)
	require.False(t, signer.IsPQCEnabled())
}

func TestSignAndVerifyEd25519(t *testing.T) {
	signer, err := NewPQCSigner(nil)
	require.NoError(t, err)

	data := []byte("test message for signing")

	sig, err := signer.Sign(data)
	require.NoError(t, err)
	require.NotNil(t, sig)
	require.Equal(t, AlgorithmEd25519, sig.Algorithm)
	require.Len(t, sig.Value, Ed25519SignatureSize)

	// Verify
	valid, err := signer.Verify(data, sig)
	require.NoError(t, err)
	require.True(t, valid)
}

func TestSignAndVerifyWithContext(t *testing.T) {
	signer, err := NewPQCSigner(nil)
	require.NoError(t, err)

	data := []byte("test message")
	context := "helm.pqc.test"

	sig, err := signer.SignWithContext(data, context)
	require.NoError(t, err)
	require.NotNil(t, sig)

	// Verify with same context
	valid, err := signer.VerifyWithContext(data, sig, context)
	require.NoError(t, err)
	require.True(t, valid)

	// Verify with wrong context should fail
	valid, err = signer.VerifyWithContext(data, sig, "wrong-context")
	require.NoError(t, err)
	require.False(t, valid)
}

func TestSignAndVerifyClassicalOnly(t *testing.T) {
	config := &PQCConfig{
		KeyID:     "ed25519-only",
		EnablePQC: false,
	}

	signer, err := NewPQCSigner(config)
	require.NoError(t, err)

	data := []byte("test message for ed25519 signing")

	sig, err := signer.Sign(data)
	require.NoError(t, err)
	require.Equal(t, AlgorithmEd25519, sig.Algorithm)

	// Verify
	valid, err := signer.Verify(data, sig)
	require.NoError(t, err)
	require.True(t, valid)
}

// TestRealMLKEMEncapsulateDecapsulate tests real FIPS 203 ML-KEM-768.
func TestRealMLKEMEncapsulateDecapsulate(t *testing.T) {
	// Create sender and receiver
	receiver, err := NewPQCSigner(nil)
	require.NoError(t, err)

	sender, err := NewPQCSigner(nil)
	require.NoError(t, err)

	// Get receiver's ML-KEM public key
	receiverPubKey := receiver.MLKEMPublicKey()
	require.NotNil(t, receiverPubKey)
	require.Len(t, receiverPubKey, MLKEMPublicKeySize)

	// Sender encapsulates to receiver
	encap, err := sender.Encapsulate(receiverPubKey)
	require.NoError(t, err)
	require.NotNil(t, encap)
	require.Equal(t, AlgorithmMLKEM768, encap.Algorithm)
	require.Len(t, encap.Ciphertext, MLKEMCiphertextSize)
	require.Len(t, encap.SharedSecret, MLKEMSharedSecretSize)

	// Receiver decapsulates
	sharedSecret, err := receiver.Decapsulate(encap.Ciphertext)
	require.NoError(t, err)
	require.Len(t, sharedSecret, MLKEMSharedSecretSize)

	// CRITICAL: Shared secrets must match (this is the proof it's working!)
	require.Equal(t, encap.SharedSecret, sharedSecret)
}

func TestEncapsulateToSelf(t *testing.T) {
	signer, err := NewPQCSigner(nil)
	require.NoError(t, err)

	encap, err := signer.EncapsulateToSelf()
	require.NoError(t, err)
	require.NotNil(t, encap)

	// Decapsulate our own encapsulation
	ss, err := signer.Decapsulate(encap.Ciphertext)
	require.NoError(t, err)
	require.Equal(t, encap.SharedSecret, ss)
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
	require.Contains(t, err.Error(), "PQC not enabled")
}

func TestDecapsulateInvalidCiphertext(t *testing.T) {
	signer, err := NewPQCSigner(nil)
	require.NoError(t, err)

	// Wrong size ciphertext should fail
	_, err = signer.Decapsulate([]byte("too short"))
	require.Error(t, err)
}

func TestPublicKeys(t *testing.T) {
	signer, err := NewPQCSigner(nil)
	require.NoError(t, err)

	keys := signer.PublicKeys()
	require.Contains(t, keys, "ed25519")
	require.Contains(t, keys, "ml-kem-768")

	// Verify key lengths (hex encoded)
	require.Len(t, keys["ed25519"], Ed25519PublicKeySize*2)
	require.Len(t, keys["ml-kem-768"], MLKEMPublicKeySize*2)
}

func TestPublicKeysClassicalOnly(t *testing.T) {
	config := &PQCConfig{
		KeyID:     "classical",
		EnablePQC: false,
	}

	signer, err := NewPQCSigner(config)
	require.NoError(t, err)

	keys := signer.PublicKeys()
	require.Contains(t, keys, "ed25519")
	require.NotContains(t, keys, "ml-kem-768")
}

func TestSignatureTimestamp(t *testing.T) {
	signer, err := NewPQCSigner(nil)
	require.NoError(t, err)

	before := time.Now()
	sig, err := signer.Sign([]byte("data"))
	after := time.Now()

	require.NoError(t, err)
	require.True(t, sig.Timestamp.After(before) || sig.Timestamp.Equal(before))
	require.True(t, sig.Timestamp.Before(after) || sig.Timestamp.Equal(after))
}

func TestVerifyWrongAlgorithm(t *testing.T) {
	signer, err := NewPQCSigner(nil)
	require.NoError(t, err)

	sig := &Signature{
		Value:     []byte("fake"),
		Algorithm: "unknown-algo",
	}

	_, err = signer.Verify([]byte("data"), sig)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unknown algorithm")
}

func TestVerifyTamperedSignature(t *testing.T) {
	signer, err := NewPQCSigner(nil)
	require.NoError(t, err)

	data := []byte("original message")
	sig, err := signer.Sign(data)
	require.NoError(t, err)

	// Tamper with signature
	sig.Value[0] ^= 0xFF

	valid, err := signer.Verify(data, sig)
	require.NoError(t, err)
	require.False(t, valid)
}

func TestVerifyWrongData(t *testing.T) {
	signer, err := NewPQCSigner(nil)
	require.NoError(t, err)

	originalData := []byte("original message")
	sig, err := signer.Sign(originalData)
	require.NoError(t, err)

	// Verify with wrong data
	valid, err := signer.Verify([]byte("different message"), sig)
	require.NoError(t, err)
	require.False(t, valid)
}

func TestKeyIDGeneration(t *testing.T) {
	signer1, _ := NewPQCSigner(nil)
	signer2, _ := NewPQCSigner(nil)

	// Different signers should have different key IDs
	require.NotEqual(t, signer1.KeyID(), signer2.KeyID())
	require.Len(t, signer1.KeyID(), 16)
}

func TestDefaultPQCConfig(t *testing.T) {
	config := DefaultPQCConfig()
	require.True(t, config.EnablePQC)
	require.Equal(t, 365*24*time.Hour, config.KeyExpiry)
	require.NotEmpty(t, config.KeyID)
}

func TestKeyExpiry(t *testing.T) {
	config := &PQCConfig{
		KeyID:     "expiry-test",
		EnablePQC: true,
		KeyExpiry: 1 * time.Hour,
	}

	signer, err := NewPQCSigner(config)
	require.NoError(t, err)

	// Check that keys are not expired initially
	require.False(t, signer.IsExpired())

	// Check expiry time
	require.True(t, signer.ExpiresAt().After(signer.CreatedAt()))
}

func TestAlgorithmConstants(t *testing.T) {
	require.Equal(t, "ML-KEM-768", AlgorithmMLKEM768)
	require.Equal(t, "Ed25519", AlgorithmEd25519)
	require.Equal(t, "Hybrid-ML-KEM+Ed25519", AlgorithmHybrid)
}

func TestKeySizes(t *testing.T) {
	require.Equal(t, 1184, MLKEMPublicKeySize)
	require.Equal(t, 1088, MLKEMCiphertextSize)
	require.Equal(t, 32, MLKEMSharedSecretSize)
	require.Equal(t, 32, Ed25519PublicKeySize)
	require.Equal(t, 64, Ed25519PrivateKeySize)
	require.Equal(t, 64, Ed25519SignatureSize)
}

func TestKeyPairFromSigner(t *testing.T) {
	signer, err := NewPQCSigner(nil)
	require.NoError(t, err)

	// Get Ed25519 key pair
	edKP := signer.KeyPairFromSigner(AlgorithmEd25519)
	require.NotNil(t, edKP)
	require.Equal(t, AlgorithmEd25519, edKP.Algorithm)

	// Get ML-KEM key pair
	mlKP := signer.KeyPairFromSigner(AlgorithmMLKEM768)
	require.NotNil(t, mlKP)
	require.Equal(t, AlgorithmMLKEM768, mlKP.Algorithm)

	// Invalid algorithm
	nilKP := signer.KeyPairFromSigner("invalid")
	require.Nil(t, nilKP)
}

// TestFIPS203Compliance verifies the implementation matches NIST FIPS 203.
func TestFIPS203Compliance(t *testing.T) {
	t.Run("ML-KEM-768 key sizes", func(t *testing.T) {
		signer, err := NewPQCSigner(nil)
		require.NoError(t, err)

		pubKey := signer.MLKEMPublicKey()
		require.Len(t, pubKey, 1184, "ML-KEM-768 encapsulation key must be 1184 bytes per FIPS 203")
	})

	t.Run("Shared secret size", func(t *testing.T) {
		signer, err := NewPQCSigner(nil)
		require.NoError(t, err)

		encap, err := signer.EncapsulateToSelf()
		require.NoError(t, err)
		require.Len(t, encap.SharedSecret, 32, "ML-KEM shared secret must be 32 bytes per FIPS 203")
	})

	t.Run("Ciphertext size", func(t *testing.T) {
		signer, err := NewPQCSigner(nil)
		require.NoError(t, err)

		encap, err := signer.EncapsulateToSelf()
		require.NoError(t, err)
		require.Len(t, encap.Ciphertext, 1088, "ML-KEM-768 ciphertext must be 1088 bytes per FIPS 203")
	})
}

// BenchmarkMLKEM768 benchmarks real FIPS 203 ML-KEM-768 operations.
func BenchmarkMLKEM768(b *testing.B) {
	signer, _ := NewPQCSigner(nil)
	pubKey := signer.MLKEMPublicKey()

	b.Run("KeyGen", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			NewPQCSigner(nil)
		}
	})

	b.Run("Encapsulate", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			signer.Encapsulate(pubKey)
		}
	})

	b.Run("Decapsulate", func(b *testing.B) {
		encap, _ := signer.EncapsulateToSelf()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			signer.Decapsulate(encap.Ciphertext)
		}
	})
}

// BenchmarkEd25519 benchmarks Ed25519 signing operations.
func BenchmarkEd25519(b *testing.B) {
	signer, _ := NewPQCSigner(nil)
	data := []byte("benchmark message for signing")

	b.Run("Sign", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			signer.Sign(data)
		}
	})

	b.Run("Verify", func(b *testing.B) {
		sig, _ := signer.Sign(data)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			signer.Verify(data, sig)
		}
	})
}
