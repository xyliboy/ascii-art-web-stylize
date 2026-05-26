package main

import (
	"ascii-art-web-stylize/internal/handlers"
	"fmt"
	"log"
	"net/http"
)

func setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/ascii-art", handlers.AsciiArtHandler)

	// Serve CSS and JS directly from the templates directory.
	// FileServer handles caching headers automatically.
	mux.Handle("/styles.css", http.FileServer(http.Dir("templates")))
	mux.Handle("/app.js", http.FileServer(http.Dir("templates")))

	// Bonus REST API
	mux.HandleFunc("/api/banners", handlers.ApiBannersHandler)
	mux.HandleFunc("/api/ascii-art", handlers.ApiAsciiArtHandler)

	return mux
}

func main() {
	fmt.Println("Open in your browser:")
	fmt.Println("http://localhost:8080")

	if err := http.ListenAndServe(":8080", setupRoutes()); err != nil {
		log.Fatal(err)
	}
}
