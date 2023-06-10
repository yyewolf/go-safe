package internal

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestAESEncryptionBackend(t *testing.T) {
	// Generate a random key for testing
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("Failed to generate random key: %v", err)
	}

	// Initialize the encryption backend
	backend := AESEncryptionBackend{}
	if err := backend.Initialize(key); err != nil {
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
