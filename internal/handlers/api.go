// api.go handles the REST API endpoints.
// These return JSON instead of HTML - useful for other apps or tools.
package handlers

import (
	"ascii-art-web-stylize/internal/ascii"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
)

// ApiBannersHandler returns the list of available banners as JSON.
func ApiBannersHandler(w http.ResponseWriter, r *http.Request) {
	banners := []string{"standard", "shadow", "thinkertoy"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(banners)
}

// ApiAsciiArtHandler reads text and banner from the URL,
// generates the ASCII art, and returns it as JSON.
// Example: /api/ascii-art?text=Hello&banner=standard
func ApiAsciiArtHandler(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")
	banner := r.URL.Query().Get("banner")

	// No text means nothing to render.
	if text == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Fall back to standard if no banner was given.
	if banner == "" {
		banner = "standard"
	}

	bannerMap, err := ascii.LoadBanner(banner)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			http.Error(w, "Not Found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	result := ascii.Render(text, bannerMap)
	if strings.TrimSpace(result) == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"result": result})
}
