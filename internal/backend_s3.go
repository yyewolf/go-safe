package internal

import (
	"bytes"
	"errors"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Config represents the configuration for the S3 backend.
type S3Config struct {
	StorageClass string
	Prepend      string
	Bucket       string
	Config       *aws.Config
}

// S3Backend represents a backend that stores and retrieves files from Amazon S3.
type S3Backend struct {
	storageclass      string
	prepend           string
	bucket            string
	config            *aws.Config
	encryptionBackend EncryptionBackend
	s3Client          *s3.S3
}

// NewS3Backend creates a new instance of the S3Backend.
func NewS3Backend(config *S3Config, encryptionBackend EncryptionBackend) (Backend, error) {
	b := S3Backend{}
	err := b.Initialize(config, encryptionBackend)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// Initialize initializes the S3 backend with the configuration and encryption backend.
func (b *S3Backend) Initialize(cfg Config, encryptionBackend EncryptionBackend) error {
	config, ok := cfg.(*S3Config)
	if !ok {
		return errors.New("config is not of type S3Config")
	}

	// Check that Config is not nil
	if config.Config == nil {
		return errors.New("config cannot be nil")
	}

	// Check that bucket is not empty
	if config.Bucket == "" {
		return errors.New("bucket cannot be empty")
	}

	// Create a new session with the AWS region
	sess, err := session.NewSession(config.Config)
	if err != nil {
		return err
	}

	// Create a new S3 client
	b.s3Client = s3.New(sess)

	b.storageclass = config.StorageClass
	b.prepend = config.Prepend
	b.bucket = config.Bucket
	b.config = config.Config
	b.encryptionBackend = encryptionBackend

	return nil
}

// Store stores a file in S3 with the specified key and encrypted data.
func (b *S3Backend) Store(key string, data []byte) error {
	// Encrypt the data using the encryption backend
	encryptedData, err := b.encryptionBackend.Encrypt(data)
	if err != nil {
		return err
	}

	key = filepath.Join(b.prepend, key)

	// Upload the encrypted data to S3
	_, err = b.s3Client.PutObject(&s3.PutObjectInput{
		Bucket:       aws.String(b.bucket),
		Key:          aws.String(key),
		Body:         bytes.NewReader(encryptedData),
		StorageClass: aws.String(b.storageclass),
	})
	if err != nil {
		return err
	}

	return nil
}

// Retrieve retrieves a file from S3 with the specified key and returns its decrypted data.
func (b *S3Backend) Retrieve(key string) ([]byte, error) {
	key = filepath.Join(b.prepend, key)
	// Download the encrypted data from S3
	resp, err := b.s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the encrypted data from the response
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	encryptedData := buf.Bytes()

	// Decrypt the data using the encryption backend
	decryptedData, err := b.encryptionBackend.Decrypt(encryptedData)
	if err != nil {
		return nil, err
	}

	return decryptedData, nil
}

// Delete deletes a file from S3 with the specified key.
func (b *S3Backend) Delete(key string) error {
	_, err := b.s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	return nil
}
