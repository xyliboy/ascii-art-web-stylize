package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSetupRoutesHome(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	setupRoutes().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestSetupRoutesNotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/missing-page", nil)
	w := httptest.NewRecorder()

	setupRoutes().ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestSetupRoutesAsciiArtPost(t *testing.T) {
	body := strings.NewReader("text=Hello&banner=standard")
	req := httptest.NewRequest(http.MethodPost, "/ascii-art", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	setupRoutes().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestSetupRoutesAPI(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/banners", nil)
	w := httptest.NewRecorder()

	setupRoutes().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
