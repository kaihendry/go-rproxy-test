package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewSingleHostReverseProxy(url *url.URL) *httputil.ReverseProxy {
	rp := httputil.NewSingleHostReverseProxy(url)
	oldDirector := rp.Director
	rp.Director = func(r *http.Request) {
		oldDirector(r)
		r.Host = url.Host
	}
	return rp
}

func main() {
	serverUrl, err := url.Parse("http://m.bbc.com")
	if err != nil {
		log.Fatal("URL failed to parse")
	}
	reverseProxy := NewSingleHostReverseProxy(serverUrl)
	http.Handle("/", reverseProxy)
	if err = http.ListenAndServe(":9000", nil); err != nil {
		log.Fatal(err)
	}
}
