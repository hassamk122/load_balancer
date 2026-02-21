package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

const (
	port = 8080
)

var serverPool = &ServerPool{}

// It will only fail if no backends are available
// if a req attempts exceeds more than 3 attempts it shows error for that backend
// other wise it pools to next available peer(backend)
func loadBalancer(res http.ResponseWriter, req *http.Request) {

	attempts := GetAttemptsFromContext(req)
	if attempts > 3 {
		log.Printf("%s(%s) Max attempts reached, terminating\n", req.RemoteAddr, req.URL.Path)
		http.Error(res, "Service not available", http.StatusServiceUnavailable)
		return
	}

	peer := serverPool.GetNextPeer()
	if peer != nil {
		log.Printf("Request from %s routed to backend %s \n", req.RemoteAddr, peer.URL)
		peer.ReverseProxy.ServeHTTP(res, req)
		return
	}

	http.Error(res, "Service not available", http.StatusServiceUnavailable)
}

// runs every 20 seconds
// checks for backends health
func healthCheck() {
	ticker := time.NewTicker(time.Second * 20)
	defer ticker.Stop()

	for range ticker.C {
		log.Printf("Starting health check ...")
		serverPool.HealthCheck()
		log.Printf("Health check complete.")
	}
}

// iterates over provided urls and creates reverse proxy
// for each setups a reverse proxy
// creates a backend struct to store its info related to its status and proxy
// reverse proxy error handler retries failed request
// marks backend dead if max retries (3) reached
// if no error adds backend to server pool
func createBackends(backendsUrls []string) {
	for _, rawUrl := range backendsUrls {
		u, err := url.Parse(rawUrl)
		if err != nil {
			log.Fatal(err)
		}

		func(u *url.URL) {
			reverseProxy := httputil.NewSingleHostReverseProxy(u)

			backend := &Backend{
				URL:          u,
				Alive:        true,
				ReverseProxy: reverseProxy,
			}

			reverseProxy.ErrorHandler = func(res http.ResponseWriter, req *http.Request, err error) {
				log.Printf("[%s] %s\n", u.Host, err.Error())
				retries := GetRetryFromContext(req)
				if retries < 3 {
					ctx := context.WithValue(req.Context(), Retry, retries+1)
					reverseProxy.ServeHTTP(res, req.WithContext(ctx))
					return
				}

				serverPool.MarkBackendStatus(u, false)

				attempts := GetAttemptsFromContext(req)
				ctx := context.WithValue(req.Context(), Attempts, attempts+1)
				loadBalancer(res, req.WithContext(ctx))
			}

			serverPool.backends = append(serverPool.backends, backend)
		}(u)

	}
}

func main() {
	backendsUrls := []string{
		"http://localhost:6969",
		"http://localhost:4200",
		"http://localhost:8989",
	}

	createBackends(backendsUrls)

	serverPool.HealthCheck()

	go healthCheck()

	server := http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(loadBalancer),
	}

	log.Println("started at port :8080")
	log.Fatal(server.ListenAndServe())

}
