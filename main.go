package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/igorsilvestre/simple-go-server/pkg/external"
	"log"
	"net/http"
	"os"

  "github.com/joho/godotenv"
)

// CORS Middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func main() {
  if os.Getenv("RAILWAY_ENVIRONMENT") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file")
		}
	}

	r := mux.NewRouter()

	// Register external routes
	external.RegisterExternalRoutes(r)

	// Main routes
	r.HandleFunc("/", handler).Methods("GET", "OPTIONS")

	// Apply the CORS middleware
	r.Use(corsMiddleware)

	log.Println("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
