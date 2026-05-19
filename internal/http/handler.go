package web

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"ascii-art-web/internal/render"
)

var allowedBanners = map[string]struct{}{
	"standard":   {},
	"shadow":     {},
	"thinkertoy": {},
}

type Renderer interface {
	Render(text, banner string) (string, error)
}

type Handler struct {
	templatesDir string
	renderer     Renderer
}

type pageData struct {
	Text   string
	Banner string
	Result string
	Error  string
}

type statusPageData struct {
	Code    int
	Title   string
	Message string
}

// NewHandler builds the main HTTP entrypoint for the web application.
// It connects templates with the rendering service behind one router.
func NewHandler(templatesDir string, renderer Renderer) *Handler {
	return &Handler{
		templatesDir: templatesDir,
		renderer:     renderer,
	}
}

// ServeHTTP dispatches each request to the correct route handler.
// It also centralizes method validation and route-level error handling.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		if r.Method != http.MethodGet {
			h.writeStatusError(w, http.StatusBadRequest, fmt.Sprintf("bad request: method %s is not allowed for %s", r.Method, r.URL.Path))
			return
		}
		h.handleHome(w)
	case "/styles.css":
		if r.Method != http.MethodGet {
			h.writeStatusError(w, http.StatusBadRequest, fmt.Sprintf("bad request: method %s is not allowed for %s", r.Method, r.URL.Path))
			return
		}
		h.serveFile(w, r, "styles.css", "text/css; charset=utf-8", "stylesheet")
	case "/app.js":
		if r.Method != http.MethodGet {
			h.writeStatusError(w, http.StatusBadRequest, fmt.Sprintf("bad request: method %s is not allowed for %s", r.Method, r.URL.Path))
			return
		}
		h.serveFile(w, r, "app.js", "application/javascript; charset=utf-8", "script")
	case "/ascii-art":
		if r.Method != http.MethodPost {
			h.writeStatusError(w, http.StatusBadRequest, fmt.Sprintf("bad request: method %s is not allowed for %s", r.Method, r.URL.Path))
			return
		}
		h.handleASCIIArt(w, r)
	default:
		h.writeStatusError(w, http.StatusNotFound, fmt.Sprintf("not found: route %s does not exist", r.URL.Path))
	}
}

// handleHome renders the empty landing page with the default banner selected.
// It serves the main form without invoking the ASCII renderer.
func (h *Handler) handleHome(w http.ResponseWriter) {
	data := pageData{Banner: "standard"}
	if err := h.renderPage(w, http.StatusOK, data); err != nil {
		h.writeTemplateError(w, err)
	}
}

// serveFile serves a static file from the templates directory with the appropriate content type.
// It reports a clear error if the file is missing, including the file type in the message.
func (h *Handler) serveFile(w http.ResponseWriter, r *http.Request, filename, contentType, fileType string) {
	filePath := filepath.Join(h.templatesDir, filename)
	if _, err := os.Stat(filePath); err != nil {
		h.writeStatusError(w, http.StatusNotFound, fmt.Sprintf("not found: %s templates/%s is missing", fileType, filename))
		return
	}

	w.Header().Set("Content-Type", contentType)
	http.ServeFile(w, r, filePath)
}

// handleASCIIArt parses the form, validates inputs, and renders the output page.
// It maps rendering and resource failures to the correct HTTP status codes.
func (h *Handler) handleASCIIArt(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderFormError(w, http.StatusBadRequest, pageData{
			Banner: "standard",
			Error:  "Bad request: the submitted form could not be parsed.",
		})
		return
	}

	text := r.FormValue("text")
	banner := r.FormValue("banner")
	if !isAllowedBanner(banner) {
		h.renderFormError(w, http.StatusBadRequest, pageData{
			Text:   text,
			Banner: "standard",
			Error:  fmt.Sprintf("Bad request: unsupported banner %q. Choose standard, shadow, or thinkertoy.", banner),
		})
		return
	}

	result, err := h.renderer.Render(text, banner)
	if err != nil {
		switch {
		case errors.Is(err, os.ErrNotExist):
			h.renderFormError(w, http.StatusNotFound, pageData{
				Text:   text,
				Banner: banner,
				Error:  fmt.Sprintf("Not found: banner resource for %q could not be loaded.", banner),
			})
		case errors.Is(err, render.ErrUnsupportedCharacter):
			h.renderFormError(w, http.StatusBadRequest, pageData{
				Text:   text,
				Banner: banner,
				Error:  "Bad request: the submitted text contains unsupported characters.",
			})
		default:
			h.renderFormError(w, http.StatusInternalServerError, pageData{
				Text:   text,
				Banner: banner,
				Error:  "Internal server error: failed to render ASCII art.",
			})
		}
		return
	}

	data := pageData{
		Text:   text,
		Banner: banner,
		Result: result,
	}
	if err := h.renderPage(w, http.StatusOK, data); err != nil {
		h.writeTemplateError(w, err)
	}
}

// renderPage executes the main template with page data and a chosen status code.
// It is used for both the empty state and form-related success or error states.
func (h *Handler) renderPage(w http.ResponseWriter, statusCode int, data pageData) error {
	page, err := template.ParseFiles(filepath.Join(h.templatesDir, "index.html"))
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)
	return page.Execute(w, data)
}

// renderStatusPage executes the standalone status template for route-level failures.
// It keeps 404 and 500 pages separate from the form result template.
func (h *Handler) renderStatusPage(w http.ResponseWriter, statusCode int, data statusPageData) error {
	page, err := template.ParseFiles(filepath.Join(h.templatesDir, "error.html"))
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)
	return page.Execute(w, data)
}

// renderFormError re-renders the main page while preserving form feedback.
// It falls back to template error handling if page rendering fails.
func (h *Handler) renderFormError(w http.ResponseWriter, statusCode int, data pageData) {
	if err := h.renderPage(w, statusCode, data); err != nil {
		h.writeTemplateError(w, err)
	}
}

// writeTemplateError translates template failures into stable HTTP responses.
// It falls back to route-style status pages instead of exposing raw template errors.
func (h *Handler) writeTemplateError(w http.ResponseWriter, err error) {
	statusCode := http.StatusInternalServerError
	message := "internal server error: the HTML template could not be rendered"
	if errors.Is(err, os.ErrNotExist) {
		statusCode = http.StatusNotFound
		message = "not found: required HTML template is missing"
	}

	h.writeStatusError(w, statusCode, message)
}

// writeStatusError renders a standalone status page for transport or route errors.
// If that template also fails, it falls back to the standard net/http error writer.
func (h *Handler) writeStatusError(w http.ResponseWriter, statusCode int, message string) {
	data := statusPageData{
		Code:    statusCode,
		Title:   http.StatusText(statusCode),
		Message: message,
	}

	if err := h.renderStatusPage(w, statusCode, data); err != nil {
		http.Error(w, message, statusCode)
	}
}

// isAllowedBanner enforces the small whitelist of supported banner names.
// This keeps route validation explicit before the renderer is called.
func isAllowedBanner(banner string) bool {
	_, ok := allowedBanners[banner]
	return ok
}
