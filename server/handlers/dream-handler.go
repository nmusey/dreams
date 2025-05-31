package handlers

import (
	"dreams/models"
	"dreams/services"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type DreamHandler struct {
	db           *gorm.DB
	aiService    *services.AIService
	queueService *services.QueueService
}

func NewDreamHandler(db *gorm.DB, aiService *services.AIService, queueService *services.QueueService) *DreamHandler {
	return &DreamHandler{
		db:           db,
		aiService:    aiService,
		queueService: queueService,
	}
}

func (h *DreamHandler) RegisterRoutes() {
	http.HandleFunc("GET /api/dreams", h.HandleGetAll)
	http.HandleFunc("GET /api/dreams/{id}", h.HandleGetById)
	http.HandleFunc("POST /api/dreams", h.HandleCreate)
	http.HandleFunc("POST /api/dreams/{id}/generate-image", h.HandleGenerateImage)
	http.HandleFunc("PUT /api/dreams/{id}", h.HandleUpdate)
	http.HandleFunc("DELETE /api/dreams/{id}", h.HandleDelete)
}

func (h *DreamHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	var dreams []models.Dream
	if err := h.db.Find(&dreams).Error; err != nil {
		log.Printf("Error fetching dreams: %v", err)
		http.Error(w, "Failed to fetch dreams", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(dreams); err != nil {
		log.Printf("Error encoding dreams: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *DreamHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var dream models.Dream
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &dream); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.db.Create(&dream).Error; err != nil {
		log.Printf("Error creating dream: %v", err)
		http.Error(w, "Failed to create dream", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(dream); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *DreamHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/dreams/")
	idStr = strings.TrimSuffix(idStr, "/")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid dream ID", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var updateData map[string]interface{}
	if err := json.Unmarshal(body, &updateData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var existingDream models.Dream
	if err := h.db.First(&existingDream, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Dream not found", http.StatusNotFound)
		} else {
			log.Printf("Error finding dream: %v", err)
			http.Error(w, "Failed to find dream", http.StatusInternalServerError)
		}
		return
	}

	if err := h.db.Model(&existingDream).Updates(updateData).Error; err != nil {
		log.Printf("Error updating dream: %v", err)
		http.Error(w, "Failed to update dream", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(existingDream); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *DreamHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/dreams/")
	idStr = strings.TrimSuffix(idStr, "/")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid dream ID", http.StatusBadRequest)
		return
	}

	var existingDream models.Dream
	if err := h.db.First(&existingDream, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Dream not found", http.StatusNotFound)
		} else {
			log.Printf("Error finding dream: %v", err)
			http.Error(w, "Failed to find dream", http.StatusInternalServerError)
		}
		return
	}

	if err := h.db.Delete(&existingDream).Error; err != nil {
		log.Printf("Error deleting dream: %v", err)
		http.Error(w, "Failed to delete dream", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *DreamHandler) HandleGetById(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/dreams/")
	idStr = strings.TrimSuffix(idStr, "/")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid dream ID", http.StatusBadRequest)
		return
	}

	var dream models.Dream
	if err := h.db.First(&dream, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Dream not found", http.StatusNotFound)
		} else {
			log.Printf("Error finding dream: %v", err)
			http.Error(w, "Failed to find dream", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(dream); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GenerateImageResponse is the response for the generate image endpoint
type GenerateImageResponse struct {
	Message       string `json:"message"`
	QueuePosition int    `json:"queuePosition,omitempty"`
}

func (h *DreamHandler) HandleGenerateImage(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleGenerateImage: Received request")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("HandleGenerateImage: Method not allowed")
		return
	}

	// Get the dream ID from the URL
	path := strings.TrimPrefix(r.URL.Path, "/api/dreams/")
	path = strings.TrimSuffix(path, "/generate-image")
	idStr := path
	if idStr == "" {
		http.Error(w, "Missing dream ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid dream ID", http.StatusBadRequest)
		return
	}

	// Get only the necessary fields from the database
	var dream models.Dream
	if err := h.db.Select("id, dream, image_url").First(&dream, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Dream not found", http.StatusNotFound)
			return
		}
		log.Printf("Error fetching dream %d: %v", id, err)
		http.Error(w, "Failed to fetch dream", http.StatusInternalServerError)
		return
	}

	log.Printf("HandleGenerateImage: Attempting to enqueue request for dream ID: %d", dream.ID)
	
	// Enqueue the image generation request
	position, err := h.queueService.EnqueueRequest(dream)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to enqueue request: %v", err)
		log.Printf("HandleGenerateImage: %s", errMsg)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   "Image generation already in progress",
			"message": errMsg,
		})
		return
	}

	log.Printf("HandleGenerateImage: Successfully enqueued request for dream ID: %d, position: %d", dream.ID, position)

	// Return the queue position to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(GenerateImageResponse{
		Message:       "Image generation queued successfully",
		QueuePosition: position,
	}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// HandleCheckImageStatus checks the status of an image generation request
func (h *DreamHandler) HandleCheckImageStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the dream ID from the URL
	path := strings.TrimPrefix(r.URL.Path, "/api/dreams/")
	path = strings.TrimSuffix(path, "/status")
	idStr := path
	if idStr == "" {
		http.Error(w, "Missing dream ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid dream ID", http.StatusBadRequest)
		return
	}

	// Get only the necessary fields from the database
	var result struct {
		ID       uint
		ImageURL string
	}
	
	// Use a more efficient query with only the fields we need
	if err := h.db.Model(&models.Dream{}).
		Select("id, image_url").
		Where("id = ?", id).
		First(&result).Error; err != nil {
		
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Dream not found", http.StatusNotFound)
			return
		}
		log.Printf("Error fetching dream status %d: %v", id, err)
		http.Error(w, "Failed to fetch dream status", http.StatusInternalServerError)
		return
	}

	// If the image URL is set, return the image URL
	if result.ImageURL != "" {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"status":   "completed",
			"imageUrl": result.ImageURL,
		}); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
		return
	}

	// Check if the dream is in the queue
	position, isInQueue := h.queueService.GetQueuePosition(uint(id))
	if isInQueue {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted) // 202 Accepted
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"status":        "processing",
			"message":       "Image generation in progress",
			"queuePosition": position,
		}); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
		return
	}

	// If we get here, the dream is not in the queue and has no image
	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
