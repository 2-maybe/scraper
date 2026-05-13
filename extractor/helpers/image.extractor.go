package helpers

import (
	"net/url"

	"github.com/2-maybe/crawler/extractor/matcher"
	"github.com/PuerkitoBio/goquery"
)

func imageExtractor(basePath *url.URL, htmlBody *goquery.Document) []string {
	var urls []string
	var key matcher.CurrMatcherSignature = matcher.Img
	matcher, ok := matcher.PrecompiledMatchers[key]
	if !ok {
		return nil
	}

	htmlBody.FindMatcher(matcher).Each(func(_ int, s *goquery.Selection) {
		val, ok := s.Attr("alt")
		if !ok || val == "" {
			return
		}

		urls = append(urls, val)
	})
	return urls
}
