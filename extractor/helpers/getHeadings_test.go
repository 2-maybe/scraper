package helpers

import (
	"bytes"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func Test_GetHeadingFromHTMLBasic(t *testing.T) {
	inputBody := []byte("<html><body><h1>Test Title</h1></body></html>")
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(inputBody))
	if err != nil {
		t.Fatalf("Err Starting test due to parse err %v\n", err)
	}
	actual := getHeadingFromHTML(doc)
	expected := "Test Title"

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}
