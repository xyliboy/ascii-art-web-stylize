package main

import (
	"ascii-art-web-stylize/internal/handlers"
	"net/http"
)

func main() {
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/ascii-art", handlers.AsciiArtHandler)

	// Serve CSS and JS directly from the templates directory.
	// FileServer handles caching headers automatically.
	http.Handle("/styles.css", http.FileServer(http.Dir("templates")))
	http.Handle("/app.js", http.FileServer(http.Dir("templates")))

	// Bonus REST API
	http.HandleFunc("/api/banners", handlers.ApiBannersHandler)
	http.HandleFunc("/api/ascii-art", handlers.ApiAsciiArtHandler)

	http.ListenAndServe(":8080", nil)
}
