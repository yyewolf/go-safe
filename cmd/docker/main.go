package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yyewolf/go-safe/internal"
)

var backupDir string
var interval int

func main() {
	// Initialize Viper
	viper.AutomaticEnv()     // Read environment variables
	viper.SetEnvPrefix("GS") // Set environment variable prefix

	// Set default values for flags
	viper.SetDefault("BACKUP_DIR", "/backup")
	viper.SetDefault("AES_KEY_LOCATION", "/aes.key")

	// Create and configure the Cobra command
	rootCmd := &cobra.Command{
		Use: "go-safe",
		Run: func(cmd *cobra.Command, args []string) {
			// Get the values from Viper
			s3AccessID := viper.GetString("S3_ACCESS_ID")
			s3AccessKey := viper.GetString("S3_ACCESS_KEY")
			s3BucketName := viper.GetString("S3_BUCKET_NAME")
			s3Endpoint := viper.GetString("S3_ENDPOINT")
			s3Region := viper.GetString("S3_REGION")
			s3Dir := viper.GetString("S3_DIR")
			backupDir = viper.GetString("BACKUP_DIR")
			aesKeyLocation := viper.GetString("AES_KEY_LOCATION")
			interval = viper.GetInt("INTERVAL")

			// Validate and process the values
			if s3AccessID == "" || s3AccessKey == "" || s3BucketName == "" || s3Endpoint == "" {
				fmt.Println("Missing required S3 configuration")
				cmd.Usage()
				os.Exit(1)
			}

			// Check key file permissions and existence
			st, err := os.Stat(aesKeyLocation)
			if err != nil {
				fmt.Printf("Failed to stat key file: %v\n", err)
				os.Exit(1)
			}

			// Check if the key file is readable by anyone other than the owner
			if st.Mode()&0004 != 0 {
				fmt.Println("key file is readable by others")
				os.Exit(1)
			}

			// Check if the key file is writable by anyone other than the owner
			if st.Mode()&0002 != 0 {
				fmt.Println("key file is writable by others")
				os.Exit(1)
			}

			// Check that backup directory exists
			st, err = os.Stat(backupDir)
			if err != nil {
				fmt.Printf("Failed to stat backup directory: %v\n", err)
				os.Exit(1)
			}

			// Check that backup directory is a directory
			if !st.IsDir() {
				fmt.Println("Backup directory is not a directory")
				os.Exit(1)
			}

			// Read the key file
			aesKey, err := os.ReadFile(aesKeyLocation)
			if err != nil {
				fmt.Printf("Failed to read key file: %v\n", err)
				os.Exit(1)
			}

			// Configure encryption backend
			encryptionBackend, err := internal.NewAESEncryptionBackend(aesKey)
			if err != nil {
				fmt.Printf("Failed to configure encryption backend: %v\n", err)
				os.Exit(1)
			}

			// Configure S3 backend
			s3Config := &internal.S3Config{
				Prepend: s3Dir,
				Bucket:  s3BucketName,
				Config: aws.NewConfig().
					WithEndpoint(s3Endpoint).
					WithCredentials(
						credentials.NewStaticCredentials(
							s3AccessID,
							s3AccessKey,
							"",
						),
					),
			}

			if s3Region != "" {
				s3Config.Config = s3Config.Config.WithRegion(s3Region)
			}

			s3Backend, err := internal.NewS3Backend(s3Config, encryptionBackend)
			if err != nil {
				fmt.Printf("Failed to configure S3 backend: %v\n", err)
				os.Exit(1)
			}

			dbFile := filepath.Join(backupDir, "db.gosafe")

			loadDatabase(dbFile)

			fmt.Println("Starting backup service in '", backupDir, "'...")
			worker(s3Backend)
		},
	}

	// Bind flags to Viper
	rootCmd.Flags().String("access-id", "", "S3 access ID")
	rootCmd.Flags().String("access-key", "", "S3 access key")
	rootCmd.Flags().String("bucket-name", "", "S3 bucket name")
	rootCmd.Flags().String("endpoint", "", "S3 endpoint")
	rootCmd.Flags().String("region", "", "S3 region")
	rootCmd.Flags().String("s3-dir", "", "S3 directory (will store under a directory in S3)")
	rootCmd.Flags().String("backup-dir", "", "Backup directory")
	rootCmd.Flags().String("aes-key-location", "", "AES key location")
	rootCmd.Flags().Int("interval", 60, "Backup interval in seconds")

	// Bind flags to environment variables
	viper.BindPFlag("S3_ACCESS_ID", rootCmd.Flags().Lookup("access-id"))
	viper.BindPFlag("S3_ACCESS_KEY", rootCmd.Flags().Lookup("access-key"))
	viper.BindPFlag("S3_BUCKET_NAME", rootCmd.Flags().Lookup("bucket-name"))
	viper.BindPFlag("S3_ENDPOINT", rootCmd.Flags().Lookup("endpoint"))
	viper.BindPFlag("S3_REGION", rootCmd.Flags().Lookup("region"))
	viper.BindPFlag("S3_DIR", rootCmd.Flags().Lookup("s3-dir"))
	viper.BindPFlag("BACKUP_DIR", rootCmd.Flags().Lookup("backup-dir"))
	viper.BindPFlag("AES_KEY_LOCATION", rootCmd.Flags().Lookup("aes-key-location"))
	viper.BindPFlag("INTERVAL", rootCmd.Flags().Lookup("interval"))

	// Execute the Cobra command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func worker(b internal.Backend) {
	duration := time.Duration(interval) * time.Second

	for {
		// Walk the backup directory and upload any new or modified files
		err := filepath.Walk(backupDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip if it's not a readable file
			if !info.Mode().IsRegular() || info.Mode()&0400 == 0 {
				return nil
			}

			savePath := path
			// Remove the backup directory from the path and add s3Dir if it's set
			if strings.HasPrefix(path, backupDir) {
				savePath = path[len(backupDir):]
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
			path = filepath.Join(backupDir, path)
			_, err := os.Stat(path)
			if err != nil {
				// File does not exist, so delete it from the database
				fmt.Println("Deleting", path, "...")

				savePath := path
				// Remove the backup directory from the path
				if strings.HasPrefix(path, backupDir) {
					savePath = path[len(backupDir):]
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
		data, err := os.ReadFile(filepath.Join(backupDir, "db.gosafe"))
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
