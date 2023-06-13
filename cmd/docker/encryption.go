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

	if config.ECIES.PublicKeyLocation != "" {
		return eciesPublicEncryptionBackend()
	}

	if config.HPKE.ClientPublicKeyLocation != "" && config.HPKE.ClientSecretKeyLocation != "" && config.HPKE.ServerPublicKeyLocation != "" {
		return hpkeEncryptionBackend()
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

func eciesPublicEncryptionBackend() encryption.EncryptionBackend {
	// Check key file permissions and existence
	st, err := os.Stat(config.ECIES.PublicKeyLocation)
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
	publicKey, err := os.ReadFile(config.ECIES.PublicKeyLocation)
	if err != nil {
		fmt.Printf("Failed to read public key file: %v\n", err)
		os.Exit(1)
	}

	// Configure encryption backend
	encryptionBackend, err := encryption.NewEciesEncryptionBackend(string(publicKey), "")
	if err != nil {
		fmt.Printf("Failed to configure encryption backend: %v\n", err)
		os.Exit(1)
	}

	return encryptionBackend
}

func hpkeEncryptionBackend() encryption.EncryptionBackend {
	// Check key file permissions and existence
	if st, err := os.Stat(config.HPKE.ClientPublicKeyLocation); err != nil || st.Mode() != 0600 && st.Mode() != 0400 {
		if err != nil {
			fmt.Printf("Failed to stat client public key file: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Client public key file permissions are too open")
		os.Exit(1)
	}

	// Check key file permissions and existence
	if st, err := os.Stat(config.HPKE.ClientSecretKeyLocation); err != nil || st.Mode() != 0600 && st.Mode() != 0400 {
		if err != nil {
			fmt.Printf("Failed to stat client secret key file: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Client secret key file permissions are too open")
		os.Exit(1)
	}

	// Check key file permissions and existence
	if st, err := os.Stat(config.HPKE.ServerPublicKeyLocation); err != nil || st.Mode() != 0600 && st.Mode() != 0400 {
		if err != nil {
			fmt.Printf("Failed to stat server public key file: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Server public key file permissions are too open")
		os.Exit(1)
	}

	// Read the key files
	clientPublicKey, err := os.ReadFile(config.HPKE.ClientPublicKeyLocation)
	if err != nil {
		fmt.Printf("Failed to read client public key file: %v\n", err)
		os.Exit(1)
	}

	clientSecretKey, err := os.ReadFile(config.HPKE.ClientSecretKeyLocation)
	if err != nil {
		fmt.Printf("Failed to read client secret key file: %v\n", err)
		os.Exit(1)
	}

	serverPublicKey, err := os.ReadFile(config.HPKE.ServerPublicKeyLocation)
	if err != nil {
		fmt.Printf("Failed to read server public key file: %v\n", err)
		os.Exit(1)
	}

	// Configure encryption backend
	encryptionBackend, err := encryption.NewHPKEBackend(clientPublicKey, clientSecretKey, serverPublicKey, nil, []byte(config.HPKE.PresharedKey), []byte(config.HPKE.PresharedKeyID))
	if err != nil {
		fmt.Printf("Failed to configure encryption backend: %v\n", err)
		os.Exit(1)
	}

	return encryptionBackend
}
