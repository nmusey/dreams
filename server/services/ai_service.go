package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"dreams/services/storage"
)

type AIService struct {
	client          *http.Client
	host            string
	endpoint        string
	storageProvider storage.StorageProvider
	model           string
}

// StorageProvider defines the interface for storage operations
type StorageProvider interface {
	SaveImage(ctx context.Context, data []byte, filename string) (string, error)
	GetImageURL(filename string) string
	GetBasePath() string
}

func NewAIService(host string, endpoint string, model string, storageProvider storage.StorageProvider) *AIService {
	// Create a custom HTTP client with extended timeouts and retries
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	return &AIService{
		host:            host,
		endpoint:        endpoint,
		model:           model,
		client:          client,
		storageProvider: storageProvider,
	}
}

func (s *AIService) GenerateImage(dreamContent string) (string, error) {
	// Create a prompt for InvokeAI
	prompt := fmt.Sprintf(`
A surreal dream-like scene featuring:
- %s
- Style: dreamy and surreal
- Mood: mysterious and captivating
- Use vibrant colors and imaginative elements
- Composition: balanced and visually interesting

Generate this as a high-quality PNG image.`, dreamContent)

	req := ImageGenerationRequest{
		Model:  s.model,
		Prompt: prompt,
	}

	// Convert request to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Build the full URL using host and endpoint
	url := fmt.Sprintf("%s%s", s.host, s.endpoint)

	// Create HTTP request
	resp, err := s.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	delay := 1 * time.Second
	maxRetries := 3
	attempts := 0

	for attempts < maxRetries {
		attempts++
		if err == nil {
			break
		}
		time.Sleep(delay)
		delay *= 2 // Exponential backoff
		resp, err = s.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	}

	if err != nil {
		return "", fmt.Errorf("failed after %d attempts: %w", maxRetries, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("AI service returned status %d: %s", resp.StatusCode, string(body))
	}

	// Read response
	var response struct {
		Images []struct {
			Base64 string `json:"base64"`
		} `json:"images"`
		Error string `json:"error,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	if response.Error != "" {
		return "", fmt.Errorf("AI service error: %s", response.Error)
	}

	if len(response.Images) == 0 {
		return "", fmt.Errorf("no images returned from AI service")
	}

	// Save image
	filename, err := s.saveImage(response.Images[0].Base64)
	if err != nil {
		return "", fmt.Errorf("error saving image: %w", err)
	}

	return filename, nil
}

type ImageGenerationRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type ImageGenerationResponse struct {
	ImageData string `json:"image_data"` // Base64-encoded image data
	Error     string `json:"error,omitempty"`
	Status    string `json:"status,omitempty"`
}

// saveImage saves the image data to the configured storage provider
func (s *AIService) saveImage(imageData string) (string, error) {
	// Generate unique filename
	filename := fmt.Sprintf("image_%d_%d.png", time.Now().Unix(), rand.Int63())

	// Decode base64 image data
	decoded, err := base64.StdEncoding.DecodeString(imageData)
	if err != nil {
		return "", fmt.Errorf("error decoding base64 image: %w", err)
	}

	// Save image using the storage provider
	_, err = s.storageProvider.SaveImage(context.Background(), decoded, filename)
	if err != nil {
		return "", fmt.Errorf("error saving image: %w", err)
	}

	// Return the URL to the saved image
	return s.storageProvider.GetImageURL(filename), nil
}
