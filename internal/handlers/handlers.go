// This package contains the web handlers.
// handlers.go deals with the HTML pages the user sees.
package handlers

import (
	"ascii-art-web-stylize/internal/ascii"
	"errors"
	"html/template"
	"net/http"
	"os"
	"strings"
)

// PageData holds what we pass to the HTML template.
// Result is the ASCII art output, Error is any message we want to show.
type PageData struct {
	Result string
	Error  string
}

// renderTemplate loads the HTML page and fills it with data.
// Returns os.ErrNotExist if the template file is missing.
func renderTemplate(w http.ResponseWriter, data PageData) error {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return os.ErrNotExist
		}
		return err
	}
	return tmpl.Execute(w, data)
}

// HomeHandler serves the main page when the user visits /
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if err := renderTemplate(w, PageData{}); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			http.Error(w, "Not Found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// AsciiArtHandler handles the form submission.
// It reads the text and banner the user picked, generates the ASCII art,
// and sends the result back to the page.
func AsciiArtHandler(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")
	banner := r.FormValue("banner")

	// No text means nothing to render.
	if text == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// If the user didn't pick a banner, use standard as default.
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
		w.WriteHeader(http.StatusBadRequest)
		renderTemplate(w, PageData{Error: "Bad Request: input contains no supported characters. Please use standard ASCII characters (A-Z, 0-9, symbols)."})
		return
	}

	if err := renderTemplate(w, PageData{Result: result}); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			http.Error(w, "Not Found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
