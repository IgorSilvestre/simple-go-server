package external

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// MapTilerResponse represents the top-level response from MapTiler Geocoding API
type MapTilerResponse struct {
	Type     string            `json:"type"`
	Features []MapTilerFeature `json:"features"`
	Query    interface{}       `json:"query,omitempty"`
}

// MapTilerFeature represents a single feature in the response
type MapTilerFeature struct {
	Type       string             `json:"type"`
	Properties MapTilerProperties `json:"properties"`
	Geometry   MapTilerGeometry   `json:"geometry"`
	BBox       []float64          `json:"bbox,omitempty"`
}

// MapTilerProperties contains the detailed location data
type MapTilerProperties struct {
	Name          string  `json:"name"`
	Label         string  `json:"label"`
	Score         float64 `json:"score,omitempty"`
	HouseNumber   string  `json:"housenumber,omitempty"`
	Street        string  `json:"street,omitempty"`
	Neighbourhood string  `json:"neighbourhood,omitempty"`
	Suburb        string  `json:"suburb,omitempty"`
	District      string  `json:"district,omitempty"`
	Postcode      string  `json:"postcode,omitempty"`
	City          string  `json:"city,omitempty"`
	County        string  `json:"county,omitempty"`
	State         string  `json:"state,omitempty"`
	Country       string  `json:"country,omitempty"`
	CountryCode   string  `json:"country_code,omitempty"`
	Region        string  `json:"region,omitempty"`
	RegionCode    string  `json:"region_code,omitempty"`
	Formatted     string  `json:"formatted,omitempty"`
	AddressLine1  string  `json:"address_line1,omitempty"`
	AddressLine2  string  `json:"address_line2,omitempty"`
	Category      string  `json:"category,omitempty"`
	Timezone      string  `json:"timezone,omitempty"`
	Result_type   string  `json:"result_type,omitempty"`
	Rank          float64 `json:"rank,omitempty"`
	PlaceType     string  `json:"place_type,omitempty"`
}

// MapTilerGeometry represents the geographic location
type MapTilerGeometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

// Function to fetch geocoding data from MapTiler API
func fetchMapTilerGeocodingData(address string) (*MapTilerResponse, error) {
	// URL encode the address
	encodedAddress := url.PathEscape(address)

	// Base URL for the MapTiler geocoding API
	baseURL := fmt.Sprintf("https://api.maptiler.com/geocoding/%s.json", encodedAddress)

	// Create URL with query parameters
	params := url.Values{}
	params.Add("autocomplete", "false")
	params.Add("fuzzyMatch", "true")
	params.Add("limit", "3")
	params.Add("key", os.Getenv("MAPTILER_API_KEY"))

	// Construct the full URL
	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Create and execute the request
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("non-200 response from MapTiler (%d): %s", resp.StatusCode, string(body))
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Unmarshal the response
	var result MapTilerResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	// Check if we got any results
	if len(result.Features) == 0 {
		return &result, fmt.Errorf("no geocoding results found for address: %s", address)
	}

	return &result, nil
}

// MapTilerGeocodingHandler handles requests to the MapTiler geocoding endpoint
func MapTilerGeocodingHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the address from query parameters
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Address parameter is required", http.StatusBadRequest)
		return
	}

	// Fetch geocoding data
	data, err := fetchMapTilerGeocodingData(address)

	// Set response content type for GeoJSON
	w.Header().Set("Content-Type", "application/geo+json")

	if err != nil {
		// If we have some data but there was an error, include both
		if data != nil {
			// Set content type back to regular JSON for error responses
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": err.Error(),
				"data":  data,
			})
		} else {
			// Just return the error if we have no data
			w.WriteHeader(http.StatusInternalServerError)
			// Set content type back to regular JSON for error responses
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": err.Error(),
			})
		}
		return
	}

	// Ensure the response is a valid GeoJSON by setting the type to "FeatureCollection"
	if data.Type == "" {
		data.Type = "FeatureCollection"
	}

	// Write the successful response in GeoJSON format
	json.NewEncoder(w).Encode(data)
}
