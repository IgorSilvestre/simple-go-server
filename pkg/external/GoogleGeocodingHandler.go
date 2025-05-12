package external

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

// GeocodingResponse represents the response structure from the Google Geocoding API
type GeocodingResponse struct {
	Results []GeocodingResult `json:"results"`
	Status  string            `json:"status"`
}

// GeocodingResult represents a single result in the geocoding response
type GeocodingResult struct {
	FormattedAddress  string             `json:"formatted_address"`
	Geometry          GeometryData       `json:"geometry"`
	PlaceID           string             `json:"place_id"`
	Types             []string           `json:"types"`
	AddressComponents []AddressComponent `json:"address_components"`
}

// GeometryData contains location information
type GeometryData struct {
	Location     LatLngData   `json:"location"`
	LocationType string       `json:"location_type"`
	Viewport     ViewportData `json:"viewport"`
}

// LatLngData contains latitude and longitude
type LatLngData struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// ViewportData contains the viewport information
type ViewportData struct {
	Northeast LatLngData `json:"northeast"`
	Southwest LatLngData `json:"southwest"`
}

// AddressComponent represents a component of an address
type AddressComponent struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

// GeocodingHandler handles geocoding requests and returns location data
func googleGeocodingHandler(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Missing 'address' query parameter", http.StatusBadRequest)
		return
	}

	// Create a cache key based on the address
	cacheKey := "google_geocoding:" + address

	// Check if the data is in the cache
	if cachedData, found := GlobalCache.Get(cacheKey); found {
		// Use the cached data
		geocodingData := cachedData.(*GeocodingResponse)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(geocodingData)
		return
	}

	geocodingData, err := getGeocodingData(address)
	if err != nil {
		http.Error(w, "Error fetching geocoding data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Store the data in the cache (24 hours expiration)
	GlobalCache.Set(cacheKey, geocodingData, 24*time.Hour)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(geocodingData)
}

// getGeocodingData fetches geocoding data from the Google Geocoding API
func getGeocodingData(address string) (*GeocodingResponse, error) {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("Google API key is not set")
	}

	endpoint := "https://maps.googleapis.com/maps/api/geocode/json"
	params := url.Values{}
	params.Add("address", address)
	params.Add("key", apiKey)

	apiURL := endpoint + "?" + params.Encode()

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Google API error: %s", string(body))
	}

	var geocodingResponse GeocodingResponse
	if err := json.Unmarshal(body, &geocodingResponse); err != nil {
		return nil, err
	}

	if geocodingResponse.Status != "OK" && geocodingResponse.Status != "ZERO_RESULTS" {
		return nil, fmt.Errorf("Google API error: %s", geocodingResponse.Status)
	}

	return &geocodingResponse, nil
}
