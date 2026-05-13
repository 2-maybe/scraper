package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/2-maybe/crawler/extractor/helpers"
)

// generateSite builds a fake site with `depth` levels, each page linking to
// `width` child pages.  Total pages = 1 + width + width^2 + ... + width^depth.
// Keep width*depth small (e.g. 4x3 = 85 pages) or the benchmark gets slow.
func generateSite(depth, width int) map[string]string {
	pages := make(map[string]string)

	var build func(path string, d int)
	build = func(path string, d int) {
		if d == 0 {
			pages[path] = `<html><body><p>leaf</p></body></html>`
			return
		}
		links := ""
		for i := 0; i < width; i++ {
			child := fmt.Sprintf("%s/%d", path, i)
			links += fmt.Sprintf(`<a href="%s">link</a>`, child)
			build(child, d-1)
		}
		pages[path] = fmt.Sprintf(`<html><body>%s</body></html>`, links)
	}

	build("", depth) // paths like "", "/0", "/0/1", ...
	return pages
}

// newMockServer returns an httptest.Server that serves the pre-built site map.
// Every unknown path returns 404 so the crawler's error path is exercised too.
func newMockServer(pages map[string]string) *httptest.Server {
	mux := http.NewServeMux()
	for path, body := range pages {
		b := body // capture
		p := path
		if p == "" {
			p = "/"
		}
		mux.HandleFunc(p, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprint(w, b)
		})
	}
	return httptest.NewServer(mux)
}

// crawlToCompletion wires up a fresh configureConcurrency, fires goroutines for
// the root page's links, and blocks until all workers are done — same logic as
// startCrawler but without log.Fatal so the benchmark can run repeatedly.
func crawlToCompletion(rootURL string, client *http.Client) int {
	rootPage, err := networkReqForPage(rootURL, client)
	if err != nil {
		return 0
	}

	c := &configureConcurrency{
		totalDispached: 0,
		pages:          make(map[string]helpers.Page),
		workerSignal:   make(chan struct{}, 4096),
		mu:             &sync.Mutex{},
		errSig:         make(chan error, 1<<20), // large enough to never block
	}

	// dedicated drainer — never competes with workerSignal in the main select
	go func() {
		for range c.errSig {
		}
	}()

	for _, link := range rootPage.HyperLinks {
		c.mu.Lock()
		c.totalDispached++
		c.mu.Unlock()
		go bsfCrawl(link, client, c, rootURL)
	}

	c.mu.Lock()
	if c.totalDispached == 0 {
		c.mu.Unlock()
		return 0
	}
	c.mu.Unlock()

	for range c.workerSignal {
		c.mu.Lock()
		n := c.totalDispached
		c.mu.Unlock()
		if n == 0 {
			return len(c.pages)
		}
	}
	return len(c.pages)
}

// sharedClient is reused across sub-benchmarks to avoid TLS/dial overhead.
var sharedClient = &http.Client{
	Transport: &http.Transport{
		MaxConnsPerHost:     50,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     30 * time.Second,
	},
	Timeout: 10 * time.Second,
}

// ── benchmarks ───────────────────────────────────────────────────────────────

// Small site: 1 + 3 + 9 = 13 pages.  Low overhead, high iteration count.
func BenchmarkCrawl_Small(b *testing.B) {
	pages := generateSite(2, 3)
	srv := newMockServer(pages)
	defer srv.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		crawlToCompletion(srv.URL+"/", sharedClient)
	}
}

// Medium site: 1 + 4 + 16 + 64 = 85 pages.  Representative real-world crawl.
func BenchmarkCrawl_Medium(b *testing.B) {
	pages := generateSite(3, 4)
	srv := newMockServer(pages)
	defer srv.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		crawlToCompletion(srv.URL+"/", sharedClient)
	}
}

// Large site: 1 + 5 + 25 + 125 + 625 = 781 pages.  Stresses goroutine fan-out
// and the isSeen map under heavy parallel writes.
func BenchmarkCrawl_Large(b *testing.B) {
	pages := generateSite(4, 5)
	srv := newMockServer(pages)
	defer srv.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		crawlToCompletion(srv.URL+"/", sharedClient)
	}
}

// Throughput in pages/sec — reported via b.ReportMetric.
func BenchmarkCrawl_PagesPerSec(b *testing.B) {
	pages := generateSite(3, 4) // 85 pages
	srv := newMockServer(pages)
	defer srv.Close()

	total := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		total += crawlToCompletion(srv.URL+"/", sharedClient)
	}
	b.ReportMetric(float64(total)/b.Elapsed().Seconds(), "pages/sec")
}

// DuplicateLinks stresses the isSeen hot path specifically: every page links
// back to the root + its children, so most goroutines will bail early on
// the isSeen check rather than doing real work.
func BenchmarkCrawl_DuplicateLinks(b *testing.B) {
	mux := http.NewServeMux()
	// root links to /a, /b, /c and also to itself repeatedly
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<html><body>
			<a href="/">home</a><a href="/">home2</a><a href="/">home3</a>
			<a href="/a">a</a><a href="/b">b</a><a href="/c">c</a>
			<a href="/a">a-dup</a><a href="/b">b-dup</a>
		</body></html>`)
	})
	for _, p := range []string{"/a", "/b", "/c"} {
		pp := p
		mux.HandleFunc(pp, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, `<html><body><a href="/">back</a></body></html>`)
		})
	}
	srv := httptest.NewServer(mux)
	defer srv.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		crawlToCompletion(srv.URL+"/", sharedClient)
	}
}
