package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// this test server should be used to create test backends
// you just have to provide a port no to run this backend at that port
// for example in this program i using 3 backends at port 8989, 6969, 4200
// so you need to create 3 test servers with ports 8989, 6969, 4200
func main() {
	port := os.Args[1]
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Response from port %s", port)
	})
	log.Printf("Backend running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
