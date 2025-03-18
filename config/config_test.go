package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTempConfig(t *testing.T, content string) string {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "config_test_*.yaml")
	require.NoError(t, err, "Failed to create temp file")

	_, err = tmpFile.Write([]byte(content))
	require.NoError(t, err, "Failed to write to temp file")

	err = tmpFile.Close()
	require.NoError(t, err, "Failed to close temp file")

	return tmpFile.Name()
}

func TestLoadValidConfig(t *testing.T) {
	yamlContent := `
server:
  host: "127.0.0.1"
  listen_port: "8080"

proxy:
  - endpoint: "/api"
    destination_url: "https://example.com"
  - endpoint: "/test"
    destination_url: "https://test.com"
`
	configFile := createTempConfig(t, yamlContent)
	defer func() {
		_ = os.Remove(configFile)
	}()

	cfg, err := LoadConfig(configFile)
	require.NoError(t, err, "Expected no error while loading valid config")

	assert.Equal(t, "127.0.0.1", cfg.Server.Host, "Host should match")
	assert.Equal(t, "8080", cfg.Server.ListenPort, "ListenPort should match")

	assert.Len(t, cfg.Proxy, 2, "Expected two proxy rules")
	assert.Equal(t, "/api", cfg.Proxy[0].Endpoint)
	assert.Equal(t, "https://example.com", cfg.Proxy[0].DestinationURL)
}

func TestLoadConfig_MissingServer(t *testing.T) {
	yamlContent := `
proxy:
  - endpoint: "/api"
    destination_url: "https://example.com"
`
	configFile := createTempConfig(t, yamlContent)
	defer func() {
		_ = os.Remove(configFile)
	}()

	_, err := LoadConfig(configFile)
	assert.Error(t, err, "Expected error due to missing server configuration")
	assert.Contains(t, err.Error(), "server configuration is missing")
}

func TestLoadConfig_MissingListenPort(t *testing.T) {
	yamlContent := `
server:
  host: "127.0.0.1"

proxy:
  - endpoint: "/api"
    destination_url: "https://example.com"
`
	configFile := createTempConfig(t, yamlContent)
	defer func() {
		_ = os.Remove(configFile)
	}()

	_, err := LoadConfig(configFile)
	assert.Error(t, err, "Expected error due to missing listen_port")
	assert.Contains(t, err.Error(), "server.listen_port cannot be empty")
}

func TestLoadConfig_MissingProxyRules(t *testing.T) {
	yamlContent := `
server:
  host: "127.0.0.1"
  listen_port: "8080"
`
	configFile := createTempConfig(t, yamlContent)
	defer func() {
		_ = os.Remove(configFile)
	}()

	_, err := LoadConfig(configFile)
	assert.Error(t, err, "Expected error due to missing proxy rules")
	assert.Contains(t, err.Error(), "at least one proxy rule must be defined")
}

func TestLoadConfig_InvalidDestinationURL(t *testing.T) {
	yamlContent := `
server:
  host: "127.0.0.1"
  listen_port: "8080"

proxy:
  - endpoint: "/api"
    destination_url: "invalid-url"
`
	configFile := createTempConfig(t, yamlContent)
	defer func() {
		_ = os.Remove(configFile)
	}()

	_, err := LoadConfig(configFile)
	assert.Error(t, err, "Expected error due to invalid destination URL")
	assert.Contains(t, err.Error(), "invalid URL in proxy rule")
}

func TestLoadConfig_DuplicateEndpoints(t *testing.T) {
	yamlContent := `
server:
  host: "127.0.0.1"
  listen_port: "8080"

proxy:
  - endpoint: "/api"
    destination_url: "https://example.com"
  - endpoint: "/api"
    destination_url: "https://duplicate.com"
`
	configFile := createTempConfig(t, yamlContent)
	defer func() {
		_ = os.Remove(configFile)
	}()

	_, err := LoadConfig(configFile)
	assert.Error(t, err, "Expected error due to duplicate endpoints")
	assert.Contains(t, err.Error(), "duplicate proxy endpoint: /api")
}

func TestLoadConfig_EmptyEndpoint(t *testing.T) {
	yamlContent := `
server:
  host: "127.0.0.1"
  listen_port: "8080"

proxy:
  - endpoint: ""
    destination_url: "https://example.com"
`
	configFile := createTempConfig(t, yamlContent)
	defer func() {
		_ = os.Remove(configFile)
	}()

	_, err := LoadConfig(configFile)
	assert.Error(t, err, "Expected error due to empty endpoint")
	assert.Contains(t, err.Error(), "proxy rule endpoint cannot be empty")
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := LoadConfig("non_existent_file.yaml")
	assert.Error(t, err, "Expected error due to missing file")
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestLoadConfig_EmptyFile(t *testing.T) {
	configFile := createTempConfig(t, "")
	defer func() {
		_ = os.Remove(configFile)
	}()

	_, err := LoadConfig(configFile)
	assert.Error(t, err, "Expected error due to empty file")
	assert.Contains(t, err.Error(), "server configuration is missing")
}
