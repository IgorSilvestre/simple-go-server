package external

import (
  "encoding/json"
  "fmt"
  "io"
  "net/http"
  "net/url"
  "os"

  "github.com/google/uuid"
)

func AddressAutocompleteHandler(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("q")
    if query == "" {
        http.Error(w, "Missing 'q' query parameter", http.StatusBadRequest)
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

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(suggestions)
}

func generateSessionToken() string {
    return uuid.New().String()
}

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

type AutocompleteResponse struct {
    Predictions []Prediction `json:"predictions"`
    Status      string       `json:"status"`
    ErrorMessage string      `json:"error_message,omitempty"`
}

type Prediction struct {
    Description string `json:"description"`
    PlaceID     string `json:"place_id"`
}
