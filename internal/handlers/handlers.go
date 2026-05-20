// Package handlers contains the HTTP handlers for the web interface.
package handlers

import (
	"ascii-art-web-stylize/internal/ascii"
	"errors"
	"html/template"
	"net/http"
	"os"
	"strings"
)

// PageData is what we pass to the HTML template on every render.
type PageData struct {
	Result string
	Error  string
	Text   string
	Banner string
}

// renderTemplate sets the status code and renders index.html with the given data.
// Content-Type must be set before WriteHeader — once the header is written it cannot change.
func renderTemplate(w http.ResponseWriter, status int, data PageData) error {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return os.ErrNotExist
		}
		return err
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	return tmpl.Execute(w, data)
}

// HomeHandler serves GET /
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if err := renderTemplate(w, http.StatusOK, PageData{}); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			http.Error(w, "404 Not Found: template missing", http.StatusNotFound)
		} else {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// AsciiArtHandler handles POST /ascii-art.
func AsciiArtHandler(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")
	banner := r.FormValue("banner")

	if text == "" {
		renderTemplate(w, http.StatusBadRequest, PageData{Error: "400 Bad Request - invalid input: please enter some text.", Banner: banner})
		return
	}

	// Fall back to standard if the form sent no banner value.
	if banner == "" {
		banner = "standard"
	}

	bannerMap, err := ascii.LoadBanner(banner)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			renderTemplate(w, http.StatusNotFound, PageData{Error: "404 Not Found: banner \"" + banner + "\" does not exist.", Text: text})
		} else {
			renderTemplate(w, http.StatusInternalServerError, PageData{Error: "500 Internal Server Error.", Text: text})
		}
		return
	}

	result := ascii.Render(text, bannerMap)

	// An empty result means every character was outside the supported ASCII range.
	if strings.TrimSpace(result) == "" {
		renderTemplate(w, http.StatusBadRequest, PageData{Error: "400 Bad Request: input contains no supported characters. Please use standard ASCII characters (A-Z, 0-9, symbols)."})
		return
	}

	if err := renderTemplate(w, http.StatusOK, PageData{Result: result, Text: text, Banner: banner}); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			http.Error(w, "404 Not Found", http.StatusNotFound)
		} else {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		}
	}
}
