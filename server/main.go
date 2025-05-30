package main

import (
	"log"
	"net/http"
	"os"

	"dreams/handlers"
	"dreams/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming request: %s %s", r.Method, r.URL.Path)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Get database connection string from environment variable
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=db user=postgres password=postgres dbname=dreams port=5432 sslmode=disable"
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(&models.Dream{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize dream handler
	dreamHandler := handlers.NewDreamHandler(db)

	// Create a new mux for handling routes
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("GET /api/dreams", dreamHandler.HandleGetAll)
	mux.HandleFunc("POST /api/dreams", dreamHandler.HandleCreate)
	mux.HandleFunc("PUT /api/dreams/{id}", dreamHandler.HandleUpdate)
	mux.HandleFunc("DELETE /api/dreams/{id}", dreamHandler.HandleDelete)

	// Wrap the mux with CORS middleware
	handler := corsMiddleware(mux)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
