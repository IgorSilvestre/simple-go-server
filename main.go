package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/likexian/whois"
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

func whoisHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    domain := vars["domain"]
    if domain == "" {
        http.Error(w, "Domain parameter is required", http.StatusBadRequest)
        return
    }

    // Perform the WHOIS lookup
    result, err := whois.Whois(domain)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error fetching WHOIS data: %v", err), http.StatusInternalServerError)
        return
    }

    // Write the result to the response
    w.Header().Set("Content-Type", "text/plain")
    _, writeErr := w.Write([]byte(result))
    if writeErr != nil {
        log.Printf("Error writing response: %v", writeErr)
    }
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/whois/{domain}", whoisHandler).Methods("GET", "OPTIONS")
    r.HandleFunc("/", handler).Methods("GET", "OPTIONS")

    // Apply the CORS middleware
    r.Use(corsMiddleware)

    log.Println("Starting server on port 8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
