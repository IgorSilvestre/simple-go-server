package external

import (
    "github.com/gorilla/mux"
)

func RegisterExternalRoutes(r *mux.Router) {
    subrouter := r.PathPrefix("/external").Subrouter()
    subrouter.HandleFunc("/whois/{domain}", WhoisHandler).Methods("GET", "OPTIONS")
    subrouter.HandleFunc("/autocomplete-address", AddressAutocompleteHandler).Methods("GET", "OPTIONS")
    subrouter.HandleFunc("/send-email", SendEmailHandler).Methods("POST", "OPTIONS")
}

