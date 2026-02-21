package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

const (
	port = 8080
)

// It will only fail if no backends are available
func loadBalancer(res http.ResponseWriter, req *http.Request) {

	attempts := GetAttemptsFromContext(req)
	if attempts > 3 {
		log.Printf("%s(%s) Max attempts reached, terminating\n", req.RemoteAddr, req.URL.Path)
		http.Error(res, "Service not available", http.StatusServiceUnavailable)
		return
	}

	serverPool := ServerPool{}
	peer := serverPool.GetNextPeer()
	if peer != nil {
		peer.ReverseProxy.ServeHTTP(res, req)
		return
	}

	http.Error(res, "Service not available", http.StatusServiceUnavailable)
}

func main() {
	u, err := url.Parse("http://localhost:8080")
	if err != nil {
		os.Exit(1)
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(u)

	reverseProxy.ErrorHandler = func(res http.ResponseWriter, req *http.Request, err error) {
		log.Printf("[%s] %s\n", u.Host, err.Error())
		retries := GetRetryFromContext(req)
		if retries < 3 {
			select {
			case <-time.After(10 * time.Millisecond):
				ctx := context.WithValue(req.Context(), Retry, retries+1)
				reverseProxy.ServeHTTP(res, req.WithContext(ctx))
			}
			return
		}

		// after 3 attempts mark it as down
		serverPool := ServerPool{}
		serverPool.MarkBackendStatus(u, false)

		// if same req routing  for few attempts with diff backend, increase attempts
		attempts := GetAttemptsFromContext(req)
		log.Printf("%s(%s) Attempting retry %d\n", req.RemoteAddr, req.URL.Path, attempts)
		ctx := context.WithValue(req.Context(), Attempts, attempts+1)
		loadBalancer(res, req.WithContext(ctx))
	}
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(loadBalancer),
	}
	http.Handle("/", reverseProxy)
	http.ListenAndServe(":8080", nil)
}
