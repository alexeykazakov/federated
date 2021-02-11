package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/alexeykazakov/proxy/pkg/log"
)

var severCount = 0

// These constant is used to define server
const (
	SERVER1 = "https://api.sandbox.x8i5.p1.openshiftapps.com:6443"
	SERVER2 = "https://www.facebook.com"
	SERVER3 = "https://www.yahoo.com"
	PORT    = "6443"
)

// Serve a reverse proxy for a given url
func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	// parse the url
	url, _ := url.Parse(target)

	// create the reverse proxy
	//proxy := httputil.NewSingleHostReverseProxy(url)
	proxy := NewReverseProxy(url)

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}

// Log the typeform payload and redirect url
func logRequestPayload(proxyURL string) {
	log.Info(nil, fmt.Sprintf("proxy_url: %s", proxyURL))
}

// Balance returns one of the servers based using round-robin algorithm
func getProxyURL() string {
	var servers = []string{SERVER1}

	server := servers[severCount]
	severCount++

	// reset the counter and start from the beginning
	if severCount >= len(servers) {
		severCount = 0
	}

	return server
}

// Given a request send it to the appropriate url
func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	url := getProxyURL()

	logRequestPayload(url)

	serveReverseProxy(url, res, req)
}

func main() {
	log.Init("proxy-service")

	// start server
	http.HandleFunc("/", handleRequestAndRedirect)

	//log.Fatal(http.ListenAndServe(":"+PORT, nil))
	////http.HandleFunc("/hello", HelloServer)

	// listen concurrently to allow for graceful shutdown
	srv := &http.Server{Addr: ":" + PORT, Handler: nil}
	go func() {
		//return server.ListenAndServeTLS("server.crt", "server.key")
		if err := srv.ListenAndServe(); err != nil {
			log.Error(nil, err, err.Error())
		}
	}()

	gracefulShutdown(srv, 10*time.Second)
}

func NewReverseProxy(target *url.URL) *httputil.ReverseProxy {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = "https"
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	return &httputil.ReverseProxy{Director: director}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func gracefulShutdown(hs *http.Server, timeout time.Duration) {
	// For a channel used for notification of just one signal value, a buffer of
	// size 1 is sufficient.
	stop := make(chan os.Signal, 1)

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C) or SIGTERM
	// (Ctrl+/). SIGKILL, SIGQUIT will not be caught.
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	sigReceived := <-stop
	log.Infof(nil, "Signal received: %+v", sigReceived.String())

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	log.Infof(nil, "Shutdown with timeout: %s", timeout.String())
	if err := hs.Shutdown(ctx); err != nil {
		log.Errorf(nil, err, "Shutdown error")
	} else {
		log.Info(nil, "Server stopped.")
	}
}
