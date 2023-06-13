package encryption

import (
	"errors"
	"fmt"
	"strings"

	ecies "github.com/yyewolf/go-ecies/v2"
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

	// Remove any trailing \r\n or \n or \r or spaces
	config.PrivateKey = strings.Trim(config.PrivateKey, "\r\n ")
	config.PublicKey = strings.Trim(config.PublicKey, "\r\n ")

	privateKey, err := ecies.NewPrivateKeyFromHex(config.PrivateKey)
	if err == nil {
		e.privateKey = privateKey
	} else {
		fmt.Println("Failed to load private key: ", err)
	}

	publicKey, err := ecies.NewPublicKeyFromHex(config.PublicKey)
	if err == nil {
		e.publicKey = publicKey
	} else {
		fmt.Println("Failed to load public key: ", err)
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
