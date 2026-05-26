package handlers

import (
	"ascii-art-web-stylize/internal/ascii"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

func init() {
	os.Chdir("../..")
}

func postAsciiArt(t *testing.T, text, banner string) *httptest.ResponseRecorder {
	t.Helper()
	form := url.Values{"text": {text}, "banner": {banner}}
	req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	AsciiArtHandler(w, req)
	return w
}

func mustRender(t *testing.T, text, bannerName string) string {
	t.Helper()
	banner, err := ascii.LoadBanner(bannerName)
	if err != nil {
		t.Fatalf("LoadBanner(%q): %v", bannerName, err)
	}
	return ascii.Render(text, banner)
}

// normalize trims trailing spaces from each line so golden comparisons are
// not sensitive to banner rows that are padded with trailing whitespace.
func normalize(s string) string {
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = strings.TrimRight(l, " ")
	}
	return strings.Join(lines, "\n")
}

// --- HTTP unit tests ---

func TestHomeHandlerOK(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	HomeHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestHomeHandlerNotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/other", nil)
	w := httptest.NewRecorder()
	HomeHandler(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Back to Home") {
		t.Errorf("expected custom 404 page, got %q", w.Body.String())
	}
}

func TestAsciiArtHandlerValidInput(t *testing.T) {
	w := postAsciiArt(t, "Hello", "standard")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAsciiArtHandlerEmptyInput(t *testing.T) {
	w := postAsciiArt(t, "", "standard")
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestAsciiArtHandlerUnknownBanner(t *testing.T) {
	w := postAsciiArt(t, "Hello", "fakebanner")
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "banner &#34;fakebanner&#34; does not exist") {
		t.Errorf("expected custom 404 message, got %q", w.Body.String())
	}
}

// --- Golden tests from docs/GOLDEN_TESTS.md ---

func TestGolden1Standard(t *testing.T) {
	got := normalize(mustRender(t, "{123}\n<Hello> (World)!", "standard"))
	want := normalize("   __                     __\n" +
		"  / /  _   ____    _____  \\ \\\n" +
		" | |  / | |___ \\  |___ /   | |\n" +
		"/ /   | |   __) |   |_ \\    \\ \\\n" +
		"\\ \\   | |  / __/   ___) |   / /\n" +
		" | |  |_| |_____| |____/   | |\n" +
		"  \\_\\                     /_/\n" +
		"\n" +
		"\n" +
		"   __  _    _          _   _          __            __ __          __                 _       _  __    _\n" +
		"  / / | |  | |        | | | |         \\ \\          / / \\ \\        / /                | |     | | \\ \\  | |\n" +
		" / /  | |__| |   ___  | | | |   ___    \\ \\        | |   \\ \\  /\\  / /    ___    _ __  | |   __| |  | | | |\n" +
		"< <   |  __  |  / _ \\ | | | |  / _ \\    > >       | |    \\ \\/  \\/ /    / _ \\  | '__| | |  / _` |  | | | |\n" +
		" \\ \\  | |  | | |  __/ | | | | | (_) |  / /        | |     \\  /\\  /    | (_) | | |    | | | (_| |  | | |_|\n" +
		"  \\_\\ |_|  |_|  \\___| |_| |_|  \\___/  /_/         | |      \\/  \\/      \\___/  |_|    |_|  \\__,_|  | | (_)\n" +
		"                                                   \\_\\                                           /_/\n" +
		"\n")
	if got != want {
		t.Errorf("golden test 1 failed\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestGolden2Standard(t *testing.T) {
	got := normalize(mustRender(t, "123??", "standard"))
	want := normalize("                     ___    ___\n" +
		" _   ____    _____  |__ \\  |__ \\\n" +
		"/ | |___ \\  |___ /     ) |    ) |\n" +
		"| |   __) |   |_ \\    / /    / /\n" +
		"| |  / __/   ___) |  |_|    |_|\n" +
		"|_| |_____| |____/   (_)    (_)\n" +
		"\n" +
		"\n")
	if got != want {
		t.Errorf("golden test 2 failed\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestGolden3Shadow(t *testing.T) {
	got := normalize(mustRender(t, "$% \"=", "shadow"))
	want := normalize("                        _|  _|\n" +
		"  _|   _|_|    _|       _|  _|\n" +
		"_|_|_| _|_|  _|                _|_|_|_|_|\n" +
		"_|_|       _|\n" +
		"  _|_|   _|  _|_|              _|_|_|_|_|\n" +
		"_|_|_| _|    _|_|\n" +
		"  _|\n" +
		"\n")
	if got != want {
		t.Errorf("golden test 3 failed\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestGolden4Thinkertoy(t *testing.T) {
	got := normalize(mustRender(t, "123 T/fs#R", "thinkertoy"))
	want := normalize("\n" +
		"  0    --  o-o        o-O-o     o  o-o      | |  o--o\n" +
		" /|   o  o    |         |      /   |       -O-O- |   |\n" +
		"o |     /   oo          |     o   -O-  o-o  | |  O-Oo\n" +
		"  |    /      |         |    /     |    \\  -O-O- |  \\\n" +
		"o-o-o o--o o-o          o   o      o   o-o  | |  o   o\n" +
		"\n" +
		"\n")
	if got != want {
		t.Errorf("golden test 4 failed\ngot:\n%s\nwant:\n%s", got, want)
	}
}
