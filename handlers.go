package main

import (
	"html/template"
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
