package main

import (
	"log"
	"net"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type Backend struct {
	URL          *url.URL
	Alive        bool
	mutex        sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

// Avoiding race conditions by using mutex
// using RLock for checking loc which basically is read lock
// simple lock to block both read and write
func (b *Backend) isAlive() bool {
	b.mutex.RLock()
	alive := b.Alive
	b.mutex.RUnlock()
	return alive
}

func (b *Backend) SetAlive(alive bool) {
	b.mutex.Lock()
	b.Alive = alive
	b.mutex.Unlock()
}

// allows us to recover dead backends or identify them
// we ping backends with fixed intervals to check status
// to ping we try to establish tcp conn if server responsds we mark it as alive
func isBackendAlive(u *url.URL) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		log.Println("site unreachable , error :", err)
		return false
	}
	conn.Close()
	return true
}
