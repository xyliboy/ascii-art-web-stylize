// main.go is the entry point of the app.
// It connects each URL path to the right handler and starts the server.
package main

import (
	"ascii-art-web-stylize/internal/handlers"
	"net/http"
)

func main() {
	// Web pages
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/ascii-art", handlers.AsciiArtHandler)

	// Static files like CSS and JS
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// REST API endpoints
	http.HandleFunc("/api/banners", handlers.ApiBannersHandler)
	http.HandleFunc("/api/ascii-art", handlers.ApiAsciiArtHandler)

	// Start the server on port 8080
	http.ListenAndServe(":8080", nil)
}
