package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

// DATA MODELS
// WikiChange maps to the specific fields extracted from the Wikipedia JSON payload
type WikiChange struct {
	Bot       bool   `json:"bot"`
	User      string `json:"user"`
	ServerURL string `json:"server_url"`
}

// AppStats holds all application metrics and a Mutex for concurrency control
type AppStats struct {
	mu            sync.Mutex
	TotalMessages int `json:"total_messages"`
	Bots          int `json:"bots"`
	NonBots       int `json:"non_bots"`
	UniqueUsers   int `json:"distinct_users"`

	// Tracks the total number of distinct server URLs
	ServerCounts int `json:"server_counts"`

	// Maps each distinct server URL to its occurrence count
	ServerList map[string]int `json:"server_list"`

	// Internal map used for deduplicating users; omitted from JSON output
	usersMap map[string]bool `json:"-"`
}

// Global stats instance. Maps must be initialized via make() before assignment.
var stats = AppStats{
	ServerList: make(map[string]int),
	usersMap:   make(map[string]bool),
}

// ENDPOINT HANDLERS
func statsHandler(w http.ResponseWriter, r *http.Request) {
	// Lock the mutex to ensure safe reading of the stats object across threads
	stats.mu.Lock()
	defer stats.mu.Unlock()

	// Map lengths dynamically determine the distinct counts
	stats.UniqueUsers = len(stats.usersMap)
	stats.ServerCounts = len(stats.ServerList)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// BACKGROUND PROCESSOR
func fetchWikiData() {
	req, err := http.NewRequest("GET", "https://stream.wikimedia.org/v2/stream/recentchange", nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "BackendWorkshopApp/1.0 (your@email.com)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Failed to connect to stream: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("Wikipedia Stream Connection Status: %d %s", resp.StatusCode, resp.Status)

	// Scanner reads the continuous stream line-by-line as data is pushed by the server
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "data: ") {
			jsonData := strings.TrimPrefix(line, "data: ")

			var change WikiChange

			// Unmarshal converts the raw JSON string into the WikiChange struct
			if err := json.Unmarshal([]byte(jsonData), &change); err != nil {
				log.Printf("JSON Parse Error: %v\nPayload: %s", err, jsonData)
				continue
			}

			// Lock the mutex before writing to shared maps and counters
			stats.mu.Lock()

			stats.TotalMessages++

			if change.Bot {
				stats.Bots++
			} else {
				stats.NonBots++
			}

			if change.User != "" {
				stats.usersMap[change.User] = true
			}

			if change.ServerURL != "" {
				stats.ServerList[change.ServerURL]++
			}

			stats.mu.Unlock()
		}
	}
}

// MAIN ENTRY POINT
func main() {
	// Execute the stream consumer in a background goroutine to prevent blocking
	go fetchWikiData()

	http.HandleFunc("/", statsHandler)

	port := ":7000"
	fmt.Printf("Server listening on port %s...\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
