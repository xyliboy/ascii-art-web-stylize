package main

import (
	"fmt"
	"log"
	"net/http"

	"ascii-art-web/internal/banner"
	web "ascii-art-web/internal/http"
	"ascii-art-web/internal/render"
)

func main() {
	loader := banner.NewLoader("testdata/banners")
	renderer := render.NewService(loader)
	handler := web.NewHandler("templates", renderer)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	log.Println("ascii-art-web listening on http://localhost:8080")
	fmt.Println("Open in your browser:")
	fmt.Println("http://localhost:8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
