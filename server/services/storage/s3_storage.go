package storage

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// s3Storage implements StorageProvider for S3-compatible storage
type s3Storage struct {
	client     *s3.Client
	bucketName string
	publicURL  string
}

// NewS3Storage creates a new S3 storage provider
func NewS3Storage(cfg Config) (StorageProvider, error) {
	// Create AWS config
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKey,
			cfg.SecretKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
			o.UsePathStyle = true // Needed for S3-compatible services like MinIO
		}
	})

	// Set public URL (either custom endpoint or standard S3 URL)
	publicURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com", cfg.BucketName, cfg.Region)
	if cfg.Endpoint != "" {
		publicURL = strings.TrimSuffix(cfg.Endpoint, "/") + "/" + cfg.BucketName
	}

	return &s3Storage{
		client:     client,
		bucketName: cfg.BucketName,
		publicURL:  publicURL,
	}, nil
}

// SaveImage uploads the image data to S3
func (s *s3Storage) SaveImage(ctx context.Context, imageData []byte, filename string) (string, error) {
	// Ensure the filename is safe and has an extension
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = ".png" // Default to png if no extension
		filename += ext
	}

	// Create a unique filename with timestamp
	timestamp := time.Now().Unix()
	uniqueFilename := fmt.Sprintf("%d_%s", timestamp, filename)

	// Upload the file to S3
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(uniqueFilename),
		Body:        bytes.NewReader(imageData),
		ContentType: aws.String("image/" + strings.TrimPrefix(ext, ".")),
		ACL:         types.ObjectCannedACLPublicRead, // Make the object publicly accessible
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	return uniqueFilename, nil
}

// GetImageURL returns the public URL for the image
func (s *s3Storage) GetImageURL(filename string) string {
	return fmt.Sprintf("%s/%s", s.publicURL, filename)
}
