package encryption

import (
	"bytes"
	"crypto/rand"
	"testing"

	ecies "github.com/yyewolf/go-ecies/v2"
)

func TestECIESEncryptionBackend(t *testing.T) {
	priv, err := ecies.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate ECIES key: %v", err)
	}
	pub := priv.PublicKey

	// Initialize the encryption backend
	backend, err := NewEciesEncryptionBackend(pub.Hex(), priv.Hex())
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
