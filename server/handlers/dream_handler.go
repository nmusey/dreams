package handlers

import (
	"dreams/models"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type DreamHandler struct {
	db *gorm.DB
}

func NewDreamHandler(db *gorm.DB) *DreamHandler {
	return &DreamHandler{db: db}
}

func (h *DreamHandler) RegisterRoutes() {
	http.HandleFunc("GET /api/dreams", h.HandleGetAll)
	http.HandleFunc("POST /api/dreams", h.HandleCreate)
	http.HandleFunc("PUT /api/dreams/{id}", h.HandleUpdate)
	http.HandleFunc("DELETE /api/dreams/{id}", h.HandleDelete)
}

func (h *DreamHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	var dreams []models.Dream
	log.Printf("Fetching all dreams")
	if err := h.db.Find(&dreams).Error; err != nil {
		log.Printf("Error fetching dreams: %v", err)
		http.Error(w, "Failed to fetch dreams", http.StatusInternalServerError)
		return
	}
	log.Printf("Found %d dreams", len(dreams))
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(dreams); err != nil {
		log.Printf("Error encoding dreams: %v", err)
		http.Error(w, "Failed to encode dreams", http.StatusInternalServerError)
		return
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
	log.Printf("Received request body: %s", string(body))

	if err := json.Unmarshal(body, &dream); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Creating dream: %+v", dream)
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
		return
	}
}

func (h *DreamHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	// Extract dream ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/api/dreams/")
	idStr = strings.TrimSuffix(idStr, "/")
	log.Printf("Received request for dream ID: %s", idStr)

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("Invalid dream ID: %v", err)
		http.Error(w, "Invalid dream ID", http.StatusBadRequest)
		return
	}

	var dream models.Dream
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Printf("Received request body: %s", string(body))

	if err := json.Unmarshal(body, &dream); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Attempting to update dream %d with content: %s", id, dream.Dream)

	// First check if the dream exists
	var existingDream models.Dream
	if err := h.db.First(&existingDream, id).Error; err != nil {
		log.Printf("Error finding dream: %v", err)
		http.Error(w, "Dream not found", http.StatusNotFound)
		return
	}

	// Update the dream
	existingDream.Dream = dream.Dream
	if err := h.db.Save(&existingDream).Error; err != nil {
		log.Printf("Error updating dream: %v", err)
		http.Error(w, "Failed to update dream", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully updated dream %d", id)

	// Return the updated dream
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(existingDream); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *DreamHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	// Extract dream ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/api/dreams/")
	idStr = strings.TrimSuffix(idStr, "/")
	log.Printf("Received request for dream ID: %s", idStr)

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("Invalid dream ID: %v", err)
		http.Error(w, "Invalid dream ID", http.StatusBadRequest)
		return
	}

	log.Printf("Attempting to delete dream %d", id)

	// First check if the dream exists
	var existingDream models.Dream
	if err := h.db.First(&existingDream, id).Error; err != nil {
		log.Printf("Error finding dream: %v", err)
		http.Error(w, "Dream not found", http.StatusNotFound)
		return
	}

	if err := h.db.Delete(&existingDream).Error; err != nil {
		log.Printf("Error deleting dream: %v", err)
		http.Error(w, "Failed to delete dream", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully deleted dream %d", id)
	w.WriteHeader(http.StatusNoContent)
}
