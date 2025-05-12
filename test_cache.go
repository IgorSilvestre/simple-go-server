package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	// Test the WhoisHandler
	testWhoisCache()

	// Test the GoogleGeocodingHandler
	testGoogleGeocodingCache()

	// Test the GeoapifyGeocodingHandler
	testGeoapifyGeocodingCache()

	// Test the NominatimGeocodingHandler
	testNominatimGeocodingCache()

	// Test the MapTilerGeocodingHandler
	testMapTilerGeocodingCache()

	// Test the AddressAutocompleteHandler
	testAddressAutocompleteCache()
}

func testWhoisCache() {
	fmt.Println("Testing WhoisHandler cache...")

	// Make the first request
	start := time.Now()
	resp1, err := http.Get("http://localhost:8080/external/whois/google.com")
	if err != nil {
		fmt.Printf("Error making first request: %v\n", err)
		return
	}
	defer resp1.Body.Close()
	duration1 := time.Since(start)

	// Make the second request (should be cached)
	start = time.Now()
	resp2, err := http.Get("http://localhost:8080/external/whois/google.com")
	if err != nil {
		fmt.Printf("Error making second request: %v\n", err)
		return
	}
	defer resp2.Body.Close()
	duration2 := time.Since(start)

	// Compare response times
	fmt.Printf("First request: %v\n", duration1)
	fmt.Printf("Second request: %v\n", duration2)
	fmt.Printf("Cache speedup: %.2fx\n", float64(duration1)/float64(duration2))
	fmt.Println()
}

func testGoogleGeocodingCache() {
	fmt.Println("Testing GoogleGeocodingHandler cache...")

	// Make the first request
	start := time.Now()
	resp1, err := http.Get("http://localhost:8080/external/geocode?address=1600+Amphitheatre+Parkway,+Mountain+View,+CA")
	if err != nil {
		fmt.Printf("Error making first request: %v\n", err)
		return
	}
	defer resp1.Body.Close()
	duration1 := time.Since(start)

	// Make the second request (should be cached)
	start = time.Now()
	resp2, err := http.Get("http://localhost:8080/external/geocode?address=1600+Amphitheatre+Parkway,+Mountain+View,+CA")
	if err != nil {
		fmt.Printf("Error making second request: %v\n", err)
		return
	}
	defer resp2.Body.Close()
	duration2 := time.Since(start)

	// Compare response times
	fmt.Printf("First request: %v\n", duration1)
	fmt.Printf("Second request: %v\n", duration2)
	fmt.Printf("Cache speedup: %.2fx\n", float64(duration1)/float64(duration2))
	fmt.Println()
}

func testGeoapifyGeocodingCache() {
	fmt.Println("Testing GeoapifyGeocodingHandler cache...")

	// Make the first request
	start := time.Now()
	resp1, err := http.Get("http://localhost:8080/external/geocode-geoapify?address=1600+Amphitheatre+Parkway,+Mountain+View,+CA")
	if err != nil {
		fmt.Printf("Error making first request: %v\n", err)
		return
	}
	defer resp1.Body.Close()
	duration1 := time.Since(start)

	// Make the second request (should be cached)
	start = time.Now()
	resp2, err := http.Get("http://localhost:8080/external/geocode-geoapify?address=1600+Amphitheatre+Parkway,+Mountain+View,+CA")
	if err != nil {
		fmt.Printf("Error making second request: %v\n", err)
		return
	}
	defer resp2.Body.Close()
	duration2 := time.Since(start)

	// Compare response times
	fmt.Printf("First request: %v\n", duration1)
	fmt.Printf("Second request: %v\n", duration2)
	fmt.Printf("Cache speedup: %.2fx\n", float64(duration1)/float64(duration2))
	fmt.Println()
}

func testNominatimGeocodingCache() {
	fmt.Println("Testing NominatimGeocodingHandler cache...")

	// Make the first request
	start := time.Now()
	resp1, err := http.Get("http://localhost:8080/external/geocode-nominatim?address=1600+Amphitheatre+Parkway,+Mountain+View,+CA")
	if err != nil {
		fmt.Printf("Error making first request: %v\n", err)
		return
	}
	defer resp1.Body.Close()
	duration1 := time.Since(start)

	// Make the second request (should be cached)
	start = time.Now()
	resp2, err := http.Get("http://localhost:8080/external/geocode-nominatim?address=1600+Amphitheatre+Parkway,+Mountain+View,+CA")
	if err != nil {
		fmt.Printf("Error making second request: %v\n", err)
		return
	}
	defer resp2.Body.Close()
	duration2 := time.Since(start)

	// Compare response times
	fmt.Printf("First request: %v\n", duration1)
	fmt.Printf("Second request: %v\n", duration2)
	fmt.Printf("Cache speedup: %.2fx\n", float64(duration1)/float64(duration2))
	fmt.Println()
}

func testMapTilerGeocodingCache() {
	fmt.Println("Testing MapTilerGeocodingHandler cache...")

	// Make the first request
	start := time.Now()
	resp1, err := http.Get("http://localhost:8080/external/geocode-maptiler?address=1600+Amphitheatre+Parkway,+Mountain+View,+CA")
	if err != nil {
		fmt.Printf("Error making first request: %v\n", err)
		return
	}
	defer resp1.Body.Close()
	duration1 := time.Since(start)

	// Make the second request (should be cached)
	start = time.Now()
	resp2, err := http.Get("http://localhost:8080/external/geocode-maptiler?address=1600+Amphitheatre+Parkway,+Mountain+View,+CA")
	if err != nil {
		fmt.Printf("Error making second request: %v\n", err)
		return
	}
	defer resp2.Body.Close()
	duration2 := time.Since(start)

	// Compare response times
	fmt.Printf("First request: %v\n", duration1)
	fmt.Printf("Second request: %v\n", duration2)
	fmt.Printf("Cache speedup: %.2fx\n", float64(duration1)/float64(duration2))
	fmt.Println()
}

func testAddressAutocompleteCache() {
	fmt.Println("Testing AddressAutocompleteHandler cache...")

	// Make the first request
	start := time.Now()
	resp1, err := http.Get("http://localhost:8080/external/autocomplete-address?q=1600+Amphitheatre+Parkway")
	if err != nil {
		fmt.Printf("Error making first request: %v\n", err)
		return
	}
	defer resp1.Body.Close()
	duration1 := time.Since(start)

	// Make the second request (should be cached)
	start = time.Now()
	resp2, err := http.Get("http://localhost:8080/external/autocomplete-address?q=1600+Amphitheatre+Parkway")
	if err != nil {
		fmt.Printf("Error making second request: %v\n", err)
		return
	}
	defer resp2.Body.Close()
	duration2 := time.Since(start)

	// Compare response times
	fmt.Printf("First request: %v\n", duration1)
	fmt.Printf("Second request: %v\n", duration2)
	fmt.Printf("Cache speedup: %.2fx\n", float64(duration1)/float64(duration2))
	fmt.Println()
}
