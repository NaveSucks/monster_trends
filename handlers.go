package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

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
	resp, err := http.Get("http://localhost:8080/offers")
	if err != nil {
		log.Println("Error fetching offers:", err)
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
