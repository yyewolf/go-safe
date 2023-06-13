package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/yyewolf/go-safe/storage"
)

// Create and configure the Cobra command
var rootCmd = &cobra.Command{
	Use: "go-safe",
	Run: func(cmd *cobra.Command, args []string) {
		// Check if in export mode
		if config.Export {
			export()
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
		if st, err := os.Stat(config.Backup.Dir); err != nil || !st.IsDir() {
			fmt.Println("Backup directory does not exist or is not a directory")
			os.Exit(1)
		}

		dbFile := filepath.Join(config.Backup.Dir, "db.gosafe")

		loadDatabase(dbFile)

		fmt.Println("Starting backup service in '", config.Backup.Dir, "'...")
		worker(s3Backend)
	},
}

func main() {
	// Execute the Cobra command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func worker(b storage.StorageBackend) {
	duration := time.Duration(config.Interval) * time.Second

	for {
		// Walk the backup directory and upload any new or modified files
		err := filepath.Walk(config.Backup.Dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip if it's not a readable file
			if !info.Mode().IsRegular() || info.Mode()&0400 == 0 {
				return nil
			}

			savePath := path
			// Remove the backup directory from the path and add s3Dir if it's set
			if strings.HasPrefix(path, config.Backup.Dir) {
				savePath = path[len(config.Backup.Dir):]
				if savePath[0] == filepath.Separator {
					savePath = savePath[1:]
				}
			}

			if savePath == "db.gosafe" {
				return nil
			}

			// Read the file
			data, err := os.ReadFile(path)
			if err != nil {
				fmt.Printf("Failed to read %s: %v\n", path, err)
				return nil
			}

			// Check if the file is already in the database
			_, ok := database[savePath]
			if !ok {
				// File is not in the database, so upload it
				fmt.Println("Uploading", path, "...")

				err = b.Store(savePath, data)
				if err != nil {
					fmt.Printf("Failed to upload %s: %v\n", path, err)
					return nil
				}

				// SHA256 sum the file
				sum := sha256.Sum256(data)
				digest := hex.EncodeToString(sum[:])

				// Add the file to the database
				database[savePath] = &File{
					Sum: digest,
				}

				return nil
			}

			// File is in the database, so check if it has been modified
			sum := sha256.Sum256(data)
			digest := hex.EncodeToString(sum[:])

			if digest != database[savePath].Sum {
				// File has been modified, so upload it
				fmt.Println("Uploading", path, " (modified)...")

				err = b.Store(savePath, data)
				if err != nil {
					fmt.Printf("Failed to upload %s: %v\n", path, err)
				}

				database[savePath].Sum = digest
			}

			return nil
		})

		if err != nil {
			fmt.Printf("Failed to walk backup directory: %v\n", err)
		}

		// Delete any files that have been deleted
		for path := range database {
			// Check if the file exists
			path = filepath.Join(config.Backup.Dir, path)
			_, err := os.Stat(path)
			if err != nil {
				// File does not exist, so delete it from the database
				fmt.Println("Deleting", path, "...")

				savePath := path
				// Remove the backup directory from the path
				if strings.HasPrefix(path, config.Backup.Dir) {
					savePath = path[len(config.Backup.Dir):]
					if savePath[0] == filepath.Separator {
						savePath = savePath[1:]
					}
				}

				err = b.Delete(savePath)
				if err != nil {
					fmt.Printf("Failed to delete %s: %v\n", path, err)
				}

				delete(database, path)
			}
		}

		// Save the database
		err = saveDatabase()
		if err != nil {
			fmt.Printf("Failed to save database: %v\n", err)
		}

		// Check if the database has been modified
		data, err := os.ReadFile(filepath.Join(config.Backup.Dir, "db.gosafe"))
		if err != nil {
			fmt.Printf("Failed to read database: %v\n", err)
		}

		sum := sha256.Sum256(data)
		digest := hex.EncodeToString(sum[:])

		if digest != databaseDigest {
			databaseDigest = digest
			// Database has been modified, so upload it
			fmt.Println("Uploading database...")
			err = b.Store("db.gosafe", data)
			if err != nil {
				fmt.Printf("Failed to upload database: %v\n", err)
			}
		}

		time.Sleep(duration)
	}
}
