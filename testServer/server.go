package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Args[1]
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Response from port %s", port)
	})
	log.Printf("Backend running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
