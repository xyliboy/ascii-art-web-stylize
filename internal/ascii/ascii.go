// This package handles everything related to ASCII art:
// loading banner files and turning text into ASCII art.
package ascii

import (
	"fmt"
	"os"
	"strings"
)

// LoadBanner reads a banner file by name (e.g. "standard") from the banners/
// directory and returns a map from rune to its 8-line ASCII-art representation.
func LoadBanner(name string) (map[rune][8]string, error) {
	path := fmt.Sprintf("banners/%s.txt", name)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open banner %q: %w", name, err)
	}

	// Normalise line endings.
	content := strings.ReplaceAll(string(data), "\r\n", "\n")
	lines := strings.Split(content, "\n")

	// The file starts with one blank line, then every character occupies 8
	// lines followed by one blank separator → block size = 9.
	// Characters range from ASCII 32 (space) to ASCII 126 (~).
	banner := make(map[rune][8]string)
	for i := 0; i <= 94; i++ {
		start := 1 + i*9 // skip the leading blank line
		if start+7 >= len(lines) {
			break
		}
		var block [8]string
		for j := 0; j < 8; j++ {
			block[j] = lines[start+j]
		}
		banner[rune(32+i)] = block
	}
	return banner, nil
}

// Render takes an input string and renders it using the provided banner.
// It splits the input on real newline characters (\n from the web form).
// Characters not present in the banner are silently skipped.
func Render(input string, banner map[rune][8]string) string {
	// Split on real newline characters (Enter key in the web form).
	segments := strings.Split(input, "\n")

	var sb strings.Builder
	for segIdx, segment := range segments {
		if segment != "" {
			// Render 8 rows for this segment.
			for row := 0; row < 8; row++ {
				for _, ch := range segment {
					block, ok := banner[ch]
					if !ok {
						continue
					}
					sb.WriteString(block[row])
				}
				sb.WriteByte('\n')
			}
		}
		// Add a blank separator line between segments (not after the last one).
		if segIdx < len(segments)-1 {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}
