package config

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadFromEnv(t *testing.T) {
	require.NoError(t, os.Setenv("PROXIER_ADDRESS", "127.0.0.1:8081"))
	require.NoError(t, os.Setenv("PROXIER_ENABLE_FASTHTTP", "true"))
	require.NoError(t, os.Setenv("PROXIER_RULES", "[{\"endpoint\":\"/foo\",\"destination\":"+
		"\"https://httpbin.org/get\"}, {\"endpoint\":\"/bar\",\"destination\":\"https://google.com\"}]"))

	cfg := LoadFromEnv()
	require.NotNil(t, cfg)
	require.NoError(t, cfg.BasicCheck())

	assert.Equal(t, "127.0.0.1:8081", cfg.Address)
	assert.Equal(t, true, cfg.EnableFastHTTP)

	err := json.Unmarshal([]byte(cfg.RawRules), &cfg.ParsedRules)
	require.NoError(t, err)

	assert.Len(t, cfg.ParsedRules, 2)
}
