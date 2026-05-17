package main

import "testing"

func TestNewServerBuildsConfiguredHTTPServer(t *testing.T) {
	server := newServer(":9090", "templates", "testdata/banners")

	if server == nil {
		t.Fatal("expected server to be created")
	}
	if server.Addr != ":9090" {
		t.Fatalf("expected addr %q, got %q", ":9090", server.Addr)
	}
	if server.Handler == nil {
		t.Fatal("expected server handler to be configured")
	}
}
