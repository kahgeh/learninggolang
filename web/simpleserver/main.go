package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(w, "<html><body>Hello world</body></html>")
	})
	http.ListenAndServe(":8080", nil)
}
