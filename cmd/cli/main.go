package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yyewolf/go-safe/storage"
)

var backupDir string

// Create and configure the Cobra command
var rootCmd = &cobra.Command{
	Use: "go-safe",
	Run: func(cmd *cobra.Command, args []string) {
		if config.ECIES.GenKey {
			eciesGenKey()
			os.Exit(0)
		}

		if config.HPKE.GenKey {
			hpkeGenKey()
			os.Exit(0)
		}

		// Configure encryption backend
		encryptionBackend := encryptionBackend()
		if encryptionBackend == nil {
			fmt.Println("No encryption backend configured")
			os.Exit(1)
		}

		// Configure storage backend
		s3Backend := storageBackend(encryptionBackend)
		if s3Backend == nil {
			fmt.Println("No storage backend configured")
			os.Exit(1)
		}

		// Check that backup directory exists and is a directory
		if st, err := os.Stat(backupDir); err != nil || !st.IsDir() {
			fmt.Println("Backup directory does not exist or is not a directory")
			os.Exit(1)
		}

		fmt.Println("Attempting to download in '", backupDir, "'...")
		downloader(s3Backend)
	},
}

func main() {
	// Execute the Cobra command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func downloader(b storage.StorageBackend) {
	// Download db.gosafe from S3
	data, err := b.Retrieve("db.gosafe")
	if err != nil {
		fmt.Printf("Failed to download db.gosafe from S3: %v\n", err)
		os.Exit(1)
	}

	err = json.Unmarshal(data, &database)
	if err != nil {
		fmt.Printf("Failed to unmarshal db.gosafe: %v\n", err)
		os.Exit(1)
	}

	// Download all files from S3
	for path, file := range database {
		savePath := filepath.Join(backupDir, path)
		data, err := b.Retrieve(path)
		if err != nil {
			fmt.Printf("Failed to download %s from S3: %v\n", path, err)
		}

		// Check the SHA256 hash of the file
		hash := sha256.Sum256(data)
		if hex.EncodeToString(hash[:]) != file.Sum {
			fmt.Printf("Failed to verify hash of %s, be careful\n", path)
		}

		// Check if we need to create any directories
		dir := filepath.Dir(savePath)
		st, err := os.Stat(dir)
		if err != nil {
			// mkdir -p
			err = os.MkdirAll(dir, 0700)
			if err != nil {
				fmt.Printf("Failed to create directory: %v\n", err)
				os.Exit(1)
			}
		} else if !st.IsDir() {
			fmt.Printf("Failed to create directory: %s is not a directory\n", dir)
			os.Exit(1)
		}

		// Write the file to disk
		err = os.WriteFile(savePath, data, 0600)
		if err != nil {
			fmt.Printf("Failed to write %s to disk: %v\n", savePath, err)
		}
	}
}
