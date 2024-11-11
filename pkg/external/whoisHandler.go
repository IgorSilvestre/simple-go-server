package external

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "regexp"
    "strings"

    "github.com/gorilla/mux"
    "github.com/likexian/whois"
)

// Function to fetch WHOIS data using the JSONWHOIS API
func fetchWhoisData(domain string) (map[string]interface{}, error) {
    req, err := http.NewRequest("GET", "https://jsonwhois.com/api/v1/whois?domain="+domain, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Add("Accept", "application/json")
    req.Header.Add("Authorization", "Token token="+os.Getenv("JSONWHOIS_APIKEY"))

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("non-200 response: %d", resp.StatusCode)
    }

    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    // Extract and parse key-value pairs from "raw" field
    if raw, ok := result["raw"].(string); ok {
        parsedData := parseRawData(raw)
        // Merge parsed data into the main result
        for key, value := range parsedData {
            result[key] = value
        }
    }

    // Check if essential data is present
    if _, ownerExists := result["owner"]; !ownerExists {
        if _, ownerIDExists := result["ownerid"]; !ownerIDExists {
            return result, fmt.Errorf("Essential WHOIS data missing from JSONWHOIS API response")
        }
    }

    return result, nil
}

// Function to fetch WHOIS data using the alternative method
func fetchWhoisDataAlternative(domain string) (map[string]interface{}, error) {
    // Perform the WHOIS lookup
    result, err := whois.Whois(domain)
    if err != nil {
        return nil, fmt.Errorf("Error fetching WHOIS data: %v", err)
    }

    // Initialize data map and include the raw WHOIS data
    data := make(map[string]interface{})
    data["raw"] = result

    // Check if the response contains "Permission denied" or similar messages
    if isPermissionDenied(result) {
        return data, fmt.Errorf("Permission denied by the WHOIS server for domain: %s", domain)
    }

    // Parse data["raw"] using parseRawData
    parsedData := parseRawData(result)

    // If no data was parsed, return an error
    if len(parsedData) == 0 {
        return data, fmt.Errorf("No WHOIS data available for domain: %s", domain)
    }

    // Merge parsedData into data
    for key, value := range parsedData {
        data[key] = value
    }

    return data, nil
}

// Function to parse key-value pairs from the "raw" data
func parseRawData(raw string) map[string]string {
    // Regular expression to match key-value pairs
    re := regexp.MustCompile(`(?m)^\s*([\w\-\. ]+):\s*(.+)$`)
    data := make(map[string]string)
    matches := re.FindAllStringSubmatch(raw, -1)

    for _, match := range matches {
        key := strings.TrimSpace(match[1])
        value := strings.TrimSpace(match[2])
        // Handle multiple values for the same key
        if existingValue, exists := data[key]; exists {
            data[key] = existingValue + ", " + value
        } else {
            data[key] = value
        }
    }

    return data
}

// Helper function to check if the WHOIS server denied permission
func isPermissionDenied(raw string) bool {
    lowerRaw := strings.ToLower(raw)
    return strings.Contains(lowerRaw, "permission denied") ||
        strings.Contains(lowerRaw, "access denied") ||
        strings.Contains(lowerRaw, "quota exceeded") ||
        strings.Contains(lowerRaw, "refused")
}

func WhoisHandler(w http.ResponseWriter, r *http.Request) {
    // Use mux.Vars to extract the domain parameter
    vars := mux.Vars(r)
    domain := vars["domain"]
    if domain == "" {
        http.Error(w, "Domain parameter is required", http.StatusBadRequest)
        return
    }

    var data map[string]interface{}
    var err error

    // Try the first method
    data, err = fetchWhoisData(domain)
    if err != nil {
        // Log the error (optional)
        fmt.Printf("Error with JSONWHOIS API: %v\n", err)

        // First method failed, try the alternative method
        data, err = fetchWhoisDataAlternative(domain)
        if err != nil {
            // Log the error (optional)
            fmt.Printf("Error with alternative WHOIS method: %v\n", err)
            // Return the data with the error message
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(map[string]interface{}{
                "error": err.Error(),
                "data":  data,
            })
            return
        }
    }

    // Write the response in JSON format
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}
