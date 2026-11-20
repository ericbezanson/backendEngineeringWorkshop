package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test functions must start with "Test" followed by a capital letter
func TestStatsHandler(t *testing.T) {
	// 1. Create a fake HTTP GET request to our endpoint
	req, err := http.NewRequest("GET", "/stats", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// 2. Create a ResponseRecorder
	// This acts as a fake browser/client that captures the handler's response
	rr := httptest.NewRecorder()

	// 3. Call the handler directly, passing in our fake browser and request
	statsHandler(rr, req)

	// 4. Assert the HTTP Status Code is 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// 5. Assert the Content-Type header is application/json
	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("Handler returned wrong content type: got %v want %v", contentType, expectedContentType)
	}

	// 6. Assert the response body is valid JSON
	// We do this by attempting to unmarshal the fake browser's captured body into our struct
	var responseStats AppStats
	if err := json.Unmarshal(rr.Body.Bytes(), &responseStats); err != nil {
		t.Errorf("Handler returned invalid JSON: %v", err)
	}
}
