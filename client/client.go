package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

func main() {
	concurrentReqs := 20
	url := "http://localhost:8080"

	var waitGroup sync.WaitGroup
	waitGroup.Add(concurrentReqs)

	for i := 0; i < concurrentReqs; i++ {
		go func(i int) {
			defer waitGroup.Done()
			response, err := http.Get(url)
			if err != nil {
				log.Printf("request %d failed : %v\n", i, err)
				return
			}
			log.Printf("Request %d: Status %d\n", i, response.StatusCode)
			response.Body.Close()
		}(i)
	}

	waitGroup.Wait()
	fmt.Println("All requests done")
}
