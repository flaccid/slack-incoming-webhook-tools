package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

type Proxy struct {
	proxy *httputil.ReverseProxy
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.proxy.ServeHTTP(w, r)
}

func logRequest(r *http.Request) {
	// we only log in debug mode due to exposure of token in request uri
	log.WithFields(log.Fields{
	  "method": r.Method,
	}).Debug(r.URL.String())
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"healthy": true}`)
	logRequest(r)
}

func Serve(targetUrl string, listenPort string, debug bool) {
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	proxyHandler := new(Proxy)
	u, err := url.Parse(targetUrl)
	if err != nil {
		log.Errorf("error parsing url")
		os.Exit(1)
	}

	proxyHandler.proxy = &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.Host = u.Host
			r.URL.Scheme = u.Scheme
			r.URL.Host = u.Host
			r.URL.Path = u.Path
			logRequest(r)
		},
	}

	http.Handle("/", proxyHandler)
	http.HandleFunc("/health", healthCheckHandler)
	log.Info("starting http server, listening on :" + listenPort)
	log.Debug("proxying requests to ", targetUrl)
	log.Fatal(http.ListenAndServe(":"+listenPort, nil))
}
