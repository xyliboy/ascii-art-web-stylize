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

func NewService(loader *banner.Loader) *Service {
	return &Service{loader: loader}
}

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

func normalizeInput(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	text = strings.ReplaceAll(text, `\n`, "\n")
	return text
}
