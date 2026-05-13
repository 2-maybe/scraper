package main

import (
	"fmt"
	"net/http"
	"net/url"

	urlsanitizer "github.com/2-maybe/crawler/urlSanitizer"
)

func checkIdenticalUrl(baseUrl string, currUrl string) bool {
	z, _ := url.Parse(baseUrl)
	k, _ := url.Parse(currUrl)
	if z.Hostname() != k.Hostname() {
		return false
	}
	return true
}

// root url is fking sanitized
func bsfCrawl(rootName string, client *http.Client, concurreny *configureConcurrency, baseUrl string) {

	defer func() {
		concurreny.mu.Lock()
		concurreny.totalDispached = (concurreny.totalDispached - 1)
		concurreny.mu.Unlock()
		concurreny.workerSignal <- struct{}{}
	}()

	if !checkIdenticalUrl(baseUrl, rootName) {
		concurreny.errSig <- fmt.Errorf("Non identical base url")
		return
	}

	if concurreny.isSeen(rootName) {
		concurreny.errSig <- fmt.Errorf("Root Url Already seen")
		return
	}

	rootPage, err := networkReqForPage(rootName, client)

	if err != nil {
		concurreny.errSig <- fmt.Errorf("%v\n", err.Error())
		return
	}

	href := make([]string, 0, outerLinks)
	href = append(href, rootPage.HyperLinks...)

	for len(href) > 0 {
		curr := href[0]
		href = href[1:]

		uri, err := urlsanitizer.NormalizeUrl(curr)

		if err != nil {
			concurreny.errSig <- err
			continue
		}

		if !checkIdenticalUrl(baseUrl, uri) {
			concurreny.errSig <- fmt.Errorf("Non identical url %s\n", (baseUrl + " " + uri))
			continue
		}

		if concurreny.isSeen(uri) {
			concurreny.errSig <- fmt.Errorf("Already seen %v\n", uri)
			continue
		}

		newPage, err := networkReqForPage(uri, client)
		if err != nil {
			concurreny.errSig <- err
			continue
		}

		concurreny.markSeen(uri, newPage)

		href = append(href, newPage.HyperLinks...)
		fmt.Println("Curr Page\t", newPage)
	}
}
