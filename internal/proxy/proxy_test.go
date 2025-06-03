package proxy

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProxyHandler(t *testing.T) {
	endpoint := "/test"
	destination := "https://example.com"

	_, handler, err := HTTPHandler(endpoint, destination)

	require.NoError(t, err, "Expected no error while creating proxy handler")
	assert.NotNil(t, handler, "Proxy handler should not be nil")
}

func TestNewProxyHandler_InvalidURL(t *testing.T) {
	endpoint := "/invalid"
	destination := "://invalid-url"

	_, _, err := HTTPHandler(endpoint, destination)

	assert.Error(t, err, "Expected error for invalid URL")
}

func TestProxyHandler(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/mockpath", r.URL.Path, "Expected correct proxied path")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Mock response"))
	}))
	defer mockServer.Close()

	endpoint := "/proxy"
	destination := mockServer.URL
	_, handler, err := HTTPHandler(endpoint, destination)
	require.NoError(t, err, "Proxy handler creation should not fail")

	proxyServer := httptest.NewServer(handler)
	defer proxyServer.Close()

	url := fmt.Sprintf("%s/mockpath", proxyServer.URL)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, url, http.NoBody)
	require.NoError(t, err, "Creating request should not fail")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "Proxy request should not fail")
	defer func() {
		_ = resp.Body.Close()
	}()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected 200 OK from upstream server")
}

func TestProxyHandler_Unreachable(t *testing.T) {
	endpoint := "/proxy"
	destination := "http://127.0.0.1:9999" // Unreachable port

	_, handler, err := HTTPHandler(endpoint, destination)
	require.NoError(t, err, "Proxy handler creation should not fail")

	proxyServer := httptest.NewServer(handler)
	defer proxyServer.Close()

	url := fmt.Sprintf("%s/test", proxyServer.URL)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, url, http.NoBody)
	require.NoError(t, err, "Creating request should not fail")

	client := &http.Client{}
	resp, err := client.Do(req)
	defer func() {
		_ = resp.Body.Close()
	}()

	require.NoError(t, err, "Expected no transport error, proxy should return 502")
	assert.NotNil(t, resp, "Expected a response, since proxy should return 502 Bad Gateway")
	assert.Equal(t, http.StatusBadGateway, resp.StatusCode, "Expected HTTP 502 Bad Gateway")
}
