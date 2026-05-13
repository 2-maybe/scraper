package helpers

import (
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type Page struct {
	PageName   string
	Heading    string
	FirstPara  string
	HyperLinks []string
	ImgAlt     []string
}

func ExtractBody(doc *goquery.Document, url *url.URL) Page {
	header := getHeadingFromHTML(doc)
	paragraph := getFirstParagraphFromHTML(doc)
	hyperLinks := hefExtractor(url, doc)
	altImgUrl := imageExtractor(url, doc)

	currPage := Page{
		PageName:   (url.Hostname() + url.Path),
		Heading:    header,
		FirstPara:  paragraph,
		HyperLinks: hyperLinks,
		ImgAlt:     altImgUrl,
	}
	return currPage
}
