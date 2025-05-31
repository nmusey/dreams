package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"dreams/handlers"
	"dreams/models"
	"dreams/services"
	"dreams/services/storage"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config holds all configuration for the application
type Config struct {
	DatabaseURL string
	Port        string
	AIApiHost   string
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
		DatabaseURL: getEnv("DATABASE_URL", "host=db user=postgres password=postgres dbname=dreams port=5432 sslmode=disable"),
		Port:        getEnv("PORT", "8080"),
		AIApiHost:   getEnv("AI_API_HOST", "http://localhost:11434"),
		AIModelName: getEnv("AI_MODEL_NAME", "llava"),

		// Storage configuration
		StorageType:    storageType,
		LocalDirectory: getEnv("STORAGE_LOCAL_DIR", filepath.Join(cwd, "images")),
		S3Bucket:       os.Getenv("AWS_S3_BUCKET"),
		S3Region:       os.Getenv("AWS_REGION"),
		S3AccessKey:    os.Getenv("AWS_ACCESS_KEY_ID"),
		S3SecretKey:    os.Getenv("AWS_SECRET_ACCESS_KEY"),
		S3Endpoint:     os.Getenv("AWS_S3_ENDPOINT"),
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

// CORS middleware
var corsMiddleware = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Vary", "Origin")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
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

	aiService := services.NewAIService(config.AIApiHost, config.AIModelName, storageProvider)

	queueService := services.NewQueueService(aiService, db)
	queueService.Start()

	dreamHandler := handlers.NewDreamHandler(db, aiService, queueService)

	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/api/dreams", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			dreamHandler.HandleGetAll(w, r)
		case http.MethodPost:
			dreamHandler.HandleCreate(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/dreams/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/dreams/")
		
		// Handle generate-image endpoint
		if strings.HasSuffix(path, "/generate-image") {
			dreamHandler.HandleGenerateImage(w, r)
			return
		}

		// Handle status check endpoint
		if strings.HasSuffix(path, "/status") {
			dreamHandler.HandleCheckImageStatus(w, r)
			return
		}

		// Extract ID for other operations
		parts := strings.Split(path, "/")
		if len(parts) == 0 || parts[0] == "" {
			http.NotFound(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			dreamHandler.HandleGetById(w, r)
		case http.MethodPut:
			dreamHandler.HandleUpdate(w, r)
		case http.MethodDelete:
			dreamHandler.HandleDelete(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Wrap the mux with CORS middleware
	handler := corsMiddleware(mux)

	// Set up file server for local images if using local storage
	if config.StorageType == storage.StorageTypeLocal {
		fs := http.FileServer(http.Dir(config.LocalDirectory))
		http.Handle("/images/", http.StripPrefix("/images/", fs))
	}

	log.Printf("Server starting on port %s...\n", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, handler))
}
