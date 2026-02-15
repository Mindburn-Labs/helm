//go:build go1.24

package pqc

import (
	"crypto/ed25519"
	"crypto/mlkem"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// Algorithm identifiers per HELM specification and NIST FIPS.
const (
	AlgorithmMLKEM768 = "ML-KEM-768" // NIST FIPS 203
	AlgorithmEd25519  = "Ed25519"    // RFC 8032
	AlgorithmHybrid   = "Hybrid-ML-KEM+Ed25519"

	// ML-KEM-768 sizes (NIST FIPS 203)
	MLKEMPublicKeySize    = 1184 // ML-KEM-768 encapsulation key bytes
	MLKEMPrivateKeySize   = 64   // ML-KEM-768 seed bytes (actual dk is larger internally)
	MLKEMCiphertextSize   = 1088 // ML-KEM-768 ciphertext bytes
	MLKEMSharedSecretSize = 32   // Shared secret bytes (crypto/mlkem.SharedKeySize)

	// Ed25519 sizes (RFC 8032)
	Ed25519PublicKeySize  = ed25519.PublicKeySize  // 32 bytes
	Ed25519PrivateKeySize = ed25519.PrivateKeySize // 64 bytes
	Ed25519SignatureSize  = ed25519.SignatureSize  // 64 bytes
)

// KeyPair represents a cryptographic key pair.
type KeyPair struct {
	PublicKey  []byte    `json:"public_key"`
	PrivateKey []byte    `json:"private_key"` // For ML-KEM, this is the seed
	Algorithm  string    `json:"algorithm"`
	KeyID      string    `json:"key_id"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at,omitempty"`
}

// Signature represents a cryptographic signature.
type Signature struct {
	Value     []byte    `json:"value"`
	Algorithm string    `json:"algorithm"`
	KeyID     string    `json:"key_id"`
	Timestamp time.Time `json:"timestamp"`
}

// EncapsulatedKey represents a key encapsulation result from ML-KEM.
type EncapsulatedKey struct {
	Ciphertext   []byte `json:"ciphertext"`
	SharedSecret []byte `json:"shared_secret"`
	Algorithm    string `json:"algorithm"`
}

// PQCSigner implements hybrid post-quantum cryptography.
// Uses ML-KEM-768 for key encapsulation and Ed25519 for signatures.
type PQCSigner struct {
	mu            sync.RWMutex
	mlkemDecapKey *mlkem.DecapsulationKey768 // Real ML-KEM-768 decapsulation key
	mlkemEncapKey *mlkem.EncapsulationKey768 // Real ML-KEM-768 encapsulation key
	mlkemSeed     []byte                     // Seed for ML-KEM key regeneration
	ed25519Pub    ed25519.PublicKey          // Real Ed25519 public key
	ed25519Priv   ed25519.PrivateKey         // Real Ed25519 private key
	enablePQC     bool
	keyID         string
	createdAt     time.Time
	expiresAt     time.Time
}

// PQCConfig configures the PQC signer.
type PQCConfig struct {
	KeyID     string
	EnablePQC bool
	KeyExpiry time.Duration
}

// DefaultPQCConfig returns production defaults.
func DefaultPQCConfig() *PQCConfig {
	return &PQCConfig{
		KeyID:     generateKeyID(),
		EnablePQC: true,
		KeyExpiry: 365 * 24 * time.Hour, // 1 year
	}
}

// NewPQCSigner creates a new hybrid PQC signer with real FIPS 203 ML-KEM.
func NewPQCSigner(config *PQCConfig) (*PQCSigner, error) {
	if config == nil {
		config = DefaultPQCConfig()
	}

	now := time.Now()
	signer := &PQCSigner{
		enablePQC: config.EnablePQC,
		keyID:     config.KeyID,
		createdAt: now,
		expiresAt: now.Add(config.KeyExpiry),
	}

	// Generate real Ed25519 keys (always available)
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("ed25519 keygen failed: %w", err)
	}
	signer.ed25519Pub = pub
	signer.ed25519Priv = priv

	if config.EnablePQC {
		// Generate real FIPS 203 ML-KEM-768 keys
		decapKey, err := mlkem.GenerateKey768()
		if err != nil {
			return nil, fmt.Errorf("ML-KEM-768 keygen failed: %w", err)
		}
		signer.mlkemDecapKey = decapKey
		signer.mlkemEncapKey = decapKey.EncapsulationKey()
		signer.mlkemSeed = decapKey.Bytes() // Store seed for serialization
	}

	return signer, nil
}

// NewPQCSignerFromKeys creates a signer from existing keys.
func NewPQCSignerFromKeys(ed25519Priv ed25519.PrivateKey, mlkemSeed []byte, keyID string) (*PQCSigner, error) {
	if len(ed25519Priv) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("invalid ed25519 private key size")
	}

	signer := &PQCSigner{
		enablePQC:   true,
		keyID:       keyID,
		createdAt:   time.Now(), // Metadata lost in simple rehydration, generic timestamp
		ed25519Priv: ed25519Priv,
		ed25519Pub:  ed25519Priv.Public().(ed25519.PublicKey),
	}

	if len(mlkemSeed) > 0 {
		// Rehydrate ML-KEM
		if len(mlkemSeed) != MLKEMPrivateKeySize { // 64 bytes for seed
			return nil, fmt.Errorf("invalid ML-KEM seed size: got %d, want %d", len(mlkemSeed), MLKEMPrivateKeySize)
		}

		decapKey, err := mlkem.NewDecapsulationKey768(mlkemSeed)
		if err != nil {
			return nil, fmt.Errorf("failed to parse ML-KEM key: %w", err)
		}
		signer.mlkemDecapKey = decapKey
		signer.mlkemEncapKey = decapKey.EncapsulationKey()
		signer.mlkemSeed = mlkemSeed
	} else {
		signer.enablePQC = false
	}

	return signer, nil
}

// Sign produces an Ed25519 signature using the real crypto/ed25519 package.
// Note: ML-KEM is for key encapsulation, not signing. Ed25519 is used for signatures.
func (s *PQCSigner) Sign(data []byte) (*Signature, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Use real Ed25519 signing
	sig := ed25519.Sign(s.ed25519Priv, data)

	return &Signature{
		Value:     sig,
		Algorithm: AlgorithmEd25519,
		KeyID:     s.keyID,
		Timestamp: time.Now(),
	}, nil
}

// SignWithContext produces an Ed25519 signature with domain separation.
func (s *PQCSigner) SignWithContext(data []byte, context string) (*Signature, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Create domain-separated message: H(context || data)
	h := sha256.New()
	h.Write([]byte(context))
	h.Write(data)
	message := h.Sum(nil)

	sig := ed25519.Sign(s.ed25519Priv, message)

	return &Signature{
		Value:     sig,
		Algorithm: AlgorithmEd25519,
		KeyID:     s.keyID,
		Timestamp: time.Now(),
	}, nil
}

// Verify verifies an Ed25519 signature using the real crypto/ed25519 package.
func (s *PQCSigner) Verify(data []byte, sig *Signature) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	switch sig.Algorithm {
	case AlgorithmEd25519, AlgorithmHybrid:
		return ed25519.Verify(s.ed25519Pub, data, sig.Value), nil
	default:
		return false, fmt.Errorf("unknown algorithm: %s", sig.Algorithm)
	}
}

// VerifyWithContext verifies an Ed25519 signature with domain separation.
func (s *PQCSigner) VerifyWithContext(data []byte, sig *Signature, context string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Recreate domain-separated message
	h := sha256.New()
	h.Write([]byte(context))
	h.Write(data)
	message := h.Sum(nil)

	return ed25519.Verify(s.ed25519Pub, message, sig.Value), nil
}

// Encapsulate performs real FIPS 203 ML-KEM-768 key encapsulation.
// This generates a shared secret and ciphertext that can be decapsulated
// by the holder of the corresponding decapsulation key.
func (s *PQCSigner) Encapsulate(recipientPubKey []byte) (*EncapsulatedKey, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.enablePQC {
		return nil, fmt.Errorf("PQC not enabled")
	}

	// Parse recipient's encapsulation key
	encapKey, err := mlkem.NewEncapsulationKey768(recipientPubKey)
	if err != nil {
		return nil, fmt.Errorf("invalid ML-KEM-768 encapsulation key: %w", err)
	}

	// Perform real FIPS 203 ML-KEM-768 encapsulation
	// Note: Go's Encapsulate returns (sharedKey, ciphertext)
	sharedSecret, ciphertext := encapKey.Encapsulate()

	return &EncapsulatedKey{
		Ciphertext:   ciphertext,
		SharedSecret: sharedSecret,
		Algorithm:    AlgorithmMLKEM768,
	}, nil
}

// EncapsulateToSelf performs ML-KEM encapsulation using our own encapsulation key.
// Useful for generating ephemeral shared secrets.
func (s *PQCSigner) EncapsulateToSelf() (*EncapsulatedKey, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.enablePQC || s.mlkemEncapKey == nil {
		return nil, fmt.Errorf("ML-KEM not enabled")
	}

	// Note: Go's Encapsulate returns (sharedKey, ciphertext)
	sharedSecret, ciphertext := s.mlkemEncapKey.Encapsulate()

	return &EncapsulatedKey{
		Ciphertext:   ciphertext,
		SharedSecret: sharedSecret,
		Algorithm:    AlgorithmMLKEM768,
	}, nil
}

// Decapsulate performs real FIPS 203 ML-KEM-768 key decapsulation.
// Returns the shared secret corresponding to the ciphertext.
func (s *PQCSigner) Decapsulate(ciphertext []byte) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.enablePQC || s.mlkemDecapKey == nil {
		return nil, fmt.Errorf("ML-KEM not enabled")
	}

	// Perform real FIPS 203 ML-KEM-768 decapsulation
	sharedSecret, err := s.mlkemDecapKey.Decapsulate(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("ML-KEM-768 decapsulation failed: %w", err)
	}

	return sharedSecret, nil
}

// PublicKeys returns all public keys for this signer.
func (s *PQCSigner) PublicKeys() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make(map[string]string)

	if s.ed25519Pub != nil {
		keys["ed25519"] = hex.EncodeToString(s.ed25519Pub)
	}
	if s.mlkemEncapKey != nil {
		keys["ml-kem-768"] = hex.EncodeToString(s.mlkemEncapKey.Bytes())
	}

	return keys
}

// Ed25519PublicKey returns the Ed25519 public key bytes.
func (s *PQCSigner) Ed25519PublicKey() []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return []byte(s.ed25519Pub)
}

// MLKEMPublicKey returns the ML-KEM-768 encapsulation key bytes.
func (s *PQCSigner) MLKEMPublicKey() []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.mlkemEncapKey == nil {
		return nil
	}
	return s.mlkemEncapKey.Bytes()
}

// KeyID returns the signer's key identifier.
func (s *PQCSigner) KeyID() string {
	return s.keyID
}

// IsPQCEnabled returns whether PQC algorithms are active.
func (s *PQCSigner) IsPQCEnabled() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.enablePQC && s.mlkemDecapKey != nil
}

// IsExpired returns whether the keys have expired.
func (s *PQCSigner) IsExpired() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return time.Now().After(s.expiresAt)
}

// CreatedAt returns when the signer was created.
func (s *PQCSigner) CreatedAt() time.Time {
	return s.createdAt
}

// ExpiresAt returns when the signer expires.
func (s *PQCSigner) ExpiresAt() time.Time {
	return s.expiresAt
}

// generateKeyID creates a secure random key identifier.
func generateKeyID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// Fallback to timestamp-based ID
		return fmt.Sprintf("key-%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)[:16]
}

// ---- Legacy Compatibility Types ----
// These maintain API compatibility with previous simulated implementation.

// KeyPairFromSigner extracts a KeyPair from the signer for serialization.
func (s *PQCSigner) KeyPairFromSigner(algorithm string) *KeyPair {
	s.mu.RLock()
	defer s.mu.RUnlock()

	switch algorithm {
	case AlgorithmEd25519:
		return &KeyPair{
			PublicKey:  []byte(s.ed25519Pub),
			PrivateKey: s.ed25519Priv.Seed(),
			Algorithm:  AlgorithmEd25519,
			KeyID:      s.keyID,
			CreatedAt:  s.createdAt,
			ExpiresAt:  s.expiresAt,
		}
	case AlgorithmMLKEM768:
		if s.mlkemEncapKey == nil {
			return nil
		}
		return &KeyPair{
			PublicKey:  s.mlkemEncapKey.Bytes(),
			PrivateKey: s.mlkemSeed,
			Algorithm:  AlgorithmMLKEM768,
			KeyID:      s.keyID,
			CreatedAt:  s.createdAt,
			ExpiresAt:  s.expiresAt,
		}
	default:
		return nil
	}
}
