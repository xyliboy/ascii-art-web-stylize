package ascii

import (
	"os"
	"strings"
	"testing"
)

func init() {
	os.Chdir("../..")
}

// --- LoadBanner ---

func TestLoadBannerStandard(t *testing.T) {
	banner, err := LoadBanner("standard")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(banner) != 95 {
		t.Errorf("expected 95 characters, got %d", len(banner))
	}
}

func TestLoadBannerShadow(t *testing.T) {
	banner, err := LoadBanner("shadow")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(banner) != 95 {
		t.Errorf("expected 95 characters, got %d", len(banner))
	}
}

func TestLoadBannerThinkertoy(t *testing.T) {
	banner, err := LoadBanner("thinkertoy")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(banner) != 95 {
		t.Errorf("expected 95 characters, got %d", len(banner))
	}
}

func TestLoadBannerInvalid(t *testing.T) {
	_, err := LoadBanner("fakebanner")
	if err == nil {
		t.Error("expected error for invalid banner, got nil")
	}
}

func TestLoadBannerEmpty(t *testing.T) {
	_, err := LoadBanner("")
	if err == nil {
		t.Error("expected error for empty banner name, got nil")
	}
}

// --- Render ---

func TestRenderSingleWord(t *testing.T) {
	banner, _ := LoadBanner("standard")
	result := Render("Hi", banner)
	if result == "" {
		t.Error("expected non-empty result, got empty string")
	}
}

func TestRenderProducesEightRows(t *testing.T) {
	banner, _ := LoadBanner("standard")
	result := Render("A", banner)
	// Single segment with no trailing newline after last row → 8 newlines
	lines := strings.Split(result, "\n")
	// Last element is empty string after the trailing newline
	if len(lines) != 9 {
		t.Errorf("expected 9 elements (8 rows + trailing newline), got %d", len(lines))
	}
}

func TestRenderMultiLine(t *testing.T) {
	banner, _ := LoadBanner("standard")
	result := Render("Hi\nWorld", banner)
	if result == "" {
		t.Error("expected non-empty result for multi-line input")
	}
	// Two segments: each 8 rows, separated by a blank line → at least 17 lines
	lines := strings.Split(result, "\n")
	if len(lines) < 17 {
		t.Errorf("expected at least 17 lines for two segments, got %d", len(lines))
	}
}

func TestRenderEmptyString(t *testing.T) {
	banner, _ := LoadBanner("standard")
	result := Render("", banner)
	// Empty input: one empty segment with no rows rendered, one blank separator
	// Result should be a single newline or just the separator
	if strings.TrimSpace(result) != "" {
		t.Errorf("expected blank output for empty input, got %q", result)
	}
}

func TestRenderNewlineOnly(t *testing.T) {
	banner, _ := LoadBanner("standard")
	result := Render("\n", banner)
	// Two empty segments → one blank separator line between them
	if strings.TrimSpace(result) != "" {
		t.Errorf("expected blank output for newline-only input, got %q", result)
	}
}

func TestRenderUnsupportedCharacter(t *testing.T) {
	banner, _ := LoadBanner("standard")
	// Characters outside ASCII 32-126 are silently skipped
	result := Render("A€B", banner)
	// Should still produce output for A and B
	if result == "" {
		t.Error("expected non-empty result when input has unsupported characters mixed with valid ones")
	}
}

func TestRenderAllBanners(t *testing.T) {
	banners := []string{"standard", "shadow", "thinkertoy"}
	for _, name := range banners {
		banner, err := LoadBanner(name)
		if err != nil {
			t.Errorf("banner %q failed to load: %v", name, err)
			continue
		}
		result := Render("Hello", banner)
		if result == "" {
			t.Errorf("banner %q: expected non-empty result, got empty string", name)
		}
	}
}

func TestRenderSpaceCharacter(t *testing.T) {
	banner, _ := LoadBanner("standard")
	result := Render(" ", banner)
	// Space is ASCII 32, should be in the banner and produce 8 rows
	if result == "" {
		t.Error("expected non-empty result for space character")
	}
}

func TestRenderNumbers(t *testing.T) {
	banner, _ := LoadBanner("standard")
	result := Render("123", banner)
	if result == "" {
		t.Error("expected non-empty result for numeric input")
	}
}

func TestRenderSpecialCharacters(t *testing.T) {
	banner, _ := LoadBanner("standard")
	result := Render("!@#", banner)
	if result == "" {
		t.Error("expected non-empty result for special characters in ASCII range")
	}
}
