package web

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

var statusPageTemplate = template.Must(template.New("status").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{{.Code}} {{.Title}} - ASCII Art Web</title>
  <link rel="stylesheet" href="/styles.css">
</head>
<body>
  <main class="page-shell status-page-shell">
    <section class="hero status-hero">
      <div class="hero-copy status-main-card">
        <p class="eyebrow">ASCII Art Web</p>
        <p class="status-code">{{.Code}}</p>
        <h1 class="status-title">{{.Title}}</h1>
        <p class="status-message">{{.Message}}</p>
        <a class="status-link" href="/">Back To Home</a>
      </div>
      <aside class="hero-card status-side-card" aria-label="Error summary">
        <p class="panel-kicker">Generator Status</p>
        <h2>{{.Title}}</h2>
        <p class="intro">The app is still running with the same interface, background, and typography as the main generator.</p>
        <p class="intro intro-secondary">Return home to create ASCII art with the standard, shadow, or thinkertoy banner.</p>
      </aside>
    </section>
  </main>
  <script>
    (function () {
      const rootElement = document.documentElement;
      window.addEventListener("pointermove", function (event) {
        const centerX = window.innerWidth / 2;
        const centerY = window.innerHeight / 2;
        const shiftX = ((event.clientX - centerX) / centerX) * 18;
        const shiftY = ((event.clientY - centerY) / centerY) * 18;

        rootElement.style.setProperty("--star-shift-x", shiftX.toFixed(2) + "px");
        rootElement.style.setProperty("--star-shift-y", shiftY.toFixed(2) + "px");
        rootElement.style.setProperty("--star-counter-x", (shiftX * -0.6).toFixed(2) + "px");
        rootElement.style.setProperty("--star-counter-y", (shiftY * -0.6).toFixed(2) + "px");
      });
    }());
  </script>
</body>
</html>`))

func NewHandler(templatesDir string, renderer Renderer) *Handler {
	return &Handler{
		templatesDir: templatesDir,
		renderer:     renderer,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/" && r.Method == http.MethodGet:
		h.handleHome(w, r)
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

func (h *Handler) handleHome(w http.ResponseWriter, _ *http.Request) {
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
		Title:   strings.ToUpper(http.StatusText(statusCode)),
		Message: message,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)
	if err := statusPageTemplate.Execute(w, data); err != nil {
		http.Error(w, message, statusCode)
	}
}

func isAllowedBanner(banner string) bool {
	_, ok := allowedBanners[banner]
	return ok
}
