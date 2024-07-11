package main

import (
    "fmt"
    "net/http"
    "log"
    "github.com/gorilla/mux"
    "github.com/likexian/whois"
)

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
    w.Write([]byte(result))
}

func main () {
    r := mux.NewRouter()
    
    r.HandleFunc("/whois/{domain}", whoisHandler).Methods("GET")

    r.HandleFunc("/", handler)

    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    } else {
        log.Println("Server started on port 8080")
    }
}
