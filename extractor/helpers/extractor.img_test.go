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

func TestImg_Extractor(t *testing.T) {
	matcher.PreCompile()
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
			name: "Img Test",
			input: transFromInput([]byte(`
				<body>
				<img src="" alt="Thunder"/>
				<img src="" alt="Apple"/>
				<img src="" alt="Cosmos"/>
				<img src="" alt="Earth"/>
				<body/>
				`)),
			expected: []string{"Thunder", "Apple", "Cosmos", "Earth"},
			basePath: "https://www.apple.com/images",
		},
		{
			name: "Img Test",
			input: transFromInput([]byte(`
				<body>
				<img src="" alt="Earth"/>
				<img src="" alt=""/>
				<img src="" alt="Cosmos"/>
				<img src="" alt=""/>
				<body/>
				`)),
			expected: []string{"Earth", "Cosmos"},
			basePath: "https://www.apple.com/images",
		},
		{
			name: "Img Test",
			input: transFromInput([]byte(`
				<body>
				<img src="" alt=""/>
				<img src="" alt=""/>
				<img src="" alt=""/>
				<img src="" alt=""/>
				<body/>
				`)),
			expected: nil,
			basePath: "https://www.apple.com/images",
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(c *testing.T) {
			p, err := url.Parse(v.basePath)
			if err != nil {
				c.Fatalf("Internal Err on test")
			}
			data := imageExtractor(p, v.input)
			if !reflect.DeepEqual(data, v.expected) {
				c.Errorf("Expected %s\n got %v\n", v.expected, data)
			}
		})
	}
}
