package external

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// NominatimGeocodingResult represents a single result from the Nominatim API
type NominatimGeocodingResult struct {
	PlaceID     int      `json:"place_id"`
	Licence     string   `json:"licence"`
	OsmType     string   `json:"osm_type"`
	OsmID       int      `json:"osm_id"`
	Lat         string   `json:"lat"`
	Lon         string   `json:"lon"`
	Class       string   `json:"class"`
	Type        string   `json:"type"`
	PlaceRank   int      `json:"place_rank"`
	Importance  float64  `json:"importance"`
	AddressType string   `json:"addresstype"`
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	BoundingBox []string `json:"boundingbox"`
}

// NominatimGeocodingHandler handles geocoding requests using the Nominatim API
func NominatimGeocodingHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the address from query parameters
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Address parameter is required", http.StatusBadRequest)
		return
	}

	// Fetch geocoding data from Nominatim
	results, err := fetchNominatimGeocodingData(address)
	if err != nil {
		http.Error(w, "Error fetching geocoding data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response content type and return the results
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// fetchNominatimGeocodingData fetches geocoding data from the Nominatim API
func fetchNominatimGeocodingData(address string) ([]NominatimGeocodingResult, error) {
	// Base URL for the Nominatim API
	baseURL := "https://nominatim.openstreetmap.org/search"

	// Create URL with query parameters
	params := url.Values{}
	params.Add("q", address)
	params.Add("format", "json")

	// Construct the full URL
	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Create a client with custom headers (Nominatim requires a User-Agent)
	client := &http.Client{}
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set a User-Agent as required by Nominatim's usage policy
	req.Header.Set("User-Agent", "simple-go-server")

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("non-200 response from Nominatim (%d): %s", resp.StatusCode, string(body))
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Unmarshal the response
	var results []NominatimGeocodingResult
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	// Check if we got any results
	if len(results) == 0 {
		return results, fmt.Errorf("no geocoding results found for address: %s", address)
	}

	return results, nil
}
