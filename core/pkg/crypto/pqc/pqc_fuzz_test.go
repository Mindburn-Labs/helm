//go:build go1.24

package pqc

import (
	"testing"
)

// FuzzPQCSign tests the PQC signer with random inputs.
// Run: go test -fuzz=FuzzPQCSign -fuzztime=30s ./core/pkg/crypto/pqc/
func FuzzPQCSign(f *testing.F) {
	// Seed corpus
	f.Add([]byte("hello world"))
	f.Add([]byte(""))
	f.Add([]byte{0x00, 0x01, 0x02, 0x03})
	f.Add([]byte("The quick brown fox jumps over the lazy dog"))
	f.Add(make([]byte, 1024))       // Large input
	f.Add([]byte{0xff, 0xfe, 0xfd}) // Binary data

	f.Fuzz(func(t *testing.T, data []byte) {
		signer, err := NewPQCSigner(nil)
		if err != nil {
			t.Fatalf("failed to create signer: %v", err)
		}

		// Sign should not panic
		sig, err := signer.Sign(data)
		if err != nil {
			t.Fatalf("sign failed: %v", err)
		}

		// Verify should succeed
		valid, err := signer.Verify(data, sig)
		if err != nil {
			t.Fatalf("verify failed: %v", err)
		}
		if !valid {
			t.Error("signature verification failed")
		}

		// Modifying data should invalidate signature
		if len(data) > 0 {
			modifiedData := make([]byte, len(data))
			copy(modifiedData, data)
			modifiedData[0] ^= 0xff // Flip bits
			valid, err := signer.Verify(modifiedData, sig)
			if err == nil && valid {
				t.Error("modified data should not verify")
			}
		}
	})
}

// FuzzPQCSignWithContext tests context-aware signing.
func FuzzPQCSignWithContext(f *testing.F) {
	f.Add([]byte("test data"), "context1")
	f.Add([]byte("another test"), "")
	f.Add([]byte{0x00}, "special\ncontext")

	f.Fuzz(func(t *testing.T, data []byte, context string) {
		signer, err := NewPQCSigner(nil)
		if err != nil {
			t.Fatalf("failed to create signer: %v", err)
		}

		// Sign with context
		sig, err := signer.SignWithContext(data, context)
		if err != nil {
			// Some contexts might be rejected
			return
		}

		// Verify with same context should succeed
		valid, err := signer.VerifyWithContext(data, sig, context)
		if err != nil {
			t.Fatalf("verify with context failed: %v", err)
		}
		if !valid {
			t.Error("context signature verification failed")
		}

		// Verify with different context should fail
		if context != "" {
			valid, _ := signer.VerifyWithContext(data, sig, "wrong-context")
			if valid {
				t.Error("different context should not verify")
			}
		}
	})
}

// FuzzMLKEMEncapsulation tests ML-KEM key encapsulation.
func FuzzMLKEMEncapsulation(f *testing.F) {
	// Seed with public key bytes
	signer, _ := NewPQCSigner(nil)
	pubKey := signer.MLKEMPublicKey()
	f.Add(pubKey)

	f.Fuzz(func(t *testing.T, randomPubKey []byte) {
		sender, err := NewPQCSigner(nil)
		if err != nil {
			t.Fatalf("failed to create sender: %v", err)
		}

		// Try encapsulation - may fail for invalid public keys
		encap, err := sender.Encapsulate(randomPubKey)
		if err != nil {
			// Expected for malformed keys
			return
		}

		// If encapsulation succeeded with valid key, check basic properties
		if len(encap.Ciphertext) != MLKEMCiphertextSize {
			t.Errorf("unexpected ciphertext size: %d", len(encap.Ciphertext))
		}
		if len(encap.SharedSecret) != MLKEMSharedSecretSize {
			t.Errorf("unexpected shared secret size: %d", len(encap.SharedSecret))
		}
	})
}

// FuzzMLKEMDecapsulation tests ML-KEM decapsulation with random ciphertext.
func FuzzMLKEMDecapsulation(f *testing.F) {
	// Generate valid ciphertext for seed
	signer, err := NewPQCSigner(nil)
	if err != nil {
		f.Logf("failed to create signer for seed: %v", err)
		return
	}
	encap, err := signer.EncapsulateToSelf()
	if err != nil {
		f.Logf("failed to encapsulate for seed: %v", err)
		return
	}
	f.Add(encap.Ciphertext)

	f.Fuzz(func(t *testing.T, ciphertext []byte) {
		receiver, err := NewPQCSigner(nil)
		if err != nil {
			t.Fatalf("failed to create receiver: %v", err)
		}

		// Decapsulation may fail for malformed ciphertext
		_, err = receiver.Decapsulate(ciphertext)
		if err != nil {
			// Expected for malformed ciphertext
			return
		}
		// If it succeeds, the shared secret was computed
		// (may not match any sender's because ciphertext is random)
	})
}
