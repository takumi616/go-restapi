package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "REST API sample in Go")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /test", testHandler)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("APP_PORT"), mux))
}
