package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"dreams/handlers"
	"dreams/models"
	"dreams/services"
	"dreams/services/storage"

	"github.com/rs/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config holds all configuration for the application
type Config struct {
	DatabaseURL string
	Port        string
	AIApiHost   string
	AIEndpoint  string
	AIModelName string

	// Storage configuration
	StorageType    storage.StorageType
	LocalDirectory string
	S3Bucket       string
	S3Region       string
	S3AccessKey    string
	S3SecretKey    string
	S3Endpoint     string
}

// loadConfig loads configuration from environment variables with defaults
func loadConfig() Config {
	// Default to local storage
	storageType := storage.StorageTypeLocal
	if os.Getenv("STORAGE_TYPE") == "s3" {
		storageType = storage.StorageTypeS3
	}

	// Get the current working directory for local storage
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}

	return Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:localhost:5432/dreams?sslmode=disable"),
		Port:        getEnv("PORT", "8080"),
		AIApiHost:   getEnv("AI_API_HOST", "http://localhost:11434"),
		AIEndpoint:  getEnv("AI_API_ENDPOINT", "/api/generate"),
		AIModelName: getEnv("AI_MODEL_NAME", "stable-diffusion-1.5"),
		StorageType:    storageType,
		LocalDirectory: getEnv("LOCAL_DIRECTORY", filepath.Join(cwd, "images")),
		S3Bucket:       getEnv("S3_BUCKET", ""),
		S3Region:       getEnv("S3_REGION", "us-east-1"),
		S3AccessKey:    getEnv("S3_ACCESS_KEY", ""),
		S3SecretKey:    getEnv("S3_SECRET_KEY", ""),
		S3Endpoint:     getEnv("S3_ENDPOINT", ""),
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Create CORS middleware
var corsMiddleware = func() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://localhost:3000"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Requested-With", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           3600,
		Debug:            false,
	})
}

func main() {
	config := loadConfig()

	// Initialize storage provider
	storageCfg := storage.Config{
		Type:           config.StorageType,
		LocalDirectory: config.LocalDirectory,
		BucketName:     config.S3Bucket,
		Region:         config.S3Region,
		AccessKey:      config.S3AccessKey,
		SecretKey:      config.S3SecretKey,
		Endpoint:       config.S3Endpoint,
	}

	storageProvider, err := storage.NewStorage(storageCfg)
	if err != nil {
		log.Fatalf("Failed to initialize storage provider: %v", err)
	}

	// Create local directory if using local storage
	if config.StorageType == storage.StorageTypeLocal {
		if err := os.MkdirAll(config.LocalDirectory, 0755); err != nil {
			log.Fatalf("Failed to create local storage directory: %v", err)
		}
	}

	db, err := gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&models.Dream{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	aiService := services.NewAIService(config.AIApiHost, config.AIEndpoint, config.AIModelName, storageProvider)

	queueService := services.NewQueueService(aiService, db)
	queueService.Start()

	dreamHandler := handlers.NewDreamHandler(db, aiService, queueService)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/dreams", dreamHandler.HandleGetAll)
	mux.HandleFunc("POST /api/dreams", dreamHandler.HandleCreate)
	mux.HandleFunc("GET /api/dreams/{id}", dreamHandler.HandleGetById)
	mux.HandleFunc("PUT /api/dreams/{id}", dreamHandler.HandleUpdate)
	mux.HandleFunc("DELETE /api/dreams/{id}", dreamHandler.HandleDelete)
	mux.HandleFunc("POST /api/dreams/{id}/generate-image", dreamHandler.HandleGenerateImage)
	mux.HandleFunc("GET /api/dreams/{id}/status", dreamHandler.HandleCheckImageStatus)

	if config.StorageType == storage.StorageTypeLocal {
		fs := http.FileServer(http.Dir(config.LocalDirectory))
		mux.Handle("/images/", http.StripPrefix("/images/", fs))
	}

	// Wrap the mux with CORS middleware
	c := corsMiddleware()
	handler := c.Handler(mux)

	log.Printf("Server starting on port %s...\n", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, handler))
}
