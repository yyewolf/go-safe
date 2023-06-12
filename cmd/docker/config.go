package main

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
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
		PublicKeyLocation string `mapstructure:"public-key-location"`
	} `mapstructure:"ecies"`

	Interval int `mapstructure:"interval"`
}

var config Config

func init() {
	// Initialize Viper
	godotenv.Load() // Load environment variables from .env file

	viper.EnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_")) // Replace "." with "_" in environment variables
	viper.SetEnvPrefix("GS")                                      // Set environment variable prefix

	viper.AutomaticEnv() // Read environment variables

	// Set default values for flags
	viper.SetDefault("backup.dir", "/backup")
	viper.SetDefault("s3.storage-class", "STANDARD")

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
	rootCmd.Flags().String("ecies.public-key-location", "", "ECIES public key location")

	// Encryption related
	rootCmd.MarkFlagsMutuallyExclusive("aes.key-location", "ecies.public-key-location")

	// Misc
	rootCmd.Flags().String("backup.dir", "", "Backup directory")
	rootCmd.Flags().Int("interval", 60, "Backup interval in seconds")

	// Bind flags to environment variables
	viper.BindPFlags(rootCmd.Flags())

	rootCmd.SetGlobalNormalizationFunc(func(f *pflag.FlagSet, name string) pflag.NormalizedName {
		replacer := strings.NewReplacer("-", "_", ".", "_")
		viper.BindEnv(name, fmt.Sprintf("GS_%s", replacer.Replace(strings.ToUpper(name))))
		return pflag.NormalizedName(name)
	})

	viper.ReadInConfig()

	// Unmarshal co nfig
	viper.Unmarshal(&config)
}
