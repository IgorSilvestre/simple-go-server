package external

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

// GeoapifyResponse represents the top-level response from Geoapify Geocoding API
type GeoapifyResponse struct {
	Type     string            `json:"type"`
	Features []GeoapifyFeature `json:"features"`
	Query    GeoapifyQuery     `json:"query"`
}

// GeoapifyFeature represents a single feature in the response
type GeoapifyFeature struct {
	Type       string             `json:"type"`
	Properties GeoapifyProperties `json:"properties"`
	Geometry   GeoapifyGeometry   `json:"geometry"`
	BBox       []float64          `json:"bbox,omitempty"`
}

// GeoapifyProperties contains the detailed location data
type GeoapifyProperties struct {
	Datasource    map[string]interface{} `json:"datasource"`
	Country       string                 `json:"country"`
	CountryCode   string                 `json:"country_code"`
	State         string                 `json:"state"`
	County        string                 `json:"county"`
	City          string                 `json:"city"`
	Postcode      string                 `json:"postcode"`
	Suburb        string                 `json:"suburb,omitempty"`
	Street        string                 `json:"street"`
	Housenumber   string                 `json:"housenumber,omitempty"`
	Lon           float64                `json:"lon"`
	Lat           float64                `json:"lat"`
	StateCode     string                 `json:"state_code,omitempty"`
	ResultType    string                 `json:"result_type"`
	Formatted     string                 `json:"formatted"`
	AddressLine1  string                 `json:"address_line1"`
	AddressLine2  string                 `json:"address_line2"`
	Category      string                 `json:"category,omitempty"`
	Timezone      map[string]interface{} `json:"timezone,omitempty"`
	PlusCode      string                 `json:"plus_code,omitempty"`
	PlusCodeShort string                 `json:"plus_code_short,omitempty"`
	Rank          map[string]interface{} `json:"rank,omitempty"`
	PlaceId       string                 `json:"place_id,omitempty"`
	// Using map for nested objects that have variable structure
}

// GeoapifyGeometry represents the geographic location
type GeoapifyGeometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

// GeoapifyQuery represents the query information from the response
type GeoapifyQuery struct {
	Text   string                 `json:"text"`
	Parsed map[string]interface{} `json:"parsed,omitempty"`
}

// Function to fetch geocoding data from Geoapify API
func fetchGeocodingData(address string) (*GeoapifyResponse, error) {
	// Base URL for the Geoapify geocoding API
	baseURL := "https://api.geoapify.com/v1/geocode/search"

	// Create URL with query parameters
	params := url.Values{}
	params.Add("text", address)
	params.Add("apiKey", os.Getenv("GEOAPIFY_API_KEY"))

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
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("non-200 response from Geoapify (%d): %s", resp.StatusCode, string(body))
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Unmarshal the response
	var result GeoapifyResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	// Check if we got any results
	if len(result.Features) == 0 {
		return &result, fmt.Errorf("no geocoding results found for address: %s", address)
	}

	return &result, nil
}

// GeocodeGeoapifyHandler handles requests to the Geoapify geocoding endpoint
func GeoapifyGeocodingHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the address from query parameters
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Address parameter is required", http.StatusBadRequest)
		return
	}

	// Fetch geocoding data
	data, err := fetchGeocodingData(address)

	// Set response content type
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		// If we have some data but there was an error, include both
		if data != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": err.Error(),
				"data":  data,
			})
		} else {
			// Just return the error if we have no data
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": err.Error(),
			})
		}
		return
	}

	// Write the successful response in JSON format
	json.NewEncoder(w).Encode(data)
}
