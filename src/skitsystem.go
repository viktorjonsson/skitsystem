package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

var nrOfRequests = 0

func main() {
	http.HandleFunc("/hello", hello)
	log.Print("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func hello(writer http.ResponseWriter, request *http.Request) {
	log.Printf("Received request on %q", html.EscapeString(request.URL.Path))
	nrOfRequests++
	fmt.Fprint(writer, "Hello World!")
	log.Printf("Returning response from %q", html.EscapeString(request.URL.Path))
	log.Printf("Received requests: %v", nrOfRequests)
}
