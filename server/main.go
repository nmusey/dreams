package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"dreams/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming request: %s %s", r.Method, r.URL.Path)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {
	// Configure database connection
	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable"

	log.Printf("Connecting to database with DSN: %s", dsn)

	// Configure GORM logger
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		},
	)

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(&models.Dream{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Setup routes
	http.HandleFunc("/api/dreams", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			var dreams []models.Dream
			log.Printf("Fetching all dreams")
			if err := db.Find(&dreams).Error; err != nil {
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

		case http.MethodPost:
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
			if err := db.Create(&dream).Error; err != nil {
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

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
