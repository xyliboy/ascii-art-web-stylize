package banner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadStandardBanner(t *testing.T) {
	loader := NewLoader(filepath.Join("..", "..", "testdata", "banners"))

	font, err := loader.Load("standard")
	if err != nil {
		t.Fatalf("load standard: %v", err)
	}

	lines, ok := font.Glyph('A')
	if !ok {
		t.Fatalf("expected glyph for A")
	}
	if len(lines) != characterHeight {
		t.Fatalf("expected %d lines, got %d", characterHeight, len(lines))
	}
}

func TestLoadMissingBanner(t *testing.T) {
	loader := NewLoader(filepath.Join("..", "..", "testdata", "banners"))

	_, err := loader.Load("missing")
	if err == nil {
		t.Fatalf("expected missing banner error")
	}
	if !os.IsNotExist(err) {
		t.Fatalf("expected os.ErrNotExist-compatible error, got %v", err)
	}
}

func TestLoadRejectsInvalidStructure(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "broken.txt"), []byte("bad\nbanner\n"), 0o644); err != nil {
		t.Fatalf("write broken banner: %v", err)
	}

	loader := NewLoader(dir)
	_, err := loader.Load("broken")
	if err == nil {
		t.Fatalf("expected invalid banner structure error")
	}
}
