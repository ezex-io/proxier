package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadFromEnv(t *testing.T) {
	require.NoError(t, os.Setenv("EZEX_PROXIER_ADDRESS", "127.0.0.1:8081"))
	require.NoError(t, os.Setenv("EZEX_PROXIER_ENABLE_FASTHTTP", "true"))
	require.NoError(t, os.Setenv("EZEX_PROXIER_PROXY_RULES", "/foo|https://httpbin.org/get,/bar|https://google.com"))

	cfg := LoadFromEnv()
	require.NotNil(t, cfg)
	require.NoError(t, cfg.BasicCheck())

	assert.Equal(t, "127.0.0.1:8081", cfg.Address)
	assert.Equal(t, true, cfg.EnableFastHTTP)
	assert.Len(t, cfg.RawRules, 2)
}
