package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":9999", nil))
}

func indexHandler(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(writer, "URL.Path = %q\n", req.URL.Path)
}

func helloHandler(writer http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		fmt.Fprintf(writer, "Header[%q] = %q\n", k, v)
	}
}
