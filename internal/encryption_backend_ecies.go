package internal

import (
	"errors"

	ecies "github.com/ecies/go/v2"
)

// EciesEncryptionConfig represents the configuration for an ECIES encryption backend.
type EciesEncryptionConfig struct {
	PublicKey  []byte
	PrivateKey []byte
}

// EciesEncryptionBackend represents an encryption backend using ECIES.
type EciesEncryptionBackend struct {
	publicKey  *ecies.PublicKey
	privateKey *ecies.PrivateKey
}

// NewEciesEncryptionBackend creates a new ECIES encryption backend.
func NewEciesEncryptionBackend(publicKey []byte, privateKey []byte) (*EciesEncryptionBackend, error) {
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

	privateKey := ecies.NewPrivateKeyFromBytes(config.PrivateKey)
	e.privateKey = privateKey

	if len(config.PublicKey) != 0 {
		publicKey, err := ecies.NewPublicKeyFromBytes(config.PublicKey)
		if err == nil {
			e.publicKey = publicKey
		}
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
