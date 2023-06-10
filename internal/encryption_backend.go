package internal

// EncryptionConfig represents the configuration for an encryption backend.
type EncryptionConfig interface{}

// EncryptionBackend represents an interface for encryption operations.
type EncryptionBackend interface {
	// Initialize the encryption backend with any necessary configuration.
	Initialize(config EncryptionConfig) error

	// Encrypt encrypts the provided data using the public key (for assymetrical encryption).
	Encrypt(data []byte) ([]byte, error)

	// Decrypt decrypts the provided encrypted data using the private key (for assymetrical encryption).
	Decrypt(encryptedData []byte) ([]byte, error)
}
