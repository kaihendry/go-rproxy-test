package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strconv"
)

var changes = regexp.MustCompile(`Trump|Corbyn`)
var plainHttp = regexp.MustCompile(`http:`)

// var scriptTags = regexp.MustCompile(`<script.*</script>`)
var scriptTags = regexp.MustCompile(`<script[^>]*>[\s\S]*?</script>`)

func newSingleHostReverseProxy(url *url.URL) *httputil.ReverseProxy {
	rp := httputil.NewSingleHostReverseProxy(url)
	oldDirector := rp.Director
	rp.Director = func(r *http.Request) {
		oldDirector(r)
		r.Host = url.Host
		delete(r.Header, "X-Forwarded-For")
		log.Println(r.URL.Path)
	}
	rp.ModifyResponse = func(resp *http.Response) (err error) {

		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return fmt.Errorf("gzip.NewReader: %v", err)
		}

		b, err := ioutil.ReadAll(reader)
		// log.Println(string(b))
		if err != nil {
			return err
		}
		err = resp.Body.Close()
		if err != nil {
			return err
		}
		b = changes.ReplaceAll(b, []byte(`Larry`))
		b = plainHttp.ReplaceAll(b, nil)
		b = scriptTags.ReplaceAll(b, nil)

		//		body := ioutil.NopCloser(bytes.NewReader(b))

		var buffer bytes.Buffer

		writer := gzip.NewWriter(&buffer)
		if _, err := writer.Write(b); err != nil {
			return fmt.Errorf("writer.Write: %v", err)
		}
		if err := writer.Close(); err != nil {
			return fmt.Errorf("writer.Close: %v", err)
		}

		resp.Body = ioutil.NopCloser(bytes.NewReader(buffer.Bytes()))

		resp.ContentLength = int64(len(b))
		resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
		return nil
	}
	return rp
}

func main() {
	serverURL, err := url.Parse("http://www.bbc.co.uk")
	if err != nil {
		log.Fatal("URL failed to parse")
	}
	reverseProxy := newSingleHostReverseProxy(serverURL)
	http.Handle("/", reverseProxy)
	if err = http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
