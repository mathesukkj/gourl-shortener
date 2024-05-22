package tests

import (
	"testing"
	"unicode"

	"github.com/mathesukkj/gourl-shortener/app"
)

func TestShortenUrl(t *testing.T) {
	for i := 0; i < 10000; i++ {
		shortenedUrl := app.ShortenUrl("https://example.com")

		if len(shortenedUrl) < 5 || len(shortenedUrl) > 10 {
			t.Errorf(
				"Shortened URL length out of bounds: got %d, want between 5 and 10",
				len(shortenedUrl),
			)
		}

		for _, ch := range shortenedUrl {
			if !unicode.IsLetter(ch) {
				t.Errorf("Shortened URL contains invalid character: %c", ch)
			}
		}

	}
}
