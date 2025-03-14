package test

import (
	"link-shortener/internal/short_link_generator"
	"testing"
)

func TestValidateGeneratedShortLinkLength(t *testing.T) {
	const expectedLen = 10

	got := short_link_generator.GenerateShortLink()

	if len(got) != expectedLen {
		t.Errorf("ожидалась длина 10, получено %d", len(got))
	}
}
