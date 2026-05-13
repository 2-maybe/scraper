package main

import (
	"fmt"
	"log"
	"os"

	"github.com/2-maybe/crawler/extractor/matcher"
	urlsanitizer "github.com/2-maybe/crawler/urlSanitizer"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Insufficent Args")
	}

	if len(os.Args) > 2 {
		log.Fatalf("Too many args")
	}

	site := os.Args[1]
	url, err := urlsanitizer.NormalizeUrl(site)
	fmt.Println("Curr args ", site, " Curr url ", url)

	if err != nil {
		log.Fatalf("%v\n", err)
	}

	matcher.PreCompile()
	startCrawler(url)
}
