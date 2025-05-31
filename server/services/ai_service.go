package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"dreams/services/storage"
	"github.com/google/uuid"
)

type AIService struct {
	baseURL        string
	model          string
	client         *http.Client
	storage        storage.StorageProvider
	imageDirectory string
}

func NewAIService(baseURL string, model string, storageProvider storage.StorageProvider) *AIService {
	// Create a custom HTTP client with extended timeouts and retries
	client := &http.Client{
		// Set a very long timeout since we'll be handling timeouts via context
		Timeout: 30 * time.Minute,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
			DisableKeepAlives:   false,
			ForceAttemptHTTP2:   true,
		},
	}

	return &AIService{
		baseURL: baseURL,
		model:   model,
		client:  client,
		storage: storageProvider,
	}
}

type ImageGenerationRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type ImageGenerationResponse struct {
	ImageData string `json:"image_data"` // Base64-encoded image data
	Error    string `json:"error,omitempty"`
	Status   string `json:"status,omitempty"`
}

// saveImage saves the image data to the configured storage provider
func (s *AIService) saveImage(imageData string) (string, error) {
	// Decode the base64 image data
	decoded, err := base64.StdEncoding.DecodeString(imageData)
	if err != nil {
		return "", fmt.Errorf("error decoding base64 image: %v", err)
	}

	// Generate a unique filename with timestamp
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s.png", timestamp, uuid.New().String())

	// Save the image using the storage provider
	_, err = s.storage.SaveImage(context.Background(), decoded, filename)
	if err != nil {
		return "", fmt.Errorf("error saving image to storage: %v", err)
	}

	// Return the URL or path to the saved image
	return s.storage.GetImageURL(filename), nil
}

// GenerateImage sends a request to the AI service to generate an image based on the dream content
// and saves the resulting image to disk. Returns the path to the saved image.
func (s *AIService) GenerateImage(dreamContent string) (string, error) {
	// Create the prompt for the AI model
	prompt := fmt.Sprintf(`Generate a dream-like image based on the following description. The image should be a visual representation of the dream, capturing its mood, setting, and key elements. The image should be in a surreal, dreamy style with vibrant colors and imaginative elements.

Dream Description:
%s

Instructions:
1. Focus on the main elements and atmosphere of the dream.
2. Use a color palette that matches the mood (e.g., soft pastels for peaceful dreams, dark and moody for nightmares).
3. Include any specific objects, creatures, or landscapes mentioned.
4. The style should be artistic and dream-like, not photorealistic.

Respond with a base64-encoded PNG image and nothing else.`, dreamContent)

	// Use the chat completion API with LLaVA
	reqBody := map[string]interface{}{
		"model": "llava:latest",
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"stream": false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	url := fmt.Sprintf("%s/api/chat", s.baseURL)
	log.Printf("Sending request to: %s", url)

	// Create a context with a 10-minute timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Create a new request with the context
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request to AI service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("AI service returned non-200 status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	log.Printf("AI Service Response Status: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		log.Printf("AI service error response: %s", string(body))
		return "", fmt.Errorf("AI service returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var response struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
		Error string `json:"error,omitempty"`
	}

	// First try to parse as JSON
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("Failed to parse JSON response, trying to extract base64 directly")
		// Try to save the raw response as it might be the image data
		return s.saveImage(strings.TrimSpace(string(body)))
	}

	if response.Error != "" {
		return "", fmt.Errorf("AI service error: %s", response.Error)
	}

	// Extract content from message
	content := response.Message.Content
	if content == "" {
		return "", fmt.Errorf("empty response content from AI service")
	}

	// Clean up the response - remove markdown code blocks if present
	if strings.Contains(content, "```") {
		parts := strings.Split(content, "```")
		if len(parts) >= 2 {
			content = parts[1]
		}
	}

	// Remove any non-base64 characters
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	// Save the image using the storage provider
	return s.saveImage(content)
}
