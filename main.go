package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

// Get env var or default
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Serve a reverse proxy for a given url
func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

	proxy.ServeHTTP(res, req)
}

// Given a request send it to the appropriate url
func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	proxyUrl := os.Getenv("PROXY_URL")
	serveReverseProxy(proxyUrl, res, req)
}

func main() {
	proxyPort := getEnv("PROXY_PORT", "1337")

	http.HandleFunc("/", handleRequestAndRedirect)
	if err := http.ListenAndServe(proxyPort, nil); err != nil {
		panic(err)
	}
}
