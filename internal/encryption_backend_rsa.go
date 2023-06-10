package internal

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

// RSAEncryptionConfig represents the configuration for an RSA encryption backend.
type RSAEncryptionConfig struct {
	PrivateKey []byte
	PublicKey  []byte
}

// RSAEncryptionBackend represents an encryption backend using RSA.
type RSAEncryptionBackend struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// NewRSAEncryptionBackend creates a new instance of the RSAEncryptionBackend.
func NewRSAEncryptionBackend(config *RSAEncryptionConfig) (EncryptionBackend, error) {
	b := RSAEncryptionBackend{}
	err := b.Initialize(config)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// Initialize initializes the RSA encryption backend with the private key.
func (e *RSAEncryptionBackend) Initialize(config EncryptionConfig) error {
	if config == nil {
		return errors.New("config not set")
	}

	cfg, ok := config.(*RSAEncryptionConfig)
	if !ok {
		return errors.New("config is not of type RSAEncryptionConfig")
	}

	privateKey, err := parsePrivateKey(cfg.PrivateKey)
	if err == nil {
		e.privateKey = privateKey
	}
	fmt.Println(err)

	publicKey, err := parsePublicKey(cfg.PublicKey)
	if err == nil {
		e.publicKey = publicKey
	}
	return nil
}

// Encrypt encrypts the provided data using RSA encryption.
func (e *RSAEncryptionBackend) Encrypt(data []byte) ([]byte, error) {
	if e.publicKey == nil {
		return nil, errors.New("public key not set")
	}

	// Separate the message into chunks
	msgLen := len(data)
	chunkSize := e.publicKey.Size() - 2*sha256.Size - 2
	chunks := make([][]byte, 0)
	for i := 0; i < msgLen; i += chunkSize {
		end := i + chunkSize
		if end > msgLen {
			end = msgLen
		}
		chunks = append(chunks, data[i:end])
	}

	// Encrypt each chunk
	encryptedChunks := make([][]byte, len(chunks))
	for i, chunk := range chunks {
		encryptedChunk, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, e.publicKey, chunk, nil)
		if err != nil {
			return nil, err
		}
		encryptedChunks[i] = encryptedChunk
	}

	// Concatenate the encrypted chunks
	encryptedData := make([]byte, 0)
	for _, encryptedChunk := range encryptedChunks {
		encryptedData = append(encryptedData, encryptedChunk...)
	}

	return encryptedData, nil
}

// Decrypt decrypts the provided encrypted data using RSA decryption.
func (e *RSAEncryptionBackend) Decrypt(encryptedData []byte) ([]byte, error) {
	if e.privateKey == nil {
		return nil, errors.New("private key not set")
	}

	// Separate the message into chunks
	msgLen := len(encryptedData)
	chunkSize := e.privateKey.Size()
	chunks := make([][]byte, 0)
	for i := 0; i < msgLen; i += chunkSize {
		end := i + chunkSize
		if end > msgLen {
			end = msgLen
		}
		chunks = append(chunks, encryptedData[i:end])
	}

	// Decrypt each chunk
	decryptedChunks := make([][]byte, len(chunks))
	for i, chunk := range chunks {
		decryptedChunk, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, e.privateKey, chunk, nil)
		if err != nil {
			return nil, err
		}
		decryptedChunks[i] = decryptedChunk
	}

	// Concatenate the decrypted chunks
	decryptedData := make([]byte, 0)
	for _, decryptedChunk := range decryptedChunks {
		decryptedData = append(decryptedData, decryptedChunk...)
	}

	return decryptedData, nil
}

// parsePrivateKey parses the private key from a PEM-encoded string.
func parsePrivateKey(privateKeyPEM []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return nil, errors.New("failed to parse private key")
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	privateKeyRSA, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("failed to parse private key")
	}
	return privateKeyRSA, nil
}

// parsePublicKey parses the public key from a PEM-encoded string.
func parsePublicKey(publicKeyPEM []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(publicKeyPEM)
	if block == nil {
		return nil, errors.New("failed to parse public key")
	}
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	publicKeyRSA, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("failed to parse public key")
	}
	return publicKeyRSA, nil
}
