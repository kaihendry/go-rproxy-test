package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func newSingleHostReverseProxy(url *url.URL) *httputil.ReverseProxy {
	rp := httputil.NewSingleHostReverseProxy(url)
	oldDirector := rp.Director
	rp.Director = func(r *http.Request) {
		oldDirector(r)
		r.Host = url.Host
		delete(r.Header, "X-Forwarded-For")
	}
	return rp
}

func main() {
	serverURL, err := url.Parse("http://www.bbc.co.uk")
	// serverURL, err := url.Parse("http://localhost:4001")
	// serverURL, err := url.Parse("https://requestb.in/169gk011")
	if err != nil {
		log.Fatal("URL failed to parse")
	}
	reverseProxy := newSingleHostReverseProxy(serverURL)
	http.Handle("/", reverseProxy)
	if err = http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
