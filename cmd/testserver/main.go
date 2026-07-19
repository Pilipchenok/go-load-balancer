package main

import (
	"os"
	"net/http"
	"fmt"
	"log"
)

func main() {
	port := os.Args[1]
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "response from backend on port %s\n", port)
	})
	log.Printf("test server on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
