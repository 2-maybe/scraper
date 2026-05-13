package main

import (
	"sync"

	"github.com/2-maybe/crawler/extractor/helpers"
)

type configureConcurrency struct {
	totalDispached int
	workerSignal   chan struct{}
	mu             *sync.Mutex
	pages          map[string]helpers.Page
	errSig         chan error
}

func (c *configureConcurrency) markSeen(key string, val helpers.Page) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.pages[key] = val
}

func (c *configureConcurrency) isSeen(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.pages[key]
	return ok
}
