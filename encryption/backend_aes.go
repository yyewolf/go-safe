package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// AESEncryptionConfig represents the configuration for an AES encryption backend.
type AESEncryptionConfig struct {
	Key []byte
}

// AESEncryptionBackend represents an encryption backend using AES.
type AESEncryptionBackend struct {
	key []byte
}

// NewAESEncryptionBackend creates a new AES encryption backend.
func NewAESEncryptionBackend(key []byte) (*AESEncryptionBackend, error) {
	b := AESEncryptionBackend{}
	err := b.Initialize(&AESEncryptionConfig{
		Key: key,
	})
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// Initialize initializes the AES encryption backend with the encryption key.
func (e *AESEncryptionBackend) Initialize(cfg EncryptionConfig) error {
	// Check the configuration type
	config, ok := cfg.(*AESEncryptionConfig)
	if !ok {
		return errors.New("invalid AES encryption configuration")
	}
	key := config.Key
	// Check the key length
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return errors.New("invalid AES key length")
	}

	e.key = key
	return nil
}

// Encrypt encrypts the provided data using AES encryption.
func (e *AESEncryptionBackend) Encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	// Create a new GCM cipher with the block and a random nonce
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Seal the plaintext using the AES GCM cipher
	encryptedData := aesgcm.Seal(nil, nonce, data, nil)

	// Prepend the nonce to the encrypted data
	encryptedData = append(nonce, encryptedData...)

	return encryptedData, nil
}

// Decrypt decrypts the provided encrypted data using AES decryption.
func (e *AESEncryptionBackend) Decrypt(encryptedData []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	// Split the encrypted data into the nonce and the actual ciphertext
	nonce := encryptedData[:12]
	ciphertext := encryptedData[12:]

	// Create a new GCM cipher with the block and the received nonce
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Decrypt the ciphertext using the AES GCM cipher
	decryptedData, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return decryptedData, nil
}
