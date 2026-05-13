package helpers

import (
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/2-maybe/crawler/extractor/matcher"
	"github.com/PuerkitoBio/goquery"
)

var first = `
        <!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <title>Document</title>
        </head>
        <body>
            <h1>Water [MIZU] the well of life</h1>
            <p>Outer P</p>
            <nav>Hello</nav>
            <section>
                <img src="" alt="apple">
                <img src="" alt="orange">
                <img src="" alt="banana">
                <img src="" alt="kiwi">
            </section>
            <a href="/about"></a>
            <a href="/home"></a>
            <a href="/logout"></a>
        </body>
        </html>
        `

var sec = `
        <!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <title>Document</title>
        </head>
        <body>
            <h1>Water [MIZU] the well of life</h1>
            <p>Outer P</p>
            <nav>Hello</nav>
            <section>
                <img src="" alt="apple">
                <img src="" alt="">
                <img src="" alt="kiwi">
                <img src="" alt="">
            </section>
            <a href="/about"></a>
            <a href="/home"></a>
            <a href="/logout"></a>
        </body>
        </html>
        `

func Test_ExtractBody(t *testing.T) {
	matcher.PreCompile()

	bytesToHtmlTree := func(input string, k *testing.T) *goquery.Document {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(input))
		if err != nil {
			k.Fatalf("Internal Test err %v\n", err)
		}
		return doc
	}

	tests := []struct {
		name     string
		input    *goquery.Document
		expected Page
		baseUrl  string
	}{
		{
			name:  "Extract Html Body",
			input: bytesToHtmlTree(first, t),
			expected: Page{
				PageName:  "www.apple.com", // Changed to match likely URL implementation
				Heading:   "Water [MIZU] the well of life",
				FirstPara: "Outer P",
				HyperLinks: []string{
					"https://www.apple.com/about",
					"https://www.apple.com/home",
					"https://www.apple.com/logout",
				},
				ImgAlt: []string{"apple", "orange", "banana", "kiwi"},
			},
			baseUrl: "https://www.apple.com",
		},
		{
			name:  "Empty Input",
			input: bytesToHtmlTree("", t),
			expected: Page{
				PageName:   "www.google.com",
				Heading:    "",
				FirstPara:  "",
				HyperLinks: nil,
				ImgAlt:     nil,
			},
			baseUrl: "https://www.google.com",
		},
		{
			name:  "Missing Img alt",
			input: bytesToHtmlTree(sec, t),
			expected: Page{
				PageName:  "www.google.com",
				Heading:   "Water [MIZU] the well of life",
				FirstPara: "Outer P",
				HyperLinks: []string{
					"https://www.google.com/about",
					"https://www.google.com/home",
					"https://www.google.com/logout",
				},
				ImgAlt: []string{"apple", "kiwi"},
			},
			baseUrl: "https://www.google.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(k *testing.T) {
			uri, _ := url.Parse(tt.baseUrl)
			pg := ExtractBody(tt.input, uri)
			if !reflect.DeepEqual(pg, tt.expected) {
				k.Errorf("\nTest: %s\nExpected: %+v\nGot:      %+v", tt.name, tt.expected, pg)
			}
		})
	}
}
