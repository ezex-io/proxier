package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ezex-io/proxier/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var log = slog.Default()

var proxyRules = []*config.ProxyRule{
	{Endpoint: "/test", DestinationURL: "https://example.com"},
	{Endpoint: "/mock", DestinationURL: "https://mockapi.com"},
}

var serverConfig = &config.ServerConfig{
	Host:       "127.0.0.1",
	ListenPort: "8080",
}

func TestNewServer(t *testing.T) {
	srv, err := NewHTTP(log, serverConfig, proxyRules)

	require.NoError(t, err, "Expected no error while creating server")
	assert.NotNil(t, srv, "httpServer instance should not be nil")
}

func TestRootEndpoint(t *testing.T) {
	srv, err := NewHTTP(log, serverConfig, proxyRules)
	require.NoError(t, err)

	sv, ok := srv.(*httpServer)
	require.True(t, ok)

	testServer := httptest.NewServer(sv.httpServer.Handler)
	defer testServer.Close()

	url := fmt.Sprintf("%s/", testServer.URL)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, url, http.NoBody)
	require.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() {
		_ = resp.Body.Close()
	}()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected HTTP 200 OK")
}

func TestLivezEndpoint(t *testing.T) {
	srv, err := NewHTTP(log, serverConfig, proxyRules)
	require.NoError(t, err)

	sv, ok := srv.(*httpServer)
	require.True(t, ok)

	testServer := httptest.NewServer(sv.httpServer.Handler)
	defer testServer.Close()

	url := fmt.Sprintf("%s/livez", testServer.URL)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, url, http.NoBody)
	require.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() {
		_ = resp.Body.Close()
	}()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected HTTP 200 OK")
}

func TestProxyRoutes(t *testing.T) {
	srv, err := NewHTTP(log, serverConfig, proxyRules)
	require.NoError(t, err)

	sv, ok := srv.(*httpServer)
	require.True(t, ok)

	testServer := httptest.NewServer(sv.httpServer.Handler)
	defer testServer.Close()

	for _, rule := range proxyRules {
		t.Run(fmt.Sprintf("Proxy Route %s", rule.Endpoint), func(t *testing.T) {
			url := fmt.Sprintf("%s%s", testServer.URL, rule.Endpoint)
			req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, url, http.NoBody)
			require.NoError(t, err)

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer func() {
				_ = resp.Body.Close()
			}()

			assert.Contains(t, []int{
				http.StatusOK, http.StatusBadGateway,
				http.StatusNotFound, http.StatusForbidden,
			}, resp.StatusCode,
				"Unexpected response status")
		})
	}
}

func TestServerStartStop(t *testing.T) {
	srv, err := NewHTTP(log, serverConfig, proxyRules)
	require.NoError(t, err)

	go srv.Start()

	time.Sleep(1 * time.Second)

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	srv.Stop(ctx)
}
