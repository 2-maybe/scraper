package helpers

import (
	"bytes"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func Test_GetFirstParagraphFromHTMLMainPriority(t *testing.T) {
	inputBody := []byte(`<html><body>i
	<p>Outside paragraph.</p>
	<main>
	<p>Inner paragraph</p>
	</main>
	</body></html>`)
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(inputBody))
	if err != nil {
		t.Fatalf("Err starting test %v\n", err)
	}
	actual := getFirstParagraphFromHTML(doc)
	expected := "Outside paragraph."

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}
