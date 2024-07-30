package external

import (
    "github.com/gorilla/mux"
)

func RegisterExternalRoutes(r *mux.Router) {
    subrouter := r.PathPrefix("/external").Subrouter()
    subrouter.HandleFunc("/whois/{domain}", WhoisHandler).Methods("GET", "OPTIONS")
}

