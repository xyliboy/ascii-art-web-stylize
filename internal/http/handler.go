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

func NewHandler(templatesDir string, renderer Renderer) *Handler {
	return &Handler{
		templatesDir: templatesDir,
		renderer:     renderer,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/" && r.Method == http.MethodGet:
		h.handleHome(w)
	case r.URL.Path == "/styles.css" && r.Method == http.MethodGet:
		h.handleStyles(w, r)
	case r.URL.Path == "/ascii-art" && r.Method == http.MethodPost:
		h.handleASCIIArt(w, r)
	case r.URL.Path == "/" || r.URL.Path == "/ascii-art":
		h.writeStatusError(w, http.StatusBadRequest, fmt.Sprintf("bad request: method %s is not allowed for %s", r.Method, r.URL.Path))
	case r.URL.Path == "/styles.css":
		h.writeStatusError(w, http.StatusBadRequest, fmt.Sprintf("bad request: method %s is not allowed for %s", r.Method, r.URL.Path))
	default:
		h.writeStatusError(w, http.StatusNotFound, fmt.Sprintf("not found: route %s does not exist", r.URL.Path))
	}
}

func (h *Handler) handleHome(w http.ResponseWriter) {
	data := pageData{Banner: "standard"}
	if err := h.renderPage(w, http.StatusOK, data); err != nil {
		h.writeTemplateError(w, err)
	}
}

func (h *Handler) handleStyles(w http.ResponseWriter, r *http.Request) {
	cssPath := filepath.Join(h.templatesDir, "styles.css")
	if _, err := os.Stat(cssPath); err != nil {
		h.writeStatusError(w, http.StatusNotFound, "not found: stylesheet templates/styles.css is missing")
		return
	}

	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	http.ServeFile(w, r, cssPath)
}

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

func (h *Handler) renderPage(w http.ResponseWriter, statusCode int, data pageData) error {
	page, err := template.ParseFiles(filepath.Join(h.templatesDir, "index.html"))
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)
	return page.Execute(w, data)
}

func (h *Handler) renderStatusPage(w http.ResponseWriter, statusCode int, data statusPageData) error {
	page, err := template.ParseFiles(filepath.Join(h.templatesDir, "error.html"))
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)
	return page.Execute(w, data)
}

func (h *Handler) renderFormError(w http.ResponseWriter, statusCode int, data pageData) {
	if err := h.renderPage(w, statusCode, data); err != nil {
		h.writeTemplateError(w, err)
	}
}

func (h *Handler) writeTemplateError(w http.ResponseWriter, err error) {
	statusCode := http.StatusInternalServerError
	message := "internal server error: the HTML template could not be rendered"
	if errors.Is(err, os.ErrNotExist) {
		statusCode = http.StatusNotFound
		message = "not found: required HTML template is missing"
	}

	h.writeStatusError(w, statusCode, message)
}

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

func isAllowedBanner(banner string) bool {
	_, ok := allowedBanners[banner]
	return ok
}
