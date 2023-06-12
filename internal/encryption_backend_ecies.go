package internal

import (
	"errors"

	ecies "github.com/ecies/go/v2"
)

// EciesEncryptionConfig represents the configuration for an ECIES encryption backend.
type EciesEncryptionConfig struct {
	PublicKey  string
	PrivateKey string
}

// EciesEncryptionBackend represents an encryption backend using ECIES.
type EciesEncryptionBackend struct {
	publicKey  *ecies.PublicKey
	privateKey *ecies.PrivateKey
}

// NewEciesEncryptionBackend creates a new ECIES encryption backend.
func NewEciesEncryptionBackend(publicKey string, privateKey string) (*EciesEncryptionBackend, error) {
	b := EciesEncryptionBackend{}
	err := b.Initialize(&EciesEncryptionConfig{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	})
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// Initialize initializes the ECIES encryption backend with the encryption key.
func (e *EciesEncryptionBackend) Initialize(cfg EncryptionConfig) error {
	// Check the configuration type
	config, ok := cfg.(*EciesEncryptionConfig)
	if !ok {
		return errors.New("invalid ECIES encryption configuration")
	}

	privateKey, err := ecies.NewPrivateKeyFromHex(config.PrivateKey)
	if err == nil {
		e.privateKey = privateKey
	}

	publicKey, err := ecies.NewPublicKeyFromHex(config.PublicKey)
	if err == nil {
		e.publicKey = publicKey
	}

	return nil
}

// Encrypt encrypts the data using ECIES.
func (e *EciesEncryptionBackend) Encrypt(data []byte) ([]byte, error) {
	return ecies.Encrypt(e.publicKey, data)
}

// Decrypt decrypts the data using ECIES.
func (e *EciesEncryptionBackend) Decrypt(data []byte) ([]byte, error) {
	return ecies.Decrypt(e.privateKey, data)
}
