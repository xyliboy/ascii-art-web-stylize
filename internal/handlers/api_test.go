package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestApiBannersHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/banners", nil)
	w := httptest.NewRecorder()

	ApiBannersHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var banners []string
	if err := json.NewDecoder(w.Body).Decode(&banners); err != nil {
		t.Fatalf("expected valid JSON, got error: %v", err)
	}

	if len(banners) != 3 {
		t.Fatalf("expected 3 banners, got %d", len(banners))
	}
}

func TestApiAsciiArtHandler_ValidInput(t *testing.T) {
	form := url.Values{}
	form.Set("text", "Hello")
	form.Set("banner", "standard")

	req := httptest.NewRequest(http.MethodGet, "/api/ascii-art?text=Hello&banner=standard", nil)
	w := httptest.NewRecorder()

	ApiAsciiArtHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var result map[string]string
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("expected valid JSON, got error: %v", err)
	}

	if strings.TrimSpace(result["result"]) == "" {
		t.Fatal("expected non-empty result")
	}
}

func TestApiAsciiArtHandler_EmptyText(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/ascii-art?text=&banner=standard", nil)
	w := httptest.NewRecorder()

	ApiAsciiArtHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestApiAsciiArtHandler_UnknownBanner(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/ascii-art?text=Hello&banner=invalid", nil)
	w := httptest.NewRecorder()

	ApiAsciiArtHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
