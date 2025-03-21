package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/valyala/fasthttp"
)

func HTTPHandler(endpoint string, destination string) (string, http.HandlerFunc, error) {
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

func FastHTTPHandler(endpoint string, destination string) (string, fasthttp.RequestHandler, error) {
	targetURL, err := url.Parse(destination)
	if err != nil {
		return "", nil, fmt.Errorf("invalid destination URL %s: %w", destination, err)
	}

	client := &fasthttp.HostClient{
		Addr:  targetURL.Host,
		IsTLS: targetURL.Scheme == "https",
	}

	handler := func(ctx *fasthttp.RequestCtx) {
		originalPath := string(ctx.Path())
		if !strings.HasPrefix(originalPath, endpoint) {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.SetBodyString("Invalid endpoint")

			return
		}

		trimmedPath := strings.TrimPrefix(originalPath, endpoint)
		if !strings.HasPrefix(trimmedPath, "/") {
			trimmedPath = "/" + trimmedPath
		}

		fullProxyPath := targetURL.Path + trimmedPath
		log.Printf("[Proxy] %s -> %s%s", originalPath, targetURL.String(), trimmedPath)

		req := &ctx.Request
		resp := &ctx.Response

		req.URI().SetScheme(targetURL.Scheme)
		req.URI().SetHost(targetURL.Host)
		req.URI().SetPath(fullProxyPath)

		if err := client.Do(req, resp); err != nil {
			ctx.SetStatusCode(fasthttp.StatusBadGateway)
			ctx.SetBodyString("Proxy error: " + err.Error())
		}
	}

	return endpoint, handler, nil
}
