package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/2-maybe/crawler/extractor/helpers"
	"github.com/PuerkitoBio/goquery"
)

// const imgSize = 4096 * 60
const outerLinks = 4096 * 60

func networkReqForPage(link string, client *http.Client) (helpers.Page, error) {
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		log.Fatalf("Failed to create req%v\n", err)
	}

	req.Header.Set("User-Agent", "DemoCrawler/1.0")

	res, err := client.Do(req)
	if err != nil {
		return helpers.Page{}, err
	}

	defer res.Body.Close()

	if res.StatusCode > 399 {
		return helpers.Page{}, fmt.Errorf("StatusCode above 399")
	}

	contentType := res.Header.Get("Content-type")

	if !strings.HasPrefix(contentType, "text/html") {
		return helpers.Page{}, fmt.Errorf("invalid content type: %s", contentType)
	}

	doc, docParseErr := goquery.NewDocumentFromReader(res.Body)
	if docParseErr != nil {
		return helpers.Page{}, docParseErr
	}

	uri, err := url.Parse(link)
	if err != nil {
		return helpers.Page{}, err
	}

	newPage := helpers.ExtractBody(doc, uri)
	return newPage, nil
}

func startCrawler(link string) {
	client := &http.Client{
		Transport: &http.Transport{
			MaxConnsPerHost:     10,
			MaxIdleConnsPerHost: 2,
			MaxIdleConns:        500,
			IdleConnTimeout:     time.Second * 30,
			ForceAttemptHTTP2:   true,
			TLSHandshakeTimeout: time.Second * 5,
		},
		Timeout: time.Second * 30,
	}

	rootPage, err := networkReqForPage(link, client)

	if err != nil {
		fmt.Println("Like this throws fatal cause its root \nroot needs to succes")
		fmt.Printf("RECORDED err %v\n", err)
		log.Fatalf("%v\n", err)
	}

	var mu *sync.Mutex = &sync.Mutex{}
	// var wg *sync.WaitGroup = &sync.WaitGroup{}

	concurrecy := configureConcurrency{
		totalDispached: 0,
		pages:          make(map[string]helpers.Page),
		workerSignal:   make(chan struct{}, 4096),
		mu:             mu,
		errSig:         make(chan error, 4096),
	}

	fmt.Println("Hey ", rootPage.HyperLinks)

	for _, value := range rootPage.HyperLinks {
		concurrecy.totalDispached++
		fmt.Println("Start Workers ", concurrecy.totalDispached)
		go bsfCrawl(value, client, &concurrecy, link)
	}

	for {
		select {
		case <-concurrecy.workerSignal:
			concurrecy.mu.Lock()
			if concurrecy.totalDispached == 0 {
				concurrecy.mu.Unlock()
				return
			}
			concurrecy.mu.Unlock()
		case msg := <-concurrecy.errSig:
			fmt.Printf("Received msg %v\n", msg)
		}
	}

}
