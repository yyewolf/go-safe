package internal

type Config interface {
}

type Backend interface {
	// Initialize the backend with any necessary configuration.
	Initialize(config Config, encryptionBackend EncryptionBackend) error

	// Store a file with the specified key and encrypted data.
	Store(key string, data []byte) error

	// Retrieve a file with the specified key and return its decrypted data.
	Retrieve(key string) ([]byte, error)

	// Delete a file with the specified key.
	Delete(key string) error
}
