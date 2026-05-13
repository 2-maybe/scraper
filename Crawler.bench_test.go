package main

import (
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/2-maybe/crawler/extractor/helpers"
)

// ── helpers ──────────────────────────────────────────────────────────────────

func newConcurrency() *configureConcurrency {
	return &configureConcurrency{
		totalDispached: 0,
		pages:          make(map[string]helpers.Page),
		workerSignal:   make(chan struct{}, 1024),
		mu:             &sync.Mutex{},
		errSig:         make(chan error, 1024),
	}
}

func seedPages(c *configureConcurrency, n int) {
	for i := 0; i < n; i++ {
		key := fmt.Sprintf("https://example.com/page-%d", i)
		c.pages[key] = helpers.Page{}
	}
}

// ── checkIdenticalUrl ────────────────────────────────────────────────────────

func BenchmarkCheckIdenticalUrl_Same(b *testing.B) {
	base := "https://askhimalaya.com/tours"
	curr := "https://askhimalaya.com/day1"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		checkIdenticalUrl(base, curr)
	}
}

func BenchmarkCheckIdenticalUrl_Diff(b *testing.B) {
	base := "https://askhimalaya.com/tours"
	curr := "https://external-site.com/page"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		checkIdenticalUrl(base, curr)
	}
}

// ── isSeen / markSeen ────────────────────────────────────────────────────────

// Cold map — key is never present.
func BenchmarkIsSeen_Miss(b *testing.B) {
	c := newConcurrency()
	seedPages(c, 10_000)
	probe := "https://askhimalaya.com/never-seen"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.isSeen(probe)
	}
}

// Warm map — key is always present.
func BenchmarkIsSeen_Hit(b *testing.B) {
	c := newConcurrency()
	seedPages(c, 10_000)
	probe := "https://example.com/page-5000"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.isSeen(probe)
	}
}

func BenchmarkMarkSeen(b *testing.B) {
	c := newConcurrency()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("https://example.com/page-%d", i)
		c.markSeen(key, helpers.Page{})
	}
}

// ── contended isSeen (parallel readers) ──────────────────────────────────────

func BenchmarkIsSeen_Parallel(b *testing.B) {
	c := newConcurrency()
	seedPages(c, 10_000)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("https://example.com/page-%d", i%10_000)
			c.isSeen(key)
			i++
		}
	})
}

// ── contended markSeen (parallel writers) ────────────────────────────────────

func BenchmarkMarkSeen_Parallel(b *testing.B) {
	c := newConcurrency()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("https://example.com/page-%d", i)
			c.markSeen(key, helpers.Page{})
			i++
		}
	})
}

// ── mixed read/write (realistic goroutine pattern) ───────────────────────────

func BenchmarkConcurrency_MixedReadWrite(b *testing.B) {
	c := newConcurrency()
	seedPages(c, 1_000)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("https://example.com/page-%d", i%2_000)
			if !c.isSeen(key) {
				c.markSeen(key, helpers.Page{})
			}
			i++
		}
	})
}

// ── networkReqForPage (integration, skipped in unit runs) ────────────────────
// Run with:  go test -bench=BenchmarkNetworkReq -benchtime=20x -tags integration

func BenchmarkNetworkReq(b *testing.B) {
	client := &http.Client{}
	target := "https://askhimalaya.com/"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := networkReqForPage(target, client)
		if err != nil {
			b.Logf("request error (expected on offline runs): %v", err)
		}
	}
}
