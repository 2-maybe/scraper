package helpers

import (
	"net/url"

	"github.com/2-maybe/crawler/extractor/matcher"
	"github.com/PuerkitoBio/goquery"
)

func hefExtractor(basePath *url.URL, htmlBody *goquery.Document) []string {
	var urls []string

	var key matcher.CurrMatcherSignature = matcher.Href
	match, ok := matcher.PrecompiledMatchers[key]
	if !ok {
		return nil
	}

	htmlBody.FindMatcher(match).Each(func(_ int, s *goquery.Selection) {
		val, ok := s.Attr("href")
		if !ok || val == "" {
			return
		}

		joinedPath, err := url.Parse(val)
		if err != nil {
			return
		}

		newUri := basePath.ResolveReference(joinedPath)

		urls = append(urls, newUri.String())
	})

	return urls
}
