package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// localStorage implements StorageProvider for local filesystem storage
type localStorage struct {
	baseDir string
}

// NewLocalStorage creates a new local filesystem storage provider
func NewLocalStorage(baseDir string) (StorageProvider, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base directory: %w", err)
	}

	return &localStorage{
		baseDir: baseDir,
	}, nil
}

// SaveImage saves the image data to the local filesystem
func (s *localStorage) SaveImage(ctx context.Context, imageData []byte, filename string) (string, error) {
	// Ensure the filename is safe and has an extension
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = ".png" // Default to png if no extension
		filename += ext
	}

	// Create the full file path
	filePath := filepath.Join(s.baseDir, filename)

	// Write the file
	if err := os.WriteFile(filePath, imageData, 0644); err != nil {
		return "", fmt.Errorf("failed to save image to %s: %w", filePath, err)
	}

	return filename, nil
}

// GetImageURL returns the relative path to the image
func (s *localStorage) GetImageURL(filename string) string {
	// For local storage, we just return the relative path
	// The HTTP server will serve files from the images directory
	return "/images/" + filename
}
