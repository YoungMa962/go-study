package main

import (
	"fmt"
	"net-study/proxy"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "Hello Young!")
	})
	_ = http.ListenAndServe(":8080", new(proxy.Proxy))
}
