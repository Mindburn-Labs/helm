package hsm

import (
	"context"
	"crypto"
	"crypto/ed25519"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPKCS11Provider(t *testing.T) {
	cfg := PKCS11Config{
		LibraryPath: "/usr/lib/softhsm/libsofthsm2.so",
		SlotID:      0,
		PIN:         "1234",
		TokenLabel:  "test-token",
	}

	provider, err := NewPKCS11Provider(cfg)
	require.NoError(t, err)
	require.NotNil(t, provider)

	ctx := context.Background()

	t.Run("OpenClose", func(t *testing.T) {
		require.NoError(t, provider.Open(ctx))
		require.True(t, provider.IsOpen())

		require.NoError(t, provider.Close())
		require.False(t, provider.IsOpen())
	})

	t.Run("KeyLifecycle", func(t *testing.T) {
		require.NoError(t, provider.Open(ctx))
		defer func() { _ = provider.Close() }()

		// Generate key
		handle, err := provider.GenerateKey(ctx, KeyGenOpts{
			Algorithm:   AlgorithmECDSAP256,
			Label:       "test-key",
			Extractable: false,
			Sensitive:   true,
			Usage:       KeyUsageSign | KeyUsageVerify,
			ExpiresIn:   24 * time.Hour,
		})
		require.NoError(t, err)
		require.NotEmpty(t, handle)

		// Get key info
		info, err := provider.GetKeyInfo(ctx, handle)
		require.NoError(t, err)
		require.Equal(t, "test-key", info.Label)
		require.Equal(t, AlgorithmECDSAP256, info.Algorithm)
		require.NotNil(t, info.ExpiresAt)

		// List keys
		keys, err := provider.ListKeys(ctx)
		require.NoError(t, err)
		require.Len(t, keys, 1)

		// Delete key
		require.NoError(t, provider.DeleteKey(ctx, handle))

		// Verify deletion
		_, err = provider.GetKeyInfo(ctx, handle)
		require.Equal(t, ErrKeyNotFound, err)
	})

	t.Run("SignVerify", func(t *testing.T) {
		require.NoError(t, provider.Open(ctx))
		defer func() { _ = provider.Close() }()

		handle, err := provider.GenerateKey(ctx, KeyGenOpts{
			Algorithm: AlgorithmECDSAP256,
			Label:     "sign-key",
			Usage:     KeyUsageSign | KeyUsageVerify,
		})
		require.NoError(t, err)

		digest := []byte("test message digest")

		// Sign
		sig, err := provider.Sign(ctx, handle, digest, SignOpts{
			HashAlgorithm: crypto.SHA256,
		})
		require.NoError(t, err)
		require.NotEmpty(t, sig)

		// Verify
		valid, err := provider.Verify(ctx, handle, digest, sig)
		require.NoError(t, err)
		require.True(t, valid)
	})

	t.Run("EncryptDecrypt", func(t *testing.T) {
		require.NoError(t, provider.Open(ctx))
		defer func() { _ = provider.Close() }()

		handle, err := provider.GenerateKey(ctx, KeyGenOpts{
			Algorithm: AlgorithmRSA2048,
			Label:     "enc-key",
			Usage:     KeyUsageEncrypt | KeyUsageDecrypt,
		})
		require.NoError(t, err)

		plaintext := []byte("secret data")

		// Encrypt
		ciphertext, err := provider.Encrypt(ctx, handle, plaintext)
		require.NoError(t, err)
		require.NotEqual(t, plaintext, ciphertext)

		// Decrypt
		decrypted, err := provider.Decrypt(ctx, handle, ciphertext)
		require.NoError(t, err)
		require.Equal(t, plaintext, decrypted)
	})

	t.Run("KeyWrapUnwrap", func(t *testing.T) {
		require.NoError(t, provider.Open(ctx))
		defer func() { _ = provider.Close() }()

		// Create wrapping key
		wrapKey, err := provider.GenerateKey(ctx, KeyGenOpts{
			Algorithm: AlgorithmRSA4096,
			Label:     "wrap-key",
			Usage:     KeyUsageWrap | KeyUsageUnwrap,
		})
		require.NoError(t, err)

		// Create key to wrap (must be extractable)
		targetKey, err := provider.GenerateKey(ctx, KeyGenOpts{
			Algorithm:   AlgorithmECDSAP256,
			Label:       "target-key",
			Extractable: true,
			Usage:       KeyUsageSign,
		})
		require.NoError(t, err)

		// Wrap
		wrapped, err := provider.WrapKey(ctx, targetKey, wrapKey)
		require.NoError(t, err)
		require.NotEmpty(t, wrapped)

		// Unwrap
		unwrapped, err := provider.UnwrapKey(ctx, wrapped, wrapKey, ImportOpts{
			Label: "unwrapped-key",
			Usage: KeyUsageSign,
		})
		require.NoError(t, err)
		require.NotEmpty(t, unwrapped)
	})

	t.Run("NotInitialized", func(t *testing.T) {
		// Close if open
		_ = provider.Close()

		_, err := provider.GenerateKey(ctx, KeyGenOpts{})
		require.Equal(t, ErrNotInitialized, err)

		_, err = provider.ListKeys(ctx)
		require.Equal(t, ErrNotInitialized, err)
	})

	t.Run("Capabilities", func(t *testing.T) {
		caps := provider.Capabilities()
		require.False(t, caps.SupportsPQC) // PKCS#11 doesn't support PQC yet
		require.True(t, caps.SupportsKeyWrap)
		require.Equal(t, 3, caps.FIPSLevel)
		require.Contains(t, caps.SupportedAlgorithms, AlgorithmRSA2048)
	})

	t.Run("ProviderInfo", func(t *testing.T) {
		require.Equal(t, "PKCS#11", provider.Name())
		require.NotEmpty(t, provider.Version())
	})
}

func TestSoftwareProvider(t *testing.T) {
	provider := NewSoftwareProvider()
	ctx := context.Background()

	t.Run("FullLifecycle", func(t *testing.T) {
		require.NoError(t, provider.Open(ctx))
		defer func() { _ = provider.Close() }()

		require.NotEmpty(t, provider.Version())
		require.True(t, provider.IsOpen())

		// Generate
		handle, err := provider.GenerateKey(ctx, KeyGenOpts{
			Algorithm: AlgorithmMLKEM768, // PQC supported in software
			Label:     "pqc-key",
			Usage:     KeyUsageEncrypt,
		})
		require.NoError(t, err)

		// List keys
		keys, err := provider.ListKeys(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, keys)

		// Import/Delete (stub implementation, but must remain reachable for interface coverage)
		priv := ed25519.NewKeyFromSeed(make([]byte, 32))
		imported, err := provider.ImportKey(ctx, crypto.PrivateKey(priv), ImportOpts{Label: "imported-key", Usage: KeyUsageSign})
		require.NoError(t, err)
		require.NotEmpty(t, imported)
		require.NoError(t, provider.DeleteKey(ctx, imported))

		// Wrap/Unwrap
		wrapped, err := provider.WrapKey(ctx, handle, handle)
		require.NoError(t, err)
		require.NotEmpty(t, wrapped)

		unwrapped, err := provider.UnwrapKey(ctx, wrapped, handle, ImportOpts{Label: "unwrapped-key", Usage: KeyUsageEncrypt})
		require.NoError(t, err)
		require.NotEmpty(t, unwrapped)

		// Get info
		info, err := provider.GetKeyInfo(ctx, handle)
		require.NoError(t, err)
		require.Equal(t, AlgorithmMLKEM768, info.Algorithm)
		require.True(t, info.Extractable) // Software keys are always extractable

		// Sign/Verify
		digest := []byte("test")
		sig, err := provider.Sign(ctx, handle, digest, SignOpts{})
		require.NoError(t, err)

		valid, err := provider.Verify(ctx, handle, digest, sig)
		require.NoError(t, err)
		require.True(t, valid)

		// Encrypt/Decrypt
		plaintext := []byte("secret")
		ciphertext, err := provider.Encrypt(ctx, handle, plaintext)
		require.NoError(t, err)

		decrypted, err := provider.Decrypt(ctx, handle, ciphertext)
		require.NoError(t, err)
		require.Equal(t, plaintext, decrypted)
	})

	t.Run("Capabilities", func(t *testing.T) {
		caps := provider.Capabilities()
		require.True(t, caps.SupportsPQC) // Software supports PQC
		require.Equal(t, 0, caps.FIPSLevel)
		require.Contains(t, caps.SupportedAlgorithms, AlgorithmMLKEM768)
	})

	t.Run("ProviderInfo", func(t *testing.T) {
		require.Contains(t, provider.Name(), "Software")
	})
}

func TestAlgorithmString(t *testing.T) {
	tests := []struct {
		algo Algorithm
		want string
	}{
		{AlgorithmRSA2048, "RSA-2048"},
		{AlgorithmRSA4096, "RSA-4096"},
		{AlgorithmECDSAP256, "ECDSA-P256"},
		{AlgorithmECDSAP384, "ECDSA-P384"},
		{AlgorithmEd25519, "Ed25519"},
		{AlgorithmMLKEM768, "ML-KEM-768"},
		{AlgorithmSLHDSA128, "SLH-DSA-128f"},
		{Algorithm(99), "Unknown"},
	}

	for _, tt := range tests {
		require.Equal(t, tt.want, tt.algo.String())
	}
}

func TestKeyUsageFlags(t *testing.T) {
	usage := KeyUsageSign | KeyUsageVerify | KeyUsageWrap

	require.True(t, usage&KeyUsageSign != 0)
	require.True(t, usage&KeyUsageVerify != 0)
	require.True(t, usage&KeyUsageWrap != 0)
	require.False(t, usage&KeyUsageDecrypt != 0)
}

func TestErrors(t *testing.T) {
	require.Error(t, ErrNotInitialized)
	require.Error(t, ErrSessionClosed)
	require.Error(t, ErrKeyNotFound)
	require.Error(t, ErrOperationFailed)
}
