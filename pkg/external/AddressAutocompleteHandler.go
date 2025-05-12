package external

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

var stateAbbreviations = map[string]string{
	"Acre":                "AC",
	"Alagoas":             "AL",
	"Amapá":               "AP",
	"Amazonas":            "AM",
	"Bahia":               "BA",
	"Ceará":               "CE",
	"Distrito Federal":    "DF",
	"Espírito Santo":      "ES",
	"Goiás":               "GO",
	"Maranhão":            "MA",
	"Mato Grosso":         "MT",
	"Mato Grosso do Sul":  "MS",
	"Minas Gerais":        "MG",
	"Pará":                "PA",
	"Paraíba":             "PB",
	"Paraná":              "PR",
	"Pernambuco":          "PE",
	"Piauí":               "PI",
	"Rio de Janeiro":      "RJ",
	"Rio Grande do Norte": "RN",
	"Rio Grande do Sul":   "RS",
	"Rondônia":            "RO",
	"Roraima":             "RR",
	"Santa Catarina":      "SC",
	"São Paulo":           "SP",
	"Sergipe":             "SE",
	"Tocantins":           "TO",
}

// AddressAutocompleteHandler handles autocomplete requests and returns suggestions
func AddressAutocompleteHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Missing 'q' query parameter", http.StatusBadRequest)
		return
	}

	// Create a cache key based on the query
	cacheKey := "address_autocomplete:" + query

	// Check if the data is in the cache
	if cachedData, found := GlobalCache.Get(cacheKey); found {
		// Use the cached data
		suggestions := cachedData.(*AutocompleteResponse)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(suggestions)
		return
	}

	sessionToken := r.URL.Query().Get("sessiontoken")
	if sessionToken == "" {
		sessionToken = generateSessionToken()
	}

	suggestions, err := getAutocompleteSuggestions(query, sessionToken)
	if err != nil {
		http.Error(w, "Error fetching autocomplete suggestions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Process each prediction to abbreviate the state if needed
	for i, prediction := range suggestions.Predictions {
		suggestions.Predictions[i].Description = abbreviateState(prediction.Description)
	}

	// Store the data in the cache (24 hours expiration)
	GlobalCache.Set(cacheKey, suggestions, 24*time.Hour)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(suggestions)
}

// generateSessionToken generates a new UUID token for autocomplete sessions
func generateSessionToken() string {
	return uuid.New().String()
}

// getAutocompleteSuggestions fetches suggestions from the Google Places Autocomplete API
func getAutocompleteSuggestions(input, sessionToken string) (*AutocompleteResponse, error) {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("Google API key is not set")
	}

	endpoint := "https://maps.googleapis.com/maps/api/place/autocomplete/json"
	params := url.Values{}
	params.Add("input", input)
	params.Add("key", apiKey)
	params.Add("sessiontoken", sessionToken)
	params.Add("language", "pt-BR") // Ensure the language is set to Portuguese

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

	var autocompleteResponse AutocompleteResponse
	if err := json.Unmarshal(body, &autocompleteResponse); err != nil {
		return nil, err
	}

	if autocompleteResponse.Status != "OK" && autocompleteResponse.Status != "ZERO_RESULTS" {
		return nil, fmt.Errorf("Google API error: %s", autocompleteResponse.Status)
	}

	return &autocompleteResponse, nil
}

// abbreviateState checks if the description already contains a state abbreviation.
// If it doesn't, it replaces the full state name with its abbreviation.
func abbreviateState(description string) string {
	// Check if the description already contains a state abbreviation
	for _, abbreviation := range stateAbbreviations {
		if strings.Contains(description, abbreviation) {
			// If an abbreviation is already present, return the description unchanged
			return description
		}
	}

	// If no abbreviation is found, replace the full state name with the abbreviation
	for state, abbreviation := range stateAbbreviations {
		if strings.Contains(description, state) {
			// Replace the full state name with its abbreviation
			description = strings.Replace(description, state, abbreviation, 1)
			break
		}
	}

	return description
}

// AutocompleteResponse represents the response structure from the Google Places API
type AutocompleteResponse struct {
	Predictions  []Prediction `json:"predictions"`
	Status       string       `json:"status"`
	ErrorMessage string       `json:"error_message,omitempty"`
}

// Prediction represents a single prediction in the autocomplete response
type Prediction struct {
	Description string `json:"description"`
	PlaceID     string `json:"place_id"`
}
