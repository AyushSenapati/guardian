package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func newRevesedProxy(definition *Definition, lb LB) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Director: createDirector(definition, lb),
	}
}

func createDirector(definition *Definition, lb LB) func(*http.Request) {
	return func(req *http.Request) {
		orgReqURI := req.URL.Path // org req URI for logging

		electedHost, err := lb.Elect(definition.Upstreams.Targets)
		if err != nil {
			log.Println("error: failed to elect an upstream")
			return
		}
		target, _ := url.Parse(electedHost)
		path := target.Path + req.URL.Path

		if definition.StripPath {
			listenPath := matcher.ReplaceAllString(definition.ListenPath, "")
			path = strings.Replace(path, listenPath, "", 1)
		}

		req.URL.Path = path
		req.URL.Host = target.Host
		req.URL.Scheme = target.Scheme

		// if preserve host is set don't modify the request host with
		// target host, which could lead to SSL host verification failure
		if !definition.PreserveHost {
			req.Host = target.Host
		}

		logStruct := struct {
			OriginalReqURI string
			UpstreamHost   string
			UpstreamURI    string
		}{
			OriginalReqURI: orgReqURI,
			UpstreamHost:   req.URL.Host,
			UpstreamURI:    req.URL.Path,
		}
		log.Printf("%+v", logStruct)
	}
}
