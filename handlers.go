package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// Define a constant for the fallback URL (the Docker service name)
const defaultAPIURL = "http://monster-scrape:8080/offers"

var templates = template.Must(template.ParseGlob("templates/*.html"))

// Global or package-level variable to hold the resolved URL
var apiURL string

func init() {
	_ = godotenv.Load("config.env")
	// Read the API_URL environment variable.
	// This is set in your docker-compose.yml for production/deployment.
	resolvedURL := os.Getenv("API_URL")

	// If the variable is not set (e.g., during local development without Compose),
	// fall back to the default service URL.
	if resolvedURL == "" {
		apiURL = defaultAPIURL
	} else {
		// Use the URL provided by the environment variable
		apiURL = resolvedURL
	}

	log.Printf("Trends service will target API at: %s\n", apiURL)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Monster Energy Trends",
	}
	err := templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Println("Error rendering template:", err)
		http.Error(w, "Internal Server Error", 500)
	}
}

// Offer struct matches the scraper JSON
type Offer struct {
	Discounter string `json:"discounter"`
	Price      string `json:"price"`
	Date       string `json:"date"`
}

// Proxy handler: fetch offers from scraper and return them
func offersHandler(w http.ResponseWriter, r *http.Request) {
	// Use the resolved API URL from the init function
	resp, err := http.Get(apiURL)

	if err != nil {
		log.Printf("Error fetching offers from %s: %v\n", apiURL, err)
		http.Error(w, "Failed to fetch offers", 500)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading offers response:", err)
		http.Error(w, "Failed to read offers", 500)
		return
	}

	// validate JSON
	var offers []Offer
	if err := json.Unmarshal(body, &offers); err != nil {
		log.Println("Error parsing offers JSON:", err)
		http.Error(w, "Invalid offers JSON", 500)
		return
	}

	// return as JSON to browser
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(offers)
}
