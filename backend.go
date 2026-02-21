package main

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backend struct {
	URL          *url.URL
	Alive        bool
	mutex        sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

// Avoiding race conditions by using mutex
func (b *Backend) isAlive() bool {
	b.mutex.Lock()
	alive := b.Alive
	b.mutex.Unlock()
	return alive
}

func (b *Backend) SetAlive(alive bool) {
	b.mutex.Lock()
	b.Alive = alive
	b.mutex.Unlock()
}

// allows us to recover dead backends or identify them
// we ping backends with fixed intervals to check status
func isBackendAlive(u *url.URL) error {
	return nil
}
