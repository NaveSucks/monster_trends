package main

import (
	"log"
	"net/http"
)

func main() {
	// Routes
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/api/offers", offersHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Monster Trends running on http://localhost:3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
