package services

import (
	"container/list"
	"context"
	"dreams/models"
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/gorm"
)

type QueuedImageRequest struct {
	Dream    models.Dream
	ResultCh chan string
	ErrorCh  chan error
}

type QueueService struct {
	queue     *list.List
	mu        sync.Mutex
	aiService *AIService
	isRunning bool
	db        *gorm.DB
	// Track active requests by dream ID
	activeRequests map[uint]*QueuedImageRequest
}

func NewQueueService(aiService *AIService, db *gorm.DB) *QueueService {
	qs := &QueueService{
		queue:          list.New(),
		aiService:      aiService,
		db:             db,
		activeRequests: make(map[uint]*QueuedImageRequest),
	}
	return qs
}

// Start starts the queue processor
func (qs *QueueService) Start() {
	qs.isRunning = true
	go qs.processQueue()
	log.Println("Queue processor started")
}

// EnqueueRequest adds a new image generation request to the queue and returns the position in the queue
func (qs *QueueService) EnqueueRequest(dream models.Dream) (int, error) {
	qs.mu.Lock()
	defer qs.mu.Unlock()

	// Check if there's already a request for this dream
	if _, exists := qs.activeRequests[dream.ID]; exists {
		return -1, fmt.Errorf("image generation already in progress for dream %d", dream.ID)
	}

	// Create a new request
	request := &QueuedImageRequest{
		Dream:    dream,
		ResultCh: make(chan string, 1),
		ErrorCh:  make(chan error, 1),
	}

	// Add to active requests and queue
	qs.activeRequests[dream.ID] = request
	position := qs.queue.Len()
	qs.queue.PushBack(request)
	log.Printf("Enqueued dream %d (position: %d, queue length: %d)", 
		dream.ID, position+1, qs.queue.Len())

	// Return the position in the queue (1-based index)
	return position + 1, nil
}

func (qs *QueueService) processQueue() {
	qs.isRunning = true

	for qs.isRunning {
		// Get the next request atomically
		var request *QueuedImageRequest
		qs.mu.Lock()
		if element := qs.queue.Front(); element != nil {
			request = qs.queue.Remove(element).(*QueuedImageRequest)
		}
		qs.mu.Unlock()

		if request == nil {
			time.Sleep(1 * time.Second) // Prevent busy waiting
			continue
		}

		log.Printf("Processing dream %d", request.Dream.ID)

		// Process the request in a goroutine
		go func(req *QueuedImageRequest) {
			defer func() {
				qs.mu.Lock()
				delete(qs.activeRequests, req.Dream.ID)
				qs.mu.Unlock()
				close(req.ResultCh)
				close(req.ErrorCh)
			}()

			// Generate the image with a timeout
			imagePath, err := qs.aiService.GenerateImage(req.Dream.Dream)
			if err != nil {
				req.ErrorCh <- fmt.Errorf("error generating image: %w", err)
				return
			}

			// Update the dream with the generated image URL
			err = qs.db.Transaction(func(tx *gorm.DB) error {
				// Set a timeout for the database operation
				dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				// Use the transaction with timeout context
				tx = tx.WithContext(dbCtx)
				if err := tx.Model(&models.Dream{}).
					Where("id = ?", req.Dream.ID).
					Update("image_url", imagePath).Error; err != nil {
					return fmt.Errorf("failed to update dream with image URL: %w", err)
				}
				return nil
			})

			if err != nil {
				req.ErrorCh <- fmt.Errorf("failed to update dream: %w", err)
				return
			}

			// Send the result
			req.ResultCh <- imagePath
		}(request)

		// Wait for either the result or an error
		select {
		case <-request.ResultCh:
			log.Printf("Successfully processed dream %d", request.Dream.ID)
		case err := <-request.ErrorCh:
			log.Printf("Error processing dream %d: %v", request.Dream.ID, err)
		}
	}
}

// GetQueuePosition returns the position of a dream in the queue and a boolean indicating if it's in the queue
func (qs *QueueService) GetQueuePosition(dreamID uint) (int, bool) {
	qs.mu.Lock()
	defer qs.mu.Unlock()

	// First check active requests (currently processing)
	if _, exists := qs.activeRequests[dreamID]; exists {
		return 0, true // Currently processing, so position is 0
	}

	// Then check the queue
	position := 0
	for e := qs.queue.Front(); e != nil; e = e.Next() {
		position++
		if req, ok := e.Value.(*QueuedImageRequest); ok && req.Dream.ID == dreamID {
			return position, true
		}
	}

	// Not found in queue
	return -1, false
}

// Stop stops the queue processor
func (qs *QueueService) Stop() {
	qs.isRunning = false
}
