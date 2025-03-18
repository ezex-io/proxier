package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func NewProxyHandler(endpoint string, destination string) (string, http.HandlerFunc, error) {
	targetURL, err := url.Parse(destination)
	if err != nil {
		return "", nil, fmt.Errorf("invalid destination URL %s: %w", destination, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	proxy.Director = func(r *http.Request) {
		originalPath := r.URL.Path
		trimmedPath := strings.TrimPrefix(originalPath, endpoint)

		r.URL.Path = targetURL.Path + trimmedPath
		r.URL.Scheme = targetURL.Scheme
		r.URL.Host = targetURL.Host
		r.Host = targetURL.Host

		log.Printf("[Proxy] %s -> %s%s", originalPath, targetURL.String(), trimmedPath)
	}

	return endpoint, proxy.ServeHTTP, nil
}
