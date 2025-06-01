package storage

import (
	"context"
	"fmt"
)

// StorageProvider defines the interface for different storage implementations
type StorageProvider interface {
	// SaveImage saves image data and returns the URL or path to access it
	SaveImage(ctx context.Context, imageData []byte, filename string) (string, error)
	// GetImageURL returns the URL or path to access the image
	GetImageURL(filename string) string
}

// StorageType represents the type of storage to use
type StorageType string

const (
	// StorageTypeLocal represents local filesystem storage
	StorageTypeLocal StorageType = "local"
	// StorageTypeS3 represents AWS S3 or S3-compatible storage
	StorageTypeS3 StorageType = "s3"
)

// Config holds configuration for storage providers
type Config struct {
	Type           StorageType
	LocalDirectory string // For local storage
	BucketName     string // For S3 storage
	Region         string // For S3 storage
	Endpoint       string // For S3-compatible storage (optional)
	AccessKey      string // For S3 storage
	SecretKey      string // For S3 storage
	UseSSL         bool   // For S3 storage
}

// NewStorage creates a new storage provider based on the configuration
func NewStorage(cfg Config) (StorageProvider, error) {
	switch cfg.Type {
	case StorageTypeLocal:
		return NewLocalStorage(cfg.LocalDirectory)
	case StorageTypeS3:
		return NewS3Storage(cfg)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", cfg.Type)
	}
}
