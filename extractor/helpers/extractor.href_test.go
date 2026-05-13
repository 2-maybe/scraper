package helpers

import (
	"bytes"
	"log"
	"net/url"
	"reflect"
	"testing"

	"github.com/2-maybe/crawler/extractor/matcher"
	"github.com/PuerkitoBio/goquery"
)

func TestHref_Extractor(t *testing.T) {
	matcher.PreCompile()

	transFromInput := func(input []byte) *goquery.Document {
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(input))
		if err != nil {
			log.Fatalf("Err starting test %v\n", err)
		}
		return doc
	}

	tests := []struct {
		name     string
		input    *goquery.Document
		expected []string
		basePath string
	}{
		{
			name: "Vlid resurlt input",
			input: transFromInput([]byte(`
				<!DOCTYPE html>
				<html lang="en">

				<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>Document</title>
				</head>

				<body>
				<h1>Hello This is official apple page</h1>

				<li>
				<ol>TOP</ol>
				<ol>OF </ol>
				<ol>THE</ol>
				<ol>GAME</ol>
				<ol>CAPITALISM</ol>
				</li>

				<a href="/in/mac/">Mac Page</a>
				<a href="/in/watch/">Watch Page</a>
				<a href="/in/services/">Services Page</a>
				<a href="/in/apple-arcade/">Services Page</a>
				</body>

				</html>
				`)),
			expected: []string{
				"https://www.apple.com/in/mac/",
				"https://www.apple.com/in/watch/",
				"https://www.apple.com/in/services/",
				"https://www.apple.com/in/apple-arcade/",
			},

			basePath: "https://www.apple.com",
		},
		{
			name: "No Src Present",
			input: transFromInput([]byte(`
				<!DOCTYPE html>
				<html lang="en">

				<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>Document</title>
				</head>

				<body>
				<h1>Hello</h1>
				<h1>World</h1>
				<h1>How</h1>
				<h1>Are</h1>
				<h1>You</h1>
				<h1>Are</h1>
				<h1>You</h1>
				<h1>Fine</h1>
				</body>

				</html>
				`)),
			expected: nil,
			basePath: "https://www.apple.com",
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(c *testing.T) {
			p, err := url.Parse(v.basePath)
			if err != nil {
				c.Fatalf("Internal Err on test")
			}
			data := hefExtractor(p, v.input)
			if !reflect.DeepEqual(data, v.expected) {
				c.Errorf("Extected %s\n got %v\n", v.expected, data)
			}
		})
	}
}
