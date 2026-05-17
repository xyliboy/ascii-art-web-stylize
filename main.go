package main

import (
	"fmt"
	"log"
	"net/http"

	"ascii-art-web/internal/banner"
	web "ascii-art-web/internal/http"
	"ascii-art-web/internal/render"
)

// main builds the HTTP server and starts listening for browser requests.
// It keeps the entrypoint small and delegates setup to a helper for testability.
func main() {
	server := newServer(":8080", "templates", "testdata/banners")

	log.Println("ascii-art-web listening on http://localhost:8080")
	fmt.Println("Open in your browser:")
	fmt.Println("http://localhost:8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

// newServer wires the banner loader, renderer, and HTTP handler into one server.
// It keeps startup configuration in one place so the setup can be tested directly.
func newServer(addr, templatesDir, bannersDir string) *http.Server {
	loader := banner.NewLoader(bannersDir)
	renderer := render.NewService(loader)
	handler := web.NewHandler(templatesDir, renderer)

	return &http.Server{
		Addr:    addr,
		Handler: handler,
	}
}
