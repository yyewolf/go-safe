package encryption

import (
	"encoding/json"
	"errors"

	hpke "github.com/jedisct1/go-hpke-compact"
)

type HPKEConfig struct {
	UtilityPKey []byte
	UtilitySKey []byte
	YourPKey    []byte
	YourSKey    []byte

	PresharedKey   []byte
	PresharedKeyID []byte
}

// HPKEBackend is an implementation of the EncryptionBackend interface using go-hpke-compact.
type HPKEBackend struct {
	suite     *hpke.Suite
	client    hpke.KeyPair
	server    hpke.KeyPair
	preshared *hpke.Psk
}

func NewHPKEBackend(clientPublicKey, clientSecretKey, serverPublicKey, serverPrivateKey, presharedKey, presharedKeyID []byte) (*HPKEBackend, error) {
	b := &HPKEBackend{}

	cfg := HPKEConfig{
		UtilityPKey: clientPublicKey,
		UtilitySKey: clientSecretKey,
		YourPKey:    serverPublicKey,
		YourSKey:    serverPrivateKey,

		PresharedKey:   presharedKey,
		PresharedKeyID: presharedKeyID,
	}

	err := b.Initialize(&cfg)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Initialize initializes the encryption backend with the provided configuration.
func (b *HPKEBackend) Initialize(config EncryptionConfig) error {
	cfg, ok := config.(*HPKEConfig)
	if !ok {
		return errors.New("invalid GoHPKECompact encryption configuration")
	}

	suite, err := hpke.NewSuite(hpke.KemX25519HkdfSha256, hpke.KdfHkdfSha256, hpke.AeadChaCha20Poly1305)
	if err != nil {
		return err
	}

	b.suite = suite
	b.client.PublicKey = cfg.UtilityPKey
	b.client.SecretKey = cfg.UtilitySKey
	b.server.PublicKey = cfg.YourPKey
	b.server.SecretKey = cfg.YourSKey
	b.preshared = &hpke.Psk{
		Key: cfg.PresharedKey,
		ID:  cfg.PresharedKeyID,
	}

	return nil
}

// Encrypt encrypts the provided plaintext.
func (b *HPKEBackend) Encrypt(plaintext []byte) ([]byte, error) {
	clientCtx, ss, err := b.suite.NewAuthenticatedClientContext(b.client, b.server.PublicKey, []byte("go-safe"), b.preshared)
	if err != nil {
		return nil, err
	}

	ciphertext, err := clientCtx.EncryptToServer(plaintext, nil)
	if err != nil {
		return nil, err
	}

	out := struct {
		EncryptedData []byte `json:"ed"`
		SharedSecret  []byte `json:"ss"`
	}{
		EncryptedData: ciphertext,
		SharedSecret:  ss,
	}

	return json.Marshal(out)
}

// Decrypt decrypts the provided ciphertext.
func (b *HPKEBackend) Decrypt(ciphertext []byte) ([]byte, error) {
	var in struct {
		EncryptedData []byte `json:"ed"`
		SharedSecret  []byte `json:"ss"`
	}

	err := json.Unmarshal(ciphertext, &in)
	if err != nil {
		return nil, err
	}

	serverCtx, err := b.suite.NewAuthenticatedServerContext(b.client.PublicKey, in.SharedSecret, b.server, []byte("go-safe"), b.preshared)
	if err != nil {
		return nil, err
	}

	return serverCtx.DecryptFromClient(in.EncryptedData, nil)
}
