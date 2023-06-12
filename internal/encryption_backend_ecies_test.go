package internal

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestECIESEncryptionBackend(t *testing.T) {
	pub := []byte{'\x04', '\x05', '\xa0', '\xa3', '\x0c', '\xa8', '\xd7', '\x1f', '\x7c', '\xfd', '\x54', '\xd4', '\x55', '\xa9', '\xdd', '\x5c', '\x27', '\x23', '\xdf', '\x9c', '\xcf', '\x9b', '\xb1', '\xf8', '\x95', '\x85', '\x39', '\x1f', '\x1d', '\x08', '\xf9', '\xb2', '\x77', '\x6c', '\x52', '\x0b', '\x8a', '\xaa', '\x39', '\xc0', '\x53', '\x66', '\xc0', '\xd7', '\x52', '\x5e', '\xd4', '\xdc', '\xa7', '\x7c', '\x26', '\x59', '\x21', '\x6b', '\x32', '\x66', '\xca', '\x4c', '\x0b', '\x78', '\x2c', '\xd0', '\xce', '\x7c', '\x01'}

	// Initialize the encryption backend
	backend, err := NewEciesEncryptionBackend(pub, nil)
	if err != nil {
		t.Fatalf("Failed to initialize encryption backend: %v", err)
	}

	// Test small file encryption and decryption
	smallData := []byte("This is a small file.")
	encryptedSmallData, err := backend.Encrypt(smallData)
	if err != nil {
		t.Fatalf("Failed to encrypt small file: %v", err)
	}

	priv := []byte{'\x93', '\xce', '\xc8', '\x5e', '\x4d', '\xe9', '\xda', '\x39', '\x60', '\x47', '\x6e', '\x00', '\xe4', '\x73', '\x48', '\x82', '\x82', '\xe7', '\x47', '\x6b', '\x03', '\x49', '\x1b', '\x0d', '\xb4', '\xbf', '\x82', '\xb9', '\xf5', '\xf8', '\x68', '\x1f'}
	backend, err = NewEciesEncryptionBackend(nil, priv)
	if err != nil {
		t.Fatalf("Failed to initialize encryption backend: %v", err)
	}

	decryptedSmallData, err := backend.Decrypt(encryptedSmallData)
	if err != nil {
		t.Fatalf("Failed to decrypt small file: %v", err)
	}
	if !bytes.Equal(smallData, decryptedSmallData) {
		t.Fatal("Small file encryption and decryption failed: data mismatch")
	}

	backend, err = NewEciesEncryptionBackend(pub, priv)
	if err != nil {
		t.Fatalf("Failed to initialize encryption backend: %v", err)
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
