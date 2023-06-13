package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/yyewolf/go-safe/encryption"
	"github.com/yyewolf/go-safe/storage"
)

func storageBackend(encryptionBackend encryption.EncryptionBackend) storage.StorageBackend {
	if config.S3.AccessID != "" {
		return s3Backend(encryptionBackend)
	}
	return nil
}

func s3Backend(encryptionBackend encryption.EncryptionBackend) storage.StorageBackend {
	// Configure S3 backend
	s3Config := &storage.S3Config{
		StorageClass: config.S3.StorageClass,
		Prepend:      config.S3.Dir,
		Bucket:       config.S3.BucketName,
		Config: aws.NewConfig().
			WithCredentials(
				credentials.NewStaticCredentials(
					config.S3.AccessID,
					config.S3.AccessKey,
					"",
				),
			),
	}

	if config.S3.Region != "" {
		s3Config.Config = s3Config.Config.WithRegion(config.S3.Region)
	}

	if config.S3.Endpoint != "" {
		s3Config.Config = s3Config.Config.WithEndpoint(config.S3.Endpoint)
	}

	s3Backend, err := storage.NewS3Backend(s3Config, encryptionBackend)
	if err != nil {
		fmt.Printf("Failed to configure S3 backend: %v\n", err)
		os.Exit(1)
	}

	return s3Backend
}
