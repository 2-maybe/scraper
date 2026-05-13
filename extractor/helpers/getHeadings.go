package helpers

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getHeadingFromHTML(Html_body *goquery.Document) string {
	head := Html_body.Find("h1").First().Text()
	head = strings.TrimSpace(head)
	return head
}
