package web

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ascii-art-web/internal/render"
)

type stubRenderer struct {
	result string
	err    error
	calls  int
	text   string
	banner string
}

func (s *stubRenderer) Render(text, banner string) (string, error) {
	s.calls++
	s.text = text
	s.banner = banner
	if s.err != nil {
		return "", s.err
	}
	return s.result, nil
}

func TestGetHomePage(t *testing.T) {
	templatesDir := writeTemplate(t)
	renderer := &stubRenderer{}
	handler := NewHandler(templatesDir, renderer)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "<form") {
		t.Fatalf("expected response body to include form")
	}
	if renderer.calls != 0 {
		t.Fatalf("renderer should not be called on GET /")
	}
}

func TestPostASCIIArtSuccess(t *testing.T) {
	templatesDir := writeTemplate(t)
	renderer := &stubRenderer{result: "ASCII RESULT"}
	handler := NewHandler(templatesDir, renderer)

	req := httptest.NewRequest(
		http.MethodPost,
		"/ascii-art",
		strings.NewReader("text=hello&banner=standard"),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if renderer.calls != 1 {
		t.Fatalf("expected renderer to be called once, got %d", renderer.calls)
	}
	if renderer.text != "hello" || renderer.banner != "standard" {
		t.Fatalf("renderer received unexpected input: %q %q", renderer.text, renderer.banner)
	}
	if !strings.Contains(rec.Body.String(), "ASCII RESULT") {
		t.Fatalf("expected response to include rendered result")
	}
}

func TestPostASCIIArtAllowsEmptyInput(t *testing.T) {
	templatesDir := writeTemplate(t)
	renderer := &stubRenderer{result: ""}
	handler := NewHandler(templatesDir, renderer)

	req := httptest.NewRequest(
		http.MethodPost,
		"/ascii-art",
		strings.NewReader("text=&banner=shadow"),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if renderer.calls != 1 {
		t.Fatalf("expected renderer to be called once, got %d", renderer.calls)
	}
	if renderer.text != "" || renderer.banner != "shadow" {
		t.Fatalf("renderer received unexpected input: %q %q", renderer.text, renderer.banner)
	}
	if !strings.Contains(rec.Body.String(), "Ready") {
		t.Fatalf("expected empty success page to show ready/info state")
	}
}

func TestPostASCIIArtRejectsInvalidBanner(t *testing.T) {
	templatesDir := writeTemplate(t)
	renderer := &stubRenderer{}
	handler := NewHandler(templatesDir, renderer)

	req := httptest.NewRequest(
		http.MethodPost,
		"/ascii-art",
		strings.NewReader("text=hello&banner=invalid"),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if renderer.calls != 0 {
		t.Fatalf("renderer should not be called for invalid banner")
	}
	if !strings.Contains(rec.Body.String(), `unsupported banner &#34;invalid&#34;`) {
		t.Fatalf("expected detailed banner error message, got %q", rec.Body.String())
	}
}

func TestMissingRouteReturnsNotFound(t *testing.T) {
	templatesDir := writeTemplate(t)
	handler := NewHandler(templatesDir, &stubRenderer{})

	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), ">404<") || !strings.Contains(rec.Body.String(), "Not Found") {
		t.Fatalf("expected standalone 404 error page, got %q", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "route /missing does not exist") {
		t.Fatalf("expected detailed missing route message, got %q", rec.Body.String())
	}
}

func TestWrongMethodReturnsBadRequest(t *testing.T) {
	templatesDir := writeTemplate(t)
	handler := NewHandler(templatesDir, &stubRenderer{})

	req := httptest.NewRequest(http.MethodGet, "/ascii-art", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "method GET is not allowed for /ascii-art") {
		t.Fatalf("expected detailed method error message, got %q", rec.Body.String())
	}
}

func TestPostHomeReturnsBadRequest(t *testing.T) {
	templatesDir := writeTemplate(t)
	handler := NewHandler(templatesDir, &stubRenderer{})

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "method POST is not allowed for /") {
		t.Fatalf("expected detailed method error message, got %q", rec.Body.String())
	}
}

func TestMissingTemplateReturnsNotFound(t *testing.T) {
	handler := NewHandler(t.TempDir(), &stubRenderer{})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "required HTML template is missing") {
		t.Fatalf("expected missing template message, got %q", rec.Body.String())
	}
}

func TestRendererMissingBannerReturnsNotFound(t *testing.T) {
	templatesDir := writeTemplate(t)
	renderer := &stubRenderer{err: os.ErrNotExist}
	handler := NewHandler(templatesDir, renderer)

	req := httptest.NewRequest(
		http.MethodPost,
		"/ascii-art",
		strings.NewReader("text=hello&banner=standard"),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `banner resource for &#34;standard&#34; could not be loaded`) {
		t.Fatalf("expected detailed banner resource error, got %q", rec.Body.String())
	}
}

func TestRendererUnexpectedErrorReturnsInternalServerError(t *testing.T) {
	templatesDir := writeTemplate(t)
	renderer := &stubRenderer{err: errors.New("boom")}
	handler := NewHandler(templatesDir, renderer)

	req := httptest.NewRequest(
		http.MethodPost,
		"/ascii-art",
		strings.NewReader("text=hello&banner=standard"),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "failed to render ASCII art") {
		t.Fatalf("expected detailed internal error message, got %q", rec.Body.String())
	}
}

func TestUnsupportedCharacterReturnsBadRequest(t *testing.T) {
	templatesDir := writeTemplate(t)
	renderer := &stubRenderer{err: render.ErrUnsupportedCharacter}
	handler := NewHandler(templatesDir, renderer)

	req := httptest.NewRequest(
		http.MethodPost,
		"/ascii-art",
		strings.NewReader("text=hello🙂&banner=standard"),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "contains unsupported characters") {
		t.Fatalf("expected unsupported character error message, got %q", rec.Body.String())
	}
}

func TestGetStylesReturnsCSS(t *testing.T) {
	templatesDir := writeTemplate(t)
	handler := NewHandler(templatesDir, &stubRenderer{})

	req := httptest.NewRequest(http.MethodGet, "/styles.css", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); !strings.Contains(got, "text/css") {
		t.Fatalf("expected CSS content type, got %q", got)
	}
}

func TestMissingStylesheetReturnsNotFound(t *testing.T) {
	templatesDir := t.TempDir()
	templatePath := filepath.Join(templatesDir, "index.html")
	if err := os.WriteFile(templatePath, []byte("<html></html>"), 0o644); err != nil {
		t.Fatalf("write template: %v", err)
	}

	handler := NewHandler(templatesDir, &stubRenderer{})

	req := httptest.NewRequest(http.MethodGet, "/styles.css", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "stylesheet templates/styles.css is missing") {
		t.Fatalf("expected detailed stylesheet error, got %q", rec.Body.String())
	}
}

func TestStylesWrongMethodReturnsBadRequest(t *testing.T) {
	templatesDir := writeTemplate(t)
	handler := NewHandler(templatesDir, &stubRenderer{})

	req := httptest.NewRequest(http.MethodPost, "/styles.css", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "method POST is not allowed for /styles.css") {
		t.Fatalf("expected detailed method error message, got %q", rec.Body.String())
	}
}

func TestMissingErrorTemplateFallsBackToHTTPError(t *testing.T) {
	templatesDir := t.TempDir()
	templatePath := filepath.Join(templatesDir, "index.html")
	if err := os.WriteFile(templatePath, []byte("<html></html>"), 0o644); err != nil {
		t.Fatalf("write template: %v", err)
	}

	handler := NewHandler(templatesDir, &stubRenderer{})
	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "not found: route /missing does not exist") {
		t.Fatalf("expected plain fallback error body, got %q", rec.Body.String())
	}
}

func writeTemplate(t *testing.T) string {
	t.Helper()

	templatesDir := t.TempDir()
	templatePath := filepath.Join(templatesDir, "index.html")
	errorTemplatePath := filepath.Join(templatesDir, "error.html")
	stylesPath := filepath.Join(templatesDir, "styles.css")
	templateBody := `<!DOCTYPE html>
<html>
<body>
<form method="post" action="/ascii-art">
<textarea name="text">{{.Text}}</textarea>
<input type="radio" name="banner" value="standard">
{{if .Error}}<p class="error">{{.Error}}</p>{{end}}
{{if and (not .Result) (not .Error)}}<h2>Ready</h2>{{end}}
<pre>{{.Result}}</pre>
</form>
</body>
</html>`
	errorTemplateBody := `<!DOCTYPE html>
<html>
<body>
<p>{{.Code}}</p>
<h1>{{.Title}}</h1>
<p>{{.Message}}</p>
</body>
</html>`
	stylesBody := "body { font-family: monospace; }\n"

	if err := os.WriteFile(templatePath, []byte(templateBody), 0o644); err != nil {
		t.Fatalf("write template: %v", err)
	}
	if err := os.WriteFile(errorTemplatePath, []byte(errorTemplateBody), 0o644); err != nil {
		t.Fatalf("write error template: %v", err)
	}
	if err := os.WriteFile(stylesPath, []byte(stylesBody), 0o644); err != nil {
		t.Fatalf("write stylesheet: %v", err)
	}

	return templatesDir
}
