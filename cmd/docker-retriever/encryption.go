package main

import (
	"fmt"
	"os"

	"github.com/yyewolf/go-safe/encryption"
)

func encryptionBackend() encryption.EncryptionBackend {
	if config.AES.KeyLocation != "" {
		return aesEncryptionBackend()
	}

	if config.ECIES.PrivateKeyLocation != "" {
		return eciesPrivateEncryptionBackend()
	}

	return nil
}

func aesEncryptionBackend() encryption.EncryptionBackend {
	// Check key file permissions and existence
	st, err := os.Stat(config.AES.KeyLocation)
	if err != nil {
		fmt.Printf("Failed to stat key file: %v\n", err)
		os.Exit(1)
	}

	// Key should only be readable by the owner
	if st.Mode() != 0600 && st.Mode() != 0400 {
		fmt.Println("Key file permissions are too open")
		os.Exit(1)
	}

	// Read the key file
	aesKey, err := os.ReadFile(config.AES.KeyLocation)
	if err != nil {
		fmt.Printf("Failed to read key file: %v\n", err)
		os.Exit(1)
	}

	// Configure encryption backend
	encryptionBackend, err := encryption.NewAESEncryptionBackend(aesKey)
	if err != nil {
		fmt.Printf("Failed to configure encryption backend: %v\n", err)
		os.Exit(1)
	}

	return encryptionBackend
}

func eciesPrivateEncryptionBackend() encryption.EncryptionBackend {
	// Check key file permissions and existence
	st, err := os.Stat(config.ECIES.PrivateKeyLocation)
	if err != nil {
		fmt.Printf("Failed to stat public key file: %v\n", err)
		os.Exit(1)
	}

	// Key should only be readable by the owner
	if st.Mode() != 0600 && st.Mode() != 0400 {
		fmt.Println("Public key file permissions are too open")
		os.Exit(1)
	}

	// Read the key file
	privKey, err := os.ReadFile(config.ECIES.PrivateKeyLocation)
	if err != nil {
		fmt.Printf("Failed to read public key file: %v\n", err)
		os.Exit(1)
	}

	// Configure encryption backend
	encryptionBackend, err := encryption.NewEciesEncryptionBackend("", string(privKey))
	if err != nil {
		fmt.Printf("Failed to configure encryption backend: %v\n", err)
		os.Exit(1)
	}

	return encryptionBackend
}
