package services

import (
	"container/list"
	"dreams/models"
	"fmt"
	"log"
	"sync"
	"time"
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
}

func NewQueueService(aiService *AIService) *QueueService {
	qs := &QueueService{
		queue:     list.New(),
		aiService: aiService,
	}
	go qs.processQueue()
	return qs
}

func (qs *QueueService) EnqueueRequest(dream models.Dream) (string, error) {
	resultCh := make(chan string, 1)
	errorCh := make(chan error, 1)

	request := QueuedImageRequest{
		Dream:    dream,
		ResultCh: resultCh,
		ErrorCh:  errorCh,
	}

	qs.mu.Lock()
	qs.queue.PushBack(request)
	qs.mu.Unlock()

	log.Printf("Enqueued image generation request for dream ID: %d. Queue length: %d", dream.ID, qs.queue.Len())

	// Wait for the result with a timeout
	select {
	case result := <-resultCh:
		return result, nil
	case err := <-errorCh:
		return "", err
	case <-time.After(5 * time.Minute): // Timeout after 5 minutes
		return "", fmt.Errorf("image generation timed out")
	}
}

func (qs *QueueService) processQueue() {
	qs.isRunning = true
	for qs.isRunning {
		qs.mu.Lock()
		if qs.queue.Len() == 0 {
			qs.mu.Unlock()
			time.Sleep(1 * time.Second)
			continue
		}

		// Get the next request
		element := qs.queue.Front()
		qs.queue.Remove(element)
		qs.mu.Unlock()

		request := element.Value.(QueuedImageRequest)
		log.Printf("Processing image generation request for dream ID: %d", request.Dream.ID)

		// Generate the image
		imageURL, err := qs.aiService.GenerateImage(request.Dream.Dream)
		if err != nil {
			log.Printf("Error generating image for dream ID %d: %v", request.Dream.ID, err)
			request.ErrorCh <- err
			continue
		}

		request.ResultCh <- imageURL
		log.Printf("Successfully generated image for dream ID: %d", request.Dream.ID)
	}
}

func (qs *QueueService) Stop() {
	qs.isRunning = false
}
