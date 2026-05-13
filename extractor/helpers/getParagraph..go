package helpers

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getFirstParagraphFromHTML(htmlBody *goquery.Document) string {
	outer := htmlBody.Find("p").First().Text()
	inner := htmlBody.Find("main p").First().Text()
	if outer == "" && inner == "" {
		return ""
	}
	if outer != "" {
		outer = strings.TrimSpace(outer)
		return outer
	}

	inner = strings.TrimSpace(inner)
	return inner
}
