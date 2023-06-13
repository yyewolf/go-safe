package main

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	S3 struct {
		AccessID     string `mapstructure:"access-id"`
		AccessKey    string `mapstructure:"access-key"`
		BucketName   string `mapstructure:"bucket-name"`
		Endpoint     string `mapstructure:"endpoint"`
		Region       string `mapstructure:"region"`
		Dir          string `mapstructure:"dir"`
		StorageClass string `mapstructure:"storage-class"`
	} `mapstructure:"s3"`

	Backup struct {
		Dir string `mapstructure:"dir"`
	} `mapstructure:"backup"`

	AES struct {
		KeyLocation string `mapstructure:"key-location"`
	} `mapstructure:"aes"`

	ECIES struct {
		PrivateKeyLocation string `mapstructure:"private-key-location"`
	} `mapstructure:"ecies"`
}

var config Config

func init() {
	// Initialize Viper
	godotenv.Load() // Load environment variables from .env file
	cobra.OnInitialize(initConfig)

	// S3 Related
	rootCmd.Flags().String("s3.access-id", "", "S3 access ID")
	rootCmd.Flags().String("s3.access-key", "", "S3 access key")
	rootCmd.Flags().String("s3.bucket-name", "", "S3 bucket name")
	rootCmd.Flags().String("s3.endpoint", "", "S3 endpoint")
	rootCmd.Flags().String("s3.region", "", "S3 region")
	rootCmd.Flags().String("s3.dir", "", "S3 directory (will store under a directory in S3)")
	rootCmd.Flags().String("s3.storage-class", "", "S3 storage class")

	rootCmd.MarkFlagsRequiredTogether("s3.access-id", "s3.access-key", "s3.bucket-name", "s3.endpoint", "s3.region")

	// AES Related
	rootCmd.Flags().String("aes.key-location", "", "AES key location")

	// ECIES Related
	rootCmd.Flags().String("ecies.private-key-location", "", "ECIES private key location")

	// Encryption related
	rootCmd.MarkFlagsMutuallyExclusive("aes.key-location", "ecies.private-key-location")

	// Misc
	rootCmd.Flags().String("backup.dir", "", "Backup directory (where to save to)")

	// Bind flags to environment variables
	viper.BindPFlags(rootCmd.Flags())

	viper.AutomaticEnv() // Read environment variables

	// Load into Config struct
	viper.Unmarshal(&config)
}

func initConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.EnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_")) // Replace "." with "_" in environment variables
	viper.SetEnvPrefix("GS")                                      // Set environment variable prefix

	// Load config from file
	viper.SetConfigName("gosafe-config") // Name of config file (without extension)
	viper.AddConfigPath(".")             // Path to look for the config file in
	viper.AddConfigPath("$HOME")         // Path to look for the config file in

	viper.BindPFlags(rootCmd.Flags())
	viper.SetDefault("backup.dir", "/backup")
	viper.SetDefault("s3.storage-class", "STANDARD")

	viper.AutomaticEnv()

	viper.ReadInConfig()

	if err := viper.Unmarshal(&config); err != nil {
		cobra.CheckErr(err)
	}
}
