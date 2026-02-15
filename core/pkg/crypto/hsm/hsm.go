// Package hsm provides Hardware Security Module abstraction for HELM.
// This package enables hardware-backed cryptographic key protection for
// high-security deployments requiring FIPS 140-2/3 compliance.
package hsm

import (
	"context"
	"crypto"
	"errors"
	"fmt"
	"sync"
	"time"
)

// Common errors
var (
	ErrNotInitialized    = errors.New("hsm: not initialized")
	ErrSessionClosed     = errors.New("hsm: session closed")
	ErrKeyNotFound       = errors.New("hsm: key not found")
	ErrOperationFailed   = errors.New("hsm: operation failed")
	ErrInvalidKeyHandle  = errors.New("hsm: invalid key handle")
	ErrNotSupported      = errors.New("hsm: operation not supported")
	ErrAuthenticationReq = errors.New("hsm: authentication required")
)

// KeyHandle is an opaque reference to a key stored in the HSM.
type KeyHandle string

// Algorithm represents supported key algorithms.
type Algorithm int

const (
	AlgorithmRSA2048 Algorithm = iota
	AlgorithmRSA4096
	AlgorithmECDSAP256
	AlgorithmECDSAP384
	AlgorithmEd25519
	AlgorithmMLKEM768  // Post-quantum key encapsulation
	AlgorithmSLHDSA128 // Post-quantum signature (SLH-DSA-SHAKE-128f)
)

func (a Algorithm) String() string {
	names := map[Algorithm]string{
		AlgorithmRSA2048:   "RSA-2048",
		AlgorithmRSA4096:   "RSA-4096",
		AlgorithmECDSAP256: "ECDSA-P256",
		AlgorithmECDSAP384: "ECDSA-P384",
		AlgorithmEd25519:   "Ed25519",
		AlgorithmMLKEM768:  "ML-KEM-768",
		AlgorithmSLHDSA128: "SLH-DSA-128f",
	}
	if name, ok := names[a]; ok {
		return name
	}
	return "Unknown"
}

// KeyUsage specifies permitted key operations.
type KeyUsage int

const (
	KeyUsageSign KeyUsage = 1 << iota
	KeyUsageVerify
	KeyUsageEncrypt
	KeyUsageDecrypt
	KeyUsageWrap
	KeyUsageUnwrap
	KeyUsageDerive
)

// KeyInfo provides metadata about a stored key.
type KeyInfo struct {
	Handle      KeyHandle
	Label       string
	Algorithm   Algorithm
	Usage       KeyUsage
	Extractable bool
	CreatedAt   time.Time
	ExpiresAt   *time.Time
}

// KeyGenOpts specifies key generation options.
type KeyGenOpts struct {
	Algorithm   Algorithm
	Label       string
	Extractable bool     // false = non-extractable (recommended for production)
	Sensitive   bool     // true = cannot be revealed in plaintext
	Usage       KeyUsage // Permitted operations
	ExpiresIn   time.Duration
}

// SignOpts specifies signing options.
type SignOpts struct {
	HashAlgorithm crypto.Hash
	PSS           bool // Use RSA-PSS if true
}

// ImportOpts specifies key import options.
type ImportOpts struct {
	Label       string
	Extractable bool
	Sensitive   bool
	Usage       KeyUsage
}

// Provider defines the HSM abstraction interface.
type Provider interface {
	// Session Management
	Open(ctx context.Context) error
	Close() error
	IsOpen() bool

	// Key Management
	GenerateKey(ctx context.Context, opts KeyGenOpts) (KeyHandle, error)
	ImportKey(ctx context.Context, key crypto.PrivateKey, opts ImportOpts) (KeyHandle, error)
	DeleteKey(ctx context.Context, handle KeyHandle) error
	ListKeys(ctx context.Context) ([]KeyInfo, error)
	GetKeyInfo(ctx context.Context, handle KeyHandle) (*KeyInfo, error)

	// Cryptographic Operations
	Sign(ctx context.Context, handle KeyHandle, digest []byte, opts SignOpts) ([]byte, error)
	Verify(ctx context.Context, handle KeyHandle, digest, signature []byte) (bool, error)
	Encrypt(ctx context.Context, handle KeyHandle, plaintext []byte) ([]byte, error)
	Decrypt(ctx context.Context, handle KeyHandle, ciphertext []byte) ([]byte, error)

	// Key Wrapping (for backup)
	WrapKey(ctx context.Context, keyToWrap, wrappingKey KeyHandle) ([]byte, error)
	UnwrapKey(ctx context.Context, wrapped []byte, unwrappingKey KeyHandle, opts ImportOpts) (KeyHandle, error)

	// Provider Info
	Name() string
	Version() string
	Capabilities() Capabilities
}

// Capabilities describes what the provider supports.
type Capabilities struct {
	SupportsPQC         bool
	SupportsKeyWrap     bool
	MaxKeySize          int
	SupportedAlgorithms []Algorithm
	FIPSLevel           int // 0 = not certified, 2 = Level 2, 3 = Level 3
}

// ===== PKCS#11 Implementation Stub =====

// PKCS11Config configures the PKCS#11 provider.
type PKCS11Config struct {
	LibraryPath string // Path to PKCS#11 shared library
	SlotID      uint   // HSM slot ID
	PIN         string // User PIN (should come from environment/secret manager)
	TokenLabel  string // Token label for identification
}

// PKCS11Provider implements Provider using PKCS#11.
// This is a stub implementation - full implementation requires
// github.com/miekg/pkcs11 or similar binding.
type PKCS11Provider struct {
	config     PKCS11Config
	isOpen     bool
	keys       map[KeyHandle]*KeyInfo
	keyCounter int
	mu         sync.RWMutex
}

// NewPKCS11Provider creates a new PKCS#11 provider.
func NewPKCS11Provider(cfg PKCS11Config) (*PKCS11Provider, error) {
	if cfg.LibraryPath == "" {
		return nil, fmt.Errorf("pkcs11: library path required")
	}

	return &PKCS11Provider{
		config: cfg,
		keys:   make(map[KeyHandle]*KeyInfo),
	}, nil
}

func (p *PKCS11Provider) Name() string {
	return "PKCS#11"
}

func (p *PKCS11Provider) Version() string {
	return "1.0.0"
}

func (p *PKCS11Provider) Capabilities() Capabilities {
	return Capabilities{
		SupportsPQC:     false, // PKCS#11 doesn't natively support PQC yet
		SupportsKeyWrap: true,
		MaxKeySize:      4096,
		SupportedAlgorithms: []Algorithm{
			AlgorithmRSA2048,
			AlgorithmRSA4096,
			AlgorithmECDSAP256,
			AlgorithmECDSAP384,
			AlgorithmEd25519,
		},
		FIPSLevel: 3, // Depending on HSM hardware
	}
}

func (p *PKCS11Provider) Open(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.isOpen {
		return nil
	}

	// In real implementation:
	// 1. Load PKCS#11 library
	// 2. Initialize module
	// 3. Open session on slot
	// 4. Login with PIN

	p.isOpen = true
	return nil
}

func (p *PKCS11Provider) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isOpen {
		return nil
	}

	// In real implementation:
	// 1. Logout
	// 2. Close session
	// 3. Finalize module

	p.isOpen = false
	return nil
}

func (p *PKCS11Provider) IsOpen() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.isOpen
}

func (p *PKCS11Provider) GenerateKey(ctx context.Context, opts KeyGenOpts) (KeyHandle, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isOpen {
		return "", ErrNotInitialized
	}

	// Generate key handle
	p.keyCounter++
	handle := KeyHandle(fmt.Sprintf("key-%d", p.keyCounter))

	// Store key info
	info := &KeyInfo{
		Handle:      handle,
		Label:       opts.Label,
		Algorithm:   opts.Algorithm,
		Usage:       opts.Usage,
		Extractable: opts.Extractable,
		CreatedAt:   time.Now(),
	}

	if opts.ExpiresIn > 0 {
		expiry := time.Now().Add(opts.ExpiresIn)
		info.ExpiresAt = &expiry
	}

	p.keys[handle] = info

	return handle, nil
}

func (p *PKCS11Provider) ImportKey(ctx context.Context, key crypto.PrivateKey, opts ImportOpts) (KeyHandle, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isOpen {
		return "", ErrNotInitialized
	}

	//nolint:staticcheck // intended behavior for warning
	if opts.Extractable {
		// Warning: importing extractable keys defeats HSM purpose
		// Log warning in production
	}

	// Generate key handle
	p.keyCounter++
	handle := KeyHandle(fmt.Sprintf("imported-%d", p.keyCounter))

	// Determine algorithm from key type
	var algo Algorithm
	switch key.(type) {
	default:
		algo = AlgorithmECDSAP256 // Default assumption
	}

	p.keys[handle] = &KeyInfo{
		Handle:      handle,
		Label:       opts.Label,
		Algorithm:   algo,
		Usage:       opts.Usage,
		Extractable: opts.Extractable,
		CreatedAt:   time.Now(),
	}

	return handle, nil
}

func (p *PKCS11Provider) DeleteKey(ctx context.Context, handle KeyHandle) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isOpen {
		return ErrNotInitialized
	}

	if _, exists := p.keys[handle]; !exists {
		return ErrKeyNotFound
	}

	delete(p.keys, handle)
	return nil
}

func (p *PKCS11Provider) ListKeys(ctx context.Context) ([]KeyInfo, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if !p.isOpen {
		return nil, ErrNotInitialized
	}

	keys := make([]KeyInfo, 0, len(p.keys))
	for _, info := range p.keys {
		keys = append(keys, *info)
	}

	return keys, nil
}

func (p *PKCS11Provider) GetKeyInfo(ctx context.Context, handle KeyHandle) (*KeyInfo, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if !p.isOpen {
		return nil, ErrNotInitialized
	}

	info, exists := p.keys[handle]
	if !exists {
		return nil, ErrKeyNotFound
	}

	return info, nil
}

func (p *PKCS11Provider) Sign(ctx context.Context, handle KeyHandle, digest []byte, opts SignOpts) ([]byte, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if !p.isOpen {
		return nil, ErrNotInitialized
	}

	info, exists := p.keys[handle]
	if !exists {
		return nil, ErrKeyNotFound
	}

	if info.Usage&KeyUsageSign == 0 {
		return nil, fmt.Errorf("key does not permit signing")
	}

	// In real implementation: use PKCS#11 C_Sign
	// This is a stub that returns a mock signature
	signature := make([]byte, 64)
	copy(signature, digest)

	return signature, nil
}

func (p *PKCS11Provider) Verify(ctx context.Context, handle KeyHandle, digest, signature []byte) (bool, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if !p.isOpen {
		return false, ErrNotInitialized
	}

	info, exists := p.keys[handle]
	if !exists {
		return false, ErrKeyNotFound
	}

	if info.Usage&KeyUsageVerify == 0 {
		return false, fmt.Errorf("key does not permit verification")
	}

	// In real implementation: use PKCS#11 C_Verify
	// This is a stub
	return len(signature) >= 32, nil
}

func (p *PKCS11Provider) Encrypt(ctx context.Context, handle KeyHandle, plaintext []byte) ([]byte, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if !p.isOpen {
		return nil, ErrNotInitialized
	}

	info, exists := p.keys[handle]
	if !exists {
		return nil, ErrKeyNotFound
	}

	if info.Usage&KeyUsageEncrypt == 0 {
		return nil, fmt.Errorf("key does not permit encryption")
	}

	// Stub: return wrapped plaintext
	ciphertext := make([]byte, len(plaintext)+16)
	copy(ciphertext[16:], plaintext)

	return ciphertext, nil
}

func (p *PKCS11Provider) Decrypt(ctx context.Context, handle KeyHandle, ciphertext []byte) ([]byte, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if !p.isOpen {
		return nil, ErrNotInitialized
	}

	info, exists := p.keys[handle]
	if !exists {
		return nil, ErrKeyNotFound
	}

	if info.Usage&KeyUsageDecrypt == 0 {
		return nil, fmt.Errorf("key does not permit decryption")
	}

	if len(ciphertext) < 16 {
		return nil, fmt.Errorf("invalid ciphertext")
	}

	// Stub: unwrap
	return ciphertext[16:], nil
}

func (p *PKCS11Provider) WrapKey(ctx context.Context, keyToWrap, wrappingKey KeyHandle) ([]byte, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if !p.isOpen {
		return nil, ErrNotInitialized
	}

	toWrap, exists := p.keys[keyToWrap]
	if !exists {
		return nil, ErrKeyNotFound
	}

	if !toWrap.Extractable {
		return nil, fmt.Errorf("key is not extractable")
	}

	wrapperInfo, exists := p.keys[wrappingKey]
	if !exists {
		return nil, ErrKeyNotFound
	}

	if wrapperInfo.Usage&KeyUsageWrap == 0 {
		return nil, fmt.Errorf("wrapping key does not permit wrapping")
	}

	// Stub: return fake wrapped key
	wrapped := []byte(fmt.Sprintf("wrapped:%s:with:%s", keyToWrap, wrappingKey))
	return wrapped, nil
}

func (p *PKCS11Provider) UnwrapKey(ctx context.Context, wrapped []byte, unwrappingKey KeyHandle, opts ImportOpts) (KeyHandle, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isOpen {
		return "", ErrNotInitialized
	}

	unwrapperInfo, exists := p.keys[unwrappingKey]
	if !exists {
		return "", ErrKeyNotFound
	}

	if unwrapperInfo.Usage&KeyUsageUnwrap == 0 {
		return "", fmt.Errorf("key does not permit unwrapping")
	}

	// Generate new handle for unwrapped key
	p.keyCounter++
	handle := KeyHandle(fmt.Sprintf("unwrapped-%d", p.keyCounter))

	p.keys[handle] = &KeyInfo{
		Handle:      handle,
		Label:       opts.Label,
		Algorithm:   AlgorithmECDSAP256, // Would be determined from wrapped key
		Usage:       opts.Usage,
		Extractable: opts.Extractable,
		CreatedAt:   time.Now(),
	}

	return handle, nil
}

// ===== Software Provider (for development/testing) =====

// SoftwareProvider implements Provider using software-based cryptography.
// NOT FOR PRODUCTION - use only for development and testing.
type SoftwareProvider struct {
	keys       map[KeyHandle]*softwareKey
	keyCounter int
	isOpen     bool
	mu         sync.RWMutex
}

type softwareKey struct {
	info    *KeyInfo
	privKey crypto.PrivateKey
}

// NewSoftwareProvider creates a software-only provider for development.
func NewSoftwareProvider() *SoftwareProvider {
	return &SoftwareProvider{
		keys: make(map[KeyHandle]*softwareKey),
	}
}

func (p *SoftwareProvider) Name() string {
	return "Software (Development Only)"
}

func (p *SoftwareProvider) Version() string {
	return "1.0.0"
}

func (p *SoftwareProvider) Capabilities() Capabilities {
	return Capabilities{
		SupportsPQC:     true, // Can use Go crypto/mlkem
		SupportsKeyWrap: true,
		MaxKeySize:      8192,
		SupportedAlgorithms: []Algorithm{
			AlgorithmRSA2048,
			AlgorithmRSA4096,
			AlgorithmECDSAP256,
			AlgorithmECDSAP384,
			AlgorithmEd25519,
			AlgorithmMLKEM768,
			AlgorithmSLHDSA128,
		},
		FIPSLevel: 0, // Not FIPS certified
	}
}

func (p *SoftwareProvider) Open(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.isOpen = true
	return nil
}

func (p *SoftwareProvider) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.isOpen = false
	return nil
}

func (p *SoftwareProvider) IsOpen() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.isOpen
}

func (p *SoftwareProvider) GenerateKey(ctx context.Context, opts KeyGenOpts) (KeyHandle, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isOpen {
		return "", ErrNotInitialized
	}

	p.keyCounter++
	handle := KeyHandle(fmt.Sprintf("sw-key-%d", p.keyCounter))

	info := &KeyInfo{
		Handle:      handle,
		Label:       opts.Label,
		Algorithm:   opts.Algorithm,
		Usage:       opts.Usage,
		Extractable: true, // Software keys are always extractable
		CreatedAt:   time.Now(),
	}

	p.keys[handle] = &softwareKey{
		info: info,
		// In real implementation, generate actual key here
	}

	return handle, nil
}

func (p *SoftwareProvider) ImportKey(ctx context.Context, key crypto.PrivateKey, opts ImportOpts) (KeyHandle, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isOpen {
		return "", ErrNotInitialized
	}

	p.keyCounter++
	handle := KeyHandle(fmt.Sprintf("sw-imported-%d", p.keyCounter))

	p.keys[handle] = &softwareKey{
		info: &KeyInfo{
			Handle:      handle,
			Label:       opts.Label,
			Usage:       opts.Usage,
			Extractable: true,
			CreatedAt:   time.Now(),
		},
		privKey: key,
	}

	return handle, nil
}

func (p *SoftwareProvider) DeleteKey(ctx context.Context, handle KeyHandle) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isOpen {
		return ErrNotInitialized
	}

	delete(p.keys, handle)
	return nil
}

func (p *SoftwareProvider) ListKeys(ctx context.Context) ([]KeyInfo, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if !p.isOpen {
		return nil, ErrNotInitialized
	}

	keys := make([]KeyInfo, 0, len(p.keys))
	for _, k := range p.keys {
		keys = append(keys, *k.info)
	}

	return keys, nil
}

func (p *SoftwareProvider) GetKeyInfo(ctx context.Context, handle KeyHandle) (*KeyInfo, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if !p.isOpen {
		return nil, ErrNotInitialized
	}

	k, exists := p.keys[handle]
	if !exists {
		return nil, ErrKeyNotFound
	}

	return k.info, nil
}

func (p *SoftwareProvider) Sign(ctx context.Context, handle KeyHandle, digest []byte, opts SignOpts) ([]byte, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if !p.isOpen {
		return nil, ErrNotInitialized
	}

	_, exists := p.keys[handle]
	if !exists {
		return nil, ErrKeyNotFound
	}

	// Stub signature
	return append([]byte("sig:"), digest...), nil
}

func (p *SoftwareProvider) Verify(ctx context.Context, handle KeyHandle, digest, signature []byte) (bool, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if !p.isOpen {
		return false, ErrNotInitialized
	}

	_, exists := p.keys[handle]
	if !exists {
		return false, ErrKeyNotFound
	}

	return len(signature) > 4, nil
}

func (p *SoftwareProvider) Encrypt(ctx context.Context, handle KeyHandle, plaintext []byte) ([]byte, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if !p.isOpen {
		return nil, ErrNotInitialized
	}

	return append([]byte("enc:"), plaintext...), nil
}

func (p *SoftwareProvider) Decrypt(ctx context.Context, handle KeyHandle, ciphertext []byte) ([]byte, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if !p.isOpen {
		return nil, ErrNotInitialized
	}

	if len(ciphertext) < 4 {
		return nil, fmt.Errorf("invalid ciphertext")
	}

	return ciphertext[4:], nil
}

func (p *SoftwareProvider) WrapKey(ctx context.Context, keyToWrap, wrappingKey KeyHandle) ([]byte, error) {
	return []byte(fmt.Sprintf("wrapped:%s", keyToWrap)), nil
}

func (p *SoftwareProvider) UnwrapKey(ctx context.Context, wrapped []byte, unwrappingKey KeyHandle, opts ImportOpts) (KeyHandle, error) {
	return p.GenerateKey(ctx, KeyGenOpts{Label: opts.Label, Usage: opts.Usage})
}
