package corsproxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type corsProxyStruct struct {
	reverseproxy            *httputil.ReverseProxy
	listenport, forwardport int
}

// CorsProxy is the entry point to the HTTP proxy
type CorsProxy interface {
	Serve()
}

func (p *corsProxyStruct) Serve() {
	log.Printf("Listening on port %v; forwarding to port %v\n", p.listenport, p.forwardport)

	corsCombiner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request %v %v", r.Method, r.URL)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")
		p.reverseproxy.ServeHTTP(w, r)
	})

	http.ListenAndServe(fmt.Sprintf(":%v", p.listenport), corsCombiner)
}

// NewProxy returns a new HTTP proxy
func NewProxy(listenport, forwardport int) CorsProxy {
	proxyURL, _ := url.Parse(fmt.Sprintf("http://127.0.0.1:%v/", forwardport))
	proxy := httputil.NewSingleHostReverseProxy(proxyURL)
	return &corsProxyStruct{
		reverseproxy: proxy,
		listenport:   listenport,
		forwardport:  forwardport,
	}
}
