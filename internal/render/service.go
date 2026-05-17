package render

import (
	"errors"
	"fmt"
	"strings"

	"ascii-art-web/internal/banner"
)

var ErrUnsupportedCharacter = errors.New("unsupported character")

type Service struct {
	loader *banner.Loader
}

// NewService creates the ASCII render service around a banner loader.
// The service stays focused on rendering while the loader owns filesystem access.
func NewService(loader *banner.Loader) *Service {
	return &Service{loader: loader}
}

// Render converts input text into ASCII art for the selected banner.
// It normalizes newline forms, validates characters, and builds the final output.
func (s *Service) Render(text, bannerName string) (string, error) {
	font, err := s.loader.Load(bannerName)
	if err != nil {
		return "", err
	}

	normalizedText := normalizeInput(text)
	segments := strings.Split(normalizedText, "\n")
	if len(segments) == 1 && segments[0] == "" {
		return "", nil
	}

	var output strings.Builder
	for _, segment := range segments {
		if segment == "" {
			output.WriteByte('\n')
			continue
		}

		for row := 0; row < 8; row++ {
			for _, r := range segment {
				glyph, ok := font.Glyph(r)
				if !ok {
					return "", fmt.Errorf("%w: %q", ErrUnsupportedCharacter, string(r))
				}
				output.WriteString(glyph[row])
			}
			output.WriteByte('\n')
		}
	}

	return output.String(), nil
}

// normalizeInput aligns browser and CLI-style newline input into one format.
// It also turns literal \n sequences into actual line breaks before rendering.
func normalizeInput(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	text = strings.ReplaceAll(text, `\n`, "\n")
	return text
}
