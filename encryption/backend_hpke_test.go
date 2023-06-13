package encryption

import (
	"bytes"
	"crypto/rand"
	"testing"

	hpke "github.com/jedisct1/go-hpke-compact"
)

func TestHPKEBackend(t *testing.T) {
	suite, err := hpke.NewSuite(hpke.KemX25519HkdfSha256, hpke.KdfHkdfSha256, hpke.AeadChaCha20Poly1305)
	if err != nil {
		t.Fatalf("Failed to initialize HPKE suite: %v", err)
	}
	clientKp, err := suite.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate HPKE keypair: %v", err)
	}
	serverKp, err := suite.GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate HPKE keypair: %v", err)
	}
	psk := &hpke.Psk{
		Key: []byte("0123456789abcdef"),
		ID:  []byte("0123456789abcdef"),
	}

	// Initialize the encryption backend
	backend, err := NewHPKEBackend(clientKp.PublicKey, clientKp.SecretKey, serverKp.PublicKey, serverKp.SecretKey, psk.Key, psk.ID)
	if err != nil {
		t.Fatalf("Failed to initialize encryption backend: %v", err)
	}

	// Test small file encryption and decryption
	smallData := []byte("This is a small file.")
	encryptedSmallData, err := backend.Encrypt(smallData)
	if err != nil {
		t.Fatalf("Failed to encrypt small file: %v", err)
	}

	decryptedSmallData, err := backend.Decrypt(encryptedSmallData)
	if err != nil {
		t.Fatalf("Failed to decrypt small file: %v", err)
	}
	if !bytes.Equal(smallData, decryptedSmallData) {
		t.Fatal("Small file encryption and decryption failed: data mismatch")
	}

	// Test large file encryption and decryption
	largeData := make([]byte, 10*1024*1024) // 1 MB
	if _, err := rand.Read(largeData); err != nil {
		t.Fatalf("Failed to generate random data for large file: %v", err)
	}
	encryptedLargeData, err := backend.Encrypt(largeData)
	if err != nil {
		t.Fatalf("Failed to encrypt large file: %v", err)
	}
	decryptedLargeData, err := backend.Decrypt(encryptedLargeData)
	if err != nil {
		t.Fatalf("Failed to decrypt large file: %v", err)
	}
	if !bytes.Equal(largeData, decryptedLargeData) {
		t.Fatal("Large file encryption and decryption failed: data mismatch")
	}
}
