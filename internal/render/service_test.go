package render

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"ascii-art-web/internal/banner"
)

func TestRenderStandardExactSample(t *testing.T) {
	service := newTestService()

	got, err := service.Render("123??", "standard")
	if err != nil {
		t.Fatalf("render standard: %v", err)
	}

	want := "                     ___    ___   \n" +
		" _   ____    _____  |__ \\  |__ \\  \n" +
		"/ | |___ \\  |___ /     ) |    ) | \n" +
		"| |   __) |   |_ \\    / /    / /  \n" +
		"| |  / __/   ___) |  |_|    |_|   \n" +
		"|_| |_____| |____/   (_)    (_)   \n" +
		"                                  \n" +
		"                                  \n"

	if got != want {
		t.Fatalf("unexpected output\nwant:\n%s\ngot:\n%s", want, got)
	}
}

func TestRenderSupportsNewlineSequences(t *testing.T) {
	service := newTestService()

	got, err := service.Render("Hi\\nGo", "standard")
	if err != nil {
		t.Fatalf("render with newline sequence: %v", err)
	}

	if count := strings.Count(got, "\n"); count != 16 {
		t.Fatalf("expected 16 output lines, got %d", count)
	}
}

func TestRenderSupportsActualNewlines(t *testing.T) {
	service := newTestService()

	got, err := service.Render("Hi\n\nGo", "standard")
	if err != nil {
		t.Fatalf("render with actual newlines: %v", err)
	}

	if !strings.Contains(got, "\n\n") {
		t.Fatalf("expected blank line between rendered blocks")
	}
}

func TestRenderEmptyInputReturnsEmptyOutput(t *testing.T) {
	service := newTestService()

	got, err := service.Render("", "standard")
	if err != nil {
		t.Fatalf("render empty input: %v", err)
	}
	if got != "" {
		t.Fatalf("expected empty output, got %q", got)
	}
}

func TestRenderUnsupportedCharacter(t *testing.T) {
	service := newTestService()

	_, err := service.Render("γεια", "standard")
	if err == nil {
		t.Fatalf("expected unsupported character error")
	}
	if !errors.Is(err, ErrUnsupportedCharacter) {
		t.Fatalf("expected ErrUnsupportedCharacter, got %v", err)
	}
}

func newTestService() *Service {
	loader := banner.NewLoader(filepath.Join("..", "..", "testdata", "banners"))
	return NewService(loader)
}
