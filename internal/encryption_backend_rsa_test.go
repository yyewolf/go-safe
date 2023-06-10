package internal_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/yyewolf/go-safe/internal"
)

func TestEncryptionBackend_RSA(t *testing.T) {
	// Generate RSA key pair for testing
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		t.Fatalf("failed to generate RSA key pair: %v", err)
	}

	pkey, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		t.Fatalf("failed to marshal private key: %v", err)
	}

	// Convert private key to PEM format
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: pkey,
	})

	// Convert public key to PEM format
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("failed to marshal public key: %v", err)
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	// Create the RSAEncryptionConfig
	config := &internal.RSAEncryptionConfig{
		PrivateKey: privateKeyPEM,
		PublicKey:  publicKeyPEM,
	}

	// Create an instance of the RSAEncryptionBackend
	backend := internal.RSAEncryptionBackend{}
	err = backend.Initialize(config)
	if err != nil {
		t.Fatalf("failed to initialize RSAEncryptionBackend: %v", err)
	}

	// Test encryption and decryption
	originalData := []byte("Hello, RSA 4096 encryption!")
	encryptedData, err := backend.Encrypt(originalData)
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}
	decryptedData, err := backend.Decrypt(encryptedData)
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	// Check if the decrypted data matches the original data
	if string(decryptedData) != string(originalData) {
		t.Fatalf("decrypted data does not match original data")
	}
}
