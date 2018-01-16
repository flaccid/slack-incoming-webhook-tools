package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	log "github.com/sirupsen/logrus"
)

type Proxy struct {
	proxy *httputil.ReverseProxy
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.proxy.ServeHTTP(w, r)
}

func Serve(targetUrl string, listenPort string) {
	proxyHandler := new(Proxy)
	u, err := url.Parse(targetUrl)
	if err != nil {
		log.Errorf("error parsing url")
		os.Exit(1)
	}

	proxyHandler.proxy = &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.Host = u.Host
			req.URL.Scheme = u.Scheme
			req.URL.Host = u.Host
			req.URL.Path = u.Path
		},
	}

	http.Handle("/", proxyHandler)
	log.Info("starting http server, listening on :" + listenPort)
	log.Debug("proxying requests to ", targetUrl)
	log.Fatal(http.ListenAndServe(":"+listenPort, proxyHandler))
}
