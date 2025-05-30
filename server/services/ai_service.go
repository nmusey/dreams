package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type AIService struct {
	baseURL string
}

func NewAIService() *AIService {
	baseURL := os.Getenv("OLLAMA_HOST")
	if baseURL == "" {
		baseURL = "http://localhost:11434" // Default to localhost if not specified
	}
	return &AIService{
		baseURL: baseURL,
	}
}

type ImageGenerationRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type ImageGenerationResponse struct {
	Response string `json:"response"`
}

func (s *AIService) GenerateImage(dreamContent string) (string, error) {
	// Create a surreal, dream-like prompt based on the dream content
	prompt := fmt.Sprintf("A surreal, dream-like artistic interpretation of the following dream: %s. Style: ethereal, mystical, dreamscape, using soft colors and flowing forms", dreamContent)

	reqBody := ImageGenerationRequest{
		Model:  "llava", // Using llava model for image generation
		Prompt: prompt,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/generate", s.baseURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result ImageGenerationResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %v", err)
	}

	// The response from Ollama will be a base64 encoded image
	// We'll return this directly as it can be used in an HTML img tag with data:image/jpeg;base64,
	return fmt.Sprintf("data:image/jpeg;base64,%s", result.Response), nil
}
