package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yyewolf/go-safe/internal"
)

var backupDir string

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
				fmt.Println("Key file is readable by others")
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
				// mkdir -p
				err = os.MkdirAll(backupDir, 0700)
				if err != nil {
					fmt.Printf("Failed to create backup directory: %v\n", err)
					os.Exit(1)
				}
				st, err = os.Stat(backupDir)
				if err != nil {
					fmt.Printf("Failed to stat backup directory: %v\n", err)
					os.Exit(1)
				}
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

			fmt.Println("Starting to retrieve in '", backupDir, "'...")
			downloader(s3Backend)
		},
	}

	// Bind flags to Viper
	rootCmd.Flags().String("access-id", "", "S3 access ID")
	rootCmd.Flags().String("access-key", "", "S3 access key")
	rootCmd.Flags().String("bucket-name", "", "S3 bucket name")
	rootCmd.Flags().String("endpoint", "", "S3 endpoint")
	rootCmd.Flags().String("region", "", "S3 region")
	rootCmd.Flags().String("s3-dir", "", "S3 directory")
	rootCmd.Flags().String("backup-dir", "", "Backup directory")
	rootCmd.Flags().String("aes-key-location", "", "AES key location")

	// Bind flags to environment variables
	viper.BindPFlag("S3_ACCESS_ID", rootCmd.Flags().Lookup("access-id"))
	viper.BindPFlag("S3_ACCESS_KEY", rootCmd.Flags().Lookup("access-key"))
	viper.BindPFlag("S3_BUCKET_NAME", rootCmd.Flags().Lookup("bucket-name"))
	viper.BindPFlag("S3_ENDPOINT", rootCmd.Flags().Lookup("endpoint"))
	viper.BindPFlag("S3_REGION", rootCmd.Flags().Lookup("region"))
	viper.BindPFlag("S3_DIR", rootCmd.Flags().Lookup("s3-dir"))
	viper.BindPFlag("BACKUP_DIR", rootCmd.Flags().Lookup("backup-dir"))
	viper.BindPFlag("AES_KEY_LOCATION", rootCmd.Flags().Lookup("aes-key-location"))

	// Execute the Cobra command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func downloader(b internal.Backend) {
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
