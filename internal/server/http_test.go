package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ezex-io/gopkg/logger"
	"github.com/ezex-io/proxier/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	log        = logger.NewSlog(nil)
	addr       = "0.0.0.0:8080"
	proxyRules = []*config.Rules{
		{
			Endpoint:    "/proxy",
			Destination: "https://example.com/proxy",
		},
		{
			Endpoint:    "/foo",
			Destination: "https://example.com/foo",
		},
	}
)

func TestNewServer(t *testing.T) {
	srv, err := NewHTTP(log, addr, proxyRules)

	require.NoError(t, err, "Expected no error while creating server")
	assert.NotNil(t, srv, "httpServer instance should not be nil")
}

func TestRootEndpoint(t *testing.T) {
	srv, err := NewHTTP(log, addr, proxyRules)
	require.NoError(t, err)

	sv, ok := srv.(*httpServer)
	require.True(t, ok)

	testServer := httptest.NewServer(sv.httpServer.Handler)
	defer testServer.Close()

	resp, err := http.Get(fmt.Sprintf("%s/", testServer.URL))
	require.NoError(t, err)
	defer func() {
		_ = resp.Body.Close()
	}()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected HTTP 200 OK")
}

func TestLivezEndpoint(t *testing.T) {
	srv, err := NewHTTP(log, addr, proxyRules)
	require.NoError(t, err)

	sv, ok := srv.(*httpServer)
	require.True(t, ok)

	testServer := httptest.NewServer(sv.httpServer.Handler)
	defer testServer.Close()

	resp, err := http.Get(fmt.Sprintf("%s/livez", testServer.URL))
	require.NoError(t, err)
	defer func() {
		_ = resp.Body.Close()
	}()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected HTTP 200 OK")
}

func TestProxyRoutes(t *testing.T) {
	srv, err := NewHTTP(log, addr, proxyRules)
	require.NoError(t, err)

	sv, ok := srv.(*httpServer)
	require.True(t, ok)

	testServer := httptest.NewServer(sv.httpServer.Handler)
	defer testServer.Close()

	for _, rule := range proxyRules {
		t.Run(fmt.Sprintf("Proxy Route %s", rule), func(t *testing.T) {
			resp, err := http.Get(fmt.Sprintf("%s%s", testServer.URL, rule.Endpoint))
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
	srv, err := NewHTTP(log, addr, proxyRules)
	require.NoError(t, err)

	go srv.Start()

	time.Sleep(1 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Stop(ctx)
}
