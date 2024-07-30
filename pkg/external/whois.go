package external

import (
    "net/http"
    "fmt"
    "log"
    "github.com/likexian/whois"
    "github.com/gorilla/mux"
)

func WhoisHandler(w http.ResponseWriter, r *http.Request) {
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

