package server_test

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ezex-io/proxier/config"
	"github.com/ezex-io/proxier/internal/server"
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
	srv, err := server.New(log, serverConfig, proxyRules)

	require.NoError(t, err, "Expected no error while creating server")
	assert.NotNil(t, srv, "Server instance should not be nil")
}

func TestRootEndpoint(t *testing.T) {
	srv, err := server.New(log, serverConfig, proxyRules)
	require.NoError(t, err)

	testServer := httptest.NewServer(srv.HTTPHandler())
	defer testServer.Close()

	resp, err := http.Get(fmt.Sprintf("%s/", testServer.URL))
	require.NoError(t, err)
	defer func() {
		_ = resp.Body.Close()
	}()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected HTTP 200 OK")
}

func TestLivezEndpoint(t *testing.T) {
	srv, err := server.New(log, serverConfig, proxyRules)
	require.NoError(t, err)

	testServer := httptest.NewServer(srv.HTTPHandler())
	defer testServer.Close()

	resp, err := http.Get(fmt.Sprintf("%s/livez", testServer.URL))
	require.NoError(t, err)
	defer func() {
		_ = resp.Body.Close()
	}()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected HTTP 200 OK")
}

func TestProxyRoutes(t *testing.T) {
	srv, err := server.New(log, serverConfig, proxyRules)
	require.NoError(t, err)

	testServer := httptest.NewServer(srv.HTTPHandler())
	defer testServer.Close()

	for _, rule := range proxyRules {
		t.Run(fmt.Sprintf("Proxy Route %s", rule.Endpoint), func(t *testing.T) {
			resp, err := http.Get(fmt.Sprintf("%s%s", testServer.URL, rule.Endpoint))
			require.NoError(t, err)
			defer func() {
				_ = resp.Body.Close()
			}()

			assert.Contains(t, []int{http.StatusOK, http.StatusBadGateway, http.StatusNotFound}, resp.StatusCode,
				"Unexpected response status")
		})
	}
}

func TestServerStartStop(t *testing.T) {
	srv, err := server.New(log, serverConfig, proxyRules)
	require.NoError(t, err)

	go srv.Start()

	time.Sleep(1 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Stop(ctx)
}
