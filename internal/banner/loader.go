package banner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	firstPrintableASCII = 32
	lastPrintableASCII  = 126
	characterHeight     = 8
	blockHeight         = 9
	expectedBannerLines = 855
)

type Font struct {
	glyphs map[rune][]string
}

func (f Font) Glyph(r rune) ([]string, bool) {
	lines, ok := f.glyphs[r]
	return lines, ok
}

type Loader struct {
	baseDir string
}

func NewLoader(baseDir string) *Loader {
	return &Loader{baseDir: baseDir}
}

func (l *Loader) Load(name string) (Font, error) {
	path := filepath.Join(l.baseDir, name+".txt")
	content, err := os.ReadFile(path)
	if err != nil {
		return Font{}, err
	}

	lines := normalizeBannerLines(string(content))
	if len(lines) != expectedBannerLines {
		return Font{}, fmt.Errorf("invalid banner line count for %q: got %d want %d", name, len(lines), expectedBannerLines)
	}

	glyphs := make(map[rune][]string, lastPrintableASCII-firstPrintableASCII+1)
	for index := 0; index <= lastPrintableASCII-firstPrintableASCII; index++ {
		start := index*blockHeight + 1
		end := start + characterHeight

		character := rune(firstPrintableASCII + index)
		glyphs[character] = append([]string(nil), lines[start:end]...)
	}

	return Font{glyphs: glyphs}, nil
}

func normalizeBannerLines(content string) []string {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")

	lines := strings.Split(content, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	return lines
}
