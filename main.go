package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strconv"
)

var changes = regexp.MustCompile(`Hello`)

func newSingleHostReverseProxy(url *url.URL) *httputil.ReverseProxy {
	rp := httputil.NewSingleHostReverseProxy(url)
	oldDirector := rp.Director
	rp.Director = func(r *http.Request) {
		oldDirector(r)
		r.Host = url.Host
	}
	rp.ModifyResponse = func(resp *http.Response) (err error) {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		err = resp.Body.Close()
		if err != nil {
			return err
		}
		b = changes.ReplaceAll(b, []byte(`HOWDY`))
		body := ioutil.NopCloser(bytes.NewReader(b))
		resp.Body = body
		resp.ContentLength = int64(len(b))
		resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
		return nil
	}
	return rp
}

func main() {
	serverURL, err := url.Parse("https://hello.natalian.org")
	if err != nil {
		log.Fatal("URL failed to parse")
	}
	reverseProxy := newSingleHostReverseProxy(serverURL)
	http.Handle("/", reverseProxy)
	if err = http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
