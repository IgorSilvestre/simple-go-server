package external

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
)

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
		// Optionally, you can remove the "raw" field after parsing if needed
		// delete(result, "raw")
	}

	return result, nil
}

// Function to parse key-value pairs from the "raw" data
func parseRawData(raw string) map[string]string {
	re := regexp.MustCompile(`(?m)^\s*([\w-]+):\s+(.+)$`)
	data := make(map[string]string)
	matches := re.FindAllStringSubmatch(raw, -1)

	for _, match := range matches {
		key := strings.TrimSpace(match[1])
		value := strings.TrimSpace(match[2])
		data[key] = value
	}

	return data
}

func WhoisHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	lastSlash := strings.LastIndex(path, "/")
	if lastSlash == -1 || lastSlash == len(path)-1 {
		http.Error(w, "Domain not provided", http.StatusBadRequest)
		return
	}
	domain := path[lastSlash+1:]

	data, err := fetchWhoisData(domain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
