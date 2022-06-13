package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
)

type HttpHandler struct {
	overwriteHost *bool
	upstream *url.URL
}

func (h *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var host string
	if *h.overwriteHost{
		host = h.upstream.Host
	} else {
		host = r.Host
	}

	r.Host = host
	r.URL.Scheme = h.upstream.Scheme
	proxy := httputil.NewSingleHostReverseProxy(h.upstream)
	proxy.ModifyResponse = func(r *http.Response) error {

		if strings.Contains(r.Request.URL.Path, "/api/") {

			// Read the response body
			b, _ := ioutil.ReadAll(r.Body)

			// Compile the regex
			var re = regexp.MustCompile(`"created":"(?:.+?)",`)

			// Replace the current created datetime by a fake date
			b = re.ReplaceAllFunc(b, func(s []byte) []byte {
				return []byte(`"created":"2020-01-01",`)
			})
			buf := bytes.NewBuffer(b)

			// Replace the body with our new body
			r.Body = ioutil.NopCloser(buf)
			// Set the content length
			r.Header["Content-Length"] = []string{fmt.Sprint(buf.Len())}
			// Set Proxy Header
			r.Header.Set("X-Proxy", "go-netbox-proxy 1.0")
			return nil
		}

		return nil
	}
	// Serve request
	proxy.ServeHTTP(w, r)

}

func main() {

	// Parse flags
	addr := flag.String("addr", ":8080", "proxy listen address")
	up := flag.String("upstream", "", "upstream http address")
	overwritehost := flag.Bool("overwritehost", false, "overwrite host header")
	flag.Parse()

	// Parse upstream url
	parsedUpstream, err := url.Parse(*up)

	if err != nil {
		log.Fatal(err.Error())
	}

	// Setup the reverse proxy server
	httpHandler := &HttpHandler{}
	httpHandler.upstream = parsedUpstream
	httpHandler.overwriteHost = overwritehost
	http.Handle("/", httpHandler)
	err = http.ListenAndServe(*addr, nil)

	if err != nil {
		log.Fatal(err.Error())
	}
}
