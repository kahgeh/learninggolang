package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	count("Go", "https://golang.org")
	count("Python", "https://www.python.org/")
	time.Sleep(10 * time.Second)
}

func count(name, url string) {
	start := time.Now()
	r, err := http.Get(url)

	if err != nil {
		fmt.Printf("%s: %s", name, err)
		return
	}

	n, _ := io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()
	fmt.Printf("%s %d [%.2fs]\n", name, n, time.Since(start).Seconds())
}
