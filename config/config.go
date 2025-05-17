package config

import (
	"fmt"
	"strings"

	"github.com/ezex-io/gopkg/env"
)

type Config struct {
	Address        string
	EnableFastHTTP bool
	ProxyRules     []string
}

func LoadFromEnv() *Config {
	return &Config{
		Address:        env.GetEnv[string]("EZEX_PROXIER_ADDRESS", env.WithDefault("0.0.0.0:8080")),
		EnableFastHTTP: env.GetEnv[bool]("EZEX_PROXIER_ENABLE_FASTHTTP", env.WithDefault("false")),
		ProxyRules:     env.GetEnv[[]string]("EZEX_PROXIER_PROXY_RULES"),
	}
}

func (c *Config) BasicCheck() error {
	for _, rule := range c.ProxyRules {
		paths := strings.Split(rule, "|")
		if len(paths) != 2 {
			return fmt.Errorf("invalid proxy rule: %s", rule)
		}

		if paths[0] == "" {
			return fmt.Errorf("invalid proxy endpoint: %s", paths[0])
		}

		if paths[1] == "" {
			return fmt.Errorf("invalid proxy destination address: %s", paths[1])
		}
	}

	return nil
}
